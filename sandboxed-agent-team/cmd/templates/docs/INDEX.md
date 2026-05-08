# Documentation Map

The project's documentation lives in **four trees**, each with a
distinct purpose, owner, and lifecycle. This file is the master
pointer; each tree carries its own comprehensive index.

A shared **glossary** sits at the root, used by all four trees via
inline Markdown links.

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
augmentable. Owned by the **Architect**; committed by the
**Analyst**. Includes generic conventions (Java, Vaadin, naming,
Lombok, comments, abstraction, fix discipline), generic
architecture patterns (modules, persistence, services, security),
UI patterns, testing patterns, deployment patterns, Figma
translation, and **recipes** (step-by-step how-to guides for
recurring capabilities like passkey auth, OIDC/SSO).

## [Architecture](architecture/INDEX.md) — `docs/architecture/`

**Project-specific** architecture and design — *how this codebase
realizes the requirements*. Names classes, traces flows, captures
non-obvious design choices, and documents the patterns this project
commits to (which may be project-specific applications of generic
patterns from `docs/patterns/`). Owned by the **Architect**;
committed by the **Analyst**.

Also under this tree:
[`architecture-debt.md`](architecture/architecture-debt.md), the
append-only log of known structural debt and recommended
resolutions.

## [Guides](guides/INDEX.md) — `docs/guides/`

End-user, administrator, and operator-facing guides for using the
running system. Distinct from the three trees above — those are
about *building* the system; this tree is about *using* it. Owned
by the **Tech Writer**. Updated at release boundaries, not on every
implementation task.

## [Glossary](glossary.md) — `docs/glossary.md`

The project's canonical vocabulary. Used by requirement docs and
other documents via inline Markdown links. Curated by the
**Architect**, committed by the **Analyst**. Small by design —
holds only ambiguous, project-specific, or implementation-agnostic
terms.
