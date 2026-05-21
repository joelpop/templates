# Patterns — Index

Application-agnostic conventions, architecture patterns, recipes, and implementation
guidelines for **Vaadin 24+ with Spring Boot 3+** and Spring Data JPA. Version-
sensitive patterns carry inline "Vaadin ≥X / <X" notes; see [`README.md`](README.md)
→ "Version Compatibility" for the summary matrix. Zero project-specific content.

| Path | Description |
|------|-------------|
| [README.md](README.md) | What the kit provides and how to reuse it |

## Conventions

| Path                                                                         | Description                                                                                                                                                                                                             |
|------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [conventions/java.md](conventions/java.md)                                   | When writing Java code in this project, follow these style and idiom conventions so every class reads consistently regardless of author                                                                                 |
| [conventions/vaadin/spring.md](conventions/vaadin/spring.md)                 | When registering Vaadin components as Spring beans, use `@SpringComponent` and the appropriate scope so Vaadin's lifecycle and Spring's DI work together correctly and to avoid name collision with Vaadin's `Component` |
| [conventions/vaadin/views.md](conventions/vaadin/views.md)                   | When creating a Vaadin view, use `Composite<T>`, place it in its own sub-package, and annotate it with `@Menu` so navigation and package structure stay consistent                                                      |
| [conventions/vaadin/components.md](conventions/vaadin/components.md)         | When building Vaadin components, initialize member variables in the constructor, delegate dialogs, and use Signals for reactive state so components stay cohesive and testable                                          |
| [conventions/vaadin/layout-diagram.md](conventions/vaadin/layout-diagram.md) | Every Vaadin view and custom component class should have a layout diagram in its class Javadoc so the intended structure is visible without running the application                                                     |
| [conventions/vaadin/uimodel.md](conventions/vaadin/uimodel.md)               | Types in `{app}-uimodel` should be named for their UI context, implement `HasCaption` when used in selection components, and implement capability interfaces when generic UI components bind to their structure         |
| [conventions/vaadin/datetime.md](conventions/vaadin/datetime.md)             | When displaying dates and times, use `DateTimeUtil` with a zone and locale resolved per-call from `ClientDetailsService` so values reflect the user's local settings                                                    |
| [conventions/singular-form.md](conventions/singular-form.md)                 | Every named artifact (class, table, module, package leaf, route path) should use the singular form, not the plural                                                                                                      |
| [conventions/lombok.md](conventions/lombok.md)                               | When using Lombok on JPA entities or enums with properties, apply only the safe annotations so runtime failures from proxies, lazy loading, and bidirectional relationships are avoided                                 |
| [conventions/comments.md](conventions/comments.md)                           | When writing code, default to no comments and add one only when the WHY is non-obvious so comments signal constraints and invariants, not history                                                                       |
| [conventions/abstraction.md](conventions/abstraction.md)                     | When encountering duplication, wait for the third instance before extracting an abstraction so the right abstraction size becomes clear from real usage                                                                 |

## Writing

Documentation-writing conventions (distinct from code conventions
above). The `conventions/comments.md` rule above straddles —
it's about comments *inside code* but is about writing as well.
The entries here are about writing the documents themselves.

| Path                                                          | Description                                                                                                                                                                          |
|---------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [writing/requirements.md](writing/requirements.md)           | When writing requirements, use system-facing imperative voice and modal verbs so each statement is unambiguous, testable, and distinguishable from user stories                      |
| [writing/conventions.md](writing/conventions.md)             | When writing a convention or pattern document, open with a scope statement and lead INDEX.md entries with the obligation so readers immediately know when and why the pattern applies |

## Architecture

| Path                                                            | Description                                                                                                                                                                                                              |
|-----------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [architecture/modules.md](architecture/modules.md)             | When setting up a new application, structure it as a multi-module Maven project with runtime-scoped `*service` implementation dependencies so the UI layer cannot reference service implementation types at compile time |
| [architecture/persistence.md](architecture/persistence.md)     | When working with Spring Data JPA, follow these entity design and query patterns so lifecycle, lazy loading, and dirty checking behave predictably                                                           |
| [architecture/services.md](architecture/services.md)           | When implementing services, use a technology-prefix name, separate queries from mutations, apply `@Transactional` at the service layer, and never pass JPA entities or Vaadin types across the boundary                  |
| [architecture/security.md](architecture/security.md)           | When adding authentication and authorization, follow these patterns for BCrypt, passkey, OIDC, session management, RBAC, and security headers so the application meets baseline security requirements                    |

## Persistence

| Path                                                                              | Description                                                                                                                                                    |
|-----------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [persistence/spring-data-jpa/naming.md](persistence/spring-data-jpa/naming.md)   | When creating JPA entities, enums, repositories, or projections, use these naming conventions so the persistence-layer role of each artifact is unambiguous from its name alone |

## UI

| Path                                            | Description                                                                                                                                                          |
|-------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [ui/theming.md](ui/theming.md)                  | When styling the application, use Lumo CSS custom properties and `LumoUtility` classes so the theme stays consistent and maintainable without custom CSS              |
| [ui/navigation.md](ui/navigation.md)            | When building navigation, use `AppLayout` + `SideNav` with `@Menu`-annotated routes so the nav structure is declarative and conditionally rendered by role           |
| [ui/components.md](ui/components.md)            | When building grids, forms, and dialogs, follow these UI component patterns so interaction behavior stays consistent across views                                     |
| [ui/responsive.md](ui/responsive.md)            | When building views, follow these responsive layout patterns so the UI adapts correctly to mobile, tablet, and desktop breakpoints                                    |
| [ui/error-views.md](ui/error-views.md)          | When handling navigation and server errors, use these error view patterns so error pages are consistent and do not leak sensitive information                         |

## Testing

| Path                                              | Description                                                                                                                                                                                             |
|---------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [testing/patterns.md](testing/patterns.md)        | When writing tests, follow the pyramid structure with one test class per production class and trace every acceptance criterion to a passing test so coverage is intentional and verifiable               |

## Deployment

| Path                                                  | Description                                                                                                                                                                        |
|-------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [deployment/patterns.md](deployment/patterns.md)      | When deploying, use a fat JAR or Docker image with profile-specific properties and versioned Flyway migrations so the deployment is reproducible and rollback-safe                 |

## Figma

| Path                                                                        | Description                                                                                                                                                                  |
|-----------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [figma/figma-to-lumo-theme.md](figma/figma-to-lumo-theme.md)               | When translating a Figma design system, map color tokens and spacing to Lumo CSS custom properties so the implementation stays in sync with the design source                |
| [figma/figma-to-vaadin.md](figma/figma-to-vaadin.md)                       | When implementing a Figma design, translate layouts and component structures to Vaadin compositions using these patterns so the result matches the design intent              |

## Recipes

| Path                                    | Description                                                                                                                                                                                                                  |
|-----------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [recipes/INDEX.md](recipes/INDEX.md)    | When implementing a recurring capability (passkey auth, OIDC/SSO, multi-tenancy, etc.), follow the relevant recipe for a step-by-step walkthrough anchored to the specific framework and library versions in use             |