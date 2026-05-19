# Documentation Map

The project's documentation lives in **four trees** plus a **glossary**,
each with a distinct purpose, owner, and lifecycle. This file is the
master pointer; each tree and the glossary carry their own `INDEX.md`.

## [Requirements](reqs/INDEX.md) — `docs/reqs/`

What the system must do and which constraints it must satisfy.
IEEE 830 / ISO 29148 (SRS structure) and ISO 25010 (quality model).
Functional features, non-functional quality attributes, technical
constraints, environmental requirements, external interfaces, and
cross-cutting concerns. Owned by the **Analyst**; human-approved.

Each requirement carries status checkboxes — see "Requirement
Status" in `CLAUDE.md` for the convention.

Also under this tree: [`open-items.md`](reqs/open-items.md), the
single tracker for questions requiring human input before
implementation can proceed.

## [Patterns](patterns/INDEX.md) — `docs/patterns/`

**Project-agnostic** conventions, architecture patterns, recipes,
and implementation guidelines for the project's framework stack.
Designed to carry across projects — extractable, reusable, and
augmentable. Owned and committed by the **Architect** on
`pattern/<slug>` branches. Includes generic conventions (Java,
Vaadin, naming, Lombok, comments, abstraction, fix discipline),
generic architecture patterns (modules, persistence, services,
security), UI patterns, testing patterns, deployment patterns,
Figma translation, and **recipes** (step-by-step how-to guides for
recurring capabilities like passkey auth, OIDC/SSO).

## [Solutions](solutions/INDEX.md) — `docs/solutions/`

**Project-specific** documentation of how *this codebase* realizes
its requirements — names classes, traces flows, records non-obvious
implementation choices, and captures how the project applies
patterns from `docs/patterns/` to its specific context. Owned and
committed by the **Coder** on the task branch.

Also under this tree:
[`technical-debt.md`](solutions/technical-debt.md), the
append-only log of known structural debt and recommended
resolutions.

## [Guides](guides/INDEX.md) — `docs/guides/`

End-user, administrator, and operator-facing guides for using the
running system. Distinct from the three trees above — those are
about *building* the system; this tree is about *using* it. Owned
by the **Tech Writer**. Updated at release boundaries, not on every
implementation task.

## [Glossary](glossary/INDEX.md) — `docs/glossary/`

The project's canonical vocabulary — business terms (Analyst) and
technical terms (Architect). Used by requirement docs, patterns,
solutions, and code via inline Markdown links.
