# Developer Onboarding

> **Generated:** `<UTC_TIMESTAMP>`  <!-- ISO 8601 UTC, e.g. 2026-04-18T14:32:05Z -->
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the kit template source and
> re-run `agent-team install`.

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
- A Claude Code OAuth token or Anthropic API key. On macOS,
  `start.sh` auto-extracts an OAuth token from the Keychain; on other
  systems, export `CLAUDE_CODE_OAUTH_TOKEN` or `ANTHROPIC_API_KEY` in
  your shell config
- If the project uses an SSH Git remote: the SSH key referenced in
  `~/.ssh/config` for that remote
- If the project uses the **PR** merge method (see [Project
  Details](#project-details) below): a platform API token (Bitbucket
  app password, GitHub fine-grained PAT, or GitLab personal access
  token). Claude Code walks you through creating one if you don't
  have it ready.

### Step 2 — Run onboard

From the project root, run the kit's onboard command:

```
./team/join.sh
```

The tool auto-detects your local state and does the right thing:

- **Fresh onboarding** → builds the sandbox container from the
  project's Dockerfile, provisions SSH material if the project uses
  an SSH Git remote, and starts the sandbox.
- **Re-onboarding** (you've onboarded this workspace before) → tears
  down the existing sandbox and rebuilds from scratch. Idempotent
  and non-destructive: nothing versioned is touched.

If prompts come up, they're for information the tool can't auto-discover
(e.g., which SSH key to use when multiple keys are configured). When
the command finishes, the sandbox is running and a Claude Code session
is attached to it, ready for the agent team.

To uninstall your local state later: `./team/join.sh
--remove`.

## Daily Use

Once onboarding is complete:

1. At your host terminal (in the project directory), start the
   sandbox: `team/start.sh`. This drops you into a Claude Code
   session running inside the sandbox. The session's system prompt
   auto-loads the Lead role, so the team spawns as soon as you send
   your first message — no slash command required. Once setup
   completes, the statusline shows "Agent Team Mode" as a visible
   confirmation that you're talking to the team.
2. Describe your work to the Lead.

For detailed daily workflows — team structure, requirements and
implementation lifecycles, pausing and resuming, ending a Claude
Code session, and what to do when something goes wrong — see
[`TEAM_GUIDE.md`](TEAM_GUIDE.md).

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
  captured at startup. It can become invalid if:
  - The access token expired (~24h) and the refresh token also expired
    (weeks/months) — rare, but happens after long breaks.
  - You ran `/login` on the host while the sandbox was running — the
    new login may invalidate the token the sandbox is using.
  In either case: stop the sandbox with `docker sandbox stop <name>`
  (the name is printed by `start.sh` at startup), re-run `claude` on
  the host and `/login` if needed, then restart with
  `team/start.sh`. On macOS, `start.sh` automatically picks up the
  fresh token from the Keychain. On other systems, update
  `CLAUDE_CODE_OAUTH_TOKEN` in your shell config first.
- **macOS Keychain password prompt:** On macOS, `start.sh` reads the
  OAuth token from the Keychain. If the Keychain is locked (e.g.,
  after a reboot or corporate IT policy), macOS may prompt for your
  login password. This is expected — enter it to continue. If the
  Keychain is managed by IT and you cannot unlock it, fall back to
  exporting `CLAUDE_CODE_OAUTH_TOKEN` manually in your shell config.
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
    to the host key — `start.sh` reads this to sync keys on startup.
- **Sandbox build fails:** Check that Docker Desktop is running.
  Review the build output for version mismatches or network errors.
- **`/project:team-start` not found:** Ensure
  `.claude/commands/team-start.md` exists in the repo. Try `git pull`
  to get the latest.

### Offboarding

When you leave the project:
1. Run `team/stop.sh` to destroy your sandbox.
2. Delete `.sandbox/` and `.claude/settings.json` — these are
   gitignored and not shared.
3. Optionally delete `.claude/.tasks/` and `.claude/.progress.md` if
   no one else needs your local task history.

