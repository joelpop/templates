# Claude Code: Verified Facts vs. Common False Beliefs

Claude's training data describes earlier versions of Claude Code and is
known to confidently misstate the current tool surface, frontmatter
schema, and Agent Teams mechanics. Past sessions on this kit have spent
hours "fixing" things that weren't broken — removing real tool calls,
stripping valid frontmatter fields, redesigning around features that
were assumed not to exist.

This file records facts that have been verified against the installed
Claude Code binary, alongside the false beliefs they correct. Each
entry includes verification commands so you don't have to trust this
file either — re-verify when in doubt.

**Currency:** facts here are tied to a specific Claude Code version.
Re-verify before treating an entry as load-bearing for a new design.
Last verification: Claude Code **2.1.128** (2026-05).

**When to read this file:**

- Before "fixing" anything in `.claude/agents/*.md`, `lead.md`,
  `CLAUDE_TEAM.md`, or kit hooks/scripts that touches Agent Teams,
  subagent frontmatter, or Claude Code tool calls.
- When a code review flags a `TeamCreate` call, an
  `isolation: worktree` frontmatter line, or a teammate calling
  `Agent` as "broken."
- Whenever your training-data recall about Claude Code conflicts with
  what this kit actually does. Default to the kit; verify if uncertain.

---

## Agent Teams is a real Claude Code feature (`TeamCreate`, `TeamDelete`, `SendMessage`)

**Common false belief:** Claude Code only has subagents (the `Agent`
tool); "Agent Teams" doesn't exist as a separate feature; there is no
`TeamCreate` tool. The public docs page describes spawning a team "in
plain language," which is easy to misread as meaning no underlying
tool exists.

**Verified:** Agent Teams is real and experimental. The following are
deferred tools (they appear in `ToolSearch` but not in the default
tool list):

- `TeamCreate` — *"create a multi-agent swarm team"*. Inputs:
  `team_name` (required), `description` (optional), `agent_type`
  (optional). Errors if you're already leading a team:
  *"A leader can only manage one team at a time."*
- `TeamDelete` — *"disband a swarm team and clean up"*.
- `SendMessage` — *"send messages to agent teammates (swarm
  protocol)"*.

The `Agent` tool, when called by a team lead, accepts additional
spawned-teammate fields: `name`, `team_name`, `mode`,
`isolation: "worktree"`, `cwd`. Hook events `TeammateIdle`,
`TaskCreated`, `TaskCompleted` are emitted during team lifecycle.

The feature is gated by `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`
(env var or settings).

**How to verify:**

```bash
# Version check (Agent Teams introduced ~2.1.32, evolved since)
claude --version

# Tool names in the installed binary (works on macOS / Linux)
strings "$(realpath "$(which claude)")" \
  | grep -E '^(TeamCreate|TeamDelete|SendMessage|TeammateIdle)$'

# Public-facing docs (may describe behavior, not all tool names)
# https://code.claude.com/docs/en/agent-teams
```

**Don't "fix" by:** removing `TeamCreate` calls; replacing the team
flow with `Agent`-only natural-language spawn; removing the
`CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS` env var on the assumption it's
unused; treating `SendMessage` between teammates as nonexistent and
routing everything through the Lead.

---

## `isolation: worktree` is valid frontmatter for spawned teammates

**Common false belief:** `isolation: worktree` is not a documented
frontmatter field; it must be a typo or stale leftover from sub-agents
and should be removed from teammate definitions.

**Verified:** The `Agent`-tool schema for spawned teammates explicitly
extends with:

```
isolation: V.enum(["worktree"]).optional()
  .describe('Isolation mode. "worktree" creates a temporary git
             worktree so the agent works on an isolated copy of the
             repo.')
cwd: V.string().optional()
  .describe('Absolute path to run the agent in. Mutually exclusive
             with isolation: "worktree".')
```

The kit's teammate definitions deliberately use `isolation: worktree`
so each teammate works on its own branch in its own checkout. Removing
the field collapses every teammate into the lead's working tree and
breaks the branching model.

**How to verify:**

```bash
strings "$(realpath "$(which claude)")" \
  | grep -E 'isolation.*worktree|spawned teammate'
```

**Don't "fix" by:** removing `isolation: worktree` from agent
frontmatter; replacing it with `cwd`-only spawning; assuming the field
is silently ignored and therefore safe to delete.

---

## Teammates can spawn synchronous subagents (just not other teammates or background agents)

**Common false belief:** "Teammates cannot spawn their own teams or
subagents — only the lead manages the team," so any `Agent` call from
a teammate is broken and should be removed.

