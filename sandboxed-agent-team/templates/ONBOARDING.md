# Developer Onboarding

> **Generated:** `<UTC_TIMESTAMP>`  <!-- ISO 8601 UTC, e.g. 2026-04-18T14:32:05Z -->
> **GENERATED FILE** — do not edit directly. Edits here will be lost
> the next time this file is regenerated. To change this file, edit
> its template in the team setup kit (SANDBOXED_AGENT_TEAMS.md) and
> re-run the setup at your host terminal.

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

### Step 2 — Prompt

Start a Claude Code session in the project directory in **accept
edits mode** (press Tab until the mode selector shows "Accept Edits"
or start Claude Code with `--allowedTools Edit,Write,Read,Glob,Grep`).
This auto-approves file reads and writes — the setup creates several
files — while still prompting you for shell commands. Then say:

> Read `ONBOARDING.md` and execute the setup checklist. Ask me before
> doing anything destructive, and stop when you need my input.

Throughout setup, Claude Code will prompt you to approve shell
commands and other tool calls not covered by accept edits mode.
These are expected — approve them to keep the setup moving. The
prompt you gave still applies, so Claude Code will pause for your
input at decision points.

### Step 3 — Proceed

Claude Code takes it from here. The checklist detects your local
state and adjusts automatically:

- **Fresh onboarding** — full setup (sandbox files, agent settings,
  sandbox build, team start)
- **Re-onboarding** — asks before overwriting existing local setup

Claude Code handles most steps autonomously. It pauses to ask for
your input when it needs information it cannot discover (e.g., your
auth method if detection fails) or confirmation (e.g., which SSH key
to use). When the sandbox files are ready, it asks you to run
`.sandbox/start.sh` in a separate terminal, then continues setup
from inside the sandbox session.

## Daily Use

Once onboarding is complete:

1. At your host terminal (in the project directory), start the
   sandbox: `.sandbox/start.sh`. This drops you into a Claude Code
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
  `.sandbox/start.sh`. On macOS, `start.sh` automatically picks up the
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
  - Verify `.sandbox/ssh/` exists and contains the private key, a
    `config` file, and a `known_hosts` file.
  - If the key was rotated since onboarding, either re-run onboarding
    or manually copy the new key to `.sandbox/ssh/` and restart the
    sandbox.
  - If `known_hosts` is missing or stale, regenerate it:
    `ssh-keyscan <hostname> > .sandbox/ssh/known_hosts 2>/dev/null`
  - Verify `.sandbox/ssh.source` contains the correct absolute path
    to the host key — `start.sh` reads this to sync keys on startup.
- **Sandbox build fails:** Check that Docker Desktop is running.
  Review the build output for version mismatches or network errors.
- **`/project:team-start` not found:** Ensure
  `.claude/commands/team-start.md` exists in the repo. Try `git pull`
  to get the latest.

### Offboarding

When you leave the project:
1. Run `.sandbox/stop.sh` to destroy your sandbox.
2. Delete `.sandbox/` and `.claude/settings.json` — these are
   gitignored and not shared.
3. Optionally delete `.claude/tasks/` and `.claude/progress.md` if
   no one else needs your local task history.

---

**SETUP CHECKLIST — Agent-Executed Content Below**

---

## Setup Checklist

### Step 0 — Detect local state
Check if this developer already has local setup files:
- Does `.sandbox/Dockerfile` exist?
- Does `.claude/settings.json` exist?

If both exist, tell the human: "Your local setup already exists. Would
you like to rebuild it? This will overwrite your sandbox config and
agent settings." If the human says no, skip to Step 4.

### Step 1 — Detect developer identity and authentication method
<!-- SYNC NOTE: The auth detection logic below is duplicated in
     Setup Step 2 of the Agent Team Setup checklist. If you update
     one, update both. -->
Auto-discover the developer's git identity from `git config user.name`
and `git config user.email`. These are needed for the Dockerfile in
Step 2. If either value is missing, ask the developer to configure
them before proceeding:
```
git config user.name "Your Name"
git config user.email "your.email@example.com"
```

The sandbox also needs an API key or OAuth token to authenticate —
your host Claude Code login does not carry over. On macOS, `start.sh`
automatically extracts the latest OAuth token from the Keychain at
each startup, so no manual token management is needed. On other
systems, a manual export is required.

