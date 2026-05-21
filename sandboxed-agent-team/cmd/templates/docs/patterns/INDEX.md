# Patterns — Index

Application-agnostic conventions, architecture patterns, recipes, and implementation
guidelines for **Vaadin 24+ with Spring Boot 3+** and Spring Data JPA. Version-
sensitive patterns carry inline "Vaadin ≥X / <X" notes; see [`README.md`](README.md)
→ "Version Compatibility" for the summary matrix. Zero project-specific content.

| Path | Description |
|------|-------------|
| [README.md](README.md) | What the kit provides and how to reuse it |

## Writing

Documentation-writing conventions (distinct from code conventions below).

| Path                                                          | Description                                                                                                                                                                          |
|---------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [writing/conventions.md](writing/conventions.md)             | When writing a convention or pattern document, open with a scope statement and lead INDEX.md entries with the obligation so readers immediately know when and why the pattern applies |
| [writing/requirements.md](writing/requirements.md)           | When writing requirements, use system-facing imperative voice and modal verbs so each statement is unambiguous, testable, and distinguishable from user stories                      |

## Language

Universal language idioms — no framework assumptions.

| Path                                                          | Description                                                                                                                                                                                    |
|---------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [language/singular-form.md](language/singular-form.md)       | Every named artifact (class, table, module, package leaf, route path) should use the singular form, not the plural                                                                             |
| [language/java/conventions.md](language/java/conventions.md) | When writing Java code in this project, follow these style and idiom conventions so every class reads consistently regardless of author                                                         |
| [language/java/lombok.md](language/java/lombok.md)           | When using Lombok on JPA entities or enums with properties, apply only the safe annotations so runtime failures from proxies, lazy loading, and bidirectional relationships are avoided         |
| [language/java/comments.md](language/java/comments.md)       | When writing code, default to no comments and add one only when the WHY is non-obvious so comments signal constraints and invariants, not history                                               |
| [language/java/abstraction.md](language/java/abstraction.md) | When encountering duplication, wait for the third instance before extracting an abstraction so the right abstraction size becomes clear from real usage                                          |

## Structure

How the project and code are organized — independent of any tool or framework.

| Path                                              | Description                                                                                                                                                                                                |
|---------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [structure/modules.md](structure/modules.md)      | When setting up a new application, structure it as a multi-module Maven project with a runtime-scoped `jpaservice` dependency so the UI layer cannot reference JPA types at compile time                   |
| [structure/services.md](structure/services.md)    | When implementing services, use a technology-prefix name, separate queries from mutations, apply `@Transactional` at the service layer, and never pass JPA entities or Vaadin types across the boundary    |

## CI/CD

Build tooling and deployment operations.

| Path                                                        | Description                                                                                                                                                          |
|-------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [cicd/deployment/patterns.md](cicd/deployment/patterns.md) | When deploying, use a fat JAR or Docker image with profile-specific properties and versioned Flyway migrations so the deployment is reproducible and rollback-safe   |

## Persistence

Persistence-technology-specific patterns and naming.

| Path                                                                            | Description                                                                                                                                                                      |
|---------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [persistence/spring-data-jpa/naming.md](persistence/spring-data-jpa/naming.md)     | When creating JPA entities, enums, repositories, or projections, use these naming conventions so the persistence-layer role of each artifact is unambiguous from its name alone |
| [persistence/spring-data-jpa/patterns.md](persistence/spring-data-jpa/patterns.md) | When working with JPA, follow these entity design and query patterns so Hibernate's lifecycle, lazy loading, and dirty checking behave predictably                               |

## Security

| Path                                                                                              | Description                                                                                                                                                                           |
|---------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [security/patterns.md](security/patterns.md)                                                      | When adding authentication and authorization, follow these patterns for Spring Security config, RBAC, session management, security headers, and PII protection                        |
| [security/spring/recipes/audited-principal.md](security/spring/recipes/audited-principal.md)      | When integrating JPA auditing with Spring Security, use this recipe to wire a principal that carries the user key for `@CreatedBy`/`@LastModifiedBy`                                 |
| [security/passkey/recipes/passkey.md](security/passkey/recipes/passkey.md)                        | When adding passkey (WebAuthn) authentication, follow this recipe to configure the Spring Security integration and credential storage                                                  |
| [security/oidc/recipes/oidc-sso.md](security/oidc/recipes/oidc-sso.md)                            | When integrating an OIDC/SSO provider, follow this recipe to configure Spring Security's OAuth2 login and map claims to application roles                                             |
| [security/form-login/recipes/form-login.md](security/form-login/recipes/form-login.md)            | When adding username/password form login, follow this recipe to configure Spring Security's form login with BCrypt and entropy-based password validation                              |
| [security/form-login/recipes/conditional-auth.md](security/form-login/recipes/conditional-auth.md)| When the application must support multiple auth methods selectable at deploy time, follow this recipe to gate each method via `application.properties` flags                          |

## UI

