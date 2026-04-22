#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# stop.sh — destroy this project's sandbox VM.
#
# Prompts for confirmation unless --yes is passed. The --yes flag
# is how team/join.sh drives a forced rebuild without double-prompting
# the user; under --yes the template image is left intact so the
# subsequent rebuild can reuse Docker's cache.

set -euo pipefail

assume_yes=0
for arg in "$@"; do
    case "$arg" in
        --yes) assume_yes=1 ;;
        *)
            echo "stop.sh: unknown argument: $arg" >&2
            echo "Usage: $0 [--yes]" >&2
            exit 2
            ;;
    esac
done

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PARENT_DIR="$(basename "$(dirname "$PROJECT_DIR")")"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
TEMPLATE_IMAGE="${PARENT_DIR}-${PROJECT_NAME}-sandbox"
SANDBOX_NAME="claude-${PARENT_DIR}-${PROJECT_NAME}"

# Short-circuit if there's nothing to tear down. Avoids confusing
# "Failed to delete sandbox: VM not found" output and spurious
# prompts when stop.sh is invoked on a project that never started
# its sandbox (e.g., from join.sh on a fresh install).
if ! docker sandbox ls 2>/dev/null | grep -qw "${SANDBOX_NAME}"; then
    echo "No sandbox '${SANDBOX_NAME}' to tear down."
    exit 0
fi

if [ "$assume_yes" -eq 1 ]; then
    echo "=== Discarding sandbox '${SANDBOX_NAME}' to match the updated setup ==="
else
    echo "=== Tearing down sandbox: ${SANDBOX_NAME} ==="
    echo "WARNING: This destroys the sandbox VM and everything inside it."
    echo "         Files in ${PROJECT_DIR} are NOT deleted."
    while : ; do
        read -r -p "Continue? [yes/NO] " resp
        case "${resp:-}" in
            yes|Yes|YES) break ;;
            no|No|NO|"") echo "Aborted."; exit 0 ;;
            *) echo "Invalid response: ${resp}" ;;
        esac
    done
fi

# Clean up git worktrees created by the agent team.
if [ -d "${PROJECT_DIR}/.claude/.worktrees" ]; then
    echo "Cleaning up agent worktrees..."
    for wt in "${PROJECT_DIR}/.claude/.worktrees"/*/; do
        git -C "$PROJECT_DIR" worktree remove --force "${wt%/}" 2>/dev/null || true
    done
    git -C "$PROJECT_DIR" worktree prune
fi

docker sandbox rm "${SANDBOX_NAME}"
echo "Sandbox removed."

# Template image removal is only offered for direct user invocations.
# --yes callers (team/join.sh) keep the cached image so the rebuild
# is fast.
if [ "$assume_yes" -eq 0 ] && docker image inspect "${TEMPLATE_IMAGE}" &>/dev/null; then
    while : ; do
        read -r -p "Also remove Docker template image '${TEMPLATE_IMAGE}'? [yes/NO] " resp
        case "${resp:-}" in
            yes|Yes|YES)
                docker rmi "${TEMPLATE_IMAGE}" 2>/dev/null || true
                echo "Template image removed."
                break
                ;;
            no|No|NO|"") break ;;
            *) echo "Invalid response: ${resp}" ;;
        esac
    done
fi

echo "=== Teardown complete. ==="
if [ "$assume_yes" -eq 0 ]; then
    echo "Delete ${PROJECT_DIR} manually per your data retention policy."
fi
