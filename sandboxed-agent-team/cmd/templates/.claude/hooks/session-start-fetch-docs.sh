#!/bin/bash
# Session-start hook: inject canonical Agent Teams docs into the
# assistant's context.
#
# This kit uses Claude Code's experimental Agent Teams feature
# (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 in .claude/settings.json).
# The assistant's training-time memory has been observed to conflate
# Agent Teams with subagents, leading to confidently wrong answers
# about team mechanics. Fetching the live doc at every session start
# (and on resume / clear / compact) keeps answers grounded.
#
# Wired up by .claude/settings.json under hooks.SessionStart, with
# matchers: startup, resume, clear, compact. The script must be
# executable; the kit's installer / git tree preserves the mode bit.

set -u

# Clear any stale activation sentinel from a previous session. The
# sentinel is written by the Integrator after TeamCreate succeeds
# (see `.claude/agents/lead.md` → Team Initialization). Starting blank ensures
# the "Agent Team Mode" statusline indicator is accurate for this
# session — it lights up only after the team is actually live.
rm -f "${CLAUDE_PROJECT_DIR:-.}/.claude/.team-active"

DOC_URL="https://code.claude.com/docs/en/agent-teams.md"

echo "## Canonical Agent Teams reference (auto-injected by SessionStart hook)"
echo ""
echo "Source: $DOC_URL"
echo ""
echo "This project uses Claude Code's (currently) experimental Agent Teams feature."
echo "When answering any question about teammates, the Lead, the team"
echo "lifecycle, TeamCreate / SendMessage / TaskCreate-family tools,"
echo "the mailbox, or how teammates differ from subagents — defer to"
echo "the content below over training-time recall. Past sessions have"
echo "wasted hours and significant tokens on designs built atop"
echo "confidently-asserted mechanics that turned out to be wrong."
echo ""
echo "---"
echo ""

if ! curl -fsSL --max-time 10 "$DOC_URL"; then
  echo ""
  echo "(Fetch of $DOC_URL failed — network unavailable, doc moved, or"
  echo "URL changed. Verify against the live page before asserting any"
  echo "Agent Teams mechanics. Do not answer from training-time recall.)"
fi

exit 0