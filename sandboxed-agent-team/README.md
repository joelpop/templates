# Sandboxed & Isolated Agent Team Kit
*Using a Docker sandbox and Git worktree isolation*

## Introduction

This kit sets up a structured Claude Code agent team inside an
isolated Docker sandbox, via a single cross-platform binary.
Templates and logic ship together; variables live in the target
project and survive kit upgrades.

Two layers of isolation: a **Docker sandbox** keeps all team
activity, tools, and credentials off the host machine and
separated from other projects; **Git worktrees** give each
teammate its own working copy of the repo so teammates never
overwrite each other's in-progress work.

**Key terms:**

- **Sandbox** — the Docker-based isolated environment where the
  team runs.
- **Lead directive** — the instructions the sandboxed Claude Code
  auto-loads on session start, so the Lead takes over and the team
  comes up on your first message — no slash command required.
- **Repo-platform API token** — a single credential from whichever
  service hosts your repository (the *repo platform*: Bitbucket,
  GitHub, or GitLab). Bitbucket calls it an app password; GitHub
  and GitLab call it a personal access token (PAT). Required for
  HTTPS git access from the sandbox; also used for PR API calls.

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

Two kinds of things to have ready: host infrastructure, and
project-config values Claude Code will ask about.

**Infrastructure (host machine):**

