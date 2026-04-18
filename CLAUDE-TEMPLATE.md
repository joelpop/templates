# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **Artifact**, a Vaadin 25 + Spring Boot 4 application using Java 25. It follows a **multi-module Maven architecture** with **layered separation** where code is organized into separate Maven modules by technical concern. The application uses **MapStruct** to map between JPA interface projections, UI models, and JPA entities, ensuring complete decoupling between persistence and presentation layers.

## Tech Stack

- Maven
- Java 25
- Vaadin 25
- Spring Boot 4
- Spring Data JPA
- MapStruct (for mapping from JPA interface projections to UI model objects and from UI model objects to JPA entities)
- Vaadin TestBench (for unit tests)
- Playwright (for integration tests)
- Git

## MCP Servers

This project was preconfigured with Model Context Protocol (MCP) servers in `.mcp.json` to provide specialized tooling assistance to Claude Code. These servers were selected based on the project's toolchain and the broad, cross-platform availability of the tooling needed to run the MCP servers (http and node/npm/npx).

### Sample MPC Server File

```json
{
  "mcpServers": {
    "fetch": {
      "command": "npm",
      "args": ["exec", "--silent", "--", "fetch-mcp"]
    },
    "java": {
      "type": "http",
      "url": "https://www.javadocs.dev/mcp"
    },
    "playwright": {
      "command": "npx",
      "args": ["-y", "@playwright/mcp"]
    },
    "spring-docs": {
      "command": "npx",
      "args": ["-y", "@enokdev/springdocs-mcp@latest"]
    },
    "vaadin": {
      "type": "http",
      "url": "https://mcp.vaadin.com/docs"
    }
  }
}
```

**Note:** Node.js must be installed and `npm`/`npx` must be available on your PATH for the command-based MCP servers to work. The HTTP-based servers (java, vaadin) require no local installation.

### Fetch MCP Server (`mcp__fetch__`)

Provides URL fetching capabilities for retrieving web content.

**Available Tools:**
- `fetch_url` - Fetch content from a URL (HTML, text, or images). Returns Markdown by default.
- `fetch_youtube_transcript` - Fetch transcript from a YouTube video URL.

### IDE MCP Server (`mcp__ide__`)

Provides IDE integration for diagnostics.

**Available Tools:**
- `getDiagnostics` - Get diagnostic information, optionally filtered by URI.

### Java MCP Server (`mcp__java__`)

Provides Java/Maven utilities for dependency management and Javadoc lookup.

**Available Tools:**
- `get_latest_version` - Get the latest version of a Maven artifact from Maven Central.
- `get_javadoc_content_list` - List contents of a Javadoc JAR for a specific artifact.
- `get_javadoc_symbol_contents` - Get Javadoc content for a specific symbol/class.
- `symbol_to_artifact` - Find the Maven groupId and artifactId for a given class/package name.

### Playwright MCP Server (`mcp__playwright__`)

Provides browser automation for UI testing and interaction. Useful for testing the Vaadin application.

**Available Tools:**
- `browser_navigate` - Navigate to a URL.
- `browser_snapshot` - Capture accessibility snapshot (preferred over screenshot for actions).
- `browser_take_screenshot` - Take a screenshot of the current page.
- `browser_click` - Click on an element.
- `browser_type` - Type text into an element.
- `browser_fill_form` - Fill multiple form fields at once.
- `browser_press_key` - Press a keyboard key.
- `browser_select_option` - Select an option in a dropdown.
- `browser_hover` - Hover over an element.
- `browser_drag` - Drag and drop between elements.
- `browser_wait_for` - Wait for text, text disappearance, or a timeout.
- `browser_tabs` - List, create, close, or select browser tabs.
- `browser_evaluate` - Execute JavaScript in the browser.
- `browser_console_messages` - Get console messages.
- `browser_network_requests` - Get network requests.
- `browser_close` - Close the browser page.
- `browser_install` - Install the browser if not available.

### Spring Docs MCP Server (`mcp__spring-docs__`)

Provides access to Spring Boot documentation, guides, tutorials, and best practices.

