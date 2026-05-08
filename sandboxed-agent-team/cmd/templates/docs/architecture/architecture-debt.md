# Architecture Debt

Tracks known structural debt, deferred decisions, and recommended
resolutions for this project. Each entry describes what the debt is,
why it exists, the impact if left unresolved, and the recommended
resolution.

New entries are added as debt is identified during implementation or
review. Resolved entries are kept for historical context, marked
**Resolved** with the resolution date.

This is an **append-only** log. Entries are not deleted; their
status changes (Open → In Progress → Resolved) but the history
stays. Each entry's ID (`AD-###`) is permanent and may be referenced
from commit messages and PRs.

---

## Format

Each entry follows this template:

### AD-### — [Short title]
**Identified:** YYYY-MM-DD
**Status:** Open | In Progress | Resolved (YYYY-MM-DD)
**Affects:** [component / doc references]

**Debt:** What the structural debt is.

**Why it exists:** The trade-off or constraint that produced it.

**Impact:** Cost of leaving it unresolved (performance, security,
velocity, etc.).

**Recommended resolution:** What we should do about it and when.

**Progress (YYYY-MM-DD):** *(Optional, repeat for each progress
update.)* What changed; any partial resolution; new findings.

---

## Entries

*(No entries yet. Add as debt is identified.)*

<!--
Example:

### AD-001 — Flat package structure

**Identified:** 2026-XX-XX
**Status:** Open
**Affects:** `<base-package>.*` Java sources

**Debt:** All classes live directly in the base package with no
sub-packages. The flat layout will not survive growth into multiple
domain entities, repositories, and views.

**Why it exists:** Normal for a scaffold — the starter doesn't
impose a package layout.

**Impact:** Moving later means rewriting imports across a growing
codebase and rebasing in-flight branches. Moving early is cheap;
delay compounds.

**Recommended resolution:** Before the first real feature lands,
introduce the package layout described in
`docs/patterns/architecture/modules.md` and
`docs/patterns/conventions/naming.md`.
-->
