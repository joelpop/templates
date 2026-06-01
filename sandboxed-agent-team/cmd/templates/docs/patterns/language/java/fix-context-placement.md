# Where Fix Context Belongs

When a fix, refactor, or design decision produces context worth preserving,
choose the right home for it so it's discoverable without polluting source code.

- **Commit message body** — the primary home. One paragraph on what was wrong,
  why the fix works, what alternatives were considered.
- **PR description** — for cross-cutting context, screenshots, links to issues.
- **A `docs/solutions/` entry** — if the fix changes a pattern for this project,
  the new pattern goes in the architecture docs.
- **A new `docs/patterns/` entry** — if the fix surfaced a project-agnostic rule
  worth carrying across projects.