**Available Tools:**
- `search_spring_docs` - Search Spring Boot documentation by keywords. Filter by `docType`: "guides", "reference", "api", or "all".
- `search_spring_projects` - Search all Spring projects on spring.io/projects.
- `get_spring_project` - Get details of a specific Spring project (e.g., "spring-boot", "spring-security").
- `get_all_spring_guides` - List all Spring guides, optionally filtered by category.
- `get_spring_guide` - Get content of a specific guide with configurable detail level.
- `get_spring_reference` - Get a section of the Spring Boot reference documentation.
- `search_spring_concepts` - Search Spring Boot concepts by category (core, web, data, security, testing, production).
- `search_spring_ecosystem` - Search across the entire Spring ecosystem.
- `get_spring_tutorial` - Get step-by-step tutorials for specific features.
- `compare_spring_versions` - Compare features between Spring Boot versions.
- `get_spring_best_practices` - Get best practices by category (architecture, performance, security, testing, configuration, deployment).
- `diagnose_spring_issues` - Diagnose common Spring Boot issues and get solutions.

### Vaadin MCP Server (`mcp__vaadin__`)

Provides access to Vaadin documentation, component APIs, and best practices. **Always call `get_vaadin_primer` first** before working with Vaadin to get current information about modern Vaadin development.

**Available Tools:**
- `get_vaadin_primer` - Returns comprehensive primer about modern Vaadin development (2025+). Essential to avoid outdated assumptions.
- `search_vaadin_docs` - Search Vaadin documentation with hybrid semantic + keyword search. Specify `ui_language` as "java", "react", or "common".
- `get_full_document` - Retrieve complete documentation pages by file path.
- `get_vaadin_version` - Get the latest stable Vaadin version.
- `get_components_by_version` - List all components available in a specific Vaadin version.
- `get_component_java_api` - Get Java API documentation for a component.
- `get_component_react_api` - Get React API documentation for a component.
- `get_component_web_component_api` - Get Web Component/TypeScript API documentation.
- `get_component_styling` - Get styling/theming documentation for a component.

**Note**: Claude Code must be restarted after modifying `.mcp.json` to load MCP server changes.

## Package Structure

All Java packages follow the base package pattern: `group.artifact.*`

The Java code is organized into the following packages (object `User` and its derivatives are used below hypothetically):

- `group.artifact` - Base package; `ArtifactBase` marker interface for type-safe component scanning by Spring
- `group.artifact.common.util` - shared libraries
- `group.artifact.jpamodel.code` - JPA entity enums (e.g., `UserTypeCode`)
- `group.artifact.jpamodel.entity` - JPA entity classes (e.g., `UserEntity`)
- `group.artifact.jpamodel.projection` - JPA interface projections (e.g., `UserNameProjection`, `UserSummaryProjection`, `UserDetailProjection`)
- `group.artifact.jpaclient.config` - Spring Data configuration (e.g., `JpaConfig`)
- `group.artifact.jpaclient.repository` - Spring Data repositories
- `group.artifact.jpaservice` - Service implementations (e.g., `JpaUserService`)
- `group.artifact.jpaservice.mapper` - MapStruct mappers
- `group.artifact.uimodel.data` - UI model POJOs (e.g., `UserName`, `UserSummary`, `UserDetail`)
- `group.artifact.uimodel.type` - UI model enums (e.g., `UserType`)
- `group.artifact.service` - Service interfaces
- `group.artifact.service.util` - Service utilities (e.g., `PageRequest`)
- `group.artifact.ui` - Shared UI components (e.g., `MainLayout`, `ViewToolbar`)
- `group.artifact.ui.component` - Shared UI components
- `group.artifact.ui.view` - Views and their related classes (each view within its own package)
- `group.artifact.app` - Main Application class
- `group.artifact.app.config` - Application configuration
- `group.artifact.app.config.security` - Application security configuration

## Architecture and Organization

Multi-module Maven project with layered separation where code is organized into separate modules by technical concern. Java packages are organized into Maven modules with corresponding names (e.g., package `group.artifact.ui` would be located in module `artifact-ui`).

### Module Structure

The project is organized as a Maven multi-module build with the following modules:

**Module Inheritance Graph**

