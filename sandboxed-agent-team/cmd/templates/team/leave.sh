#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# leave.sh — tear down the current developer's local state for this
# project. Does NOT touch versioned kit files; only the developer-local
# (gitignored) artifacts. Reverses what team/join.sh did.
#
# Prompts for confirmation unless --yes is passed. The --yes flag
# is how team/uninstall.sh drives this script after its own combined
# prompt, avoiding a double prompt.

set -euo pipefail

assume_yes=0
for arg in "$@"; do
    case "$arg" in
        --yes) assume_yes=1 ;;
        *)
            echo "leave.sh: unknown argument: $arg" >&2
            echo "Usage: $0 [--yes]" >&2
            exit 2
            ;;
    esac
done

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "${PROJECT_DIR}"

while [ "$assume_yes" -eq 0 ]; do
    read -r -p "Remove developer-local state for $(basename "${PROJECT_DIR}")? [yes/NO] " resp
    case "${resp:-}" in
        yes|Yes|YES) assume_yes=1 ;;
        no|No|NO) echo "Aborted."; exit 0 ;;
        *) echo "?invalid response: ${resp}" ;;
    esac
done

echo "=== Removing developer-local state for $(basename "${PROJECT_DIR}") ==="

# Stop the sandbox if it's running. No-op if not.
if [ -x "${PROJECT_DIR}/team/stop.sh" ]; then
    "${PROJECT_DIR}/team/stop.sh" || true
fi

# Per-developer state under .sandbox/ (all dot-prefixed per kit convention).
rm -rf "${PROJECT_DIR}/.sandbox/.ssh"
rm -f  "${PROJECT_DIR}/.sandbox/.ssh.source"
rm -f  "${PROJECT_DIR}/.sandbox/.platform-api.env"
rm -f  "${PROJECT_DIR}/.sandbox/.oauth-token"
rm -f  "${PROJECT_DIR}/.sandbox/.last-directive"

# Per-developer state under .claude/ (also dot-prefixed).
rm -f  "${PROJECT_DIR}/.claude/.last-onboarded"
rm -f  "${PROJECT_DIR}/.claude/.team-active"
rm -rf "${PROJECT_DIR}/.claude/.tasks"
rm -f  "${PROJECT_DIR}/.claude/.progress.md"
rm -rf "${PROJECT_DIR}/.claude/.worktrees"

echo "Local state removed. Kit artifacts (versioned files) are untouched."
echo "Run team/join.sh to re-provision this workspace later."
