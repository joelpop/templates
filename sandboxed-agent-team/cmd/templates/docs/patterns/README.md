# Patterns Kit

Application-agnostic conventions, architecture patterns, recipes, and implementation
guidelines for **Vaadin 24+ with Spring Boot 3+ and Spring Data JPA**. Zero project-
specific content — designed to carry across projects on this stack.

## What This Provides

| Directory      | Contents                                                                                                                                     |
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| `language/`    | Universal Java idioms, Lombok guidelines, comment discipline, abstraction patterns, singular-form naming                                      |
| `structure/`   | Multi-module project organization and service-layer design                                                                                    |
| `cicd/`        | Build tooling and deployment operations                                                                                                       |
| `persistence/` | Spring Data JPA naming, entity design, and query patterns                                                                                    |
| `security/`    | Spring Security config, RBAC, session management, and auth-method recipes (passkey, OIDC, form login)                                        |
| `ui/vaadin/`   | Vaadin coding conventions, UI component patterns, and view recipes                                                                           |
| `design/`      | Translating Figma designs into Lumo themes and Vaadin component compositions                                                                 |
| `testing/`     | Unit, browserless UI, and E2E testing patterns; acceptance-criteria traceability                                                             |

## Version Compatibility

The kit targets **Vaadin 24+** and **Spring Boot 3+**. Most patterns (JPA entity hierarchy,
service layer, theming, dialogs with `Binder`, testing, deployment) work identically across
every supported version. A handful of Vaadin APIs differ between releases — those patterns
carry an inline **"Vaadin ≥X" / "Vaadin <X"** note describing each side of the split.

The matrix below summarizes the version-sensitive areas:

| Area                                                                   | Vaadin 24.0 – 24.x                                                                                                                                                                     | Vaadin 25.0                                                                                                                        | Vaadin ≥25.1                                                                                                                                                                       |
|------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Main layout + access checker                                           | `@PermitAll` on main layout; `@Layout` annotation available from 24.1. Earlier 24.0 uses `MainLayout implements RouterLayout` + explicit `layout = MainLayout.class` on each `@Route`. | `@AnonymousAllowed` **required** on `@Layout` class — `AnnotatedViewAccessChecker` does not reliably find `@PermitAll` on layouts. | Same as 25.0.                                                                                                                                                                      |
| `@Menu` annotation for nav entries                                     | Available from Vaadin 24.4+. On older 24.x, register `SideNavItem` instances manually in `MainLayout`.                                                                                 | Available.                                                                                                                         | Available.                                                                                                                                                                         |
| Vaadin Signals (cross-session reactive state)                          | Not available — use private state-management methods (fields, listeners, manual rebind).                                                                                               | Available; use for cross-session reactive data and `Binder`-external state where it adds value.                                    | **Preferred for all non-`Binder` component state management.** `Binder` is still preferred for bean-backed forms. Pre-25.1 projects fall back to private state-management methods. |
| Phone bottom tab bar via `addToNavbar(true, ...)` touch-optimized slot | Not available — compose a `Tabs` component manually at the bottom of the main layout for touch navigation.                                                                             | Available.                                                                                                                         | Available.                                                                                                                                                                         |
| Extended client details API (for timezone, etc.)                       | `UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> ...)` asynchronous callback. Cache in `VaadinSession`.                                                             | Consult current docs — may be synchronous.                                                                                         | Consult current docs.                                                                                                                                                              |
| Spring Boot compatibility                                              | Spring Boot 3.x and 4.x patterns are identical for the APIs used here. Spring Boot 3 requires Java 17+; Spring Boot 4 raises the floor further — pick the JDK to match.                | Same.                                                                                                                              | Same.                                                                                                                                                                              |
| Java version                                                           | Java 17–20: use a conventional short name (e.g., `unused`) for unused lambda parameters — `_` is not available.                                                                        | Java 21+: use `_` (unnamed variable) for unused lambda parameters.                                                                 | Java 21+.                                                                                                                                                                          |

If you're targeting a specific Vaadin version, tighten any inline "≥X / <X" split to
the single side that applies to your project.

## How to Reuse

1. Copy this directory into the new project's `docs/` tree.
2. Search for `{app}` and replace with the application prefix (e.g., `myapp`).
3. Search for `{base_package}` and replace with the root Java package (e.g., `com.example.myapp`).
4. Pick your target Vaadin and Spring Boot versions, then delete or flatten the inline
   "Vaadin ≥X / <X" notes that don't apply to your target. The matrix above identifies every
   place those notes appear.
