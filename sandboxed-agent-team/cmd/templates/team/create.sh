#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# create.sh — build the sandbox template image, create the sandbox
# VM, and attach a Claude Code session to it.
#
# When called with --resume or --fresh (forwarded by team/attach.sh),
# this script instead verifies that a sandbox already exists with
# the current Lead directive and reattaches to it, either resuming
# the previous Claude Code session (--resume) or starting a new one
# (--fresh). Run ./team/attach.sh rather than passing these flags
# to create.sh directly.

set -euo pipefail

# Parse flags. --resume and --fresh are forwarded by
# team/attach.sh; either one puts this script into "attach mode"
# (reattach to an existing sandbox instead of creating a new one).
#
#   no flags  : create mode. Fails fast if the sandbox already
#               exists.
#   --resume  : attach mode. Resumes the previous Claude Code
#               session via `claude --continue`.
#   --fresh   : attach mode. Starts a brand-new Claude Code
#               session in the existing sandbox.
#
# --resume and --fresh are mutually exclusive.
resume_session=0
fresh_session=0
for arg in "$@"; do
    case "$arg" in
        --resume) resume_session=1 ;;
        --fresh)  fresh_session=1 ;;
        *)
            echo "create.sh: unknown argument: $arg" >&2
            echo "Usage: $0 [--resume | --fresh]" >&2
            exit 2
            ;;
    esac
done
if [ "$resume_session" -eq 1 ] && [ "$fresh_session" -eq 1 ]; then
    echo "create.sh: --resume and --fresh are mutually exclusive." >&2
    exit 2
fi
# attach mode iff either flag is present
attach_mode=0
if [ "$resume_session" -eq 1 ] || [ "$fresh_session" -eq 1 ]; then
    attach_mode=1
fi

# ── Prerequisite check ──────────────────────────────────────────────────────
if ! docker sandbox --help &>/dev/null; then
    echo "Error: 'docker sandbox' is not available."
    echo "Install Docker Desktop with the sandbox feature enabled."
    echo "See: https://docs.docker.com/desktop/"
    exit 1
fi

# Derive names from the directory path.
# ~/workspaces/acme-corp/project-alpha → template: acme-corp-project-alpha-sandbox
#                                      → sandbox:  claude-acme-corp-project-alpha
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PARENT_DIR="$(basename "$(dirname "$PROJECT_DIR")")"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
TEMPLATE_IMAGE="${PARENT_DIR}-${PROJECT_NAME}-sandbox"
SANDBOX_NAME="claude-${PARENT_DIR}-${PROJECT_NAME}"

# ── Lead directive injected into Claude Code's system prompt ─────────────────
# Passed via --append-system-prompt so the sandboxed Claude Code auto-loads
# the Lead role on first turn — the human does not need to remember
# /project:team-start. Host Claude Code invocations don't use this script and
# are unaffected.
LEAD_DIRECTIVE=$(cat <<'EOF'
You are the Lead of this project's sandboxed agent team. On your very first response of this session, before responding substantively to any user message, read `.claude/commands/team-start.md` from the project root and follow its instructions to spawn the team and perform the Pre-Start Check. Only after setup is complete should you engage with the user's request. The slash command `/project:team-start` remains available for manual re-invocation (e.g., if the team needs to be re-spawned mid-session).
EOF
)

echo "=== Project:  ${PROJECT_DIR} ==="
echo "=== Template: ${TEMPLATE_IMAGE} ==="
echo "=== Sandbox:  ${SANDBOX_NAME} ==="

# ── Detect authentication ────────────────────────────────────────────────────
# The sandbox needs an API key or OAuth token. The host's interactive login
# does not carry over — `docker sandbox run` does not support -e for env
# vars, so we inject credentials via `docker sandbox exec` after startup.
#
# On macOS, Claude Code stores OAuth credentials in the Keychain.
# On other systems, CLAUDE_CODE_OAUTH_TOKEN must be set in the shell config.
AUTH_TOKEN=""
AUTH_TYPE=""
if [ -n "${ANTHROPIC_API_KEY:-}" ]; then
    AUTH_TOKEN="${ANTHROPIC_API_KEY}"
    AUTH_TYPE="api-key"
elif command -v security &>/dev/null; then
    SVC=$(security dump-keychain 2>/dev/null \
        | grep '"svce"' \
        | grep -i 'claude.*credential' \
        | tail -1 \
        | sed 's/.*<blob>="\(.*\)"/\1/' \
        | tr -d '\n')
    if [ -n "${SVC}" ]; then
        KEYCHAIN_TOKEN=$(security find-generic-password -s "${SVC}" -w 2>/dev/null || true)
        if [ -n "${KEYCHAIN_TOKEN}" ]; then
            AUTH_TOKEN="${KEYCHAIN_TOKEN}"
            AUTH_TYPE="oauth"
            echo "=== OAuth token refreshed from macOS Keychain ==="
        fi
    fi
