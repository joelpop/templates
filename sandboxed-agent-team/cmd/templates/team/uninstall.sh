#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team install`.

# uninstall.sh — remove the agent team kit from this project.
#
# Stops the sandbox, deletes kit-generated files, excises the
# @CLAUDE_TEAM.md import block from CLAUDE.md and the kit's block
# from .gitignore, commits the removal. Does NOT touch docs/ —
# that belongs to the project.
#
# A developer can run this directly (`./team/uninstall.sh`) without
# the Go installer being available. `agent-team uninstall` simply
# shells out to this script.

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "${PROJECT_DIR}"

if [ ! -f "${PROJECT_DIR}/CLAUDE_TEAM.md" ]; then
    echo "No kit installation detected here. Nothing to remove."
    exit 0
fi

read -r -p "This deletes the kit and commits the removal. Continue? [y/N] " resp
case "${resp:-}" in
    y|Y|yes|YES) ;;
    *) echo "Aborted."; exit 0 ;;
esac

# Stop the sandbox BEFORE we delete team/stop.sh.
if [ -x "${PROJECT_DIR}/team/stop.sh" ]; then
    "${PROJECT_DIR}/team/stop.sh" || true
fi

# Remove developer-local state (the same paths team/leave.sh targets).
rm -rf "${PROJECT_DIR}/.sandbox/.ssh"
rm -f  "${PROJECT_DIR}/.sandbox/.ssh.source"
rm -f  "${PROJECT_DIR}/.sandbox/.platform-api.env"
rm -f  "${PROJECT_DIR}/.sandbox/.oauth-token"
rm -f  "${PROJECT_DIR}/.sandbox/.last-directive"
rm -f  "${PROJECT_DIR}/.claude/.last-onboarded"
rm -f  "${PROJECT_DIR}/.claude/.team-active"
rm -rf "${PROJECT_DIR}/.claude/.tasks"
rm -f  "${PROJECT_DIR}/.claude/.progress.md"
rm -rf "${PROJECT_DIR}/.claude/.worktrees"

# Delete tracked kit files. --ignore-unmatch makes this idempotent.
git rm -rf --ignore-unmatch -- \
    CLAUDE_TEAM.md \
    ONBOARDING.md \
    TEAM_GUIDE.md \
    .mcp.json \
    .claude/team-variables.yaml \
    .claude/settings.json \
    .claude/commands/team-start.md \
    .sandbox \
    team \
    > /dev/null

# Excise the @CLAUDE_TEAM.md import block from CLAUDE.md.
# Markers: <!-- sandboxed-agent-team: begin --> / end -->
# Whitespace-tolerant: a Prettier run that normalizes spaces still matches.
if [ -f "${PROJECT_DIR}/CLAUDE.md" ]; then
    sed -i.bak -E '/<!-- *sandboxed-agent-team: *begin *-->/,/<!-- *sandboxed-agent-team: *end *-->/d' \
        "${PROJECT_DIR}/CLAUDE.md"
    rm -f "${PROJECT_DIR}/CLAUDE.md.bak"
    if ! grep -q '[^[:space:]]' "${PROJECT_DIR}/CLAUDE.md"; then
        git rm -f --ignore-unmatch -- CLAUDE.md > /dev/null
    else
        git add -f -- CLAUDE.md
    fi
fi

# Excise the kit's block from .gitignore (same marker pattern).
if [ -f "${PROJECT_DIR}/.gitignore" ]; then
    sed -i.bak -E '/^# *BEGIN *sandboxed-agent-team/,/^# *END *sandboxed-agent-team/d' \
        "${PROJECT_DIR}/.gitignore"
    rm -f "${PROJECT_DIR}/.gitignore.bak"
    if ! grep -q '[^[:space:]]' "${PROJECT_DIR}/.gitignore"; then
        git rm -f --ignore-unmatch -- .gitignore > /dev/null
    else
        git add -f -- .gitignore
    fi
fi

if git diff --cached --quiet; then
    echo "Nothing to commit — kit may already have been uninstalled."
else
    git commit -m "Uninstall sandboxed-agent-team kit" > /dev/null
    echo "Kit removed and removal committed."
fi