- Docker Desktop installed and running
  (https://www.docker.com/products/docker-desktop/). The sandbox
  feature is required; `create.sh` and `attach.sh` verify it at
  launch.
- Claude Code installed on the host and authenticated (`claude`
  and `/login` if you haven't — needed to execute this checklist).
- Git identity configured (`git config user.name "..."` and
  `git config user.email "..."`).
- The project directory is a Git repository (`git init` already
  run). Scenario A setup checks this up-front.
- A Git remote configured (`git remote add origin <url>`). Not
  strictly required at setup — setup proceeds either way — but
  needed before you can use the **PR** merge method or push
  branches to collaborators.

**Project config** (the installer auto-discovers from `pom.xml`
and git config where it can; you can accept or override each
prompt):

- **Team's development branch name** — auto-detected when a
  conventional name (`develop`, `dev`, `feature/develop`) is
  present, or inferred from your remote's default branch.
- **Stack details from `pom.xml`** — Java version, Vaadin version,
  Spring Boot version, JUnit version, database, build tool.
- **Merge method** — how completed work reaches the development
  branch (PR, Integrator merge, Human merge).
- **CI platform** — GitHub Actions, GitLab CI, Bitbucket Pipelines,
  Jenkins, or none.
- **Whether to include a cost report in squash-merge commit
  messages.**
- **Preference to set up your local sandbox** (run
  `./team/join.sh`) after install.
- **A repo-platform API token** (Bitbucket app password, GitHub
  fine-grained PAT, or GitLab PAT) — required unless the repo is
  public. See *"A note on the sandbox's git access"* below for
  rationale and storage. Scopes: Bitbucket — Repositories R+W and
  Pull requests R+W; GitHub fine-grained PAT — Contents R+W and
  Pull requests R+W; GitLab — `api`. Setup walks you through
  creation with direct links; leave blank for a public repo.

### Step 2 — Run setup

From the project root, run the `agent-team-install` binary — either
on `PATH` or via its explicit path:

```
agent-team-install
```

The binary is self-contained: every template is bundled via Go's
`go:embed`, so it has no runtime dependency on the kit source
directory.

The tool auto-detects the project's state:

- **No kit installed** → full setup: discovers stack details from
  `pom.xml`, prompts for what it can't auto-derive (dev branch,
  CI platform, merge method, cost-in-commit), writes the kit
  files, and offers to run `./team/join.sh`.
- **Kit already installed** → state-aware re-run: shuts down the
  sandbox, preserves your variables file (adding new placeholders,
  cleaning orphans), regenerates every generated file from current
  templates, and re-provisions an already-provisioned workstation
  to sync with the refreshed kit.

**No git operations.** The installer writes files; it does not
switch branches, stage, commit, merge, or push. Run `git status`
after install and commit on your own schedule.

To remove the kit: `./team/uninstall.sh`. This chains
`./team/leave.sh --yes` to discard local sandbox state, then
deletes kit files and excises the `CLAUDE.md` import block and
the kit's `.gitignore` block. Like install, no git operations —
review with `git status` and commit when ready.

Lifecycle commands (`join`, `leave`, `start`, `stop`, `uninstall`)
all live in `./team/` so they stay in lockstep with the kit
version committed to the project.

### Step 3 — Onboarding other developers

New developers joining a project where the kit is already
installed run:

```
./team/join.sh
```

This provisions their local sandbox without re-running setup. To
undo just the local state (leaving versioned kit files alone),
run `./team/leave.sh`.

**A note on authentication:** The sandbox is a separate
environment, so your host Claude Code login doesn't carry over.
The kit autodetects credentials and injects them into the sandbox
on every create and attach. On macOS, it extracts the OAuth token
from the Keychain automatically — no manual token management. On
other systems, a one-time export of `CLAUDE_CODE_OAUTH_TOKEN` or
`ANTHROPIC_API_KEY` in your shell config is required.

**A note on the sandbox's git access:** Docker Sandbox blocks
outbound port 22 (SSH) for isolation. Regardless of how your host
reaches the repo, teammates *inside* the sandbox use **HTTPS** for
git operations. The kit handles this transparently — your host's
git config and origin URL are never modified.

What the kit does, once per project, on your first
`./team/join.sh`:

1. Detects the platform hosting `origin` (Bitbucket / GitHub /
   GitLab). SSH aliases (e.g., `git@bitbucket-syntech:…`) are
   resolved via `~/.ssh/config`'s `HostName` directive.
2. Prompts you to create and paste a repo-platform API token
   (Bitbucket app password, GitHub fine-grained PAT, or GitLab
   PAT) with a direct link to the token-creation page for the
   detected platform.
3. Stores the token:
   - **macOS** → Keychain, service name
     `agent-team.<sandbox-name>`, one entry per project
     (multi-project isolation). Never touches the regular
     filesystem.
   - **Linux / Windows** → `.sandbox/.repo-platform-api.env`
     (mode 600, gitignored). A `join.sh` banner flags this as
     less hardened than Keychain; credential-manager integration
     (libsecret, Credential Manager) is a planned follow-up.
4. On every `create.sh` / `attach.sh`, reads the token back and
   pipes it via stdin into the sandbox — never on a command line
   visible to `ps`. Inside the sandbox the kit writes
   `~/.git-credentials` and sets
   `git config --global credential.helper store`.
5. If your host's `origin` URL is SSH-shaped (`git@host:…` or
   `ssh://git@host/…`), the sandbox's git also gets a
   `url.https://<host>/.insteadOf git@<alias>:` rewrite — so SSH
   URLs in your committed `.git/config` transparently resolve to
   HTTPS inside the sandbox.

**Token persistence:** survives `./team/destroy.sh`,
`agent-team-install` updates, and re-runs of `./team/join.sh`.
Wiped by `./team/leave.sh` and `./team/uninstall.sh` (chains
leave). On macOS, the Keychain entry persists even across
`leave.sh`; on Linux/Windows the file is wiped and the next
`join.sh` re-prompts.

**Public-repo opt-out:** press Enter without pasting at the token
prompt. Teammates will be able to read public repos but not push.

## Daily Use

Once agent team setup is complete:

**Note:** Each teammate runs as a separate Claude Code instance
with its own context window. From the Lead's terminal you can
cycle through teammates with `Shift+Down` and message any directly,
but you typically don't need to — the Lead coordinates for you.