```
spring-boot-starter-parent
└── artifact
    ├── artifact-common
    ├── artifact-jpamodel
    ├── artifact-jpaclient
    ├── artifact-uimodel
    ├── artifact-service
    ├── artifact-jpaservice
    ├── artifact-ui
    └── artifact-app
```

### Module Brief Descriptions & Dependencies

- **artifact** - Parent of all the other project modules
- **artifact-common** - Shared utility classes and enums
- **artifact-jpamodel** - JPA entities with `@Entity` annotations, code enums, interface projections
    - **Naming convention**: Entities are suffixed with `Entity` (e.g., `UserEntity`)
- **artifact-jpaclient** - Spring Data JPA repository interfaces (depends on: `artifact-jpamodel`)
    - Contains Spring configuration (e.g., `JpaConfig`)
- **artifact-uimodel** - Plain POJOs for UI layer, type enums (no persistence knowledge)
    - **Naming convention**: UI models have no suffix (e.g., `User`)
- **artifact-service** - Service interfaces that work with UI models only (depends on: `artifact-uimodel`)
- **artifact-jpaservice** - Service implementations using MapStruct to convert from JPA interface projections to UI models and from UI models to JPA entities (depends on: artifact-service & artifact-jpaclient)
    - **Important**: Must NOT depend on UI libraries (e.g., Vaadin) to maintain layer separation
- **artifact-ui** - Vaadin UI components (views, layouts, and components) (depends on: `artifact-service`)
    - Contains only UI layer code (no application infrastructure)
    - Pure UI library module with no direct persistence dependencies
- **artifact-app** - Spring Boot application entry point (main executable JAR) (depends on: `artifact-ui` & `artifact-jpaservice`)
    - Contains `Application.java` (Spring Boot main class with `@EnableVaadin` for route scanning in UI packages)
    - Contains `ArtifactBase` marker interface for type-safe component scanning
    - Contains security configuration
    - Contains `application.properties`
    - Contains frontend resources
    - Depends on `artifact-jpaservice` at runtime only (enforces layer separation)
    - Assembles and runs the complete application

**Important**: `artifact-jpaservice` must be a **runtime-only** dependency of `artifact-app`. This enforces compile-time separation--the UI layer may not reference JPA entities, repositories, or service implementations during compilation.

### Compile-Time Layer Separation

The `artifact-jpaservice` dependency in `artifact-app` has **`<scope>runtime</scope>`**. Additionally, `artifact-ui` has no direct dependency on `artifact-jpaservice` at all. This enforces:
- ✅ UI code can call service interfaces
- ✅ UI code can work with UI models
- ❌ UI code **cannot** import or reference JPA entities at compile time
- ❌ UI code **cannot** import or reference repositories at compile time
- ❌ UI code **cannot** import or reference MapStruct mappers at compile time

Additionally, service implementations must not depend on UI libraries:
- ❌ Service code **cannot** import or reference Vaadin classes

This guarantees proper layering and prevents accidental coupling between UI and persistence layers. The application layer (`artifact-app`) is responsible for wiring everything together at runtime.

### Layer Separation Details

- **Presentation Layer** (`artifact-ui`, `artifact-uimodel`)
    - Vaadin UI components module
    - Contains Vaadin views, components, and layouts
    - Vaadin views work exclusively with UI model objects (no suffix)
    - Pure library module with no direct persistence dependencies
    - Cannot reference JPA entities, repositories, or service implementations at compile time
    - Has no knowledge of Persistence layer implementation

- **Persistence Layer** (`artifact-jpamodel` + `artifact-jpaclient`)
    - Contains JPA entities with full persistence annotations
    - Repositories for data access
    - Contains `JpaConfig` class with `@EntityScan` and `@EnableJpaRepositories` annotations
    - Has no knowledge of the presentation layer

- **Service Interface Layer** (`artifact-service`)
    - Defines business operations using UI models
    - No dependencies on JPA or persistence
    - Allows for multiple implementations

