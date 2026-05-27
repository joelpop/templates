# Multi-Module Maven Structure

Pattern for organizing a Vaadin 24+ / Spring Boot 3+ / Spring Data JPA project into Maven
modules with compile-time layer separation. The module structure is the same across every
supported Vaadin and Spring Boot line — only dependency versions change. The `{app}`
placeholder represents the application name prefix (e.g., `myapp`). Replace
`{base_package}` with the root Java package.

## Module Inheritance Graph

```
spring-boot-starter-parent
└── {app} (root POM)
    ├── {app}-common
    ├── {app}-jpamodel
    ├── {app}-jpaclient
    ├── {app}-uimodel
    ├── {app}-service
    ├── {app}-jpaservice
    ├── {app}-ui
    └── {app}-app
```

## Module Responsibilities

| Module | Contents | Dependencies |
|--------|----------|--------------|
| `{app}-common` | Shared utility classes and enums used across modules | None |
| `{app}-jpamodel` | JPA entities (`@Entity`), code enums (`*Code`), interface projections (`*Projection`) | `{app}-common` |
| `{app}-jpaclient` | Spring Data JPA repositories; `JpaConfig` (`@EntityScan`, `@EnableJpaRepositories`) | `{app}-jpamodel` |
| `{app}-uimodel` | Plain POJOs for the UI layer (no persistence annotations); UI-layer enums | `{app}-common` |
| `{app}-service` | Service interfaces operating on UI model objects only; `ClientDetailsService` default methods provide `Instant` ↔ `LocalDateTime` conversion | `{app}-uimodel` |
| `{app}-jpaservice` | Service implementations using MapStruct to convert JPA projections ↔ UI models ↔ entities | `{app}-service`, `{app}-jpaclient` |
| `{app}-ui` | Vaadin views, layouts, and UI components | `{app}-service` only |
| `{app}-app` | Spring Boot entry point; assembles all modules; `application.properties`; frontend resources | `{app}-ui`, `{app}-jpaservice` (runtime only) |

## Compile-Time Layer Separation

The critical enforcement mechanism: `{app}-jpaservice` is declared as a **runtime-only**
dependency of `{app}-app`. This means the UI module can never reference JPA entities,
repositories, or service implementations at compile time — violations produce compiler errors
immediately, not runtime surprises.

```xml
<!-- In {app}-app/pom.xml -->
<dependency>
    <groupId>com.example</groupId>
    <artifactId>{app}-jpaservice</artifactId>
    <scope>runtime</scope>   <!-- enforces UI-to-persistence separation -->
</dependency>
```

What this enforces:
- UI code can call service interfaces (`{app}-service`)
- UI code can work with UI models (`{app}-uimodel`)
- UI code **cannot** import JPA entities, repositories, or mappers at compile time
- Service implementations **cannot** import Vaadin classes

## Application Entry Point and Route Scanning

`Application.java` must reside in the root package `{base_package}` so that Vaadin
automatically scans `{base_package}.ui.*` subpackages for `@Route`-annotated views.
Using `@EnableVaadin` is not required and can interfere with theme loading — avoid it.

```java
@SpringBootApplication
@StyleSheet(Lumo.STYLESHEET)
@StyleSheet(Lumo.UTILITY_STYLESHEET)
public class Application implements AppShellConfigurator {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

## vaadin.allowed-packages

`application.properties` in `{app}-app` must include:

```properties
vaadin.allowed-packages=com.vaadin,org.vaadin
```

Add the root package of any Vaadin add-on to this list when it is introduced. Without this,
add-on components may silently fail to render.

## JPA Configuration

`JpaConfig` lives in `{app}-jpaclient` and carries:

```java
@Configuration
@EntityScan(basePackages = "{base_package}.jpamodel")
@EnableJpaRepositories(basePackages = "{base_package}.jpaclient")
@EnableJpaAuditing
public class JpaConfig {

    @Bean
    public AuditorAware<UserEntity> auditorAware(EntityManager entityManager) {
        return () -> Optional.ofNullable(SecurityContextHolder.getContext().getAuthentication())
                .filter(Authentication::isAuthenticated)
                .map(Authentication::getPrincipal)
                .filter(AuditedPrincipal.class::isInstance)
                .map(AuditedPrincipal.class::cast)
                .map(AuditedPrincipal::getKey)
                .map(key -> entityManager.getReference(UserEntity.class, key));
    }
}
```

Each auth flow's principal implements `AuditedPrincipal` and carries the user's surrogate
key from its login-time validation step, so `AuditorAware` reads the key off the principal
without a per-write DB lookup. `EntityManager.getReference` returns a Hibernate proxy
holding just the key — JPA persists the FK from the proxy with no `SELECT` issued.

`spring.jpa.open-in-view=false` must be set in `application.properties`. See
`docs/patterns/persistence/spring-data-jpa/jpa-config.md` for why OSIV must be disabled.

## MapStruct Configuration

MapStruct is configured as an annotation processor in the **parent POM**:

```xml
<!-- In root {app}/pom.xml -->
<build>
    <pluginManagement>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <configuration>
                    <annotationProcessorPaths>
                        <path>
                            <groupId>org.mapstruct</groupId>
                            <artifactId>mapstruct-processor</artifactId>
                            <version>${mapstruct.version}</version>
                        </path>
                    </annotationProcessorPaths>
                </configuration>
            </plugin>
        </plugins>
    </pluginManagement>
</build>
```

All MapStruct mappers use Spring component model:

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface EmployeeMapper {
    EmployeeDetail toDetail(EmployeeDetailProjection projection);
    List<EmployeeListItem> toListItems(List<EmployeeListItemProjection> projections);

    // Updates existing entity from UI model — leaves unrelated fields untouched
    EmployeeEntity toEntity(EmployeeDetail detail, @MappingTarget EmployeeEntity entity);
}
```

Generated mapper sources appear under `target/generated-sources/annotations` and are
injected as Spring beans.

## Package Naming

Each module uses a consistent sub-package structure under `{base_package}`:

```
{base_package}                         — root; Application.java lives here
{base_package}.common.util             — shared utilities
{base_package}.jpamodel.entity         — JPA entity classes
{base_package}.jpamodel.code           — JPA entity enums (*Code suffix)
{base_package}.jpamodel.projection     — JPA interface projections (*Projection suffix)
{base_package}.jpaclient.repository    — Spring Data repositories
{base_package}.jpaclient.config        — JpaConfig (@EntityScan, @EnableJpaRepositories)
{base_package}.uimodel.data            — UI model POJOs (no suffix)
{base_package}.uimodel.type            — UI model enums (no suffix)
{base_package}.service                 — service interfaces
{base_package}.jpaservice              — service implementations
{base_package}.jpaservice.mapper       — MapStruct mappers
{base_package}.ui                      — shared UI components and MainLayout
{base_package}.ui.component            — shared reusable UI components
{base_package}.ui.view                 — views (each view in its own sub-package)
{base_package}.app.config              — application configuration
{base_package}.app.config.security     — security configuration
```

Package leaf names use the **singular** form — `…ui.view.admin.product`, not
`…ui.view.admin.products`. See `docs/patterns/conventions/singular-form.md`.

## Dependency Version Management

All dependency versions are declared in the parent POM's `<dependencyManagement>` section.
Child module POMs do not specify version numbers for managed dependencies.

All plugin versions are declared in the parent POM's `<pluginManagement>` section.
Child module POMs do not specify plugin versions inline.

This ensures version conflicts are resolved centrally and upgrade commits are minimal.
