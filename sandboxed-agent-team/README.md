# Sandboxed & Isolated Agent Team Kit
*Using a Docker sandbox and Git worktree isolation*

## Introduction

This kit sets up a structured Claude Code agent team inside an
isolated Docker sandbox for a software project, via a single
cross-platform binary. Templates and logic ship together; variables
live in the target project and survive kit upgrades.

Two layers of isolation: a **Docker sandbox** keeps all agent
activity, installed tools, and credentials off the developer's host
machine and separated from other projects; **Git worktrees** give
each teammate its own working copy of the repository so agents
never overwrite each other's in-progress work.

**In this README:**
- [Quick Start](#quick-start) — get up and running without reading the rest
- [Daily Use](#daily-use) — running the team after setup
- [Overview](#overview) — what the template provides
  - [Team Structure](#team-structure) — the Lead and seven teammates
  - [Features](#features)
    - [Capabilities](#capabilities) — isolation, status tracking, sub-task parallelism, task suspension, cost tracking, multi-developer support
    - [Workflows](#workflows) — coordination, requirements, branching, task lifecycle
    - [Guardrails](#guardrails) — quality, testing, dependency hygiene, context preservation
  - [Workspace Layout](#workspace-layout) — project directory structure
  - [Kit Contents](#kit-contents) — the nine files the template produces

## Quick Start

**Already have the agent team setup on your project?** You don't need this
template — look for `ONBOARDING.md` in the project root instead.

### Step 1 — Prepare

You'll need two kinds of things ready: infrastructure on your host
machine, and project-config values Claude Code will ask about.

**Infrastructure (host machine):**

- Docker Desktop installed and running
  (https://www.docker.com/products/docker-desktop/). The sandbox
  feature is required; `start.sh` verifies it at launch.
- Claude Code installed on the host and authenticated (run `claude`
  and `/login` if you haven't — you need this to execute this
  checklist in the first place).
- Git identity configured (`git config user.name "..."` and
  `git config user.email "..."`).
- The project directory is a Git repository (`git init` already run).
  Scenario A setup checks this up-front.
- A Git remote configured (`git remote add origin <url>`). Not
  strictly required at setup time — setup proceeds either way — but
  needed before you can use the **PR** merge method or push
  branches to collaborators.
- If the project's Git remote uses SSH: your SSH key available at
  the path referenced in `~/.ssh/config` for that remote.

**Project config (Claude Code auto-discovers most of these from
pom.xml and will ask you to confirm):**

- Java version
- Vaadin version
- Database (if any)
- CI platform
- Development branch name
- Merge method — how completed work reaches the development branch
  (PR, Integrator merge, Human merge, or your own method)
- If using the **PR** merge method: a platform API token so the Lead
  can create, read, and merge PRs. Bitbucket: app password
  (Settings → App passwords, with Repositories:Read and Pull
  requests:Read+Write). GitHub: fine-grained PAT. GitLab: personal
  access token. The setup walks you through creating one if you
  don't have it ready.

### Step 2 — Run setup

From the project root, run the `agent-team-install` binary — either
on `PATH` or via its explicit path:

```
agent-team-install
```

The binary is self-contained: every template file is bundled into
it via Go's `go:embed`, so it has no runtime dependency on the
kit source directory and can live anywhere.

The tool auto-detects the project's state and does the right thing:

- **No kit installed** → full setup: identifies the development branch
  (with fetch + freshness check), places on it, discovers stack
  details from `pom.xml`, prompts for what it can't auto-derive
  (CI platform, merge method, cost-in-commit), renders all template
  files into the project, commits them, and offers to run
  `./team/join.sh` to provision your workstation and start the team.
- **Kit already installed** → state-aware re-run: shuts down the
  sandbox, preserves your variables file (adding any new
  placeholders the current kit introduces, cleaning orphans),
  regenerates every generated file from current templates, commits
  the update, and — if your workstation is already provisioned —
  automatically re-provisions it to sync with the refreshed kit.

To remove the kit from a project: `./team/uninstall.sh`. Lifecycle
commands (`join`, `leave`, `start`, `stop`, `uninstall`) all live in
`./team/` and are the only supported way to invoke them — running
them directly keeps lifecycle commands in lockstep with the kit
version committed to the project.

### Step 3 — Onboarding other developers

New developers joining a project where the kit is already installed
run:

```
./team/join.sh
```

This provisions their local sandbox without re-running setup. To
undo just the local state (leaving versioned kit files alone), run
`./team/leave.sh`.

**A note on authentication:** Because the sandbox is a separate
environment, your host Claude Code login doesn't carry over.
`start.sh` autodetects your credentials and injects them into the
sandbox at each startup. On macOS, it extracts the OAuth token from
the macOS Keychain automatically — no manual token management needed.
On other systems, a one-time export of `CLAUDE_CODE_OAUTH_TOKEN` or
`ANTHROPIC_API_KEY` is required in your shell config.

**A note on SSH remotes:** If the project's Git remote uses SSH (e.g.,
`git@bitbucket.org:…` or a custom Host alias from `~/.ssh/config`),
the sandbox also needs the developer's SSH key, config, and
`known_hosts` to reach the remote. The setup process detects this
automatically and provisions the SSH material into `.sandbox/.ssh/`,
which is injected into the sandbox at each startup.

## Daily Use

Once agent team setup is complete:

**Note:** Teammates run as subagents within the Lead's session —
their work appears as expandable blocks in the same terminal. Each
agent does not get its own terminal pane.

1. At your host terminal (in the project directory), start the
   sandbox: `team/start.sh`. This drops you into a Claude Code
   session running inside the sandbox. The session's system prompt
   auto-loads the Lead role (see [Auto-loading Lead in sandbox
   sessions](#capabilities) under Capabilities), so the team spawns
   as soon as you send your first message — no slash command
   required. Once setup completes, the statusline shows "Agent Team
   Mode" as a visible confirmation that you're talking to the team.
2. The sandboxed Claude Code pre-authorizes common agent commands
   (`mvn`, `git`, `ls`, `chmod`, etc.) via `.claude/settings.json`'s
   allow/deny rules, so teammates can spawn agents, run builds,
   run tests, and perform routine git operations without being
   prompted. Destructive or out-of-scope operations
   (`git reset --hard`, `git push --force`, arbitrary shell through
   `curl | bash`, etc.) are explicitly denied by the same
   settings. The Lead will not implement directly — this is
   enforced by its instructions in `team-start.md`.
3. Describe what you want to the Lead. The Lead coordinates the team
   and drives the teammates through the workflows described in the
   Overview below.
4. You can switch between requirements and implementation freely.
   Requirements can be drafted for future tasks while a current task is
   being implemented. You can also switch requirements topics at any
   time — just tell the Lead. The team tracks all in-flight requirement
   branches so nothing gets lost.
5. You review and approve requirement drafts and PRs when the Lead
   presents them. You may also provide feedback, answer questions the
   team surfaces, and perform any human-in-the-loop actions (e.g.,
   hardware passkey prompts during E2E testing). You may see multiple
   Coders and Unit Testers working simultaneously in different panes —
   this is by design when the Lead splits a task into parallel subtasks.
6. The Lead reports approximate cost per task (token usage and USD
   estimate per model, plus totals) at task wrap-up. You can also
   ask the Lead for the current cost at any time.
7. You can ask agents to take screenshots of the running application
   for visual verification — tell the Lead what you want to see.
8. **If something goes wrong:**
   - Agent seems stuck or unresponsive: tell the Lead. The Lead will
     respawn the agent.
   - The Lead itself loses context mid-session: run
     `/project:team-start` at the sandbox's Claude Code prompt to
     re-invoke the Lead (the auto-load fires only at session start,
     so mid-session recovery uses the slash command).
   - Sandbox crashes: back at your host terminal, run
     `team/start.sh` to reconnect (which reopens Claude Code
     inside the sandbox). The new session auto-loads the Lead, which
     reads `progress.md` to recover state.
   - The Lead may suspend a task to work on a prerequisite it
     discovered — this is normal. It will resume the original task
     after the prerequisite is complete.
9. **Pausing and resuming:** Exiting Claude Code (`/exit` or Ctrl+D)
    ends your Claude Code session and drops you back to the shell, but the sandbox VM keeps running in the background.
    To resume: at your host terminal run `team/start.sh` again
    — it detects the existing sandbox, connects you to it, and starts a new Claude Code
    session inside it. The Lead auto-loads and reads `progress.md`
    to pick up where you left off.
10. To end a Claude Code session cleanly, tell the Lead you're
    wrapping up the session. The Lead confirms all work is merged and flags
    anything unresolved for your next Claude Code session. Then
    exit Claude Code (`/exit` or Ctrl+D) — the sandbox VM keeps
    running so you can reconnect later.
11. To end the engagement (i.e., destroy the sandbox), after ending your final Claude Code
    session, at your host terminal: `team/stop.sh`

## Overview

### Team Structure

The team has a **Lead** and seven **teammates**, each in their own Git
worktree:

- **Lead** — the main Claude Code session. Coordinates work, manages
  the lifecycle of requirements and tasks. Communicates with the
  human. Does not write files or run commands — delegates all
  operational and application work to teammates.
- **Integrator** — the Lead's operational arm. Owns task files,
  progress tracking, all git operations, PR lifecycle (via platform
  API), and cost recording. Also the default delegate for tasks that
  don't clearly map to another teammate.
- **Analyst** — owns requirement docs in `docs/` and status tracking.
- **Architect** — architecture guardian; proposes design approaches and reviews code, does not write it.
- **Coder** — implements features and fixes bugs.
- **Janitor** — linting, cleanup, dependency hygiene.
- **Unit Tester** — unit and browserless UI tests.
- **E2E Tester** — end-to-end browser tests.

### Features

#### Capabilities

- **Isolation & Infrastructure** — Each project gets its own Docker
  sandbox built from a customizable Dockerfile. One-command scripts
  handle startup and disposal. Claude Code authentication is
  autodetected and injected via environment variable. SSH keys for
  Git remote access are provisioned into the sandbox automatically.
- **Auto-loading Lead in sandbox sessions** — The sandbox's Claude
  Code session starts with the Lead role pre-configured:
  `team/start.sh` passes `--append-system-prompt` to `claude` so
  the first turn reads `team-start.md` and spawns the team
  automatically. The human does not need to remember
  `/project:team-start`. Host Claude Code sessions are unaffected
  (they don't go through `start.sh`). The `/project:team-start`
  slash command remains available as a manual re-invocation fallback
  if the Lead needs to be reset mid-session.
- **"Agent Team Mode" statusline indicator** — The sandbox's
  statusline displays "Agent Team Mode" once the Lead has completed
  the Pre-Start Check and spawned the team, giving a visible cue at
  the keyboard that the human is interacting with the team (not
  bare Claude Code). Implemented via a `statusLine` entry in
  `.claude/settings.json` that checks for a sentinel file
  (`.claude/.team-active`) written by the Lead at the end of team
  spawn. The indicator is blank before setup completes and between
  sessions; it updates each session based on the current state.
- **Status Tracking** — Requirement status checkboxes (`[ ]`/`[-]`/`[x]`)
  in `docs/` plus role-assigned plan steps in task files. A progress
  dispatcher tracks active and suspended tasks for recovery after
  context compaction.
- **Sub-Task Parallelism** — Within a single task, the Lead may split
  file-disjoint implementation work across multiple Coders, each with
  a paired Unit Tester. Phases support dependencies between subtasks.
  Roles also work in parallel where possible: the Unit Tester and
  Architect review simultaneously after Coder work is merged.
- **MCP Documentation Servers** — Agents consult MCP servers for
  authoritative framework documentation (Java, Vaadin, Spring,
  Playwright) rather than relying on training data. The Playwright
  MCP server also provides visual debugging — agents can navigate the
  running application, take screenshots, and interact with the UI.
- **Task Suspension & Resumption** — The Lead formally suspends a task
  when a prerequisite is discovered mid-work, preserving the branch and
  status. Resumption merges the latest development branch in and
  continues from the first incomplete step. Nested suspension is
  supported.
- **Cost Tracking** — At task kickoff the Integrator captures a
  `ccusage daily` baseline; at task conclusion (T.6) it runs
  `ccusage daily` again, subtracts the baseline, and formats a
  per-model + totals cost report. The Lead always reports the
  numbers to the human verbally. Recording the report in the final
  squash-merge commit message (so it persists in git history) is a
  **per-project setting** in `CLAUDE_TEAM.md`'s Branching section
  (`Include cost report in commit message: yes|no`), chosen at
  setup and changeable later by asking the Lead. Works for both
  API-key and subscription users because `ccusage` reads Claude
  Code's local JSONL session logs directly, sidestepping the
  billing-mode-dependent `/cost` command (which returns no numbers
  for subscription users). `ccusage` inside the sandbox sees only
  sandbox sessions — the human's concurrent host Claude Code work
  is naturally excluded.
- **Multi-Developer Support** — Shared context (`CLAUDE_TEAM.md`,
  `team-start.md`) is version-controlled; developer-local state (tasks,
  progress, worktrees, settings) is gitignored. An `ONBOARDING.md` is
  generated for new developers.

#### Workflows

- **Coordination Model** — Teammates message each other directly for
  routine coordination. They escalate to the Lead when a decision
  requires human input or intervention. The human only talks to the
  Lead, who coordinates but does not implement directly.
- **Requirements Management** — All requirements originate from the
  human and are documented by the Analyst. New capabilities go through
  a requirement gate (Analyst drafts, human approves). Agents must
  escalate ambiguity — guessing is forbidden. Refinements and
  preferences bypass the gate and go directly to the Coder.
  Requirement branches are per-topic or related group (e.g.,
  `requirement/authentication`), not per individual requirement — the
  Analyst freely splits, merges, and cross-references within a group.
  The human can switch requirement topics at any time. The Analyst can
  draft requirements for future tasks while the current task is being
  implemented — requirement branches and task branches are independent.
- **Branching & Merging** — Work branches off a configurable development
  branch. Requirement changes, implementation tasks, and individual
  agent roles each get dedicated branches. Agents merge (never rebase).
  All merges to the development branch are squash merges. The merge
  method (PR, Integrator merge, human merge, or custom) is configured per
  project.
- **Task Lifecycle** — Three workflows: the *requirements workflow*
  (classify, draft, approve, merge), the *implementation workflow*
  (task branch, dependency audit, per-commit review cycle, pre-PR gate,
  human validation), and the *integration merge workflow* (incorporate
  upstream changes, resolve conflicts, re-test, finalize).

#### Guardrails

- **Quality** — The Coder must consult framework documentation before
  writing UI code (training data is not authoritative). A diagnosis-first
  fix protocol with a two-attempt limit prevents spiraling. Workarounds
  require Architect approval. The Architect reviews every commit for
  incremental rot, cross-cutting drift, and cohesion decay.
- **Testing Strategy** — Tests follow a pyramid: unit, browserless UI,
  then E2E (browser-only scenarios). The Unit Tester owns all tests by
  default and delegates to the E2E Tester. Per-commit runs are targeted;
  full suites run at the pre-PR gate. Human-in-the-loop E2E steps use a
  structured pause/resume cycle.
- **Dependency & Code Hygiene** — The Janitor audits dependencies before
  every task, after every dependency change, and after every merge. CVEs
  block merging. Version upgrades follow pinning rules — patch upgrades
  are safe, minor upgrades follow pinning rules, major upgrades need
  approval. Linting and dead code detection are also Janitor-owned.
- **Dev-Branch Health** — When the development branch is broken (by
  the team's own merge or by external changes), the Lead escalates to
  the human and holds off on new work until the issue is resolved.
- **Context Preservation** — Claude Code may silently compact context,
  dropping loaded files. Every agent must re-read a defined set of files
  before starting any task. Agents in worktrees access gitignored files
  via absolute project root path.

### Workspace Layout

This template assumes the following workspace layout:

```
~/workspaces/
├── acme-corp/          # Customer A — may contain multiple projects
│   ├── project-alpha/
│   └── project-beta/
├── widgets-inc/        # Customer B
│   └── main-app/
├── vaadin/             # Your company
│   └── ...
└── personal/           # Personal projects
    └── ...
```

Infrastructure files go **inside** each project alongside the code:

```
~/workspaces/acme-corp/project-alpha/
├── .sandbox/                        # ← All sandbox infra lives here
│   ├── Dockerfile                   # Custom Docker Sandbox template
│   ├── start.sh                     # One-command startup
│   └── stop.sh                      # Sandbox disposal
├── CLAUDE.md                        # Project-owned (kit imports CLAUDE_TEAM.md here)
├── CLAUDE_TEAM.md                   # Kit-owned agent context
├── ONBOARDING.md                    # Developer onboarding (generated)
├── TEAM_GUIDE.md                    # Daily-use reference for humans (generated)
├── .claude/
│   ├── .last-onboarded              # Records when this developer last completed onboarding (compared against ONBOARDING.md's Generated: to detect staleness)
│   ├── settings.json                # Agent Teams + permissions
│   ├── progress.md                  # Dispatcher: which task is active, which are suspended
│   ├── tasks/                       # One file per active or suspended task
│   └── commands/
│       └── team-start.md            # Reusable slash command
├── docs/
│   ├── INDEX.md                     # ← Master doc index (sample, project-owned)
│   ├── agnostic/                    # Project-agnostic patterns / preferences / standing guidance
│   │   ├── patterns.md
│   │   ├── preferences.md
│   │   └── standards.md
│   └── reqs/                        # Project-specific requirements
│       ├── architecture-debt.md     # Structural debt findings
│       ├── non-functional/          # Quality attributes (ISO 25010)
│       │   ├── performance.md
│       │   ├── security/            # Auth, authz, hardening, data protection
│       │   ├── reliability.md
│       │   └── ...
│       ├── functional/
│       │   ├── cross-cutting/       # Error handling, validation, APIs, etc.
│       │   ├── data/                # Schema, migrations
│       │   └── features/            # Feature docs + supplementals
│       │       ├── feature-a.md
│       │       ├── feature-a/       # views.md, ux.md, etc.
│       │       └── ...
│       ├── external-interfaces/     # UI, software, communication interfaces
│       ├── environmental/           # Infrastructure, platforms, deployment
│       └── technical/               # Stack, build, constraints
├── .gitignore                       # Kit manages a bracketed block of developer-local patterns
└── (existing project files)
```

### Kit Contents

The kit produces these files in a target project:

| Path | Purpose | Usage |
|------|---------|-------|
| `.sandbox/Dockerfile` | Custom sandbox image for this project | Built automatically by `start.sh` |
| `team/start.sh` | One-command sandbox build + startup | Human runs at host terminal |
| `team/stop.sh` | Sandbox disposal | Human runs at host terminal |
| `docs/INDEX.md` | Sample requirement-document index | Seeded once on initial setup; project-owned thereafter (re-setup and remove leave it alone) |
| `CLAUDE_TEAM.md` | Project context for agents (kit-owned) | Imported into `CLAUDE.md` via a bracketed `@CLAUDE_TEAM.md` line |
| `CLAUDE.md` | Project-owned context file | Kit adds/removes only the bracketed import line; everything else is yours |
| `.claude/team-variables.yaml` | Per-project persisted variables | Human-readable, hand-editable, survives kit upgrades |
| `.claude/settings.json` | Agent team config and permissions | Auto-loaded by Claude Code at session start |
| `.mcp.json` | Project-scoped MCP server config | Auto-loaded by Claude Code at session start (canonical location for project-scoped MCP servers; see Claude Code docs) |
| `.claude/commands/team-start.md` | Lead's operating manual | Auto-loaded by the sandboxed Claude Code at session start (via `--append-system-prompt` in `start.sh`); also exposed as `/project:team-start` for manual re-invocation |
| `ONBOARDING.md` | Developer onboarding (generated) | New developer runs `./team/join.sh` |
| `TEAM_GUIDE.md` | Daily-use reference for humans (generated) | Human reads for workflows, troubleshooting, recovery |