- **Service Implementation Layer** (`artifact-jpaservice`)
    - Implements service interfaces
    - **Uses MapStruct for automatic mapping** from JPA interface projections to UI model object and from UI model objects to JPA entities
    - MapStruct generates implementation classes at compile time (see `target/generated-sources/annotations`)
    - Example mapper: `UserMapper` converts between `UserDetailProjection` (JPA) and `UserDetail` (UI) and between `UserDetail` (UI) and `UserEntity` (JPA)
    - **Must NOT depend on UI libraries** (e.g., Vaadin)

- **Application Layer** (`artifact-app`)
    - Spring Boot application entry point
    - Contains `Application.java` with `@SpringBootApplication` and `@EnableVaadin("group.artifact.ui")`
    - The `@EnableVaadin` annotation is required because UI classes are in a sibling package (`group.artifact.ui`), not a descendant of the application package (`group.artifact.app`)
    - Uses `scanBasePackageClasses` with marker interface `ArtifactBase` for type-safe component scanning
    - Contains `application.properties`
    - Contains frontend resources
    - **Runtime-only dependency on `artifact-jpaservice`** - enforces compile-time separation
    - Assembles all modules into executable application

## Vaadin Configuration

### Route Scanning

The `Application.java` class uses the `@EnableVaadin` annotation because Vaadin views are in a sibling package (`group.artifact.ui`) rather than a descendant package of where the Application class resides (`group.artifact.app`):

```java
@SpringBootApplication(scanBasePackageClasses = ArtifactBase.class)
// @EnableVaadin allows the UI classes to not have to be in a descendant package of the application class
@EnableVaadin("group.artifact.ui")
@StyleSheet(Lumo.STYLESHEET)
public class Application implements AppShellConfigurator {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

The `@EnableVaadin("group.artifact.ui")` annotation explicitly tells Vaadin to scan the `group.artifact.ui` package for `@Route` annotated views and components.

The `scanBasePackageClasses = ArtifactBase.class` uses a marker interface for type-safe component scanning instead of string-based package names.

### Allowed Packages

To efficiently scan for Vaadin components, including those that are core, VFC add-ons, and the application's, `application.properties` contains:
```
vaadin.allowed-packages=com.vaadin,org.vaadin,group.artifact
```
When adding new Vaadin add-ons, update this property to include their package prefixes.

## Persistence and Service Configuration

### Entity & Repository Registration

JPA entities and repositories are configured in `group.artifact.jpaclient.config.JpaConfig`:
- `@EntityScan(basePackages = "group.artifact.jpamodel")` - scans for JPA entities
- `@EnableJpaRepositories(basePackages = "group.artifact.jpaclient")` - scans for Spring Data repositories

**Note**: In Spring Boot 4, `@EntityScan` is located in `org.springframework.boot.persistence.autoconfigure` package.

- Turn off Spring Data JPA's `open-session-in-view`

### MapStruct Integration

MapStruct is used to map from JPA interface projections to UI model objects and from UI model objects to JPA entities.

MapStruct is configured in the parent POM:
- Annotation processor configured in `maven-compiler-plugin`
- Mappers are Spring components (`componentModel = MappingConstants.ComponentModel.SPRING`)

Example partial mapper interface for `UserEntity`-related objects:
```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserMapper {
    // from interface projectins to UI model
    UserDetail toDetail(UserDetailProjection detailProjection);
    List<UserSummary> toSummaryList(List<UserSummaryProjection> summaryProjection);

    // from UI model to entity (overwrites only common properties, leaves others untouched)
    UserEntity toEntity(UserDetail detail, @MappingTarget UserEntity entity);
}
```

### Fetch Queries (Selects)

Use interface projections and corresponding mappers to map to the UI model objects.

### Insert Operations

When saving a **new** UI model object:
1. Instantiate the corresponding entity.
2. Use the mapper to copy the UI model object values into it.
3. Explicitly persist the entity.

### Update Operations

When saving **edits** made to a UI model object for an existing entity:
1. First, fetch the entity (in order for it to be managed in the persistence context) via its unique identifier.
2. Have the service update (via MapStruct) corresponding values from the UI model object into the entity while leaving the entity's other values untouched.
3. The automatic committing of the transaction will persist the value changes.

## Database Configuration

- For development, `spring.jpa.hibernate.ddl-auto=create-drop` is used in `application.properties` (recreates schema on each startup)
- For production, use `spring.jpa.hibernate.ddl-auto=validate` to validate the schema without modifications

## Development Commands

### Running the Application
```bash
# Start in development mode from artifact-app module
cd artifact-app && ../mvnw spring-boot:run