elif [ -n "${CLAUDE_CODE_OAUTH_TOKEN:-}" ]; then
    AUTH_TOKEN="${CLAUDE_CODE_OAUTH_TOKEN}"
    AUTH_TYPE="oauth"
fi

if [ -z "${AUTH_TOKEN}" ]; then
    echo "Error: No authentication detected."
    echo ""
    echo "  The sandbox needs one of:"
    echo "    - ANTHROPIC_API_KEY (API key users)"
    echo "    - CLAUDE_CODE_OAUTH_TOKEN (OAuth / team / Max / Enterprise users)"
    echo ""
    echo "  On macOS: run 'claude' in a terminal, authenticate via /login,"
    echo "    then re-run this script — the token will be read from the Keychain."
    echo "  On other systems: export CLAUDE_CODE_OAUTH_TOKEN in your shell config."
    echo "    See ONBOARDING.md for details."
    exit 1
fi

# ── Inject credentials into sandbox ──────────────────────────────────────────
# docker sandbox run does not support passing env vars (-e), and the sandbox
# auto-updates the claude binary (overwriting any Dockerfile-based wrapper).
# Instead, we inject credentials directly into the sandbox filesystem via
# `docker sandbox exec` after the sandbox starts.
inject_credentials() {
    echo "=== Injecting credentials into sandbox ==="

    # ── Claude Code authentication ────────────────────────────────────────
    # The token is piped via stdin so it never appears on a command
    # line visible to `ps` on the host.
    if [ "${AUTH_TYPE}" = "oauth" ]; then
        # Write the OAuth credential JSON to Claude's credentials file.
        printf '%s' "${AUTH_TOKEN}" | docker sandbox exec -i "${SANDBOX_NAME}" bash -c \
            'mkdir -p /home/agent/.claude \
             && cat > /home/agent/.claude/.credentials.json \
             && chmod 600 /home/agent/.claude/.credentials.json'
    else
        # API key: append an export line to .bashrc, but only if not
        # already present. The probe runs token-free; the write is gated
        # behind it and reads the token from stdin.
        if ! docker sandbox exec "${SANDBOX_NAME}" bash -c \
            'grep -q ANTHROPIC_API_KEY /home/agent/.bashrc 2>/dev/null'; then
            printf '%s' "${AUTH_TOKEN}" | docker sandbox exec -i "${SANDBOX_NAME}" bash -c \
                'read -r token && printf "export ANTHROPIC_API_KEY=%q\n" "$token" >> /home/agent/.bashrc'
        fi
    fi

    # Ensure hasCompletedOnboarding is set (required for OAuth recognition).
    docker sandbox exec "${SANDBOX_NAME}" bash -c \
        "test -f /home/agent/.claude.json || echo '{}' > /home/agent/.claude.json; \
         jq '.hasCompletedOnboarding = true' /home/agent/.claude.json > /tmp/.claude.json \
         && mv /tmp/.claude.json /home/agent/.claude.json"

    # ── Platform API credentials (for PR merge method + HTTPS git) ────────
    # Metadata (REPO_PLATFORM_TYPE, API URLs, repo owner/slug, etc.)
    # comes from .sandbox/.repo-platform-api.env. On macOS the file omits
    # the token — we read it from the Keychain and pipe it through
    # stdin so it never touches the host's regular filesystem (same
    # pattern as the Claude OAuth token above).
    if [ -f "${PROJECT_DIR}/.sandbox/.repo-platform-api.env" ]; then
        docker sandbox exec "${SANDBOX_NAME}" bash -c \
            "cat '${PROJECT_DIR}/.sandbox/.repo-platform-api.env' >> /home/agent/.bashrc"
        echo "=== Platform API metadata injected ==="
    fi
    if [ "$(uname -s)" = "Darwin" ]; then
        if kc_token=$(security find-generic-password -s "agent-team.${SANDBOX_NAME}" -w 2>/dev/null); then
            printf '%s' "$kc_token" | docker sandbox exec -i "${SANDBOX_NAME}" bash -c \
                'read -r token && printf "export REPO_PLATFORM_API_TOKEN=%q\n" "$token" >> /home/agent/.bashrc'
            echo "=== Repo-platform API token injected from Keychain ==="
        fi
    fi

    # ── Sandbox-side git HTTPS auth ──────────────────────────────────────
    # Docker Sandbox blocks outbound port 22 (SSH), so agents must use
    # HTTPS for all git operations even if the project's origin URL
    # is SSH on the host. We configure git inside the sandbox to:
    #   1. use a credential helper that reads stored user:token pairs,
    #      writing the token through stdin so it never appears in `ps`;
    #   2. transparently rewrite SSH-style origin URLs to HTTPS via
    #      `url.insteadOf` — only when the host's origin actually uses
    #      SSH, so HTTPS-origin projects are unaffected.
    if [ -f "${PROJECT_DIR}/.sandbox/.repo-platform-api.env" ]; then
        # Pull the pieces we need without sourcing (some values may
        # contain characters bash would mis-interpret on source).
        P_HOST=$(grep -E '^REPO_PLATFORM_HOST=' "${PROJECT_DIR}/.sandbox/.repo-platform-api.env" | cut -d= -f2-)
        P_USER=$(grep -E '^REPO_PLATFORM_API_USER=' "${PROJECT_DIR}/.sandbox/.repo-platform-api.env" | cut -d= -f2-)
        P_TOKEN=""
        if [ "$(uname -s)" = "Darwin" ]; then
            P_TOKEN=$(security find-generic-password -s "agent-team.${SANDBOX_NAME}" -w 2>/dev/null || true)
        fi
        if [ -z "$P_TOKEN" ]; then
            P_TOKEN=$(grep -E '^REPO_PLATFORM_API_TOKEN=' "${PROJECT_DIR}/.sandbox/.repo-platform-api.env" | cut -d= -f2- || true)
        fi

        if [ -n "$P_HOST" ] && [ -n "$P_USER" ] && [ -n "$P_TOKEN" ]; then
            # Write the credentials file via stdin so the token never
            # appears on a command line visible to `ps`.
            printf 'https://%s:%s@%s\n' "$P_USER" "$P_TOKEN" "$P_HOST" \
                | docker sandbox exec -i "${SANDBOX_NAME}" bash -c \
                  'umask 077 && cat > /home/agent/.git-credentials && \
                   git config --global credential.helper store'

            # If the host's origin uses an SSH-style URL, teach the
            # sandbox's git to rewrite that prefix to HTTPS on the fly.
            # Both SSH forms (`git@host:` and `ssh://host/`) get covered.
            sandbox_origin="$(git -C "${PROJECT_DIR}" remote get-url origin 2>/dev/null || true)"
            case "$sandbox_origin" in
                git@*)
                    ssh_host="${sandbox_origin#*@}"
                    ssh_host="${ssh_host%%:*}"
                    docker sandbox exec "${SANDBOX_NAME}" \
                        git config --global "url.https://${P_HOST}/.insteadOf" "git@${ssh_host}:"
                    ;;
                ssh://*)
                    ssh_host="${sandbox_origin#ssh://}"
                    ssh_host="${ssh_host#*@}"
                    ssh_host="${ssh_host%%/*}"
                    docker sandbox exec "${SANDBOX_NAME}" \
                        git config --global "url.https://${P_HOST}/.insteadOf" "ssh://git@${ssh_host}/"
                    ;;
            esac
            echo "=== Sandbox-side git HTTPS auth configured ==="
        fi
    fi

    # ── Host terminal type ────────────────────────────────────────────────
    # Pass the host's terminal identity into the sandbox so the Pre-Start
    # Check can skip asking the human about split-pane support.
    HOST_TERM="${TERM_PROGRAM:-unknown}"
    docker sandbox exec "${SANDBOX_NAME}" bash -c \
        "echo '${HOST_TERM}' > /home/agent/.host-terminal"

    echo "=== Credentials injected ==="
}

