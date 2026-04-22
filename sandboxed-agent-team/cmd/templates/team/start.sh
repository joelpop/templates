#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

set -euo pipefail

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

# ── SSH key sync (host side) ─────────────────────────────────────────────────
# Refresh the key pair in .sandbox/.ssh/ from the host path recorded in
# .sandbox/.ssh.source, in case the developer rotated their key.
SSH_SOURCE_FILE="${PROJECT_DIR}/.sandbox/.ssh.source"
if [ -f "$SSH_SOURCE_FILE" ]; then
    SSH_KEY=$(grep -v '^#' "$SSH_SOURCE_FILE" | head -1 | tr -d '[:space:]')
    SSH_KEY="${SSH_KEY/#\~/$HOME}"
    if [ -n "$SSH_KEY" ] && [ -f "$SSH_KEY" ]; then
        SSH_DIR="${PROJECT_DIR}/.sandbox/.ssh"
        mkdir -p "$SSH_DIR"
        cp "$SSH_KEY" "$SSH_DIR/"
        [ -f "${SSH_KEY}.pub" ] && cp "${SSH_KEY}.pub" "$SSH_DIR/"
        echo "=== SSH key synced from ${SSH_KEY} ==="
    else
        echo ""
        echo "SSH key '${SSH_KEY}' (from .sandbox/.ssh.source) not found."
        echo ""
        echo "  The project declares SSH use but the key at that path is"
        echo "  missing. Git operations will fail inside the sandbox"
        echo "  until this is fixed."
        echo ""
        echo "  If you know the correct path, enter it now — this script"
        echo "  will update .sandbox/.ssh.source and re-sync. Otherwise,"
        echo "  press Enter to abort; to reconfigure SSH interactively,"
        echo "  start a Claude Code session at your host terminal in this"
        echo "  project directory and say:"
        echo "    Read ONBOARDING.md and execute the setup checklist."
        echo ""
        # Retry loop: accept empty input to abort; otherwise re-prompt
        # if the given path doesn't exist (up to 3 attempts to prevent
        # runaway loops from clipboard errors etc.).
        NEW_PATH=""
        for attempt in 1 2 3; do
            read -p "  Correct SSH key path (or Enter to abort): " NEW_PATH
            if [ -z "$NEW_PATH" ]; then
                echo "Aborting. See the note above on re-running onboarding."
                exit 1
            fi
            NEW_PATH="${NEW_PATH/#\~/$HOME}"
            if [ -f "$NEW_PATH" ]; then
                break
            fi
            echo "File '$NEW_PATH' not found."
            if [ $attempt -eq 3 ]; then
                echo "Aborting after 3 attempts. See the note above on re-running onboarding."
                exit 1
            fi
            echo "Try again (attempt $((attempt + 1)) of 3) or press Enter to abort."
        done
        echo "$NEW_PATH" > "$SSH_SOURCE_FILE"
        SSH_KEY="$NEW_PATH"
        SSH_DIR="${PROJECT_DIR}/.sandbox/.ssh"
        mkdir -p "$SSH_DIR"
        cp "$SSH_KEY" "$SSH_DIR/"
        [ -f "${SSH_KEY}.pub" ] && cp "${SSH_KEY}.pub" "$SSH_DIR/"
        echo "=== Updated .sandbox/.ssh.source and synced SSH key from ${SSH_KEY} ==="
    fi
fi

# ── Build custom template ────────────────────────────────────────────────────
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

    # ── SSH keys ──────────────────────────────────────────────────────────
    # The workspace is mounted at PROJECT_DIR inside the sandbox (same path
    # as on the host). Copy SSH material from .sandbox/.ssh/ to /home/agent/.ssh/.
    if [ -d "${PROJECT_DIR}/.sandbox/.ssh" ]; then
        docker sandbox exec "${SANDBOX_NAME}" bash -c \
            "cp '${PROJECT_DIR}/.sandbox/.ssh/'* /home/agent/.ssh/ 2>/dev/null; \
             chmod 600 /home/agent/.ssh/id_* 2>/dev/null; \
             chmod 644 /home/agent/.ssh/*.pub /home/agent/.ssh/config \
                       /home/agent/.ssh/known_hosts 2>/dev/null"
        echo "=== SSH keys injected ==="
    fi

    # ── Platform API credentials (for PR merge method) ────────────────────
    if [ -f "${PROJECT_DIR}/.sandbox/.platform-api.env" ]; then
        docker sandbox exec "${SANDBOX_NAME}" bash -c \
            "cat '${PROJECT_DIR}/.sandbox/.platform-api.env' >> /home/agent/.bashrc"
        echo "=== Platform API credentials injected ==="
    fi

    # ── Host terminal type ────────────────────────────────────────────────
    # Pass the host's terminal identity into the sandbox so the Pre-Start
    # Check can skip asking the human about split-pane support.
    HOST_TERM="${TERM_PROGRAM:-unknown}"
    docker sandbox exec "${SANDBOX_NAME}" bash -c \
        "echo '${HOST_TERM}' > /home/agent/.host-terminal"

    echo "=== Credentials injected ==="
}

# ── Start, reconnect, or recreate ──────────────────────────────────────────
# Docker bakes the sandbox's entrypoint (the claude invocation and its
# flags) at creation time. `docker sandbox run <name>` re-attaches to
# that baked-in entrypoint; extra args here would be ignored. So if
# LEAD_DIRECTIVE ever changes (e.g., because this script was
# regenerated from the setup kit), an existing sandbox would keep
# running with the stale directive.
#
# To fix that automatically, we store LEAD_DIRECTIVE alongside the
# sandbox (in .sandbox/.last-directive, gitignored via the .sandbox/
# rule). On each run we compare the current directive to the stored
# one. If they differ, we destroy the sandbox so the new-sandbox
# branch below recreates it with the updated entrypoint. The
# recreation is fast because Docker caches the template image.
EXISTING=$(docker sandbox ls 2>/dev/null | grep -w "${SANDBOX_NAME}" || true)
DIRECTIVE_FILE="${PROJECT_DIR}/.sandbox/.last-directive"
# Hash-based comparison (robust to trailing whitespace, IDE-auto-edits,
# line-ending differences, etc.). The file holds just the hash, not
# the directive text.
DIRECTIVE_HASH=$(printf '%s' "${LEAD_DIRECTIVE}" | shasum -a 256 | awk '{print $1}')
STORED_HASH=$(cat "$DIRECTIVE_FILE" 2>/dev/null || true)

if [ -n "$EXISTING" ] && [ "${DIRECTIVE_HASH}" != "${STORED_HASH}" ]; then
    if [ -z "${STORED_HASH}" ]; then
        echo "=== Sandbox predates directive tracking — recreating with current LEAD_DIRECTIVE ==="
    else
        echo "=== LEAD_DIRECTIVE has changed since last run — recreating sandbox ==="
    fi
    docker sandbox rm "${SANDBOX_NAME}"
    EXISTING=""
fi

if [ -n "$EXISTING" ]; then
    echo "=== Reconnecting to existing sandbox ==="
    inject_credentials
    # Re-attaches to the baked-in entrypoint set at creation.
    docker sandbox run "${SANDBOX_NAME}"
else
    # New sandbox: docker sandbox run blocks (it is interactive), so we
    # inject credentials from a background job that polls until the sandbox
    # appears, then runs inject_credentials.
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

    # Record the directive hash before starting so future runs can
    # detect changes even if this run is interrupted.
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
fi