1. At your host terminal (in the project directory), reattach:
   `./team/attach.sh`. The sandboxed Claude Code auto-loads the
   Lead role (see [Auto-loading Lead in sandbox
   sessions](#capabilities)), so the team comes up on your first
   message — no slash command required. The statusline shows
   "Agent Team Mode" as visible confirmation. (If `attach.sh`
   reports no sandbox — e.g., because `destroy.sh` was run — use
   `./team/create.sh` to build a fresh one.)
2. The sandboxed Claude Code pre-authorizes common shell commands
   (`mvn`, `git`, `ls`, `chmod`, etc.) via `.claude/settings.json`
   allow/deny rules, so teammates can run builds, tests, and
   routine git without prompts. Destructive operations
   (`git reset --hard`, `git push --force`, arbitrary shell via
   `curl | bash`) are explicitly denied. The Lead won't implement
   directly — enforced by its instructions in
   `.claude/agents/lead.md`.
3. Describe what you want to the Lead. The Lead coordinates the
   team through the workflows in the Overview below.
4. Switch between requirements and implementation freely.
   Requirements can be drafted for future tasks while a current
   one is being implemented; switch topics any time — just tell
   the Lead. The team tracks all in-flight requirement branches.
5. Review and approve requirement drafts and PRs when the Lead
   presents them. You may also provide feedback, answer team
   questions, and perform human-in-the-loop actions (e.g.,
   hardware passkey prompts during E2E testing). Multiple Coders
   and Unit Testers may run simultaneously when the Lead splits a
   task into parallel subtasks — by design.
6. The Lead reports approximate cost per task (per-model token
   usage and USD, plus totals) at wrap-up; ask any time for
   current cost.
7. Ask the team to screenshot the running app for visual
   verification — tell the Lead what you want to see.
8. **If something goes wrong:**
   - Teammate stuck or unresponsive: tell the Lead. The Lead
     recovers (resume first; replace from the agent definition if
     resume fails — see Teammate Recovery in
     `.claude/agents/lead.md`).
   - The Lead loses context mid-session: run
     `/project:lead-reload` at the sandbox's Claude Code prompt
     (auto-load fires only at session start, so mid-session
     recovery uses the slash command).
   - Sandbox crashes: run `./team/attach.sh` from your host
     terminal to reconnect; the new session auto-loads the Lead,
     which reads `progress.md` to recover. If the sandbox itself
     is gone, `./team/create.sh` to rebuild.
   - The Lead may suspend a task for a discovered prerequisite,
     then resume — normal.
9. **Pausing and resuming:** Exit Claude Code (`/exit`, `exit`, or
   Ctrl-D quickly twice) — the session ends, the sandbox VM keeps
   running. To resume: `./team/attach.sh` reattaches and starts a
   new Claude Code session inside it; the Lead auto-loads and
   reads `progress.md`.
10. To end a Claude Code session cleanly, tell the Lead you're
    wrapping up. The Lead confirms all work is merged and flags
    anything unresolved for next session. Then exit Claude Code
    — the sandbox VM keeps running.
11. To end the engagement (destroy the sandbox), after your final
    Claude Code session: `./team/destroy.sh` from the host
    terminal.

## Overview

### Team Structure

The team has a **Lead** and seven **teammates**, each in its own
Git worktree:

- **Lead** — the main Claude Code session. Coordinates work,
  manages the lifecycle of requirements and tasks, talks to the
  human. Doesn't write files or run commands — delegates all
  operational and application work to teammates.
- **Integrator** — the Lead's operational arm. Owns task files,
  progress tracking, all git operations, PR lifecycle (via
  platform API), cost recording. Default delegate for tasks that
  don't clearly map to another teammate.
- **Analyst** — owns requirement docs in `docs/` and status
  tracking.
- **Architect** — architecture guardian; proposes design
  approaches and reviews code, does not write it.
- **Coder** — implements features and fixes bugs; runs
  lint/format/analysis on touched files at commit; runs
  dependency audit when adding or removing a dependency.
- **Unit Tester** — unit and browserless UI tests.
- **E2E Tester** — end-to-end browser tests.
- **Tech Writer** — owns `docs/guides/` (install / deploy / user
  / admin / operator). Updates on the release cadence.

### Features

#### Capabilities

- **Isolation & Infrastructure** — Each project gets its own
  Docker sandbox built from a customizable Dockerfile. One-command
  scripts handle startup and disposal. Claude Code authentication
  is autodetected and injected via environment variable. The
  sandbox reaches the repo over HTTPS (Docker Sandbox blocks
  outbound port 22); a repo-platform API token is provisioned
  automatically for private-repo access — see *A note on the
  sandbox's git access* above.
- **Auto-loading Lead in sandbox sessions** — The sandboxed
  Claude Code starts with the Lead role pre-configured:
  `team/create.sh` passes `--append-system-prompt` to `claude` so
  the first turn reads `.claude/agents/lead.md` and brings up the
  team via `TeamCreate`. No need to remember
  `/project:lead-reload`. Host Claude Code sessions are
  unaffected (they don't go through `create.sh`/`attach.sh`).
  `/project:lead-reload` remains available as a manual fallback
  if the Lead needs to be reset mid-session.
- **"Agent Team Mode" statusline indicator** — The sandbox's
  statusline shows "Agent Team Mode" once the Lead has completed
  the Pre-Start Check and brought up the team. Implemented via a
  `statusLine` entry in `.claude/settings.json` that checks for a
  sentinel file (`.claude/.team-active`) written by the Lead
  after `TeamCreate` succeeds. Blank before setup completes and
  between sessions; updates each session based on current state.
- **Status Tracking** — Requirement status checkboxes
  (`[ ]`/`[-]`/`[x]`) in `docs/` plus role-assigned plan steps in
  task files. A progress dispatcher tracks active and suspended
  tasks for recovery after context compaction.
- **Sub-Task Parallelism** — Within a single task, the Lead may
  split file-disjoint implementation work across multiple Coders,
  each with a paired Unit Tester. Phases support dependencies
  between subtasks. Roles also work in parallel where possible:
  Unit Tester and Architect review simultaneously after Coder
  work is merged.
- **MCP Documentation Servers** — Teammates consult MCP servers
  (Java, Vaadin, Spring, Playwright) for authoritative framework
  docs rather than relying on training data. Playwright also
  provides visual debugging — navigate the running app,
  screenshot, interact with the UI.
- **Task Suspension & Resumption** — The Lead formally suspends a
  task when a prerequisite is discovered mid-work, preserving the
  branch and status. Resumption merges the latest dev branch in
  and continues from the first incomplete step. Nested suspension
  supported.
- **Cost Tracking** — At task kickoff the Integrator captures a
  `ccusage daily` baseline; at conclusion (T.6) it runs `ccusage
  daily` again, subtracts the baseline, and formats a per-model +
  totals report. The Lead always reports verbally. Recording in
  the squash-merge commit message (persistent in git history) is
  a per-project setting in `CLAUDE_TEAM.md`'s Branching section
  (`Include cost report in commit message: yes|no`), chosen at
  setup and changeable by asking the Lead. Works for both
  API-key and subscription users — `ccusage` reads Claude Code's
  local JSONL session logs directly, sidestepping the
  billing-mode-dependent `/cost` command (which returns no
  numbers for subscription users). Inside the sandbox `ccusage`
  sees only sandbox sessions — the human's concurrent host work
  is naturally excluded.
- **Multi-Developer Support** — Shared context (`CLAUDE_TEAM.md`,
  `.claude/agents/lead.md`) is version-controlled; developer-local
  state (tasks, progress, worktrees, settings) is gitignored. An
  `ONBOARDING.md` is generated for new developers.

#### Workflows

- **Coordination Model** — Teammates message each other directly
  for routine coordination, escalating to the Lead when a decision
  needs human input. The human only talks to the Lead, who
  coordinates but does not implement directly.
- **Requirements Management** — All requirements originate from
  the human; the Analyst documents them. New capabilities go
  through a requirement gate (Analyst drafts, human approves).
  Teammates must escalate ambiguity — guessing is forbidden.
  Refinements and preferences bypass the gate and go directly to
  the Coder. Requirement branches are per-topic or related group
  (e.g., `requirement/authentication`), not per individual
  requirement — the Analyst freely splits, merges, and
  cross-references within a group. The human can switch topics
  any time. The Analyst can draft for future tasks while the
  current one is being implemented — requirement branches and
  task branches are independent.
- **Branching & Merging** — Work branches off a configurable dev
  branch. Requirement changes, implementation tasks, and per-role
  teammate work each get dedicated branches. Teammates merge
  (never rebase). All merges to the dev branch are squash merges.
  Merge method (PR, Integrator merge, human merge, or custom) is
  configured per project.
- **Task Lifecycle** — Three workflows: *requirements*
  (classify, draft, approve, merge), *implementation* (task
  branch, dependency audit, per-commit review cycle, pre-PR gate,
  human validation), *integration merge* (incorporate upstream,
  resolve conflicts, re-test, finalize).

#### Guardrails

- **Quality** — The Coder consults framework documentation before
  writing UI code (training data is not authoritative). A
  diagnosis-first fix protocol with a configurable fix-attempt
  limit prevents spiraling. Workarounds require Architect
  approval. The Architect reviews every commit for incremental
  rot, cross-cutting drift, and cohesion decay.
- **Testing Strategy** — Tests follow a pyramid: unit,
  browserless UI, then E2E (browser-only scenarios). The Unit
  Tester owns all tests by default and delegates to the E2E
  Tester. Per-commit runs are targeted; full suites at the pre-PR
  gate. Human-in-the-loop E2E steps use a structured pause/resume
  cycle.
- **Dependency & Code Hygiene** — The Coder runs lint/format and
  dependency audit at commit time when touching dependencies; the
  Architect handles dead-code judgment during review; the
  Integrator runs on-demand dependency audits when the human
  requests one. CVEs block merging. Patch upgrades are safe,
  minor follow pinning rules, major need approval. External
  tooling (lint, SonarLint, OWASP) does detection; the team
  reacts to findings.
- **Dev-Branch Health** — When the dev branch is broken (by the
  team's own merge or by external changes), the Lead escalates
  and holds off new work until resolved.
- **Context Preservation** — Claude Code may silently compact
  context, dropping loaded files. Every teammate re-reads a
  defined set before starting any task. Teammates in worktrees
  access gitignored files via absolute project root path.

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
├── .sandbox/                        # Sandbox image + developer-local state
│   ├── Dockerfile                   # Custom sandbox image template (tracked)
│   ├── .repo-platform-api.env       # (gitignored) repo-platform API token + metadata (Linux/Windows only; macOS uses the Keychain)
│   ├── .oauth-token                 # (gitignored) captured Claude OAuth token
│   └── .last-directive              # (gitignored) hash of last Lead directive
├── team/                            # Lifecycle scripts (all tracked)
│   ├── README.md                    # One-page reference for these scripts
│   ├── join.sh                      # Provision workstation + launch team
│   ├── leave.sh                     # Discard developer-local state
│   ├── create.sh                    # Build sandbox image + launch fresh sandbox
│   ├── attach.sh                    # Reattach to a running sandbox (daily use)
│   ├── shell.sh                     # Drop into a bash shell inside the sandbox
│   ├── destroy.sh                   # Destroy the sandbox VM
│   └── uninstall.sh                 # Remove the kit from the project
├── CLAUDE.md                        # Project-owned; kit adds a bracketed import line
├── CLAUDE_TEAM.md                   # Kit-owned team context (imported by CLAUDE.md)
├── ONBOARDING.md                    # Developer onboarding (generated)
├── TEAM_GUIDE.md                    # Daily-use reference for humans (generated)
├── .mcp.json                        # Project-scoped MCP server config (canonical location)
├── .claude/
│   ├── settings.json                # Agent team config + permissions (tracked)
│   ├── team-variables.yaml          # Persisted kit variables (tracked)
│   ├── agents/
│   │   ├── lead.md                  # Lead's role + standing instructions (tracked)
│   │   ├── integrator.md            # Integrator role (tracked)
│   │   ├── analyst.md               # Analyst role (tracked)
│   │   ├── architect.md             # Architect role (tracked)
│   │   ├── coder.md                 # Coder role (tracked)
│   │   ├── unit-tester.md           # Unit Tester role (tracked)
│   │   ├── e2e-tester.md            # E2E Tester role (tracked)
│   │   └── tech-writer.md           # Tech Writer role (tracked)
│   ├── commands/
│   │   └── lead-reload.md           # Manual Lead-instruction-reload slash command (tracked)
│   ├── hooks/
│   │   └── session-start-fetch-docs.sh  # SessionStart hook (tracked)
│   ├── .last-onboarded              # (gitignored) marker written by team/join.sh
│   ├── .team-active                 # (gitignored) marker for statusline indicator
│   ├── .progress.md                 # (gitignored) dispatcher task log
│   ├── .tasks/                      # (gitignored) one file per active/suspended task
│   ├── .cost-log.md                 # (gitignored) optional per-task cost log; controlled by COST_IN_LOG
│   └── .worktrees/                  # (gitignored) per-teammate git worktrees
├── docs/
│   └── INDEX.md                     # Sample doc index (project-owned after initial seed)
├── .gitignore                       # Kit maintains a bracketed block of developer-local patterns
└── (existing project files)
```

The kit seeds only `docs/INDEX.md`. The broader `docs/` structure
(`agnostic/`, `reqs/`, non-functional/functional subtrees, etc.) is a
recommended organization for requirements documents, not something the
kit creates. Add and evolve it to fit the project.

### Kit Contents

The kit produces these files in a target project:

| Path | Purpose | Usage |
|------|---------|-------|
| `.sandbox/Dockerfile` | Custom sandbox image for this project | Built automatically by `create.sh` |
| `team/README.md` | One-page reference for the scripts in this directory | Developer reads when onboarding |
| `team/join.sh` | Provisions the developer's workstation (sandbox, repo-platform API token) and launches the team | Human runs at host terminal (first time on this workstation) |
| `team/leave.sh` | Discards developer-local sandbox state; preserves kit files | Human runs at host terminal |
| `team/create.sh` | Builds the sandbox image and launches a fresh sandbox | Human runs at host terminal (or invoked by `join.sh`) |
| `team/attach.sh` | Reattaches to an already-running sandbox | Human runs at host terminal (daily use) |
| `team/shell.sh` | Drops into a bash shell inside the running sandbox | Human runs at host terminal (for inspection or ad-hoc commands) |
| `team/destroy.sh` | Destroys the sandbox VM (and optionally the template image) | Human runs at host terminal |
| `team/uninstall.sh` | Removes the kit from the project (chains through `leave.sh --yes`) | Human runs at host terminal |
| `docs/INDEX.md` | Sample requirement-document index | Seeded once on initial setup; project-owned thereafter (re-setup and remove leave it alone) |
| `CLAUDE_TEAM.md` | Project context for the team (kit-owned) | Imported into `CLAUDE.md` via a bracketed `@CLAUDE_TEAM.md` line |
| `CLAUDE.md` | Project-owned context file | Kit adds/removes only the bracketed import line; everything else is yours |
| `.claude/team-variables.yaml` | Per-project persisted variables | Human-readable, hand-editable, survives kit upgrades |
| `.claude/settings.json` | Agent team config and permissions | Auto-loaded by Claude Code at session start |
| `.mcp.json` | Project-scoped MCP server config | Auto-loaded by Claude Code at session start (canonical location for project-scoped MCP servers; see Claude Code docs) |
| `.claude/agents/lead.md` | Lead's operating manual | Auto-loaded by the sandboxed Claude Code at session start (via `--append-system-prompt` in `create.sh`); also exposed as `/project:lead-reload` for manual re-invocation |
| `ONBOARDING.md` | Developer onboarding (generated) | New developer runs `./team/join.sh` |
| `TEAM_GUIDE.md` | Daily-use reference for humans (generated) | Human reads for workflows, troubleshooting, recovery |