# Or from root directory
./mvnw spring-boot:run -pl artifact-app -am
```

The application will start on port 8080 by default and automatically launch a browser window.

### Building
```bash
# Build all modules
./mvnw clean package

# Build and skip tests
./mvnw clean package -DskipTests

# Install to local Maven repository
./mvnw clean install
```

### Testing
```bash
# Run all tests across all modules
./mvnw test

# Run tests for a specific module
./mvnw test -pl artifact-jpaservice

# Run a specific test class
./mvnw test -pl artifact-jpaservice -Dtest=JpaUserServiceTest

# Run a specific test method
./mvnw test -pl artifact-jpaservice -Dtest=JpaUserServiceTest#users_are_stored_in_the_database_with_the_current_timestamp
```

## Security Configuration

### Application Security

Spring Security is to be configured using Vaadin 25-specific facilities and techniques as this is a Single Page Application (SPA)--most conventional Spring Security configuration advice will not apply.

### View Access Control

Use standard Jakarta Security annotations for view permissions.

## UI/UX

### Views

Application initially has only two views, "Login" and "About." "About" should be within the "main layout" and show the version of the application, dependency, app server, database, and browser.

All normal application views live within MainLayout

### Vaadin Error Views

Application displays exception views when there are unhandled errors:

| Num | Type                           | Purpose                                              |
|-----|--------------------------------|------------------------------------------------------|
| 400 | `HasErrorParameter`            | Interface for views that handle specific error types |
| 403 | `AccessDeniedExceptionHandler` | Handles Spring Security access denied                |
| 404 | `RouteNotFoundError`           | Base for 404 handling                                |
| 500 | `ErrorHandler`                 | Global handler for uncaught exceptions               |

#### Sensitive Information

Error views must not expose:
- Stack traces to users
- Internal system paths
- Database error details
- Security-sensitive information
- 403 errors must appear as 404 errors to the user
- 400 errors are only shown if there is no 403 error possible

Detailed error information is logged server-side only.

### Responsive Layout

Use Vaadin's responsive layout features.

Implement appropriate smooth transitions and animations.

### Composition

Views and custom composition components should extend `Composite<T>` rather than extending layout classes directly. This provides better encapsulation and cleaner APIs. Implement appropriate Vaadin "Has" interfaces to expose typical required functionality, such as for sizing and styling.

### Navigation

For routing of views in the main navigation menu, use the `@Menu` annotation.

### Theming and Styling

Use the "Aura" theme.

Use component theme variants where available to achieve desired styling.

For simple styling adjustments to components, prefer using `addClassNames()` with `LumoUtility` class names over using CSS.

## Testing

### Unit Tests

Create unit tests for each non-UI public method.

Create UI unit tests for each UI feature using TestBench UI Unit Testing.

### Integration Tests

Create integration tests for each feature using Playwright.

### Patterns

Tests use:
- `@SpringBootTest` with custom `TestConfiguration`
- `@ComponentScan(basePackages = "group.artifact")`
- `@EnableJpaRepositories(basePackages = "group.artifact.jpaclient")`
- `@Transactional` to rollback changes
- AssertJ for assertions
- H2 in-memory database

## Naming Conventions

- **UI Model Objects**: No prefix or suffix, place in `uimodel.data` package. Names aligned with UI feature.
- **UI Model Enums**: No prefix or suffix, place in `uimodel.type` package
- **Service Interfaces**: Suffix with "Service" (e.g., "UserService"), place in `service` package
- **Entities**: Suffix with "Entity" and place in `jpamodel.entity` package
- **Entity Enums**: Suffix with "Code" and place in `jpamodel.code` package
- **Interface Projections**: Suffix with "Projection" and place in `jpamodel.projection` package
- **Service Implementations**: Prefix with the technology (e.g., "JpaUserService" for JPA implementation), place in `jpaservice` package

## Maven Conventions

- Specify versions for all dependencies in the parent POM's `dependencyManagement` section
- Specify versions for all plugins in the parent POM's `pluginManagement` section

## Java Conventions

- Use `var` instead of explicit types whenever possible

## Development Guidelines

Do things the Vaadin 25 way.

When adding new features across modules:

1. **Create UI model** in `artifact-uimodel`
    - **Name without suffix**: `UserSummary`
    - Plain POJO without JPA annotations
    - Specific to the feature using it (i.e., just key/name for combo boxes, summary for dashboards, abbreviated for listings, detail for modifying, all for admin)

2. **Create JPA entity** (if it does not already exist) in `artifact-jpamodel`
    - **Base name with `Entity` suffix**: `UserEntity`
    - Add `@Entity`, `@Table`, column mappings

3. **Create JPA Interface Projection** in `artifact-jpamodel`
    - **Name with `Projection` suffix**: `UserSummaryProjection`
    - Add projection as implemented by Entity, `public class UserEntity implements UserSummaryProjection`

4. **Create repository** (if it does not already exist) in `artifact-jpaclient`
    - Extend `JpaRepository<EntityName, ID>`
    - Example: `JpaRepository<UserEntity, Long>`

5. **Create MapStruct mapper** in `artifact-jpaservice/mapper`
    - Interface with `@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)`
    - Declare conversion methods between JPA entity/projections and UI model
    - Examples:
        - `UserSummary toModel(UserSummaryProjection projection)`
        - `UserEntity toEntity(UserSummary summary)`

6. **Define service interface** in `artifact-service`
    - Methods work with UI models, not entities
    - Example: `List<UserSummary> fetchSummaryByUserType(UserType userType)`
    - **Naming convention**: Sufffix base model with `Service` (e.g., `UserService`)

7. **Implement service** in `artifact-jpaservice`
    - Inject repository and mapper
    - Use `@Service` and `@Transactional`
    - **Naming convention**: Prefix with implementation type (e.g., `JpaUserService`)
    - **Do NOT add UI library dependencies** (e.g., Vaadin)

8. **Create UI view** in `artifact-ui`
    - Each view and its related classes should be in its own sub-package under `ui.view`
    - Hypothetical example: `UserListView` would be in `group.artifact.ui.view.userlist`
    - Shared components should be in `group.artifact.ui.component` (e.g., `FilterBar`)
    - Constructor-inject service interface
    - Work with UI models only (no suffix)
    - Add `@Route`, `@PageTitle`, `@Menu` annotations

9. **Write tests** in `artifact-jpaservice/src/test/java`
    - Use `@SpringBootTest` with `TestConfiguration`
    - Test service layer with actual repository and MapStruct mapping

## Important Notes

### frontend folder

The only module that needs and uses a Vaadin Flow `src/main/frontend` folder is the `artifact-app` folder.

### MapStruct Generated Code

- MapStruct implementations are generated at compile time
- Find generated classes in `target/generated-sources/annotations`
- If mapper changes don't take effect, run `./mvnw clean compile`

### Module Dependencies

- Modules must be built in dependency order (handled automatically by Maven reactor)
- When working on a single module, use `-am` (also-make) flag to build dependencies:
  ```bash
  ./mvnw test -pl artifact-jpaservice -am
  ```

### Spring Boot 4 Notes

- `@EntityScan` is in `org.springframework.boot.persistence.autoconfigure` package.

### Vaadin Development Mode

- The `vaadin-dev` dependency is required in `artifact-app` for development mode support (hot reload, debug features)
- This dependency should be marked as `<optional>true</optional>` so it's not included in production builds
- Without this dependency, you'll get a runtime error: `vaadin-dev-server artifact is not found`

### Managed Sources

- `.gitignore` should exclude node modules and target directories, files specific to developer enviroment, such as IDE, JRebel, MCPs, etc., and OS hidden files.

## Best Practices

- Make a reasonable effort to follow best practices
- Keep the code DRY (Don't Repeat Yourself)
- Use type-safe alternatives to string-based specifiers where applicable
- Service implementations must not depend on UI libraries to maintain layer separation
