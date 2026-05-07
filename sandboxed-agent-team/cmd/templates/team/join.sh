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
#   - prompt for a repo-platform API token (if MERGE_METHOD is PR)
#   - create the sandbox and launch the team (delegating to team/create.sh)
#   - record .claude/.last-onboarded
#
# Note: no SSH keys are provisioned inside the sandbox — Docker Sandbox
# blocks outbound port 22, so the sandbox reaches the repo exclusively
# over HTTPS (with url.insteadOf rewriting SSH-style origin URLs on
# the fly; see create.sh). We still parse the host's ~/.ssh/config
# here to resolve SSH host aliases (e.g. `bitbucket-work` →
# `bitbucket.org`) for correct repo-platform detection.
#
# Idempotent: running again discards the local sandbox and rebuilds.

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "${PROJECT_DIR}"

PARENT_DIR="$(basename "$(dirname "$PROJECT_DIR")")"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
SANDBOX_NAME="claude-${PARENT_DIR}-${PROJECT_NAME}"
# Per-project Keychain service name (macOS). Keyed by SANDBOX_NAME
# so every project on this workstation gets its own entry.
KEYCHAIN_SERVICE="agent-team.${SANDBOX_NAME}"

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

# ── Parse origin to resolve SSH host aliases (for platform detection) ──────
# If origin uses an SSH-style URL, resolve any ~/.ssh/config HostName
# alias (e.g. `bitbucket-work` → `bitbucket.org`) so the repo-platform
# detection below can match against the real host. We do NOT copy SSH
# keys into the sandbox — see the note at the top of this file.
REMOTE_URL="$(git remote get-url origin 2>/dev/null || true)"
SSH_HOST=""
REAL_HOSTNAME=""
if [[ "${REMOTE_URL}" == git@* ]]; then
    REST="${REMOTE_URL#*@}"
    SSH_HOST="${REST%%:*}"
elif [[ "${REMOTE_URL}" == ssh://* ]]; then
    TMP="${REMOTE_URL#ssh://}"
    TMP="${TMP#*@}"
    SSH_HOST="${TMP%%/*}"
fi
if [ -n "${SSH_HOST}" ] && [ -f "${HOME}/.ssh/config" ]; then
    REAL_HOSTNAME=$(awk -v host="${SSH_HOST}" '
        BEGIN { in_block = 0 }
        /^[[:space:]]*Host[[:space:]]/ {
            in_block = 0
            for (i = 2; i <= NF; i++) {
                if ($i == host) { in_block = 1; break }
            }
            next
        }
        in_block == 1 && tolower($1) == "hostname" { print $2; exit }
    ' "${HOME}/.ssh/config")
fi

# ── Platform API token (for sandbox's HTTPS git access) ──────────────────────
# Docker Sandbox blocks outbound port 22, so teammate git operations
# (fetch, push, PR API) must go over HTTPS regardless of how the
# host accesses the repo. We need an app password / PAT / GitLab
# PAT to authenticate.
#
# Storage model:
#   macOS         token → Keychain only (encrypted, app-authorized).
#                 .sandbox/.repo-platform-api.env holds only metadata —
#                 no REPO_PLATFORM_API_TOKEN line. create.sh reads the
#                 Keychain and pipes the token into the sandbox via
#                 stdin, so the token never touches the host's
#                 regular filesystem.
#   Linux/Win     token → .sandbox/.repo-platform-api.env (mode 600,
#                 gitignored). Credential-manager integration
#                 (libsecret on Linux, Credential Manager on
#                 Windows) is a planned follow-up; the file is
#                 adequate for personal dev workstations but not
#                 as hardened as a proper keyring.
#
# Precedence for finding an already-saved token:
#   1. macOS Keychain (when on macOS; survives leave.sh)
#   2. .repo-platform-api.env's REPO_PLATFORM_API_TOKEN line
#   3. Prompt the user with platform-specific instructions
ENV_FILE="${PROJECT_DIR}/.sandbox/.repo-platform-api.env"
ON_MACOS=0
if [ "$(uname -s)" = "Darwin" ]; then
    ON_MACOS=1
fi

# Detect the platform hosting origin, resolving SSH aliases through
# the REAL_HOSTNAME derived earlier (if the SSH block ran) or from
# the HTTPS URL directly.
detect_platform_from_host() {
    local host="$1"
    case "$host" in
        *github.com)
            REPO_PLATFORM_TYPE="GITHUB"
            REPO_PLATFORM_HOST="github.com"
            API_URL="https://api.github.com"
            return 0 ;;
        *bitbucket.org)
            REPO_PLATFORM_TYPE="BITBUCKET"
            REPO_PLATFORM_HOST="bitbucket.org"
            API_URL="https://api.bitbucket.org"
            return 0 ;;
        *gitlab.com)
            REPO_PLATFORM_TYPE="GITLAB"
            REPO_PLATFORM_HOST="gitlab.com"
            API_URL="https://gitlab.com/api/v4"
            return 0 ;;
    esac
    return 1
}

