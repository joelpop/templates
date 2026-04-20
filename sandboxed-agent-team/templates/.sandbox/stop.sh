#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the team setup kit
# (SANDBOXED_AGENT_TEAMS.md) and re-run the setup at your host
# terminal.

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PARENT_DIR="$(basename "$(dirname "$PROJECT_DIR")")"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
TEMPLATE_IMAGE="${PARENT_DIR}-${PROJECT_NAME}-sandbox"
SANDBOX_NAME="claude-${PARENT_DIR}-${PROJECT_NAME}"

echo "=== Tearing down sandbox: ${SANDBOX_NAME} ==="
echo "WARNING: This destroys the sandbox VM and everything inside it."
echo "         Files in ${PROJECT_DIR} are NOT deleted."
read -p "Continue? (y/N) " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Clean up git worktrees created by the agent team
    if [ -d "${PROJECT_DIR}/.claude/worktrees" ]; then
        echo "Cleaning up agent worktrees..."
        for wt in "${PROJECT_DIR}/.claude/worktrees"/*/; do
            git -C "$PROJECT_DIR" worktree remove --force "${wt%/}" 2>/dev/null || true
        done
        git -C "$PROJECT_DIR" worktree prune
    fi

    docker sandbox rm "${SANDBOX_NAME}"
    echo "Sandbox removed."

    read -p "Also remove Docker template image '${TEMPLATE_IMAGE}'? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker rmi "${TEMPLATE_IMAGE}" 2>/dev/null || true
        echo "Template image removed."
    fi

    echo "=== Teardown complete. ==="
    echo "Delete ${PROJECT_DIR} manually per your data retention policy."
else
    echo "Aborted."
fi
