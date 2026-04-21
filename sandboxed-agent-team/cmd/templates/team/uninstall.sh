#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# uninstall.sh — remove the agent team kit from this project.
#
# Delegates teardown of developer-local state to team/leave.sh,
# then deletes the kit's versioned files, excises the
# @CLAUDE_TEAM.md import block from CLAUDE.md and the kit's
# block from .gitignore, and commits the removal.
#
# Does NOT touch docs/ — that belongs to the project.
#
# Does NOT remove .mcp.json or .claude/settings.json — these are
# generic Claude Code config files and the user may have added
# their own entries alongside the kit's. Leaving them alone lets
# those additions survive an uninstall. (Install still overwrites
# them, so user additions don't survive a re-install; that's a
# known asymmetry.)
#
# This script is the only programmatic way to uninstall the kit.
# Run it directly from the project root:
#
#   ./team/uninstall.sh
#
# (The agent-team-install binary only handles installs and updates;
# lifecycle commands (join, leave, start, stop, uninstall) live
# here in team/ so they stay in lockstep with the kit version
# committed to the project.)

set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "${PROJECT_DIR}"

if [ ! -f "${PROJECT_DIR}/CLAUDE_TEAM.md" ]; then
    echo "No kit installation detected here. Nothing to remove."
    exit 0
fi

# ── Safety check for in-progress work ───────────────────────────────────────
# Uninstall will delete everything under .claude/.worktrees/, .claude/.tasks/,
# and .claude/.progress.md. Surface anything non-trivial so the developer
# can save it before we discard it.
warnings=()

if [ -d "${PROJECT_DIR}/.claude/.worktrees" ]; then
    for wt in "${PROJECT_DIR}"/.claude/.worktrees/*/; do
        [ -d "$wt" ] || continue
        if [ -n "$(git -C "$wt" status --porcelain 2>/dev/null)" ]; then
            warnings+=("  uncommitted changes in worktree: ${wt#${PROJECT_DIR}/}")
        fi
    done
fi

if [ -d "${PROJECT_DIR}/.claude/.tasks" ]; then
    task_count=$(find "${PROJECT_DIR}/.claude/.tasks" -type f -name '*.md' 2>/dev/null | wc -l | tr -d ' ')
    if [ "${task_count:-0}" -gt 0 ]; then
        warnings+=("  ${task_count} task file(s) in .claude/.tasks/ (active or suspended work)")
    fi
fi

if [ -f "${PROJECT_DIR}/.claude/.progress.md" ] \
        && [ -s "${PROJECT_DIR}/.claude/.progress.md" ]; then
    warnings+=("  .claude/.progress.md has content (dispatcher task log)")
fi

if [ ${#warnings[@]} -gt 0 ]; then
    echo "⚠  Uninstall will delete in-progress work:"
    for w in "${warnings[@]}"; do
        echo "$w"
    done
    echo ""
    echo "These files are gitignored, so this is the only copy."
    echo "If you want to preserve any of it, abort now and copy it"
    echo "somewhere outside the project before re-running uninstall."
    echo ""
fi

while : ; do
    read -r -p "Uninstall the kit AND tear down your workstation's local sandbox state? [yes/NO] " resp
    case "${resp:-}" in
        yes|Yes|YES) break ;;
        no|No|NO|"") echo "Aborted."; exit 0 ;;
        *) echo "?invalid response: ${resp}" ;;
    esac
done

# Tear down developer-local state by delegating to team/leave.sh.
# --yes bypasses leave.sh's own prompt since we already confirmed
# the combined action above. leave.sh also stops the sandbox before
# removing state, so we don't need a separate team/stop.sh call.
"${PROJECT_DIR}/team/leave.sh" --yes

# Delete tracked kit files. --ignore-unmatch makes this idempotent.
# Intentionally omitted from the delete list:
#   - .mcp.json and .claude/settings.json (preserved; see top-of-file note)
#   - docs/ and its contents (project-owned; never ours to touch)
git rm -rf --ignore-unmatch -- \
    CLAUDE_TEAM.md \
    ONBOARDING.md \
    TEAM_GUIDE.md \
    .claude/team-variables.yaml \
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