REPO_PLATFORM_TYPE=""; REPO_PLATFORM_HOST=""; API_URL=""
case "${REMOTE_URL}" in
    https://*|http://*)
        host_part="${REMOTE_URL#*://}"
        host_part="${host_part%%/*}"
        detect_platform_from_host "$host_part" || true ;;
esac
# For SSH remotes, prefer REAL_HOSTNAME (resolved from ~/.ssh/config
# in the SSH block above); fall back to the raw SSH_HOST.
if [ -z "${REPO_PLATFORM_TYPE}" ] && [ -n "${REAL_HOSTNAME:-}" ]; then
    detect_platform_from_host "$REAL_HOSTNAME" || true
fi
if [ -z "${REPO_PLATFORM_TYPE}" ] && [ -n "${SSH_HOST:-}" ]; then
    detect_platform_from_host "$SSH_HOST" || true
fi

# Prompt if auto-detection failed.
if [ -z "${REPO_PLATFORM_TYPE}" ]; then
    echo ""
    echo "Could not auto-detect the platform hosting origin."
    echo "Pick one:"
    echo "  1. Bitbucket"
    echo "  2. GitHub"
    echo "  3. GitLab"
    echo "  4. Other / skip (no token provisioning)"
    while : ; do
        read -r -p "Enter number [1-4]: " plat_choice
        case "${plat_choice:-}" in
            1) detect_platform_from_host bitbucket.org; break ;;
            2) detect_platform_from_host github.com;    break ;;
            3) detect_platform_from_host gitlab.com;    break ;;
            4) REPO_PLATFORM_TYPE="SKIP"; break ;;
            *) echo "Please enter 1, 2, 3, or 4." ;;
        esac
    done
fi

# Extract owner/repo from REMOTE_URL for .repo-platform-api.env.
PATH_PART="${REMOTE_URL}"
PATH_PART="${PATH_PART#git@*:}"
PATH_PART="${PATH_PART#ssh://*/}"
PATH_PART="${PATH_PART#https://*/}"
PATH_PART="${PATH_PART#http://*/}"
PATH_PART="${PATH_PART%.git}"
OWNER="${PATH_PART%%/*}"
REPO_NAME="${PATH_PART##*/}"

if [ "${REPO_PLATFORM_TYPE}" != "SKIP" ] && [ -n "${REPO_PLATFORM_TYPE}" ]; then
    TOKEN=""
    API_USER=""

    # Precedence 1: macOS Keychain (canonical source on macOS).
    if [ "${ON_MACOS}" = "1" ]; then
        if kc_token=$(security find-generic-password -s "${KEYCHAIN_SERVICE}" -w 2>/dev/null); then
            TOKEN="$kc_token"
            # The account (-a) field in the Keychain entry is the API user.
            API_USER=$(security find-generic-password -s "${KEYCHAIN_SERVICE}" 2>&1 \
                | awk -F'"' '/"acct"<blob>=/ { print $2; exit }')
            echo "=== Repo-platform API token restored from macOS Keychain ==="
        fi
    fi

    # Precedence 2: .repo-platform-api.env file. Primary source on
    # non-macOS; fallback on macOS in case a Keychain write
    # previously failed.
    if [ -z "${TOKEN}" ] && [ -f "${ENV_FILE}" ]; then
        TOKEN="$(grep -E '^REPO_PLATFORM_API_TOKEN=' "$ENV_FILE" | cut -d= -f2- || true)"
        API_USER="$(grep -E '^REPO_PLATFORM_API_USER=' "$ENV_FILE" | cut -d= -f2- || true)"
        [ -n "${TOKEN}" ] && echo "=== Reusing repo-platform API token from .sandbox/.repo-platform-api.env ==="
    fi

    # Precedence 3: prompt the user.
    if [ -z "${TOKEN}" ]; then
        # Warn non-macOS users about the less-secure storage path
        # before they paste anything. Banner uses box-drawing
        # characters; renders correctly in any modern terminal.
        if [ "${ON_MACOS}" != "1" ]; then
            cat <<'BANNER'

┌─────────────────────────────────────────────────────────────────────────────┐
│  Note on repo token storage                                                 │
│  ─────────────────────────────────────────────────────────────────────────  │
│  On macOS, this kit stores your repo-platform API token in the Keychain —   │
│  encrypted at rest and app-authorized. However on your OS, credential       │
│  management integration isn't wired up yet, so the token will be saved to:  │
│                                                                             │
│      .sandbox/.repo-platform-api.env   (mode 600, gitignored)               │
│                                                                             │
│  That's adequate for a personal dev workstation, but not as hardened as a   │
│  proper credential store. If your compliance posture needs stronger, OS-    │
│  native integration (libsecret on Linux, Credential Manager on Windows) is  │
│  a straightforward follow-up — open an issue or contribute the integration. │
└─────────────────────────────────────────────────────────────────────────────┘

