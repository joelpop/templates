# Developer Onboarding

> **Generated:** `<UTC_TIMESTAMP>`  <!-- ISO 8601 UTC, e.g. 2026-04-18T14:32:05Z -->
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the kit template source and
> re-run `agent-team-install`.

## Introduction

This document contains the project-specific settings and a setup
checklist to prepare your local development environment for working
on this project's Claude Code agent team. It is version-controlled in
the repo, so any developer joining the project can invoke it to
recreate the same local environment.

Onboarding creates three things: a **Docker sandbox** (an isolated
environment where the agents run, built from a project-specific
Dockerfile), **authentication and SSH material** (provisioned into
the sandbox at each startup so agents can reach Claude and the Git
remote), and **agent team permissions** (`.claude/settings.json`
tailored for this project). These artifacts are developer-local and
gitignored — each developer generates their own from this file — so
credentials, SSH keys, and host-specific paths never get committed.
The end state is a running sandbox with the agent team ready for
work.

This document has two parts. Everything above the divider is
human-facing front matter — read this to understand what onboarding
provides and how to invoke it. Everything below the divider is a
setup checklist executed by the agent when you point it at this file.

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
- If the project uses an SSH Git remote: the SSH key referenced in
  `~/.ssh/config` for that remote
- A **repo-platform API token** — Bitbucket app password, GitHub
  fine-grained PAT, or GitLab PAT. Required for any private repo
  so the sandbox can reach it over HTTPS. The sandbox can't use
  SSH (Docker Sandbox blocks outbound port 22), so agents reach
  the repo over HTTPS and need a token to authenticate pushes and
  PR API calls. Bitbucket: app password with Repositories R+W and
  Pull requests R+W. GitHub: fine-grained PAT with Contents R+W
  and Pull requests R+W. GitLab: PAT with `api` scope. The
  onboarding prompt walks you through creating one with a direct
  link. If this is a public repo and you don't need the sandbox
  to push, you can leave the prompt blank.

  Token storage:

  - **macOS** — stored in the Keychain (service name
    `agent-team.<sandbox-name>`, one entry per project). Survives
    sandbox rebuilds and `./team/destroy.sh`. Wiped only by
    `./team/leave.sh` or `./team/uninstall.sh`.
  - **Linux / Windows** — stored in `.sandbox/.platform-api.env`
    (mode 600, gitignored). Onboarding will print a banner
    explaining this tradeoff and a pointer to the future
    credential-manager integration.

### Step 2 — Run onboard

From the project root, run the kit's onboard command:

```
./team/join.sh
```

The tool auto-detects your local state and does the right thing:

- **Fresh onboarding** → builds the sandbox container from the
  project's Dockerfile, provisions SSH material if the project uses
  an SSH Git remote, and starts the sandbox.
- **Re-onboarding** (you've onboarded this workspace before) →
  discards the existing sandbox and rebuilds from scratch. Your
  workstation's SSH material and repo-platform API token are preserved;
  the project's versioned files are untouched.

If prompts come up, they're for information the tool can't auto-discover
(e.g., which SSH key to use when multiple keys are configured). When
the command finishes, the sandbox is running and a Claude Code session
is attached to it, ready for the agent team.

To discard your local workstation state later (without
removing the kit from the project): `./team/leave.sh`.

## Daily Use

Once onboarding is complete:

1. **Start a work session.** At your host terminal (in the project
   directory): `./team/attach.sh`. You'll be prompted to choose
   `resume` (continue the previous Claude Code session, preserving
   context) or `fresh` (start a new session). Flags `--resume` and
   `--fresh` skip the prompt. You're dropped into Claude Code
   running inside the sandbox; the system prompt auto-loads the
   Lead role, so the team spawns on your first message — no slash
   command required. The statusline shows "Agent Team Mode" as
   visible confirmation.

   (If `attach.sh` reports no sandbox exists, run `./team/create.sh`
   to rebuild one.)

2. **Describe your work to the Lead.** The Lead coordinates the team
   and drives the workflow; you don't talk to individual teammates.

3. **End the session.** When you're wrapping up, exit Claude Code
   with `/exit`, `exit`, or Ctrl-D *quickly twice* (Claude Code
   confirms on the first press). That ends your Claude Code
   session but **the sandbox VM keeps running** — no rebuild
   needed to come back.

4. **Reconnect next day.** Same as step 1: `./team/attach.sh` →
   choose `resume` to continue where you left off.

5. **Inspect the sandbox without starting Claude Code.** Run
   `./team/shell.sh` to drop into a bash shell inside the sandbox
   (as the `agent` user, workspace as cwd). Exit with `exit` or
   Ctrl-D (single — regular shell, not Claude Code). The sandbox
   keeps running.

6. **Free sandbox resources.** `./team/destroy.sh` destroys the
   sandbox VM. Run when you're done with the project for the day
   (or longer). Your workstation-local state (SSH material, API
   token, task files) is preserved — a later `./team/create.sh`
   rebuilds the sandbox from the same state.

For detailed daily workflows — team structure, requirements and
implementation lifecycles, task suspensions, troubleshooting, and
offboarding — see [`TEAM_GUIDE.md`](TEAM_GUIDE.md).

## When the kit is updated

When you pull commits that have modified kit files (generated scripts
under `team/`, `CLAUDE_TEAM.md`, `.sandbox/Dockerfile`, etc.),
re-run `./team/join.sh` to rebuild your local sandbox against the
updated templates. Your workstation state (SSH keys, API tokens)
is preserved; only the sandbox image is rebuilt.

## Overview

### Project Details

These values were captured during the project's agent team setup.
Do not modify them unless the project's stack has changed (in which
case, ask the Lead to regenerate this file).

