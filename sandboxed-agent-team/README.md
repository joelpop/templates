# Sandboxed & Isolated Agent Team Kit
*Using a Docker sandbox and Git worktree isolation*

## Introduction

This kit sets up a structured Claude Code agent team inside an
isolated Docker sandbox for a software project, via a single
cross-platform binary. Templates and logic ship together; variables
live in the target project and survive kit upgrades.

Two layers of isolation: a **Docker sandbox** keeps all team
activity, installed tools, and credentials off the developer's host
machine and separated from other projects; **Git worktrees** give
each teammate its own working copy of the repository so teammates
never overwrite each other's in-progress work.

**Key terms you'll see throughout:**

- **Sandbox** — the Docker-based isolated environment where the
  team runs.
- **Lead directive** — the instructions the sandboxed Claude Code
  auto-loads on session start so the Lead role takes over and the
  team comes up automatically on your first message, with no slash
  command required.
- **Repo-platform API token** — a single credential from whichever
  service hosts your repository (your *repo platform*: Bitbucket,
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

You'll need two kinds of things ready: infrastructure on your host
machine, and project-config values Claude Code will ask about.

**Infrastructure (host machine):**

- Docker Desktop installed and running
  (https://www.docker.com/products/docker-desktop/). The sandbox
  feature is required; `create.sh` and `attach.sh` verify it at launch.
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

**Project config (the installer auto-discovers from `pom.xml` and
git config where it can; you'll be prompted with the discovered
value as the default and can press Enter to accept or type an
override):**

- **Team's development branch name** — auto-detected when a
  conventional name (`develop`, `dev`, `feature/develop`) is
  present, or inferred from your remote's default branch.
  Otherwise prompted.
- **Stack details from `pom.xml`**: Java version, Vaadin version,
  Spring Boot version, JUnit version, database, build tool. Each
  prompted with the detected value as default.
- **Merge method** — how completed work reaches the development
  branch (PR, Integrator merge, Human merge).
- **CI platform** — GitHub Actions, GitLab CI, Bitbucket Pipelines,
  Jenkins, or none.
- **Whether to include a cost report in squash-merge commit messages.**
- **Preference to set up your local sandbox** (run `./team/join.sh`)
  after the install completes.
- **A repo-platform API token** (Bitbucket app password, GitHub
  fine-grained PAT, or GitLab PAT) — required unless the repo is
  public. See *"A note on the sandbox's git access"* below for the
  rationale and storage details. Bitbucket: Repositories R+W and
  Pull requests R+W. GitHub fine-grained PAT: Contents R+W and Pull
  requests R+W. GitLab: `api` scope. The setup walks you through
  creating one with direct links; leave blank at the prompt for a
  public repo.

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

- **No kit installed** → full setup: discovers stack details from
  `pom.xml`, prompts for what it can't auto-derive (team's
  development branch, CI platform, merge method, cost-in-commit),
  writes the kit files into the target directory, and offers to
  run `./team/join.sh` to provision your workstation and start the
  team.
- **Kit already installed** → state-aware re-run: shuts down the
  sandbox, preserves your variables file (adding any new
  placeholders the current kit introduces, cleaning orphans),
  regenerates every generated file from current templates, and —
  if your workstation is already provisioned — automatically
  re-provisions it to sync with the refreshed kit.

**No git operations.** The installer writes files; it does not
switch branches, stage, commit, merge, or push. Run `git status`
after an install to review the changes and commit on your own
schedule.

To remove the kit from a project: `./team/uninstall.sh`. This first
runs `./team/leave.sh --yes` to discard your workstation's local
sandbox state, then deletes the kit files and excises the
CLAUDE.md import block and the kit's .gitignore block from the
working tree. Like install, it performs no git operations —
review with `git status` and commit the removal when ready.

Lifecycle commands (`join`, `leave`, `start`, `stop`, `uninstall`)
all live in `./team/` and are the only supported way to invoke
them — running them directly keeps lifecycle commands in lockstep
with the kit version committed to the project.

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
environment, your host Claude Code login doesn't carry over. The
kit autodetects your credentials and injects them into the sandbox
on every create and attach. On macOS, it extracts the OAuth token
from the macOS Keychain automatically — no manual token management
needed. On other systems, a one-time export of
`CLAUDE_CODE_OAUTH_TOKEN` or `ANTHROPIC_API_KEY` is required in
your shell config.

**A note on the sandbox's git access:** Docker Sandbox blocks
outbound port 22 (SSH) for isolation. Regardless of how your host
reaches the repo, teammates *inside* the sandbox use **HTTPS** for
git operations. The kit handles this transparently — your host's
git config and origin URL are never modified.

What the kit does, once per project, on your first `./team/join.sh`:

1. Detects the platform hosting `origin` (Bitbucket / GitHub / GitLab).
   SSH aliases (e.g., `git@bitbucket-syntech:…`) are resolved through
   `~/.ssh/config`'s `HostName` directive.
2. Prompts you to create and paste a repo-platform API token
   (Bitbucket app password, GitHub fine-grained PAT, or GitLab
   PAT). Instructions for the detected platform and a direct link
   to its token-creation page are shown.
3. Stores the token:
   - **macOS** → Keychain, service name `agent-team.<sandbox-name>`,
     one entry per project (multi-project isolation). The token
     never touches the host's regular filesystem.
   - **Linux / Windows** → `.sandbox/.repo-platform-api.env` (mode 600,
     gitignored). A banner during `join.sh` flags this storage
     model as less hardened than Keychain; credential-manager
     integration (libsecret, Credential Manager) is a planned
     follow-up.
4. On every `create.sh` / `attach.sh`, reads the token back and
   pipes it via stdin into the sandbox — the token is never put
   on a command line visible to `ps`. Inside the sandbox, the kit
   writes `~/.git-credentials` and sets
   `git config --global credential.helper store`.
5. If your host's `origin` URL is SSH-shaped (`git@host:…` or
   `ssh://git@host/…`), the sandbox's git also gets a
   `url.https://<host>/.insteadOf git@<alias>:` rewrite rule — so
   SSH URLs in your committed `.git/config` transparently resolve
   to HTTPS when teammates run git inside the sandbox.