| Path                                                                                    | Description                                                                                                                                                                                        |
|-----------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [ui/vaadin/spring.md](ui/vaadin/spring.md)                                              | When registering Vaadin components as Spring beans, use `@SpringComponent` and the appropriate scope so Vaadin's lifecycle and Spring's DI work together correctly                                  |
| [ui/vaadin/views.md](ui/vaadin/views.md)                                                | When creating a Vaadin view, use `Composite<T>`, place it in its own sub-package, and annotate it with `@Menu` so navigation and package structure stay consistent                                  |
| [ui/vaadin/component-construction.md](ui/vaadin/component-construction.md)              | When building Vaadin components, initialize member variables in the constructor, delegate dialogs, and use Signals for reactive state so components stay cohesive and testable                       |
| [ui/vaadin/layout-diagram.md](ui/vaadin/layout-diagram.md)                              | Every Vaadin view and custom component class should have a layout diagram in its class Javadoc so the intended structure is visible without running the application                                   |
| [ui/vaadin/uimodel.md](ui/vaadin/uimodel.md)                                            | Types in `{app}-uimodel` should be named for their UI context, implement `HasCaption` when used in selection components, and implement capability interfaces when generic UI components bind to their structure |
| [ui/vaadin/datetime.md](ui/vaadin/datetime.md)                                          | When displaying dates and times, use `DateTimeUtil` with a zone and locale resolved per-call from `ClientDetailsService` so values reflect the user's local settings                                |
| [ui/vaadin/theming.md](ui/vaadin/theming.md)                                            | When styling the application, use Lumo CSS custom properties and `LumoUtility` classes so the theme stays consistent and maintainable without custom CSS                                            |
| [ui/vaadin/navigation.md](ui/vaadin/navigation.md)                                      | When building navigation, use `AppLayout` + `SideNav` with `@Menu`-annotated routes so the nav structure is declarative and conditionally rendered by role                                         |
| [ui/vaadin/components.md](ui/vaadin/components.md)                                      | When building grids, forms, and dialogs, follow these UI component patterns so interaction behavior stays consistent across views                                                                    |
| [ui/vaadin/responsive.md](ui/vaadin/responsive.md)                                      | When building views, follow these responsive layout patterns so the UI adapts correctly to mobile, tablet, and desktop breakpoints                                                                   |
| [ui/vaadin/error-views.md](ui/vaadin/error-views.md)                                    | When handling navigation and server errors, use these error view patterns so error pages are consistent and do not leak sensitive information                                                        |
| [ui/vaadin/recipes/base-view.md](ui/vaadin/recipes/base-view.md)                        | When creating the base view scaffold, follow this recipe to set up `MainLayout`, `AppLayout`, and the shared navigation chrome                                                                      |
| [ui/vaadin/recipes/item-browser.md](ui/vaadin/recipes/item-browser.md)                  | When building a grid-plus-editor view, follow this recipe to implement the `ItemBrowser` pattern with quick filter, record count, and inline edit dialog                                            |
| [ui/vaadin/recipes/app-icon.md](ui/vaadin/recipes/app-icon.md)                          | When adding an application icon (favicon and PWA icon), follow this recipe                                                                                                                         |
| [ui/vaadin/recipes/view-icon.md](ui/vaadin/recipes/view-icon.md)                        | When adding icons to side-nav menu entries, follow this recipe to supply icons via `@Menu`                                                                                                         |
| [ui/vaadin/recipes/java-lit-bridge.md](ui/vaadin/recipes/java-lit-bridge.md)            | When a Vaadin view needs a custom Lit web component, follow this recipe to bridge the Java and TypeScript sides                                                                                     |

## Design

| Path                                                                              | Description                                                                                                                                                              |
|-----------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [design/figma/figma-to-lumo-theme.md](design/figma/figma-to-lumo-theme.md)       | When translating a Figma design system, map color tokens and spacing to Lumo CSS custom properties so the implementation stays in sync with the design source            |
| [design/figma/figma-to-vaadin.md](design/figma/figma-to-vaadin.md)               | When implementing a Figma design, translate layouts and component structures to Vaadin compositions using these patterns so the result matches the design intent          |

## Testing

| Path                                                                                          | Description                                                                                                                                                                             |
|-----------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [testing/patterns.md](testing/patterns.md)                                                    | When writing tests, follow the pyramid structure with one test class per production class and trace every acceptance criterion to a passing test so coverage is intentional and verifiable |
| [testing/recipes/testbench-browserless.md](testing/recipes/testbench-browserless.md)          | When writing fast UI tests without a browser, follow this recipe to set up Vaadin's browserless `ComponentTester<T>` infrastructure                                                     |
| [testing/recipes/testbench-e2e-server.md](testing/recipes/testbench-e2e-server.md)            | When running TestBench E2E tests, follow this recipe to configure the JUnit 5 server extension with concurrency limiting                                                                |
| [testing/recipes/testbench-e2e-parallel.md](testing/recipes/testbench-e2e-parallel.md)        | When parallelizing TestBench E2E tests, follow this recipe to configure `@Execution` and the concurrency limit so tests run fast without saturating the server                          |