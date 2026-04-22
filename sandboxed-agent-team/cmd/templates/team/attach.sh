#!/usr/bin/env bash

# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

# attach.sh — reattach to this project's already-running sandbox
# and drop into a Claude Code session.
#
# Fails fast if no sandbox exists (points you at create.sh) or if
# the Lead directive has changed since the sandbox was created
# (points you at destroy + create).
#
# This is a thin wrapper over create.sh, which implements both the
# create path and the reattach path behind a --skip-create flag.

exec "$(dirname "$0")/create.sh" --skip-create "$@"
