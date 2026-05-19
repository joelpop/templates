# Solutions — Index

Project-specific documentation of how this codebase realizes its
requirements — names classes, traces flows, records non-obvious
implementation choices, and captures how the project applies patterns
from `docs/patterns/` to its specific context.

This tree is authored and committed by the **Coder** (per the Directory
Ownership Rules in CLAUDE.md). New entries are added when:

- A non-obvious implementation choice is made during a task (e.g.,
  why a particular approach was taken over alternatives).
- An abstraction is recognized during implementation (e.g., the
  `ContentData` value-object pattern from
  `docs/patterns/conventions/abstraction.md`) and the project
  commits to it.
- A solution is project-specific — it isn't portable to other
  projects, but matters for this one. (Portable patterns go in
  `docs/patterns/` instead.)

| Path | Description |
|------|-------------|
| [technical-debt.md](technical-debt.md) | Tracks known structural debt, deferred decisions, and recommended resolutions for this project. Append-only log of `TD-###` entries. |

## Adding a new solutions entry

Each entry is a `.md` file in this directory or a topic-named
subdirectory if the area grows large enough to warrant it. Use a
descriptive filename (`auth-flow.md`, `multi-tenancy.md`,
`reporting-pipeline.md`) — not a number. Keep one concern per file;
cross-link with relative paths.

Suggested entry structure:

```markdown
# <Area> — Solution

## Overview
What this part of the system does and why.

## Decisions
The choices that shape this area, with rationale.

## Mechanism
How it works — classes, flows, integration points. Reference
specific source paths.

## References
Links to relevant `docs/reqs/` requirements, `docs/patterns/`
entries, and any external references.
```
