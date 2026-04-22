#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# attach.sh — reattach to this project's already-running sandbox
# and drop into a Claude Code session.
#
# Flags:
#   --resume   Resume the most recent Claude Code session (keeps
#              the conversation context from before you exited).
#   --fresh    Start a brand-new Claude Code session in the same
#              sandbox.
#
# --resume and --fresh are mutually exclusive. If neither is
# given, you'll be prompted.
#
# Fails fast if no sandbox exists (points you at create.sh) or if
# the Lead directive has changed since the sandbox was created
# (points you at destroy + create). The heavy lifting is in
# create.sh's --skip-create branch, which this script forwards to
# with the chosen session mode.

set -euo pipefail

mode=""
for arg in "$@"; do
    case "$arg" in
        --resume)
            if [ "$mode" = "fresh" ]; then
                echo "attach.sh: --resume and --fresh are mutually exclusive." >&2
                exit 2
            fi
            mode="resume"
            ;;
        --fresh)
            if [ "$mode" = "resume" ]; then
                echo "attach.sh: --resume and --fresh are mutually exclusive." >&2
                exit 2
            fi
            mode="fresh"
            ;;
        *)
            echo "attach.sh: unknown argument: $arg" >&2
            echo "Usage: $0 [--resume | --fresh]" >&2
            exit 2
            ;;
    esac
done

if [ -z "$mode" ]; then
    while : ; do
        read -r -p "Resume previous Claude Code session or start fresh? [resume/fresh] " resp
        case "${resp:-}" in
            r|R|resume|Resume|RESUME) mode="resume"; break ;;
            f|F|fresh|Fresh|FRESH)    mode="fresh";  break ;;
            *) echo "Please answer 'resume' or 'fresh'." ;;
        esac
    done
fi

exec "$(dirname "$0")/create.sh" --skip-create --"$mode"