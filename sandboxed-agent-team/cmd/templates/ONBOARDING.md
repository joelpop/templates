# Developer Onboarding

> **Generated:** `{{UTC_TIMESTAMP}}`  <!-- ISO 8601 UTC, e.g. 2026-04-18T14:32:05Z -->
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the kit template source and
> re-run `agent-team-install`.

## Introduction

Project-specific settings and a setup checklist to prepare your
local development environment for this project's Claude Code agent
team. Version-controlled, so any new developer can invoke it to
recreate the same local environment.

Onboarding creates three things, all developer-local and gitignored:

- **Docker sandbox** — isolated environment where the team runs,
  built from a project-specific Dockerfile.
- **Authentication material** — Claude Code credentials and a
  repo-platform API token, provisioned into the sandbox at each
  startup so teammates can reach Claude and the Git remote over
  HTTPS.
- **Agent team permissions** — `.claude/settings.json` tailored
  for this project.

The end state is a running sandbox with the agent team ready for
work. Each developer generates their own state from this file, so
credentials and host-specific paths never get committed.

**In this section:**
- [Quick Start](#quick-start) — invoke the setup checklist
- [Daily Use](#daily-use) — running the team after setup
- [Overview](#overview) — project details, troubleshooting, offboarding
  - [Project Details](#project-details) — values captured at original setup
  - [Troubleshooting](#troubleshooting) — common sandbox and auth issues
  - [Offboarding](#offboarding) — steps to clean up when you leave

## Quick Start

### Step 1 — Prepare

Have these ready — Claude Code auto-discovers what it can and prompts
you for the rest:

- Docker Desktop installed and running
  (https://www.docker.com/products/docker-desktop/)
- Git identity configured (`git config user.name` and
  `git config user.email`)
- A Claude Code OAuth token or Anthropic API key. On macOS, the
  kit auto-extracts an OAuth token from the Keychain on every
  create/attach; on other systems, export `CLAUDE_CODE_OAUTH_TOKEN`
  or `ANTHROPIC_API_KEY` in your shell config
- A **repo-platform API token** — Bitbucket app password, GitHub
  fine-grained PAT, or GitLab PAT. Required for private repos:
  the sandbox can't use SSH (Docker Sandbox blocks outbound port
  22), so teammates reach the repo over HTTPS and need a token to
  authenticate pushes and PR API calls. Scopes: Bitbucket app
  password — Repositories R+W and Pull requests R+W; GitHub
  fine-grained PAT — Contents R+W and Pull requests R+W; GitLab
  PAT — `api`. The onboarding prompt walks you through creating
  one with a direct link. Leave blank for a public repo if you
  don't need the sandbox to push.

  Token storage:

  - **macOS** — Keychain (service name `agent-team.<sandbox-name>`,
    one entry per project). Survives sandbox rebuilds and
    `./team/destroy.sh`. Wiped only by `./team/leave.sh` or
    `./team/uninstall.sh`.
  - **Linux / Windows** — `.sandbox/.repo-platform-api.env` (mode
    600, gitignored). Onboarding prints a banner explaining the
    tradeoff and pointing to the planned credential-manager
    integration.

### Step 2 — Run onboard

From the project root, run the kit's onboard command:

```
./team/join.sh
```

The tool auto-detects your local state:

- **Fresh onboarding** → builds the sandbox container from the
  project's Dockerfile, prompts for your repo-platform API token
  if needed, and starts the sandbox.
- **Re-onboarding** (you've onboarded this workspace before) →
  discards the existing sandbox and rebuilds. Your workstation's
  repo-platform API token is preserved; the project's versioned
  files are untouched.

Prompts only appear for information the tool can't auto-discover.
When the command finishes, the sandbox is running with a Claude
Code session attached, ready for the agent team.

To discard local workstation state later (without removing the
kit from the project): `./team/leave.sh`.

## Daily Use

Once onboarding is complete:

1. **Start a work session.** At your host terminal (in the project
   directory): `./team/attach.sh`. You'll be prompted to choose
   `resume` (continue the previous Claude Code session) or `fresh`
   (start a new one); flags `--resume` and `--fresh` skip the
   prompt. The sandbox's Claude Code auto-loads the Lead role, so
   the team comes up on your first message — no slash command
   required. The statusline shows "Agent Team Mode" as visible
   confirmation.

   (If `attach.sh` reports no sandbox exists, run `./team/create.sh`
   to rebuild one.)

2. **Describe your work to the Lead.** The Lead coordinates the
   team and drives the workflow; you don't talk to individual
   teammates.

3. **End the session.** Exit Claude Code with `/exit`, `exit`, or
   Ctrl-D *quickly twice* (Claude Code confirms on the first
   press). The Claude Code session ends but **the sandbox VM
   keeps running** — no rebuild needed to come back.

4. **Reconnect next day.** Same as step 1: `./team/attach.sh` →
   `resume` to continue where you left off.

5. **Inspect the sandbox without starting Claude Code.**
   `./team/shell.sh` drops into a bash shell inside the sandbox
   (as the `agent` user, workspace as cwd). Exit with `exit` or
   Ctrl-D. The sandbox keeps running.

6. **Free sandbox resources.** `./team/destroy.sh` destroys the
   sandbox VM. Run when you're done with the project for the day.
   Your workstation-local state (repo-platform API token, task
   files) is preserved — a later `./team/create.sh` rebuilds the
   sandbox from the same state.

For daily workflows — team structure, requirement and
implementation lifecycles, task suspensions, troubleshooting, and
offboarding — see [`TEAM_GUIDE.md`](TEAM_GUIDE.md).

## When the kit is updated

When you pull commits that modify kit files (`team/` scripts,
`CLAUDE_TEAM.md`, `.sandbox/Dockerfile`, etc.), re-run
`./team/join.sh` to rebuild your local sandbox against the
updated templates. Your workstation state (repo-platform API
token) is preserved; only the sandbox image is rebuilt.

## Overview

### Project Details

Captured during the project's agent team setup. Don't modify
unless the project's stack has changed (in which case, ask the
Lead to regenerate this file).

- **Project:** {{PROJECT_NAME}}
- **Stack:** {{STACK_SUMMARY}}
- **Build tool:** {{BUILD_TOOL}}
- **Development branch:** {{DEV_BRANCH_NAME}}
- **Existing code requirements:** {{EXISTING_CODE_REQS}}
- **Feature workflow:** {{FEATURE_WORKFLOW}}
- **Bug workflow:** {{BUG_WORKFLOW}}
- **Auth method at original agent team setup:** <API_KEY | OAUTH>
- **Git remote transport:** <SSH | HTTPS>
- **Merge method:** <PR | INTEGRATOR_MERGE | HUMAN_MERGE>
- **Repo platform:** <BITBUCKET | GITHUB | GITLAB> (if PR method)

### Troubleshooting

- **Docker not installed:** Install Docker Desktop from
  https://www.docker.com/products/docker-desktop/
- **Sandbox authentication fails:** The OAuth token is captured
  fresh on every attach. It can become invalid if the access
  token expired (~24h) and the refresh token also expired (weeks/
  months — rare, after long breaks), or if you ran `/login` on
  the host while the sandbox was running. Either way: run
  `claude` on the host and `/login` if needed, then re-run
  `./team/attach.sh`. On macOS the fresh token is read
  automatically from the Keychain; on other systems update
  `CLAUDE_CODE_OAUTH_TOKEN` in your shell config first.
- **macOS Keychain password prompt:** The kit reads the OAuth
  token from the Keychain on every attach. If the Keychain is
  locked (after a reboot or corporate IT policy), macOS prompts
  for your login password — expected. If IT manages the Keychain
  and you can't unlock it, fall back to exporting
  `CLAUDE_CODE_OAUTH_TOKEN` manually in your shell config.
- **Sandbox git push/pull/fetch fails with 401 or 403:** The
  sandbox reaches the repo over HTTPS using your repo-platform
  API token. Causes:
  - Token expired or was revoked. Re-run `./team/join.sh` to
    re-prompt.
  - Insufficient scopes (Bitbucket: Repositories R+W and Pull
    requests R+W; GitHub fine-grained PAT: Contents R+W and Pull
    requests R+W; GitLab: `api`). Recreate with correct scopes
    and re-run `./team/join.sh`.
  - macOS Keychain returned a stale token (e.g., recreated the
    app password but reused the label). Wipe the entry with
    `security delete-generic-password -s
    agent-team.<sandbox-name>` and re-run `./team/join.sh`.
- **Sandbox build fails:** Verify Docker Desktop is running;
  check the build output for version mismatches or network errors.
- **`/project:lead-reload` not found:** Ensure
  `.claude/agents/lead.md` exists in the repo. Try `git pull`.

### Offboarding

When leaving the project on this workstation, run
`./team/leave.sh`. It destroys the sandbox VM and discards all
developer-local state (repo-platform API token, in-progress task
files, git worktrees). Kit files and committed code are
untouched.

