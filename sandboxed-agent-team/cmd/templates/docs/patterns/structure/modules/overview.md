# Module Structure Overview

When designing the Maven module structure for a new Vaadin 24+ / Spring Boot 3+ project,
follow the standard 8-module inheritance graph and responsibility table so compile-time
layer separation is enforced and each module has a single, predictable role.

The `{app}` placeholder represents the application name prefix (e.g., `myapp`). Replace
`{base_package}` with the root Java package. The module structure is identical across
every supported Vaadin and Spring Boot line — only dependency versions change.

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

| Module          | Contents                                                                                                                | Dependencies                        |
|-----------------|--------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| `{app}-common`  | Shared utility classes and enums used across modules                                                                     | None                                |
| `{app}-jpamodel`| JPA entities (`@Entity`), code enums (`*Code`), interface projections (`*Projection`)                                  | `{app}-common`                      |
| `{app}-jpaclient` | Spring Data JPA repositories; `JpaConfig` (`@EntityScan`, `@EnableJpaRepositories`)                                  | `{app}-jpamodel`                    |
| `{app}-uimodel` | Plain POJOs for the UI layer (no persistence annotations); UI-layer enums                                               | `{app}-common`                      |
| `{app}-service` | Service interfaces operating on UI model objects only; `ClientDetailsService` default methods provide `Instant` ↔ `LocalDateTime` conversion | `{app}-uimodel` |
| `{app}-jpaservice` | Service implementations using MapStruct to convert JPA projections ↔ UI models ↔ entities                           | `{app}-service`, `{app}-jpaclient`  |
| `{app}-ui`      | Vaadin views, layouts, and UI components                                                                                | `{app}-service` only                |
| `{app}-app`     | Spring Boot entry point; assembles all modules; `application.properties`; frontend resources                           | `{app}-ui`, `{app}-jpaservice` (runtime only) |
