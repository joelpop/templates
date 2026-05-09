# Agent Team Guide — {{PROJECT_NAME}}

> **Generated:** {{UTC_TIMESTAMP}}
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the kit template source and
> re-run `agent-team-install`.

Day-to-day reference for working with the Claude Code agent team
on this project. For setup, see [`ONBOARDING.md`](ONBOARDING.md).

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
for you.

1. At your host terminal (in the project directory), reattach:
   `./team/attach.sh`. The sandbox's Claude Code auto-loads the
   Lead role, so the team comes up on your first message — no
   slash command required. The statusline shows "Agent Team Mode"
   as visible confirmation. (If `attach.sh` reports no sandbox
   exists, run `./team/create.sh` to rebuild one.)
2. Describe what you want to the Lead. The Lead coordinates the
   team and drives the work; it does not implement directly.
3. Switch between requirements and implementation freely.
   Requirements can be drafted for future tasks while a current
   task is being implemented; you can switch requirement topics
   any time — just tell the Lead.
4. Review and approve requirement drafts and PRs when the Lead
   presents them. You may also provide feedback, answer questions
   the team surfaces, and perform human-in-the-loop actions (e.g.,
   hardware passkey prompts during E2E testing). You may see
   multiple Coders and Unit Testers working simultaneously — by
   design when the Lead splits a task into parallel subtasks.
5. The Lead reports approximate cost per task (per-model token
   usage and USD estimate, plus totals) at wrap-up; ask any time
   for current cost. Whether the report is also recorded in the
   squash-merge commit message (persistent audit trail) is a
   per-project setting (`Include cost report in commit message:
   yes|no` in `CLAUDE.md`'s Branching section) — ask the Lead to
   flip it if needed.
6. Ask the team to take screenshots of the running app for visual
   verification — tell the Lead what you want to see.

## How Requirements Work

- All requirements originate from you; the Analyst formalizes them
  in `docs/`.
- Requirement branches are per-topic or related group (e.g.,
  `requirement/authentication`), not per individual requirement.
- New capabilities or constraints go through a requirement gate
  (Analyst drafts, you approve). Refinements and preferences go
  directly to the Coder — no gate.
- Switch topics any time. The team tracks all in-flight
  requirement branches in `progress.md` so nothing gets lost.

## How Implementation Works

- Each task has a task branch, with a sub-branch per teammate.
  Teammates merge into the task branch; the Integrator merges the
  task branch to `{{DEV_BRANCH_NAME}}`.
- Within a task, the Lead may split file-disjoint work across
  multiple Coders, each with a paired Unit Tester. Subtask
  dependencies run in phases.
- After each Coder commit: Unit Tester and Architect review in
  parallel. After all work is merged: full test suites at the
  pre-PR gate.
- The Lead may suspend a task to work on a prerequisite it
  discovered, then resume after the prerequisite lands. Normal.

## Visual Debugging

Agents can use the Playwright MCP server to interact with the
running app — navigate, screenshot, click, inspect visual state.
Requires the dev server running. Ask the Lead to "take a screenshot
of [page]" to verify progress without running the app yourself.

## If Something Goes Wrong

- **Teammate seems stuck or unresponsive:** Tell the Lead. The
  Lead will recover — resume first; spawn a replacement from the
  same agent definition if resume fails (see Teammate Recovery in
  `.claude/agents/lead.md`).
- **The Lead itself loses context mid-session:** Run
  `/project:lead-reload` at the sandbox's Claude Code prompt. The
  auto-load fires only at session start, so mid-session recovery
  uses the slash command. The Lead reads `progress.md` to recover
  state.
- **Sandbox crashes:** Run `./team/attach.sh` from your host
  terminal to reconnect; the new session auto-loads the Lead. If
  the sandbox itself is gone, run `./team/create.sh` to rebuild.
- **Sandbox authentication fails:** Run `claude` on the host and
  `/login` if needed, then re-run `./team/attach.sh` — credentials
  are picked up fresh on every attach. On macOS the new token is
  read automatically from the Keychain; on other systems update
  `CLAUDE_CODE_OAUTH_TOKEN` in your shell config first.
- **Sandbox git push/pull/fetch fails with 401 or 403:** The
  sandbox reaches the repo over HTTPS with a repo-platform API
  token. Most common causes: token expired/revoked, or scopes
  insufficient (Bitbucket: Repositories R+W and Pull requests R+W;
  GitHub fine-grained PAT: Contents R+W and Pull requests R+W;
  GitLab: `api`). Re-run `./team/join.sh` to re-prompt.
- **Dev branch is broken:** The Lead escalates. The breakage may
  be from the team's own merge or from external changes. You
  decide: wait (another team may be fixing it), fix it with this
  team, or work on something else.

## Pausing and Resuming

Exit Claude Code (`/exit`, `exit`, or Ctrl-D quickly twice) — the
session ends and you're back at the host shell, but the sandbox
VM keeps running. To resume:

1. `./team/attach.sh` from your host terminal — reattaches to the
   running sandbox and starts a new Claude Code session inside it.
2. The Lead auto-loads and reads `progress.md` to pick up where
   you left off.

## Session End

Tell the Lead you're wrapping up. The Lead confirms all work is
merged and flags anything unresolved for next session. Then exit
Claude Code (`/exit`, `exit`, or Ctrl-D quickly twice) — the
sandbox VM keeps running so you can reconnect later (see
[Pausing and Resuming](#pausing-and-resuming)).

## Engagement End

To end the engagement, after your final Claude Code session run
`./team/destroy.sh` from the host terminal to destroy the sandbox
VM. Host files remain. Delete the project directory manually per
your data retention policy.
