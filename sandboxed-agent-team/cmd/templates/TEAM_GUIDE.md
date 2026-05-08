# Agent Team Guide — {{PROJECT_NAME}}

> **Generated:** {{UTC_TIMESTAMP}}
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the kit template source and
> re-run `agent-team-install`.

This document describes how to work with the Claude Code agent team
on this project. It is your day-to-day reference — not a setup guide
(see [`ONBOARDING.md`](ONBOARDING.md) for setup).

## Team Structure

The team has a **Lead** and seven **teammates**, each in their own Git
worktree:

- **Lead** — the main Claude Code session. Coordinates work, manages
  the lifecycle of requirements and tasks. Communicates with the
  human. Does not write files or run commands.
- **Integrator** — the Lead's operational arm. Task files, progress
  tracking, all git operations, PR lifecycle, cost recording. Also
  the default delegate for tasks that don't map to another teammate.
- **Analyst** — owns requirement docs in `docs/` and status tracking.
- **Architect** — architecture guardian; proposes design approaches and
  reviews code, does not write it.
- **Coder** — implements features and fixes bugs; runs lint/format
  on touched files at commit; runs the dependency audit when
  adding or removing a dependency. The Lead may spawn up to
  `{{MAX_PARALLEL_CODERS}}` Coders in parallel for file-disjoint
  subtasks within a single task.
- **Unit Tester** — unit and browserless UI tests.
- **E2E Tester** — end-to-end browser tests (Node.js Playwright).
- **Tech Writer** — owns `docs/guides/` (install / deploy / user /
  admin / operator guides). Updates on the release cadence, not the
  per-task cadence.

You only talk to the Lead. The Lead coordinates everything else.

## Daily Use

**Note:** Each teammate runs as a separate Claude Code instance
with its own context window. From the Lead's terminal you can
cycle through teammates with `Shift+Down` and message any of them
directly, but you typically don't need to — the Lead coordinates
the team for you.

1. At your host terminal (in the project directory), reattach to
   your sandbox: `./team/attach.sh`. This drops you into a Claude
   Code session running inside the sandbox. The session's system
   prompt auto-loads the Lead role, so the team comes up
   automatically on your first message — the Lead bootstraps and
   brings up the teammates for you. No slash command required.
   Once setup completes, the statusline shows "Agent Team Mode" as
   a visible confirmation that you're talking to the team.
   (If `attach.sh` reports no sandbox exists, run
   `./team/create.sh` to rebuild one.)
2. Describe what you want to the Lead. The Lead coordinates the team
   and drives the work — it does not implement directly.
3. You can switch between requirements and implementation freely.
   Requirements can be drafted for future tasks while a current task
   is being implemented. You can switch requirements topics at any
   time — just tell the Lead.
4. You review and approve requirement drafts and PRs when the Lead
   presents them. You may also provide feedback, answer questions the
   team surfaces, and perform any human-in-the-loop actions (e.g.,
   hardware passkey prompts during E2E testing). You may see multiple
   Coders and Unit Testers working simultaneously — this is by design
   when the Lead splits a task into parallel subtasks.