Detect the auth method so `start.sh` can inject the correct
credentials into the sandbox at startup:

1. Check if `ANTHROPIC_API_KEY` is set in the current environment.
   - If set → API key method. Proceed to Step 2.
2. On macOS: check if the macOS Keychain has Claude Code credentials.
   Run:
   ```
   security dump-keychain 2>/dev/null | grep -i 'claude.*credential'
   ```
   - If a match is found → OAuth method. Tell the human: "Your OAuth
     token is in the macOS Keychain. `start.sh` will extract it
     automatically at each startup — no manual export needed." Proceed
     to Step 2.
   - If no match → go to step 3.
3. Check if `CLAUDE_CODE_OAUTH_TOKEN` is already exported in the
   current environment.
   - If set → OAuth method. Verify it is also in the human's shell
     config (`~/.zshrc` or equivalent) so it persists across sessions.
     Proceed to Step 2.
4. If neither environment variable is set, ask the human for permission
   to read `~/.claude/.credentials.json`.
   - If it contains an OAuth token → tell the human: "You're using
     OAuth but `CLAUDE_CODE_OAUTH_TOKEN` is not exported in your shell
     config. The sandbox needs this to authenticate." Walk them through:
     1. Copy the token value from `~/.claude/.credentials.json`.
     2. Add `export CLAUDE_CODE_OAUTH_TOKEN="<token>"` to `~/.zshrc`
        or equivalent.
     3. Confirm when done before proceeding.
   - If the file doesn't exist or contains no token → go to step 5.
5. If no method is detected, ask: "I couldn't detect your auth method.
   Do you authenticate via an Anthropic API key or via OAuth
   (company/team account)?" Then guide setup accordingly.

**SSH remote access:**
<!-- SYNC NOTE: The SSH detection logic below is duplicated in
     Setup Step 1 of the Agent Team Setup checklist. If you update
     one, update both. -->
Check if the project's Git remote uses SSH:
```
git remote -v
```
If any remote URL uses an SSH-style address (starts with `git@` or
`ssh://`), the sandbox needs the developer's SSH key to access the
remote.

If SSH is detected:
1. Extract the host or alias from the remote URL. For example,
   `git@bb-client:org/repo.git` → host alias is `bb-client`;
   `git@bitbucket.org:org/repo.git` → host is `bitbucket.org`.
2. Look up the host/alias in `~/.ssh/config` to find the
   `IdentityFile` and `HostName` (if different from the alias).
   If there is no matching entry in `~/.ssh/config` and the host
   in the URL is a real hostname (e.g., `bitbucket.org`), the
   default key (`~/.ssh/id_ed25519` or `~/.ssh/id_rsa`) may be
   used — ask the developer to confirm which key authenticates
   with this remote.
3. Confirm the key path with the developer.
4. Note the SSH key path, host alias, and real hostname for Step 2.

If SSH is not detected (remote uses HTTPS or no remote is configured),
skip SSH setup.

### Step 2 — Create `.sandbox/` files

Create the following three files. Replace `<GIT_USER_NAME>` and
`<GIT_USER_EMAIL>` in the Dockerfile with this developer's git identity
from Step 1. If SSH is not needed, remove `openssh-client` from the
apt-get line and the SSH directory block from the Dockerfile.
Authentication is handled by `start.sh`, not the Dockerfile.

**`.sandbox/Dockerfile`:**
```dockerfile
<FULLY_CUSTOMIZED_DOCKERFILE_CONTENT>
```

**`.sandbox/start.sh`:**
```bash
<VERBATIM_START_SH_CONTENT>
```

**`.sandbox/stop.sh`:**
```bash
<VERBATIM_TEARDOWN_SH_CONTENT>
```

Run: `chmod +x .sandbox/start.sh .sandbox/stop.sh`

If SSH was detected in Step 1, create `.sandbox/ssh/` with the
developer's SSH material:
1. Copy the SSH key pair:
   ```
   mkdir -p .sandbox/ssh
   cp <IdentityFile> .sandbox/ssh/
   cp <IdentityFile>.pub .sandbox/ssh/ 2>/dev/null
   ```
