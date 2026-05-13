# Patterns — Index

Application-agnostic conventions, architecture patterns, recipes, and implementation
guidelines for **Vaadin 24+ with Spring Boot 3+** and Spring Data JPA. Version-
sensitive patterns carry inline "Vaadin ≥X / <X" notes; see [`README.md`](README.md)
→ "Version Compatibility" for the summary matrix. Zero project-specific content.

| Path | Description |
|------|-------------|
| [README.md](README.md) | What the kit provides and how to reuse it |

## Conventions

| Path | Description |
|------|-------------|
| [conventions/java.md](conventions/java.md) | Java coding conventions: `var`, member variable init, local variable declaration, nested types, event handler naming, lambdas, unused params, null checks, DI, SOLID, access modifiers, suppressing warnings, JavaDoc |
| [conventions/vaadin/spring.md](conventions/vaadin/spring.md) | Spring bean registration, scopes, security annotations, build configuration: `@SpringComponent`, session scope, `vaadin.allowed-packages`, access annotations on layouts |
| [conventions/vaadin/views.md](conventions/vaadin/views.md) | View structure, package layout, navigation: `Composite<T>`, per-view packages, `@Menu` |
| [conventions/vaadin/components.md](conventions/vaadin/components.md) | Component construction and state: constructor init order, dialog delegation, `NonComponent` events, Signals, `Binder`, `LumoUtility` |
| [conventions/vaadin/uimodel.md](conventions/vaadin/uimodel.md) | `{app}-uimodel` patterns: `HasCaption` for enums and picker records, `HasActive`/`HasRole` capability interfaces |
| [conventions/vaadin/datetime.md](conventions/vaadin/datetime.md) | Date/time display: `ClientDetailsService` bridge, `DateTimeUtil` (short/medium/long formatting, per-call zone/locale resolution) |
| [conventions/naming.md](conventions/naming.md) | Naming conventions: entities, services, UI models, DB columns (`_key` vs `_id`), packages, methods |
| [conventions/lombok.md](conventions/lombok.md) | Lombok guidelines: safe and unsafe annotations on JPA entities (`@Data` / `@EqualsAndHashCode` / `@ToString` pitfalls, managed collection fields, `@Builder`); `@Getter` + `@RequiredArgsConstructor` for enums with properties; `@Slf4j` for SLF4J logger declarations |
| [conventions/comments.md](conventions/comments.md) | Code comment discipline: default to no comments, comments justify invariants not history, the fix-mode trap (don't paste conversational explanations into the code), where the *why* of a fix actually belongs (commit message, PR description, architecture docs) |
| [conventions/abstraction.md](conventions/abstraction.md) | Abstraction recognition: the third-instance rule; value-object patterns (`ContentData`, `PersonName`); sizing the abstraction (value object → shared base class → component-family package) with wrong-size signs for each tier; where extracted abstractions live; when to leave duplication alone |
| [conventions/fixing.md](conventions/fixing.md) | Fix discipline: diagnose-before-fix classification (trivial / localized / structural), workaround signatures (suppression annotations, swallowed exceptions, defensive casts, masking null checks, copied code), fix-attempt limit, right-thing vs working-thing, revert-before-rework, tests are part of the fix |

## Writing

Documentation-writing conventions (distinct from code conventions
above). The `conventions/comments.md` rule above straddles —
it's about comments *inside code* but is about writing as well.
The entries here are about writing the documents themselves.

| Path | Description |
|------|-------------|
| [writing/requirements.md](writing/requirements.md) | Requirement-writing form: system-facing imperative ("the system must..."), modal verbs (must/shall vs. should/may), one concept per requirement, testability, unambiguous agnostic vocabulary, distinguishing requirements from user stories and from acceptance criteria, active voice |

## Architecture

| Path | Description |
|------|-------------|
| [architecture/modules.md](architecture/modules.md) | Multi-module Maven structure, compile-time layer separation, route scanning, `vaadin.allowed-packages`, `JpaConfig`, MapStruct |
| [architecture/persistence.md](architecture/persistence.md) | JPA best practices: entity hierarchy, `equals`/`hashCode`, `@Version`, auditing, OSIV disabled, `@ManyToOne LAZY`, `@Enumerated STRING`, insert vs update, `@DynamicUpdate`, projections, cascade strategies, `@Embeddable`, batch/bulk ops, N+1 prevention, `@DataJpaTest` |
| [architecture/services.md](architecture/services.md) | Service layer patterns: query/mutation separation, `@Transactional`, dirty-checking updates, grid loading, `@Cacheable` caching, error contracts, MapStruct mapper pattern |
| [architecture/security.md](architecture/security.md) | Auth patterns: BCrypt, entropy-based password validation, passkey/WebAuthn, session management, RBAC via Jakarta Security annotations, `@AnonymousAllowed` on layout, conditional nav/component rendering, security headers, CSRF, rate limiting, SQL injection prevention, PII in logs, file upload validation |

## UI

| Path | Description |
|------|-------------|
| [ui/theming.md](ui/theming.md) | Lumo theme, `LumoUtility`, component variants, brand customization via CSS custom properties |
| [ui/navigation.md](ui/navigation.md) | `AppLayout` + `SideNav`, `@Menu` annotation, route-path grouping, `DrawerToggle`, conditional nav rendering, navigation guards, mobile bottom tab bar |
| [ui/components.md](ui/components.md) | Quick Filter, Avatar, grid standards (`ListDataProvider`, record count, row click), form standards (`Binder`, inline validation), dialogs, notification toasts, confirmation dialogs, loading indicators |
| [ui/responsive.md](ui/responsive.md) | Breakpoints (mobile/tablet/desktop), Vaadin responsive tools, full-width dialogs on mobile, reduced grid columns, server-side breakpoint detection |
| [ui/error-views.md](ui/error-views.md) | Error view types (404, 403-as-404, 500, 400), shared base / chrome, action row pattern, sensitive information protection, logging standards |

## Testing

| Path | Description |
|------|-------------|
| [testing/patterns.md](testing/patterns.md) | Testing pyramid, one test class per production class, JUnit + Mockito unit tests, Vaadin browserless UI tests, `@DataJpaTest`, `@Transactional` rollback, H2 compatibility mode, N+1 detection, Playwright E2E, **acceptance-criteria traceability** (every AC has a passing test; tests named for the behavior they verify) |

## Deployment

| Path | Description |
|------|-------------|
| [deployment/patterns.md](deployment/patterns.md) | Fat JAR, Docker, Spring profiles (dev/staging/prod), Flyway migrations, versioned/immutable/idempotent scripts, seed data, rollback documentation, structured logging, HikariCP, health check |

## Figma

| Path | Description |
|------|-------------|
| [figma/figma-to-lumo-theme.md](figma/figma-to-lumo-theme.md) | Translating a Figma design system into a Vaadin Lumo theme — color tokens, spacing scale, component variants |
| [figma/figma-to-vaadin.md](figma/figma-to-vaadin.md) | Translating Figma layouts and components into Vaadin component compositions |

## Recipes

| Path | Description |
|------|-------------|
| [recipes/INDEX.md](recipes/INDEX.md) | Step-by-step how-to guides for capabilities that recur across projects on this stack — passkey auth, OIDC/SSO integration, conditional auth via `application.properties`, multi-tenancy patterns, etc. Distinct from `conventions/`/`architecture/` (rules and idioms) — recipes are *implementation walkthroughs* anchored to the framework stack |