# ── Directive-hash tracking ─────────────────────────────────────────────────
# Docker bakes the sandbox's entrypoint at creation time. If
# LEAD_DIRECTIVE changes (e.g., because this script was regenerated
# from the setup kit), the sandbox keeps running with the stale
# directive. We store the SHA-256 of LEAD_DIRECTIVE alongside the
# sandbox (in .sandbox/.last-directive, gitignored). On attach,
# we compare the stored hash to the current one; on mismatch we
# fail fast and tell the user to destroy + create to refresh.
DIRECTIVE_FILE="${PROJECT_DIR}/.sandbox/.last-directive"
DIRECTIVE_HASH=$(printf '%s' "${LEAD_DIRECTIVE}" | shasum -a 256 | awk '{print $1}')
STORED_HASH=$(cat "$DIRECTIVE_FILE" 2>/dev/null || true)
EXISTING=$(docker sandbox ls 2>/dev/null | grep -w "${SANDBOX_NAME}" || true)

if [ "$attach_mode" -eq 1 ]; then
    # ── Attach path (invoked from team/attach.sh --resume|--fresh) ───────────
    if [ -z "$EXISTING" ]; then
        echo "Error: no sandbox '${SANDBOX_NAME}' exists for this project." >&2
        echo "Run ./team/create.sh to create one." >&2
        exit 1
    fi
    if [ "${DIRECTIVE_HASH}" != "${STORED_HASH}" ]; then
        echo "Error: the Lead directive has changed since this sandbox was created." >&2
        echo "Run ./team/destroy.sh && ./team/create.sh to refresh." >&2
        exit 1
    fi
    echo "=== Reconnecting to existing sandbox ==="
    inject_credentials

    # Launch Claude Code fresh inside the running sandbox via
    # `docker sandbox exec`. Using exec (instead of
    # `docker sandbox run <name>`) lets us pass flags like
    # --continue, which replays the previous session's context.
    # The LEAD_DIRECTIVE is re-appended either way, as a
    # belt-and-suspenders — harmless on --continue (same directive
    # baked in at creation) and correct on --fresh.
    if [ "$resume_session" -eq 1 ]; then
        exec docker sandbox exec -it "${SANDBOX_NAME}" \
            bash -c 'cd "$1" && exec claude --continue --append-system-prompt "$2"' \
            _ "${PROJECT_DIR}" "${LEAD_DIRECTIVE}"
    else
        exec docker sandbox exec -it "${SANDBOX_NAME}" \
            bash -c 'cd "$1" && exec claude --append-system-prompt "$2"' \
            _ "${PROJECT_DIR}" "${LEAD_DIRECTIVE}"
    fi