5. The Lead reports approximate cost per task (token usage and USD
   estimate per model, plus totals) at task wrap-up. You can also
   ask the Lead for the current cost at any time. Whether the
   cost report is also recorded in the final squash-merge commit
   message (for a persistent audit trail in git history) is a
   per-project setting (`Include cost report in commit message:
   yes|no` in `CLAUDE.md`'s Branching section). Ask the Lead to
   change it if you want a different default.
6. You can ask the team to take screenshots of the running
   application for visual verification — tell the Lead what you
   want to see.

## How Requirements Work

- Requirements are documented by the Analyst in `docs/`. All
  requirements originate from you — the Analyst formalizes them.
- Requirement branches are per-topic or related group (e.g.,
  `requirement/authentication`), not per individual requirement.
- New capabilities or constraints go through a requirement gate: the
  Analyst drafts, you approve. Refinements and preferences go directly
  to the Coder — no gate needed.
- You can switch topics at any time. The team tracks all in-flight
  requirement branches in `progress.md` so nothing gets lost.

## How Implementation Works

- Each task has a task branch, and each teammate gets a sub-branch
  for their work. Teammates merge into the task branch; the
  Integrator merges the task branch to `{{DEV_BRANCH_NAME}}`.
- Within a task, the Lead may split file-disjoint work across multiple
  Coders, each with a paired Unit Tester. Dependencies between
  subtasks are handled in phases.
- After each Coder commit: the Unit Tester and Architect review in
  parallel. After all work is merged: full test suites run at the
  pre-PR gate.
- The Lead may suspend a task to work on a prerequisite it discovered.
  This is normal — it will resume the original task after the
  prerequisite is complete.

## Visual Debugging

Agents can use the Playwright MCP server to interact with the running
application — navigate to pages, take screenshots, click elements,
and inspect visual state. This requires the dev server to be running.
Ask the Lead to "take a screenshot of [page]" to verify progress
without running the app yourself.

## If Something Goes Wrong

- **Teammate seems stuck or unresponsive:** Tell the Lead. The
  Lead will recover the teammate — resume first; spawn a
  replacement from the same agent definition if resume fails (see
  Teammate Recovery in `.claude/agents/lead.md`).
- **The Lead itself loses context mid-session:** Run
  `/project:lead-reload` at the sandbox's Claude Code prompt to
  re-invoke the Lead. The auto-load fires only at session start, so
  mid-session recovery uses the slash command. The Lead reads
  `progress.md` to recover state.
- **Sandbox crashes:** Back at your host terminal, run
  `./team/attach.sh` to reconnect (which reopens Claude Code
  inside the sandbox). The new session auto-loads the Lead. If the
  sandbox itself is gone, run `./team/create.sh` to rebuild it.
- **Sandbox authentication fails:** Refresh the token by running
  `claude` on the host and `/login` if needed, then re-run
  `./team/attach.sh`. This picks up fresh credentials on every
  attach. On macOS the new token is read automatically from the
  Keychain; on other systems update `CLAUDE_CODE_OAUTH_TOKEN` in
  your shell config first.
- **Sandbox git push/pull/fetch fails with 401 or 403:** The
  sandbox reaches the repo over HTTPS with a repo-platform API
  token. Most common causes: the token expired or was revoked, or
  its scopes are insufficient (Bitbucket needs Repositories R+W
  and Pull requests R+W; GitHub fine-grained PAT needs Contents
  R+W and Pull requests R+W; GitLab needs `api`). Re-run
  `./team/join.sh` to re-prompt for a fresh token.
- **Dev branch is broken:** The Lead will escalate to you. The
  breakage may be from the team's own merge or from external changes.
  You decide: wait (another team may already be fixing it), fix it
  with this team, or work on something else.

## Pausing and Resuming

Exiting Claude Code (`/exit`, `exit`, or Ctrl-D quickly twice) ends your Claude Code
session and drops you back to the host shell, but the sandbox VM keeps
running in the background. To resume:
1. At your host terminal: `./team/attach.sh` — it reattaches to
   the running sandbox and starts a new Claude Code session inside
   it.
2. The Lead auto-loads and reads `progress.md` to pick up where you
   left off.

## Session End

To end a Claude Code session cleanly, tell the Lead you're wrapping
up the session. The Lead confirms all work is merged and flags
anything unresolved for your next Claude Code session. Then exit
Claude Code (`/exit`, `exit`, or Ctrl-D quickly twice) — the sandbox VM keeps running so
you can reconnect later (see [Pausing and
Resuming](#pausing-and-resuming)).

## Engagement End

To end the engagement (i.e., destroy the sandbox), after ending
your final Claude Code session, at your host terminal run
`./team/destroy.sh` to destroy the sandbox VM. Host files remain.
Delete the project directory manually per your data retention
policy.
