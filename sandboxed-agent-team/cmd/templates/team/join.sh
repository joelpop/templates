#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# join.sh — provision this developer's local state for the project.
#
# Scope:
#   - verify the kit is installed on the current branch
#   - stop any existing sandbox
#   - provision SSH material (if origin is an SSH remote)
#   - prompt for a platform API token (if MERGE_METHOD is PR)
#   - create the sandbox and launch the team (delegating to team/create.sh)
#   - record .claude/.last-onboarded
#
# Idempotent: running again discards the local sandbox and rebuilds.

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "${PROJECT_DIR}"

# ── Verify the kit is installed on this branch ──────────────────────────────
if [ ! -f "${PROJECT_DIR}/CLAUDE_TEAM.md" ]; then
    cat >&2 <<'EOF'
Error: the agent team kit is not installed on this branch.

  Either you're on the wrong branch, or the project does not have
  the kit set up yet. Ask your project lead which branch contains
  the kit, check it out, and run team/join.sh from there.
EOF
    exit 1
fi

# ── Read the variables file ─────────────────────────────────────────────────
VAR_FILE="${PROJECT_DIR}/.claude/team-variables.yaml"
get_var() {
    local key="$1"
    [ -f "${VAR_FILE}" ] || { echo ""; return; }
    grep -E "^${key}: " "${VAR_FILE}" \
        | sed -E 's/^[^:]+: *"?([^"]*)"?.*$/\1/' \
        | head -n1
}

MERGE_METHOD="$(get_var MERGE_METHOD)"

# ── Discard any existing sandbox so we rebuild from the current setup ──────
# destroy.sh --yes is a no-op if no sandbox exists, and keeps the
# template image in place so the rebuild below reuses Docker's cache.
if [ -x "${PROJECT_DIR}/team/destroy.sh" ]; then
    "${PROJECT_DIR}/team/destroy.sh" --yes
fi

# ── SSH provisioning (if origin remote uses SSH) ────────────────────────────
REMOTE_URL="$(git remote get-url origin 2>/dev/null || true)"