2. Write an SSH config for the sandbox — the `IdentityFile` must
   reference the sandbox path (`~/.ssh/<key-filename>`), not the
   host path:
   ```
   cat > .sandbox/ssh/config << 'EOF'
   Host <alias-or-hostname>
       HostName <real-hostname>
       User git
       IdentityFile ~/.ssh/<key-filename>
       IdentitiesOnly yes
   EOF
   ```
   If the remote URL uses the real hostname directly (no alias in
   `~/.ssh/config`), write the config entry with
   `Host <real-hostname>` instead.
3. Generate `known_hosts` for the remote host:
   ```
   ssh-keyscan <real-hostname> > .sandbox/ssh/known_hosts 2>/dev/null
   ```
4. Record the host key path so `start.sh` can sync fresh keys:
   ```
   echo "<IdentityFile-absolute-path>" > .sandbox/ssh.source
   ```

If the merge method is **PR** (see Project Details above), provision
platform API access. This is required for the Lead to create, read,
and merge PRs via the repo server's REST API.

Guide the developer through creating an API token:
- **Bitbucket Cloud:** "Go to
  https://bitbucket.org/account/settings/app-passwords/ and create an
  app password with **Repositories:Read** and **Pull requests:Read +
  Write** permissions. Paste the generated password here."
- **GitHub:** "Go to Settings → Developer settings → Personal access
  tokens → Fine-grained tokens. Create a token scoped to this
  repository with **Pull requests:Read and write** and
  **Contents:Read** permissions. Paste the token here."
- **GitLab:** "Go to User Settings → Access Tokens. Create a token
  with the **api** scope. Paste the token here."

Write `.sandbox/platform-api.env` with the token and repo details
(detected from the remote URL):
```
PLATFORM_TYPE=<PLATFORM_TYPE>
PLATFORM_API_URL=<API_BASE_URL>
PLATFORM_API_USER=<USERNAME>
PLATFORM_API_TOKEN=<TOKEN>
PLATFORM_REPO_WORKSPACE=<WORKSPACE>
PLATFORM_REPO_SLUG=<REPO_SLUG>
PLATFORM_REPO_OWNER=<OWNER>
PLATFORM_REPO_NAME=<REPO_NAME>
```
Include only the fields relevant to the platform.

### Step 3 — Create `.claude/settings.json`

Create `.claude/settings.json` with the following content:

```json
<EXACT_SETTINGS_JSON_CONTENT>
```

Do NOT overwrite or modify `.claude/commands/team-start.md` — it is
versioned and shared across all developers.

### Step 4 — Build and start the sandbox

Tell the human:

> "The sandbox files are ready. Please open a new terminal in the
> project directory and run: `.sandbox/start.sh`
>
> This builds the Docker image and starts a new Claude Code session
> inside the sandbox. Let me know when the sandbox session is ready."

Wait for the human to confirm the sandbox session is running, then
tell them:

> "In the sandbox session, say: *Read `ONBOARDING.md` and continue
> from Step 5.*"

Poll for `.claude/.last-onboarded` to appear (visible via the
bidirectional sandbox mount), confirming that Onboarding Step 5 has
completed. If the file does not appear within a few minutes, ask the
human to check the sandbox terminal for errors. Common causes:
- **Sandbox build failed or hung:** Check the `start.sh` output for Docker
  build errors (network issues, missing packages, insufficient disk).
- **Authentication failed:** The OAuth token may be expired or missing.
  See the Troubleshooting section.
- **Human closed the sandbox terminal:** Re-run `.sandbox/start.sh` to
  reconnect.

### Step 5 — Start the agent team

*(This step is executed by the sandbox Claude Code session.)*

Record the current UTC timestamp as this developer's onboarding
time to `.claude/.last-onboarded`. The file contains one line
with a `Last onboarded:` label followed by the timestamp in ISO
8601 UTC format:

```bash
TS=$(date -u +%Y-%m-%dT%H:%M:%SZ)
echo "Last onboarded: $TS" > .claude/.last-onboarded
```

The session's system prompt already auto-loaded the Lead role
from `.sandbox/start.sh`'s `--append-system-prompt`, so the team
spawns automatically on the first human message. Then tell the
human: "The team is ready. Describe what you'd like to work on."

### Step 6 — End the host session

*(Back in the host session.)* Once the human confirms the sandbox
team is running, tell them: "Onboarding is complete — the sandbox
session is running your team. This host session's setup work is done.
You can close it, or keep it open for any work you want to do outside
the sandbox (note: work in the host session is not sandboxed)."