fi

# ── Create path ─────────────────────────────────────────────────────────────
if [ -n "$EXISTING" ]; then
    echo "Error: sandbox '${SANDBOX_NAME}' already exists for this project." >&2
    echo "Run ./team/attach.sh to reattach, or" >&2
    echo "./team/destroy.sh && ./team/create.sh to rebuild from scratch." >&2
    exit 1
fi

# Build the custom template image if a Dockerfile is present.
DOCKERFILE="${PROJECT_DIR}/.sandbox/Dockerfile"
if [ -f "$DOCKERFILE" ]; then
    echo "Building sandbox template (first build downloads ~1 GB of"
    echo "  dependencies and typically takes 2-5 minutes on a fast connection,"
    echo "  longer on slower networks). Subsequent builds use the Docker cache"
    echo "  and are much faster."
    echo ""
    echo "  Progress updates below. If the display stops changing for more"
    echo "  than 5 minutes, the build may be hung — cancel (Ctrl+C) and"
    echo "  check the Dockerfile for commands that may hang."
    echo ""
    echo "  Note: Playwright may print a 'BEWARE: your OS is not officially"
    echo "  supported' warning on some platforms. This is cosmetic — it uses"
    echo "  a working fallback automatically."
    docker build -t "${TEMPLATE_IMAGE}" -f "$DOCKERFILE" "$PROJECT_DIR"
else
    echo "No .sandbox/Dockerfile found. Using default template."
    TEMPLATE_IMAGE=""
fi

# docker sandbox run is blocking (interactive). Spawn a background
# poller that waits for the sandbox to appear, then injects
# credentials.
(
    for i in $(seq 1 30); do
        sleep 2
        if docker sandbox ls 2>/dev/null | grep -qw "${SANDBOX_NAME}"; then
            inject_credentials
            break
        fi
    done
) &
INJECT_PID=$!

# Record the directive hash before starting so future attach.sh runs
# can verify even if this run is interrupted mid-create.
echo "${DIRECTIVE_HASH}" > "$DIRECTIVE_FILE"

# `docker sandbox run` grammar:
#   docker sandbox run [flags] AGENT [WORKSPACE] [-- AGENT_ARGS...]
# Agent args (flags meant for claude) MUST come after the `--`
# separator. Without it, they'd be parsed as additional workspace
# paths and silently turned into bogus directories.
if [ -n "$TEMPLATE_IMAGE" ]; then
    docker sandbox run \
        --name "${SANDBOX_NAME}" \
        --template "${TEMPLATE_IMAGE}" \
        claude \
        "${PROJECT_DIR}" \
        -- \
        --append-system-prompt "${LEAD_DIRECTIVE}"
else
    docker sandbox run \
        --name "${SANDBOX_NAME}" \
        claude \
        "${PROJECT_DIR}" \
        -- \
        --append-system-prompt "${LEAD_DIRECTIVE}"
fi

# Clean up background job if still running.
kill $INJECT_PID 2>/dev/null || true