- **Project:** <PROJECT_NAME>
- **Stack:** <STACK_SUMMARY>
- **Build tool:** <BUILD_TOOL>
- **Development branch:** <DEV_BRANCH_NAME>
- **Auth method at original agent team setup:** <API_KEY | OAUTH>
- **Git remote transport:** <SSH | HTTPS>
- **Merge method:** <PR | INTEGRATOR_MERGE | HUMAN_MERGE>
- **Repo platform:** <BITBUCKET | GITHUB | GITLAB> (if PR method)

### Troubleshooting

- **Docker not installed:** Install Docker Desktop from
  https://www.docker.com/products/docker-desktop/
- **Sandbox authentication fails:** The sandbox's OAuth token is
  captured fresh on every attach. It can become invalid if:
  - The access token expired (~24h) and the refresh token also expired
    (weeks/months) — rare, but happens after long breaks.
  - You ran `/login` on the host while the sandbox was running — the
    new login may invalidate the token the sandbox is using.
  In either case: refresh the token by running `claude` on the host
  and `/login` if needed, then re-run `./team/attach.sh`. On macOS
  the fresh token is read automatically from the Keychain; on other
  systems update `CLAUDE_CODE_OAUTH_TOKEN` in your shell config
  first.
- **macOS Keychain password prompt:** On macOS the kit reads the
  OAuth token from the Keychain on every attach. If the Keychain is
  locked (e.g., after a reboot or corporate IT policy), macOS may
  prompt for your login password. This is expected — enter it to
  continue. If the Keychain is managed by IT and you cannot unlock
  it, fall back to exporting `CLAUDE_CODE_OAUTH_TOKEN` manually in
  your shell config.
- **SSH remote access fails inside the sandbox:** Git operations over
  SSH (push, pull, fetch) fail with "Permission denied" or "Host key
  verification failed."
  - Verify `.sandbox/.ssh/` exists and contains the private key, a
    `config` file, and a `known_hosts` file.
  - If the key was rotated since onboarding, either re-run onboarding
    or manually copy the new key to `.sandbox/.ssh/` and restart the
    sandbox.
  - If `known_hosts` is missing or stale, regenerate it:
    `ssh-keyscan <hostname> > .sandbox/.ssh/known_hosts 2>/dev/null`
  - Verify `.sandbox/.ssh.source` contains the correct absolute path
    to the host key — the kit reads this to sync keys on every
    create and attach.
- **Sandbox build fails:** Check that Docker Desktop is running.
  Review the build output for version mismatches or network errors.
- **`/project:team-start` not found:** Ensure
  `.claude/commands/team-start.md` exists in the repo. Try `git pull`
  to get the latest.

### Offboarding

When you leave the project on this workstation, run
`./team/leave.sh`. It destroys the sandbox VM and discards all
developer-local state (SSH material, repo-platform API token,
in-progress task files, git worktrees). Kit files and the project's
committed code are untouched.