**Verified:** The binary enforces only two restrictions on
teammate-initiated spawns:

- *"Teammates cannot spawn other teammates — the team roster is
  flat."* — i.e., no nested teams; teammates must omit the `name` /
  `team_name` / `mode` fields when calling `Agent`.
- *"In-process teammates cannot spawn background agents. Use
  `run_in_background=false` for synchronous subagents."*

Synchronous subagent calls (e.g., a Coder forking a research helper
via `Agent`) are explicitly allowed — the binary's error message even
points the way: *"To spawn a subagent instead, omit them"* (referring
to `name`/`team_name`/`mode`).

**How to verify:**

```bash
strings "$(realpath "$(which claude)")" \
  | grep -E 'cannot spawn|team roster is flat|run_in_background'
```

**Don't "fix" by:** removing `Agent` calls from teammate definitions;
forbidding teammates from forking; routing every fork through the
Lead.

---

## Public docs describe user-facing behavior, not the full tool surface

**Common false belief:** If a tool isn't named in
`code.claude.com/docs/...`, it doesn't exist; the docs are exhaustive
about Claude Code's tool surface.

**Verified:** Several real tools are deferred (loaded via
`ToolSearch`) and described in the public docs only by *behavior*, not
by name. `TeamCreate` is the canonical example: the docs say *"ask
Claude to create a team in plain language,"* which is accurate from
the human's perspective but does not imply the tool name doesn't
exist. The lead-side mechanism is `TeamCreate`.

**How to verify (when the docs are silent):** grep the binary for the
tool name or its `searchHint` string. The binary is the source of
truth for the tool surface; docs are the source of truth for intended
*usage*.

```bash
# List all deferred-tool searchHints in the binary
strings "$(realpath "$(which claude)")" \
  | grep -oE 'searchHint:"[^"]+"' | sort -u
```

**Don't "fix" by:** declaring a tool nonexistent because the public
docs don't name it; redesigning the kit around the
documented-by-behavior surface only; treating "I couldn't find it in
the docs" as proof of absence.

---

## Session resumption drops team context (lead and teammate alike)

**Common false belief:** `/resume` restores in-process teammates, or
at least re-attaches the lead to its on-disk team state. The kit can
rely on resume to pick up where a team left off.

**Verified:** Team context is held entirely in memory. The
reconnection path (`computeInitialTeamContext`) only restores context
when the *parent process* set dynamic team context at spawn time
(`setDynamicTeamContext` writes to a module-level variable; resume
doesn't repopulate it).

On `/resume`:

- The lead's `teamContext` is null. The lead has no awareness of any
  prior team.
- The on-disk config at `~/.claude/teams/<team_name>/config.json`
  persists, but `TeamDelete` cannot reach it — it reads
  `teamContext.teamName` from in-memory state, which is empty.
- Reconnection logs for orphaned cases:
  *"No teammate context set (not a teammate)."*

The public docs say `/resume` doesn't restore in-process teammates;
the actual scope is broader — the *lead* loses team awareness too,
not just the teammate processes.

**How to verify:**

```bash
strings "$(realpath "$(which claude)")" \
  | grep -E 'computeInitialTeamContext|setDynamicTeamContext|No teammate context set'
```

Force-quit a team session and inspect `~/.claude/teams/` — the
`config.json` for the orphaned team remains; a new `claude` session
in the same project starts with empty `teamContext`.

**Don't "fix" by:** assuming `/resume` returns the lead to a working
team; calling `TeamDelete` to clean up after resume (it cannot see
the prior `team_name`); leaving the kit silent about orphan state on
disk.

**Kit guard:** the lead's Pre-Start Check looks for
`~/.claude/teams/{{TEAM_NAME}}/config.json` and prompts the human to
clean it up before calling `TeamCreate`. Manual cleanup:
`rm -rf ~/.claude/teams/{{TEAM_NAME}}` — `TeamDelete` cannot recover
context once the session has ended.

---

## When this file disagrees with a fresh fetch from `code.claude.com`

If the public docs and this file disagree, neither is automatically
right. The public docs lag the binary, and this file lags the public
docs. The binary is authoritative for *what tools exist*; the public
docs are authoritative for *intended usage and policy*.

When in conflict:

1. Re-run the verification command in this file's entry against the
   currently-installed binary.
2. Fetch the relevant public doc page.
3. If the binary still has the tool/field but the docs say something
   different, prefer the binary for "does this work" and the docs for
   "how should this be used."
4. Update this file's "Last verification" line and the affected entry
   so the next session sees the current state.