**Token persistence:** survives `./team/destroy.sh` (which only
destroys the sandbox VM), survives `agent-team-install` updates,
and survives re-runs of `./team/join.sh`. Wiped by
`./team/leave.sh` (workstation-local state removal) and
`./team/uninstall.sh` (chains leave). On macOS, the Keychain entry
persists even across `leave.sh`; on Linux/Windows the file is
wiped and the next `join.sh` re-prompts.

**Public-repo opt-out:** at the token prompt, press Enter without
pasting anything. Teammates will be able to read public repos but
not push.

## Daily Use

Once agent team setup is complete:

**Note:** Each teammate runs as a separate Claude Code instance
with its own context window. From the Lead's terminal you can
cycle through teammates with `Shift+Down` and message any of them
directly, but you typically don't need to — the Lead coordinates
the team for you.

1. At your host terminal (in the project directory), reattach to
   the sandbox: `./team/attach.sh`. This drops you into a Claude
   Code session running inside the sandbox. The session's system
   prompt auto-loads the Lead role (see [Auto-loading Lead in
   sandbox sessions](#capabilities) under Capabilities), so the
   team comes up automatically on your first message — the Lead
   bootstraps and brings up the teammates for you. No slash
   command required. Once setup completes, the statusline shows
   "Agent Team Mode" as a visible confirmation that you're talking
   to the team. (If `attach.sh` tells you no sandbox exists — for
   example, because you ran `destroy.sh` last time — run
   `./team/create.sh` to build a fresh one.)
2. The sandboxed Claude Code pre-authorizes common shell commands
   (`mvn`, `git`, `ls`, `chmod`, etc.) via `.claude/settings.json`'s
   allow/deny rules, so teammates can run builds, tests, and
   routine git operations without being prompted. Destructive or
   out-of-scope operations
   (`git reset --hard`, `git push --force`, arbitrary shell through
   `curl | bash`, etc.) are explicitly denied by the same
   settings. The Lead will not implement directly — this is
   enforced by its instructions in `.claude/agents/lead.md`.
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
7. You can ask the team to take screenshots of the running
   application for visual verification — tell the Lead what you
   want to see.
8. **If something goes wrong:**
   - Teammate seems stuck or unresponsive: tell the Lead. The
     Lead will recover the teammate (resume first; replace from
     the same agent definition if resume fails — see Teammate
     Recovery in `.claude/agents/lead.md`).
   - The Lead itself loses context mid-session: run
     `/project:lead-reload` at the sandbox's Claude Code prompt to
     re-invoke the Lead (the auto-load fires only at session start,
     so mid-session recovery uses the slash command).
   - Sandbox crashes: back at your host terminal, run
     `./team/attach.sh` to reconnect (which reopens Claude Code
     inside the sandbox). The new session auto-loads the Lead,
     which reads `progress.md` to recover state. If the sandbox
     itself is gone, run `./team/create.sh` to rebuild it.
   - The Lead may suspend a task to work on a prerequisite it
     discovered — this is normal. It will resume the original task
     after the prerequisite is complete.
9. **Pausing and resuming:** Exiting Claude Code (`/exit`, `exit`, or Ctrl-D quickly twice)
    ends your Claude Code session and drops you back to the host shell,
    but the sandbox VM keeps running in the background.
    To resume: at your host terminal run `./team/attach.sh`
    — it reattaches to the running sandbox and starts a new Claude
    Code session inside it. The Lead auto-loads and reads
    `progress.md` to pick up where you left off.
10. To end a Claude Code session cleanly, tell the Lead you're
    wrapping up the session. The Lead confirms all work is merged and flags
    anything unresolved for your next Claude Code session. Then
    exit Claude Code (`/exit`, `exit`, or Ctrl-D quickly twice) — the sandbox VM keeps
    running so you can reconnect later.
11. To end the engagement (i.e., destroy the sandbox), after ending your final Claude Code
    session, at your host terminal: `./team/destroy.sh`

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
- **Coder** — implements features and fixes bugs; runs
  lint/format/analysis on touched files at commit; runs dependency
  audit when adding or removing a dependency.
- **Unit Tester** — unit and browserless UI tests.
- **E2E Tester** — end-to-end browser tests.
- **Tech Writer** — owns `docs/guides/` (install / deploy / user /
  admin / operator guides). Updates on the release cadence.

### Features

#### Capabilities

- **Isolation & Infrastructure** — Each project gets its own Docker
  sandbox built from a customizable Dockerfile. One-command scripts
  handle startup and disposal. Claude Code authentication is
  autodetected and injected via environment variable. The sandbox
  reaches the repo over HTTPS (Docker Sandbox blocks outbound port
  22); a repo-platform API token is provisioned automatically for
  private-repo access — see *A note on the sandbox's git access*
  above.
- **Auto-loading Lead in sandbox sessions** — The sandbox's Claude
  Code session starts with the Lead role pre-configured:
  `team/create.sh` passes `--append-system-prompt` to `claude` so
  the first turn reads `.claude/agents/lead.md` and brings up the team
  automatically (via `TeamCreate`). The human does not need to
  remember `/project:lead-reload`. Host Claude Code sessions are
  unaffected (they don't go through `create.sh`/`attach.sh`). The
  `/project:lead-reload` slash command remains available as a
  manual re-invocation fallback if the Lead needs to be reset
  mid-session.
- **"Agent Team Mode" statusline indicator** — The sandbox's
  statusline displays "Agent Team Mode" once the Lead has
  completed the Pre-Start Check and brought up the team, giving a
  visible cue at the keyboard that the human is interacting with
  the team (not bare Claude Code). Implemented via a `statusLine`
  entry in `.claude/settings.json` that checks for a sentinel file
  (`.claude/.team-active`) written by the Lead after `TeamCreate`
  succeeds. The indicator is blank before setup completes and
  between sessions; it updates each session based on the current
  state.
- **Status Tracking** — Requirement status checkboxes (`[ ]`/`[-]`/`[x]`)
  in `docs/` plus role-assigned plan steps in task files. A progress
  dispatcher tracks active and suspended tasks for recovery after
  context compaction.
- **Sub-Task Parallelism** — Within a single task, the Lead may split
  file-disjoint implementation work across multiple Coders, each with
  a paired Unit Tester. Phases support dependencies between subtasks.
  Roles also work in parallel where possible: the Unit Tester and
  Architect review simultaneously after Coder work is merged.
- **MCP Documentation Servers** — Teammates consult MCP servers for
  authoritative framework documentation (Java, Vaadin, Spring,
  Playwright) rather than relying on training data. The Playwright
  MCP server also provides visual debugging — teammates can navigate
  the running application, take screenshots, and interact with the
  UI.
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
  `.claude/agents/lead.md`) is version-controlled; developer-local state (tasks,
  progress, worktrees, settings) is gitignored. An `ONBOARDING.md` is
  generated for new developers.

#### Workflows

- **Coordination Model** — Teammates message each other directly for
  routine coordination. They escalate to the Lead when a decision
  requires human input or intervention. The human only talks to the
  Lead, who coordinates but does not implement directly.
- **Requirements Management** — All requirements originate from the
  human and are documented by the Analyst. New capabilities go through
  a requirement gate (Analyst drafts, human approves). Teammates must
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
  teammate roles each get dedicated branches. Teammates merge (never rebase).
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
- **Dependency & Code Hygiene** — The Coder runs lint/format and
  the dependency audit at commit time when touching dependencies;
  the Architect handles dead-code judgment during code review; the
  Integrator runs on-demand dependency audits when the human
  requests one. CVEs block merging. Version upgrades follow
  pinning rules — patch upgrades are safe, minor upgrades follow
  pinning rules, major upgrades need approval. External tooling
  (lint, SonarLint, OWASP) does the detection work; the team
  reacts to the findings.
- **Dev-Branch Health** — When the development branch is broken (by
  the team's own merge or by external changes), the Lead escalates to
  the human and holds off on new work until the issue is resolved.
- **Context Preservation** — Claude Code may silently compact context,
  dropping loaded files. Every teammate must re-read a defined set of
  files before starting any task. Teammates in worktrees access
  gitignored files via absolute project root path.

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