BANNER
        fi

        echo "=== Repo-platform API token (${REPO_PLATFORM_TYPE}) ==="
        echo ""
        echo "Teammates inside the sandbox push commits to your repo and"
        echo "drive PRs. The sandbox blocks outbound SSH, so this goes"
        echo "over HTTPS and needs a token scoped to this repo."
        echo ""
        case "${REPO_PLATFORM_TYPE}" in
            BITBUCKET)
                echo "Create a Bitbucket app password at:"
                echo "  https://bitbucket.org/account/settings/app-passwords/"
                echo ""
                echo "  Label:  agent-team:${SANDBOX_NAME}"
                echo "          (recognizable in the Bitbucket UI — revoke"
                echo "           it specifically when you leave this project)"
                echo "  Scopes: Repositories      Read, Write"
                echo "          Pull requests     Read, Write"
                echo ""
                echo "Bitbucket shows the token once — copy it before"
                echo "closing the tab."
                echo ""
                read -r -p "Bitbucket username: " API_USER ;;
            GITHUB)
                echo "Create a GitHub fine-grained PAT at:"
                echo "  https://github.com/settings/tokens?type=beta"
                echo ""
                echo "  Name:         agent-team:${SANDBOX_NAME}"
                echo "  Repository:   only this repo"
                echo "  Permissions:  Contents       Read and write"
                echo "                Pull requests  Read and write"
                echo ""
                echo "GitHub shows the token once — copy it before leaving"
                echo "the page."
                echo ""
                read -r -p "GitHub username: " API_USER ;;
            GITLAB)
                echo "Create a GitLab personal access token at:"
                echo "  https://gitlab.com/-/user_settings/personal_access_tokens"
                echo ""
                echo "  Name:    agent-team:${SANDBOX_NAME}"
                echo "  Scopes:  api"
                echo ""
                echo "GitLab shows the token once — copy it before leaving"
                echo "the page."
                echo ""
                # GitLab HTTPS basic-auth uses the literal string 'oauth2'
                # as the username alongside a PAT as the password.
                API_USER="oauth2" ;;
        esac
        echo "Leave the token blank to skip — teammates can still read"
        echo "public repos but won't be able to push or drive PRs on"
        echo "private ones."
        echo ""
        read -r -s -p "Token: " TOKEN
        echo ""

        # Save to Keychain on macOS (no opt-out; this is the secure
        # default). If the Keychain write fails, fall through to the
        # file path by flipping ON_MACOS off for the file-write block
        # below, so the user isn't left without a working token.
        if [ -n "${TOKEN}" ] && [ "${ON_MACOS}" = "1" ]; then
            if security add-generic-password \
                    -s "${KEYCHAIN_SERVICE}" \
                    -a "${API_USER:-platform-user}" \
                    -l "agent-team:${SANDBOX_NAME}" \
                    -w "${TOKEN}" \
                    -U >/dev/null 2>&1; then
                echo "Token saved to the Keychain (label: agent-team:${SANDBOX_NAME})."
            else
                echo "Warning: Keychain save failed; falling back to .repo-platform-api.env." >&2
                ON_MACOS=0
            fi
        fi
    fi

    # Write .repo-platform-api.env. Metadata always; token only on
    # non-macOS (or if a macOS Keychain save failed, which flipped
    # ON_MACOS to 0 above).
    if [ -n "${TOKEN}" ]; then
        mkdir -p "${PROJECT_DIR}/.sandbox"
        umask 077
        {
            echo "# Developer-local. Gitignored. Do NOT commit."
            echo "REPO_PLATFORM_TYPE=${REPO_PLATFORM_TYPE}"
            echo "REPO_PLATFORM_HOST=${REPO_PLATFORM_HOST}"
            echo "REPO_PLATFORM_API_URL=${API_URL}"
            echo "REPO_PLATFORM_API_USER=${API_USER}"
            if [ "${ON_MACOS}" != "1" ]; then
                echo "REPO_PLATFORM_API_TOKEN=${TOKEN}"
            fi
            echo "REPO_OWNER=${OWNER}"
            echo "REPO_NAME=${REPO_NAME}"
            echo "REPO_WORKSPACE=${OWNER}"
            echo "REPO_SLUG=${REPO_NAME}"
        } > "${ENV_FILE}"
    else
        echo "Skipping token setup (public-repo mode). Private-repo pushes from the sandbox will fail." >&2
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
