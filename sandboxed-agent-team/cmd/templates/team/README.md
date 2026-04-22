# team/

Scripts that manage this project's local sandbox and the agent team
running inside it.

## Scripts

| Script         | Purpose                                                    |
|----------------|------------------------------------------------------------|
| `join.sh`      | Provision your workstation and launch the team.            |
| `leave.sh`     | Discard your workstation's local sandbox state.            |
| `create.sh`    | Build the sandbox image and launch a fresh sandbox.        |
| `attach.sh`    | Reattach to an already-running sandbox.                    |
| `destroy.sh`   | Destroy the sandbox VM.                                    |
| `uninstall.sh` | Remove the kit from the project.                           |

## Common flows

**First time on this workstation:**

```
./team/join.sh
```

Provisions SSH material, a platform API token (if the PR merge
method is configured), and runs `create.sh` to build the sandbox
image and launch the team.

**Daily: reconnect to your sandbox:**

```
./team/attach.sh
```

Reattaches to the running sandbox and resumes your Claude Code
session. Fails fast if the sandbox was destroyed or if the Lead
directive has changed — in which case run `./team/destroy.sh`
followed by `./team/create.sh` to rebuild.

**Ending a session:**

Exit Claude Code with `/exit` or Ctrl+D. The sandbox VM keeps
running in the background. Reconnect with `./team/attach.sh`.
When you're done with the project for the day (or longer), run
`./team/destroy.sh` to free the sandbox resources.

**Leaving the project on this workstation:**

```
./team/leave.sh
```

Discards your developer-local state (SSH material, API token,
in-progress task files, worktrees) and destroys the sandbox.
Does not remove the kit's versioned files.

**Removing the kit from the project:**

```
./team/uninstall.sh
```

Runs `leave.sh` first to discard local state, then deletes the
kit's versioned files from the working tree and excises the
CLAUDE.md import block and the kit's `.gitignore` block. No git
operations — review with `git status` and commit the removal
when ready.

## See also

- `../TEAM_GUIDE.md` — daily workflows, troubleshooting, recovery.
- `../ONBOARDING.md` — first-time setup and project variables.
- `../CLAUDE_TEAM.md` — agent team configuration.
