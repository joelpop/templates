#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# shell.sh — drop into a bash shell inside this project's running
# sandbox. Useful for inspection and ad-hoc commands (e.g.,
# checking SSH config, viewing injected credentials, or running
# project build commands outside a Claude Code session).
#
# The sandbox must already be running (created via team/create.sh,
# typically via team/join.sh). Exits with a clear message if not.

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PARENT_DIR="$(basename "$(dirname "$PROJECT_DIR")")"
PROJECT_NAME="$(basename "$PROJECT_DIR")"
SANDBOX_NAME="claude-${PARENT_DIR}-${PROJECT_NAME}"

if ! docker sandbox ls 2>/dev/null | grep -qw "${SANDBOX_NAME}"; then
    echo "Error: no sandbox '${SANDBOX_NAME}' exists for this project." >&2
    echo "Run ./team/create.sh to create one." >&2
    exit 1
fi

# Start at the project workspace; the user expects to land where
# agents normally work, not wherever docker's default cwd is.
exec docker sandbox exec -it "${SANDBOX_NAME}" \
    bash -c 'cd "$1" && exec bash' _ "${PROJECT_DIR}"