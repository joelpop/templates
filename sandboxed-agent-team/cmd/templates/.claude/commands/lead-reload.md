<!-- GENERATED FILE — do not edit directly. Edits here will be lost
the next time this file is regenerated. To change this file, edit
its template in the kit template source and
re-run `agent-team-install`. -->

# Reload Lead standing instructions

You are the team's Lead. Re-read your role definition and standing
instructions at `.claude/agents/lead.md` and resume operation per
its current contents. The team itself stays up; this command only
refreshes your own context.

Use this slash command when:

- Your context has been compacted and your standing instructions
  are no longer fresh.
- The human has updated `lead.md` (e.g., a kit upgrade) and wants
  you to pick up the new content immediately.
- The team is up but you've drifted from the documented behavior
  (e.g., handling a request without applying Request Triage; doing
  a teammate's work yourself; missing a doc-first-fix routing).

Do not run TeamCreate again from this command — the team is
already up. If `TeamCreate` has not been called this session,
treat the situation as a session-start instead and follow the
full Pre-Start Check in `.claude/agents/lead.md`.