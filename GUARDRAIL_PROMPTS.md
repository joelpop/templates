# Guardrail Prompts

## Foundational

Use the Vaadin, Playwright, and JavaDoc MCP Servers. Suggest
additional MCP Servers as they are relevant.

For the tech stack, use Maven, Java 25, Vaadin 25, Spring Boot 4, Spring
Data JPA, MapStruct (for mapping between Entities and UI model
objects), Vaadin TestBench (for unit tests), Playwright (for integration
tests), and git.

## Architecture and Organization

Use a multi-module Maven project with layered separation where code is organized into separate modules by technical concern.

## Modules

Use the following Maven modules, where `artifact` is the name of the project/application:

```
artifact (parent)
├── artifact-common (shared library)
├── artifact-jpamodel (entity classes, code enums, interface projections)
├── artifact-jpaclient (repositories; depends on: artifact-jpamodel)
├── artifact-uimodel (model classes, type enums)
├── artifact-service (service interfaces; depends on: artifact-uimodel)
├── artifact-jpaservice (JPA service implementation classes;
persistence configuration; depends on: artifact-service, artifact-jpamodel, artifact-jpaclient, artifact-uimodel + MapStruct)
└── artifact-ui (views, layouts, components, frontend; depends on: artifact-service, artifact-uimodel)
└── artifact-app (application executable, configuration; assembles the
executable; depends on: artifact-ui, artifact-jpaservice [runtime scope])
```

## Packages

Organize these packages into the modules matching the artifact name.

```
group.artifact.common
  .util
group.artifact.jpamodel
  .code
  .entity
  .projection
group.artifact.jpaclient
  .repository
group.artifact.jpaservice
  .config
  .mapper
group.artifact.uimodel
  .data
  .type
group.artifact.service
group.artifact.ui
  .component
  .layout
  .view
group.artifact.app
  .config
```

## Persistence

Turn off Spring Data JPA’s `open-session-in-view`.

Configure the `EntityScan` and `EnableJpaRepositories` in the
`artifact-jpaservice` module.

Use MapStruct to map between Entities/interface projections and UI model objects.

When performing fetch queries (selects), use interface projections and corresponding mappers to map to the UI model objects.

When saving a new (insert) UI model object, instantiate the corresponding entity, use the mapper to map the UI model object values into it, then persist the entity.

When saving edits (updates) of a UI model object, first fetch the entity (in order for it to be managed), have the mapper set appropriate values from the (potentially sparse) UI model object (without nulling values that were just sparse). The automatic committing of the transaction will persist the value changes.

## UX

### Responsive Layout

Implement appropriate smooth transitions and animations.

## Conventions

Specify the versions for all dependencies and plugins in the parent POM’s
`dependencyManagement` and `pluginManagement` sections, respectively.

Use var instead of explicit types whenever possible.

### Naming

Suffix Entities with “Entity” and place them in the `jpamodel.entity` package

Suffix entity enums with “Code” and place them in the `jpamodel.code` package.

Do not apply any prefix or suffix to UI model objects.

Place UI model objects in the `uimodel.data` package.

Place UI model enums in the `uimodel.type` package.

Suffix service interfaces with “Service”, such as “TaskService” for “Task” services.

Prefix service implementation classes with the technology of their
implementation, such as “JpaTaskService” rather than “TaskServiceImp”
for a “TaskService” interface implemented with JPA.

## Best Practices

Make a reasonable effort to follow best practices.

Keep the code DRY.

## Process

/code-review:code-review

Before each commit, do a git diff and pretend you're a senior dev
doing a code review and you HATE both naive and junior implementations. What would you criticize? What legitimate edge cases are missing?




## Features

### Authentication

#### Login View

Support Passkeys through
  - Apple TouchID
  - Apple ID
  - Google
  - Facebook
  - X

Support account registration (must confirm email address via link in email)

Support password reset

#### Account View

Support logout

Support password changing

Support user detail maintenance
  - First Name
  - Last Name
  - Image
  - Email Address (must confirm first with original address, then with new address)
  - Bio

### Authorization

Roles and Grants are managed by database
