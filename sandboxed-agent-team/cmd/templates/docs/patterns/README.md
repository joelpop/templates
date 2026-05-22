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

| Area                                                         | Vaadin 24.0 – 24.x                                                                                                                                                           | Vaadin 25.0                                                                                                              | Vaadin ≥25.1                                                                                                      |
|--------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------|
| Main layout + access checker                                 | `@PermitAll` on main layout; `@Layout` annotation available from 24.1. Earlier 24.0 uses `MainLayout implements RouterLayout` + explicit `layout = MainLayout.class` on each `@Route`. | `@AnonymousAllowed` **required** on `@Layout` class — `AnnotatedViewAccessChecker` does not reliably find `@PermitAll` on layouts. | Same as 25.0.                                                                                                     |
| `@Menu` annotation for nav entries                           | Available from Vaadin 24.4+. On older 24.x, register `SideNavItem` instances manually in `MainLayout`.                                                                      | Available.                                                                                                               | Available.                                                                                                        |
| Vaadin Signals (cross-session reactive state)                | Not available — use private state-management methods (fields, listeners, manual rebind).                                                                                     | Available; use for cross-session reactive data and `Binder`-external state where it adds value.                          | **Preferred for all non-`Binder` component state management.** `Binder` is still preferred for bean-backed forms. Pre-25.1 projects fall back to private state-management methods. |
| Phone bottom tab bar via `addToNavbar(true, ...)` touch-optimized slot | Not available — compose a `Tabs` component manually at the bottom of the main layout for touch navigation.                                                          | Available.                                                                                                               | Available.                                                                                                        |
| Extended client details API (for timezone, etc.)             | `UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> ...)` asynchronous callback. Cache in `VaadinSession`.                                                   | Consult current docs — may be synchronous.                                                                               | Consult current docs.                                                                                             |

Spring Boot's patterns used here (`AuditorAware`, `BCryptPasswordEncoder`, `@Transactional`,
`@EnableJpaAuditing`, `HttpSessionRequestCache`, `spring.jpa.open-in-view`, `@Cacheable`) are
identical between Spring Boot 3.x and 4.x. Spring Boot 3 requires Java 17+; Spring Boot 4
raises the floor further — pick the JDK to match the Spring Boot line you target.

Java 21+ is assumed for the `_` unnamed-variable syntax shown in the Java conventions. On
Java 17 projects, use a conventional short name (e.g., `e` or the specific event type) for
the unused lambda parameter instead.

If you're targeting a specific Vaadin version, you can tighten any inline "≥X / <X" split to
the single side that applies to your project.

## How to Reuse

1. Copy this directory into the new project's `docs/` tree.
2. Search for `{app}` and replace with the application prefix (e.g., `myapp`).
3. Search for `{base_package}` and replace with the root Java package (e.g., `com.example.myapp`).
4. Pick your target Vaadin and Spring Boot versions, then delete or flatten the inline
   "Vaadin ≥X / <X" notes that don't apply to your target. The matrix above identifies every
   place those notes appear.
5. Create application-specific docs in the parent `docs/` tree that reference these patterns
   using relative links (e.g., `See docs/patterns/architecture/persistence.md`).

## Patterns vs. Recipes

**Patterns** establish what to do and why — the rule, the obligation, the rationale.
**Recipes** are reference implementations of those patterns — concrete, copy-adaptable
code showing *how* to satisfy the pattern in a specific technology context.

A recipe lives under a `recipes/` sub-directory of the relevant technology directory
and assumes the reader has already read the corresponding pattern. It does not
re-argue the why.

## Conventions vs. Requirements

`patterns/` contains descriptive guidance, not checkbox requirements. Application-specific
docs carry the `[ ]` checkboxes that reference these patterns.

## Writing Pattern Documents

Every convention and pattern document in this kit should follow two structural
rules so agents and developers know when and how to apply the practice.

### Open with a scope statement

The first sentence (before any section heading) should answer: "what should this
be applied to?" State the obligation directly:

> Every Vaadin view and custom component class should have a layout diagram in
> its class Javadoc.

Without a scope statement the document describes mechanics with no stated
audience or trigger. A reader scanning the file cannot tell whether it applies
to their current task.

If the document covers multiple unrelated practices, it should be split into
separate single-practice files — possibly grouped into a subdirectory — so each
file can carry its own concise, unambiguous obligation.

### Write INDEX.md entries that lead with the obligation

The INDEX.md description is the only thing an agent sees when scanning the index.
It should answer "when do I apply this?" at a glance, not enumerate topics:

**Avoid** — topic list only:
> Text-based layout diagrams in Javadoc `<pre>` blocks: placement, box labeling…

**Prefer** — obligation first, then topics:
> Every Vaadin view and custom component class should have a layout diagram in
> its class Javadoc. Defines placement, box labeling…