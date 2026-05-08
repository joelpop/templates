# Patterns Kit

This directory (`docs/patterns/`) contains application-agnostic conventions,
architecture patterns, recipes, and implementation guidelines for **Vaadin 24+
with Spring Boot 3+ and Spring Data JPA**.

The content here contains **zero references** to any specific application, domain, or
business entity. It is designed to be extracted and reused across projects.

## What This Provides

| Directory | Contents |
|-----------|----------|
| `conventions/` | Java, Vaadin, naming, and Lombok coding conventions; comment, abstraction, and fix discipline |
| `architecture/` | Module structure, persistence, service layer, and security patterns |
| `ui/` | Theming, navigation, components, responsive layout, and error views |
| `testing/` | Unit, browserless UI, and E2E testing patterns; acceptance-criteria traceability |
| `deployment/` | Fat JAR, Docker, Spring profiles, migrations, logging, health check |
| `figma/` | Translating Figma designs into Lumo themes and Vaadin component compositions |
| `recipes/` | Step-by-step how-to guides (passkey auth, OIDC/SSO, etc.) — distinct from rules and idioms above |

## Version Compatibility

The kit targets **Vaadin 24+** and **Spring Boot 3+**. Most patterns (JPA entity hierarchy,
service layer, theming, dialogs with `Binder`, testing, deployment) work identically across
every supported version. A handful of Vaadin APIs differ between releases — those patterns
carry an inline **"Vaadin ≥X" / "Vaadin <X"** note describing each side of the split.

The matrix below summarizes the version-sensitive areas:

| Area | Vaadin 24.0 – 24.x | Vaadin 25.0 | Vaadin ≥25.1 |
|------|--------------------|-------------|--------------|
| Main layout + access checker | `@PermitAll` on main layout; `@Layout` annotation available from 24.1. Earlier 24.0 uses `MainLayout implements RouterLayout` + explicit `layout = MainLayout.class` on each `@Route`. | `@AnonymousAllowed` **required** on `@Layout` class — `AnnotatedViewAccessChecker` does not reliably find `@PermitAll` on layouts. | Same as 25.0. |
| `@Menu` annotation for nav entries | Available from Vaadin 24.4+. On older 24.x, register `SideNavItem` instances manually in `MainLayout`. | Available. | Available. |
| Vaadin Signals (cross-session reactive state) | Not available — use private state-management methods (fields, listeners, manual rebind). | Available; use for cross-session reactive data and `Binder`-external state where it adds value. | **Preferred for all non-`Binder` component state management.** `Binder` is still preferred for bean-backed forms. Pre-25.1 projects fall back to private state-management methods. |
| Phone bottom tab bar via `addToNavbar(true, ...)` touch-optimized slot | Not available — compose a `Tabs` component manually at the bottom of the main layout for touch navigation. | Available. | Available. |
| Extended client details API (for timezone, etc.) | `UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> ...)` asynchronous callback. Cache in `VaadinSession`. | Consult current docs — may be synchronous. | Consult current docs. |

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

## Conventions vs. Requirements

Documents in `patterns/` are written as **conventions and patterns** — descriptive guidance,
not checkbox requirements. Application-specific docs carry the `[ ]` requirement checkboxes
that reference these patterns.
