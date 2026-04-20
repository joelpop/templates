# Agent Team Guide — <PROJECT_NAME>

> **Generated:** <DATE>
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the team setup kit (SANDBOXED_AGENT_TEAMS.md) and
> re-run the setup at your host terminal.

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
- **Coder** — implements features and fixes bugs. The Lead may spawn
  multiple Coders for parallel subtasks within a single task.
- **Janitor** — linting, cleanup, dependency hygiene.
- **Unit Tester** — unit and browserless UI tests.
- **E2E Tester** — end-to-end browser tests (Node.js Playwright).

You only talk to the Lead. The Lead coordinates everything else.

## Daily Use

**Note:** Teammates run as subagents within the Lead's session —
their work appears as expandable blocks in the same terminal. Each
agent does not get its own terminal pane.

1. At your host terminal (in the project directory), start the
   sandbox: `.sandbox/start.sh`. This drops you into a Claude Code
   session running inside the sandbox. The session's system prompt
   auto-loads the Lead role, so the team spawns as soon as you send
   your first message — no slash command required. Once setup
   completes, the statusline shows "Agent Team Mode" as a visible
   confirmation that you're talking to the team.
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
6. You can ask agents to take screenshots of the running application
   for visual verification — tell the Lead what you want to see.

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

- Each task has a task branch, and each agent gets a sub-branch for
  their work. Agents merge into the task branch; the Integrator merges
  the task branch to `<DEV_BRANCH_NAME>`.
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

- **Agent seems stuck or unresponsive:** Tell the Lead. The Lead will
  respawn the agent.
- **The Lead itself loses context mid-session:** Run
  `/project:team-start` at the sandbox's Claude Code prompt to
  re-invoke the Lead. The auto-load fires only at session start, so
  mid-session recovery uses the slash command. The Lead reads
  `progress.md` to recover state.
- **Sandbox crashes:** Back at your host terminal, run
  `.sandbox/start.sh` to reconnect (which reopens Claude Code inside
  the sandbox). The new session auto-loads the Lead.
- **Sandbox authentication fails:** At your host terminal, stop the
  sandbox with `docker sandbox stop <name>` (the name is printed by
  `start.sh` at startup), re-run `claude` on the host and `/login` if
  needed, then restart with `.sandbox/start.sh`. This can happen if
  the OAuth token expired or if you ran `/login` on the host while
  the sandbox was running. On macOS, `start.sh` automatically picks
  up the fresh token from the Keychain.
- **SSH remote access fails:** If `git push/pull/fetch` fails with
  "Permission denied" or "Host key verification failed", check that
  `.sandbox/ssh/` contains the correct key, config, and known_hosts.
  If the key was rotated, copy the new key to `.sandbox/ssh/` (or
  re-run onboarding) and restart the sandbox.
- **Dev branch is broken:** The Lead will escalate to you. The
  breakage may be from the team's own merge or from external changes.
  You decide: wait (another team may already be fixing it), fix it
  with this team, or work on something else.

## Pausing and Resuming

Exiting Claude Code (`/exit` or Ctrl+D) ends your Claude Code
session and drops you back to the shell, but the sandbox VM keeps
running in the background. To resume:
1. At your host terminal: `.sandbox/start.sh` again — it detects the
   existing sandbox, connects you to it, and starts a new Claude
   Code session inside it.
2. The Lead auto-loads and reads `progress.md` to pick up where you
   left off.

## Session End

To end a Claude Code session cleanly, tell the Lead you're wrapping
up the session. The Lead confirms all work is merged and flags
anything unresolved for your next Claude Code session. Then exit
Claude Code (`/exit` or Ctrl+D) — the sandbox VM keeps running so
you can reconnect later (see [Pausing and
Resuming](#pausing-and-resuming)).

## Engagement End

To end the engagement (i.e., destroy the sandbox), after ending
your final Claude Code session, at your host terminal run
`.sandbox/stop.sh` to destroy the sandbox VM. Host files remain.
Delete the project directory manually per your data retention
policy.
