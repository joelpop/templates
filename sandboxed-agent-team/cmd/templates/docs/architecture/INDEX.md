# Architecture — Index

Project-specific architecture and design documentation. *How this
codebase realizes the requirements* — names classes, traces flows,
records non-obvious design choices, and captures the patterns this
project commits to (which may be project-specific applications of a
generic pattern from `docs/patterns/`).

This tree is curated by the **Architect** and committed by the
**Analyst** (per the GLOSSARY AND ARCHITECTURE COMMITS rule in the
Analyst role). New entries are added when:

- The Architect proposes a structural approach at task kickoff and
  the human approves it (the approved approach lands here).
- An abstraction is recognized during implementation (e.g., the
  `ContentData` value-object pattern from
  `docs/patterns/conventions/abstraction.md`) and the project
  commits to it.
- A pattern is *project-specific* — it isn't portable to other
  projects, but matters for this one. (Portable patterns go in
  `docs/patterns/` instead.)

| Path | Description |
|------|-------------|
| [architecture-debt.md](architecture-debt.md) | Tracks known structural debt, deferred decisions, and recommended resolutions for this project. Append-only log of `AD-###` entries. |

## Adding a new architecture entry

Each entry is a `.md` file in this directory or a topic-named
subdirectory if the area grows large enough to warrant it. Use a
descriptive filename (`auth-flow.md`, `multi-tenancy.md`,
`reporting-pipeline.md`) — not a number. Keep one architectural
concern per file; cross-link with relative paths.

Suggested entry structure:

```markdown
# <Area> Architecture

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