5. Create application-specific docs in the parent `docs/` tree that reference these patterns
   using relative links (e.g., `See docs/patterns/persistence/spring-data-jpa/entity-hierarchy.md`).

## Conventions vs. Requirements

- `patterns/` contains descriptive guidance, not checkbox requirements.
- Application-specific docs carry the `[ ]` checkboxes that reference these patterns.

## Patterns vs. Recipes

- **Patterns** establish what to do and why — the rule, the obligation, the rationale.
- **Recipes** are reference implementations of those patterns — concrete, copy-adaptable
  code showing *how* to satisfy the pattern in a specific technology context.

A recipe lives under a `recipes/` sub-directory of the relevant technology directory
and assumes the reader has already read the corresponding patterns. It does not
re-argue the why.

## Writing Pattern Documents

Every convention and pattern document in this kit should follow two structural
rules so agents and developers know when and how to apply the practice.

### Open with a scope statement

The first sentence (before any section heading) states the obligation directly:
when to apply it, what to do, and why. Three criteria for a precise scope statement:

- **"When"** names a specific artifact or task — not a broad domain ("security")
  or technology area ("JPA"). If you can't name the specific thing being created
  or configured, the scope is too vague.
- **The obligation** states the actual rule — not "follow these patterns"
  (circular). It names what to do, use, or avoid.
- **"So"** names a specific failure prevented or benefit gained — not "so it
  works" or "so it is correct."

**Avoid** — domain-scoped, circular obligation, vague benefit:
> *When working on security, follow these patterns so the application is secure.*

**Prefer** — artifact-scoped, specific rule, concrete benefit:
> *When annotating a `@Route` view, add exactly one of `@AnonymousAllowed`,
> `@PermitAll`, `@RolesAllowed`, or `@DenyAll` so `AnnotatedViewAccessChecker`
> has an unambiguous access rule for every route.*

If a single precise scope statement cannot be written for the file, the file
covers too many independent practices and must be split. Each split-out file
gets its own scope statement.

### Formatting conventions

- **"Avoid" before "Preferred"** — always place the avoid example first, the preferred example second, each in a separate fenced code block.
- **Language tags** — use ` ```java `, ` ```xml `, ` ```sql `, ` ```properties `, etc. on every fenced code block. Do not omit the tag to suppress IDE warnings.

### Code examples illustrate the rule; they don't replace it

A code block inside a pattern shows what the obligation looks like in practice —
it only makes sense alongside the surrounding explanation. When a code block
becomes detailed enough to copy-adapt for a complete feature without reading the
prose, it belongs in a recipe instead.

## Writing Recipes

### Open with a scope statement

The scope statement for a recipe names a concrete implementation goal:

> *When implementing [specific task], follow this recipe to [produce specific output].*

The task is not a pattern name — it is what the developer is about to build.

### Recipes are working examples, not law

A recipe demonstrates one valid way to implement the relevant patterns. It is
not the only acceptable implementation. Projects must adapt recipes to their
actual technology choices — a recipe that uses Keycloak is a starting point for
an Okta implementation, not a mandate to use Keycloak.

### Recipes demonstrate multiple patterns

Recipes are not constrained to a single pattern. A realistic implementation
task routinely demonstrates several patterns together, and the recipe should
reflect that reality.

### Prerequisites

Reference the pattern documents the recipe demonstrates. The reader is assumed
to have read those patterns; the recipe does not re-argue the why.

### Structure

Numbered steps with copy-adaptable code blocks. Keep prose to a minimum —
the code is the primary communication.

### Don't introduce new patterns inline

If writing the recipe surfaces a pattern not yet documented, document it in a
pattern file first and reference it from the recipe. Recipes are implementations,
not the authoritative source for pattern rules.

---

## Writing INDEX.md Entries

The INDEX.md description is the only thing an agent sees when scanning the index.
Apply the same "When … use/do … so …" obligation form as the scope statement — it
should answer "when do I apply this?" at a glance, not enumerate topics:

**Avoid** — topic list only:
> Text-based layout diagrams in Javadoc `<pre>` blocks: placement, box labeling…

**Prefer** — obligation first, then topics:
> Every Vaadin view and custom component class should have a layout diagram in
> its class Javadoc. Defines placement, box labeling…