if [[ "${REMOTE_URL}" == git@* || "${REMOTE_URL}" == ssh://* ]]; then
    echo "=== Provisioning SSH material ==="

    SSH_USER=""; SSH_HOST=""
    if [[ "${REMOTE_URL}" == git@* ]]; then
        SSH_USER="${REMOTE_URL%%@*}"
        REST="${REMOTE_URL#*@}"
        SSH_HOST="${REST%%:*}"
    else
        TMP="${REMOTE_URL#ssh://}"
        if [[ "${TMP}" == *@* ]]; then
            SSH_USER="${TMP%%@*}"
            TMP="${TMP#*@}"
        else
            SSH_USER="git"
        fi
        SSH_HOST="${TMP%%/*}"
    fi

    IDENTITY_FILE=""
    REAL_HOSTNAME=""
    if [ -f "${HOME}/.ssh/config" ]; then
        block=$(awk -v host="${SSH_HOST}" '
            BEGIN { in_block = 0 }
            /^[[:space:]]*Host[[:space:]]/ {
                in_block = 0
                for (i = 2; i <= NF; i++) {
                    if ($i == host) { in_block = 1; break }
                }
                next
            }
            in_block == 1 { print }
        ' "${HOME}/.ssh/config")
        IDENTITY_FILE=$(echo "${block}" | awk 'tolower($1) == "identityfile" { print $2; exit }')
        REAL_HOSTNAME=$(echo "${block}" | awk 'tolower($1) == "hostname" { print $2; exit }')
        IDENTITY_FILE="${IDENTITY_FILE/#\~/${HOME}}"
    fi

    if [ -z "${IDENTITY_FILE}" ]; then
        for candidate in "${HOME}/.ssh/id_ed25519" "${HOME}/.ssh/id_rsa"; do
            if [ -f "${candidate}" ]; then
                IDENTITY_FILE="${candidate}"
                break
            fi
        done
    fi

    if [ -z "${IDENTITY_FILE}" ] || [ ! -f "${IDENTITY_FILE}" ]; then
        echo "  Warning: no SSH key found for ${SSH_HOST}." >&2
        echo "           Sandbox will start, but git over SSH may fail." >&2
    else
        [ -z "${REAL_HOSTNAME}" ] && REAL_HOSTNAME="${SSH_HOST}"
        KEY_NAME="$(basename "${IDENTITY_FILE}")"
        SSH_DIR="${PROJECT_DIR}/.sandbox/.ssh"
        mkdir -p "${SSH_DIR}"
        chmod 700 "${SSH_DIR}"
        cp "${IDENTITY_FILE}" "${SSH_DIR}/${KEY_NAME}"
        chmod 600 "${SSH_DIR}/${KEY_NAME}"
        if [ -f "${IDENTITY_FILE}.pub" ]; then
            cp "${IDENTITY_FILE}.pub" "${SSH_DIR}/${KEY_NAME}.pub"
            chmod 644 "${SSH_DIR}/${KEY_NAME}.pub"
        fi
        cat > "${SSH_DIR}/config" <<EOF_CONFIG
Host ${SSH_HOST}
    HostName ${REAL_HOSTNAME}
    User ${SSH_USER}
    IdentityFile ~/.ssh/${KEY_NAME}
    IdentitiesOnly yes
EOF_CONFIG
        if ssh-keyscan "${REAL_HOSTNAME}" > "${SSH_DIR}/known_hosts" 2>/dev/null \
                && [ -s "${SSH_DIR}/known_hosts" ]; then
            :
        else
            echo "  Warning: ssh-keyscan produced no output for ${REAL_HOSTNAME}." >&2
            echo "           You may need to populate ${SSH_DIR}/known_hosts manually." >&2
        fi
        echo "${IDENTITY_FILE}" > "${PROJECT_DIR}/.sandbox/.ssh.source"
        echo "  SSH material provisioned in .sandbox/.ssh/"
    fi
fi

# ── Platform API token (PR merge method only) ──────────────────────────────
ENV_FILE="${PROJECT_DIR}/.sandbox/.platform-api.env"

if [ "${MERGE_METHOD}" = "PR" ] && [ ! -f "${ENV_FILE}" ]; then
    echo "=== Platform API token ==="

    PLATFORM_TYPE=""; API_URL=""
    case "${REMOTE_URL}" in
        *github.com*)    PLATFORM_TYPE="GITHUB";    API_URL="https://api.github.com" ;;
        *bitbucket.org*) PLATFORM_TYPE="BITBUCKET"; API_URL="https://api.bitbucket.org" ;;
        *gitlab*)        PLATFORM_TYPE="GITLAB";    API_URL="https://gitlab.com/api/v4" ;;
    esac

    PATH_PART="${REMOTE_URL}"
    PATH_PART="${PATH_PART#git@*:}"
    PATH_PART="${PATH_PART#ssh://*/}"
    PATH_PART="${PATH_PART#https://*/}"
    PATH_PART="${PATH_PART#http://*/}"
    PATH_PART="${PATH_PART%.git}"
    OWNER="${PATH_PART%%/*}"
    REPO_NAME="${PATH_PART##*/}"

    if [ -z "${PLATFORM_TYPE}" ]; then
        echo "  Warning: could not detect a supported platform (GitHub/Bitbucket/GitLab)" >&2
        echo "           from origin (${REMOTE_URL}). Skipping token provisioning." >&2
    else
        echo "PR merge method is configured. Need a ${PLATFORM_TYPE} API token."
        case "${PLATFORM_TYPE}" in
            GITHUB)
                echo "  Create a fine-grained PAT at:"
                echo "    https://github.com/settings/tokens?type=beta"
                echo "  Scope it to this repo: Contents:Read, Pull requests:Read+Write."
                ;;
            BITBUCKET)
                echo "  Create an app password at:"
                echo "    https://bitbucket.org/account/settings/app-passwords/"
                echo "  Grant: Repositories:Read, Pull requests:Read+Write."
                ;;
            GITLAB)
                echo "  Create a personal access token at:"
                echo "    https://gitlab.com/-/user_settings/personal_access_tokens"
                echo "  Scope: api."
                ;;
        esac

        read -r -p "Paste the token: " TOKEN
        if [ -z "${TOKEN}" ]; then
            echo "  Skipping: no token provided." >&2
        else
            API_USER=""
            if [ "${PLATFORM_TYPE}" = "BITBUCKET" ]; then
                read -r -p "Bitbucket username (for basic auth): " API_USER
            fi
            mkdir -p "${PROJECT_DIR}/.sandbox"
            umask 077
            cat > "${ENV_FILE}" <<EOF_ENV
# Developer-local. Gitignored. Do NOT commit.
PLATFORM_TYPE=${PLATFORM_TYPE}
PLATFORM_API_URL=${API_URL}
PLATFORM_API_USER=${API_USER}
PLATFORM_API_TOKEN=${TOKEN}
PLATFORM_REPO_OWNER=${OWNER}
PLATFORM_REPO_NAME=${REPO_NAME}
PLATFORM_REPO_WORKSPACE=${OWNER}
PLATFORM_REPO_SLUG=${REPO_NAME}
EOF_ENV
            echo "  Wrote .sandbox/.platform-api.env"
        fi
    fi
fi

# ── Record onboarding timestamp ─────────────────────────────────────────────
# Must happen BEFORE create.sh because create.sh blocks on
# `docker sandbox run`, and the sandboxed Claude Code's Pre-Start
# Check reads this marker on first turn. If we wrote it after
# create.sh returned, the marker wouldn't exist yet when the
# Pre-Start Check ran, and the team would (incorrectly) tell the
# user to re-run onboarding.
mkdir -p "${PROJECT_DIR}/.claude"
date -u +%Y-%m-%dT%H:%M:%SZ > "${PROJECT_DIR}/.claude/.last-onboarded"

# ── Create the sandbox and launch the team ────────────────────────────────
"${PROJECT_DIR}/team/create.sh"

echo ""
echo "=== Join complete ==="
