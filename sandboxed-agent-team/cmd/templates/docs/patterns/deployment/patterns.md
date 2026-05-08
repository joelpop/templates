# Deployment Patterns

Fat JAR packaging, Docker, Spring profiles, database migrations, logging, and health
check patterns for Vaadin 24+ applications running on Spring Boot 3+. The deployment
surface is identical across every supported Vaadin and Spring Boot line — only the base
JRE image and Spring Boot Maven plugin versions vary by the line you target.

## Executable Fat JAR

The application artifact is a self-contained executable JAR produced by the Spring Boot
Maven plugin:

```bash
mvn clean package
java -jar {app}-app/target/{app}-app.jar
```

No additional classpath configuration is required. All dependencies are bundled inside
the JAR. This is the deployment unit for all environments.

## Docker

A `Dockerfile` in the project root builds a runnable image from the fat JAR. Pick the
Temurin tag that matches your project's Java target — Spring Boot 3 supports Java 17+,
Spring Boot 4 raises the minimum. The example below uses Java 25; substitute `21` or `17`
for older Spring Boot lines:

```dockerfile
FROM eclipse-temurin:25-jre-alpine     # or :21-jre-alpine / :17-jre-alpine
WORKDIR /app
COPY {app}-app/target/{app}-app.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

```bash
docker build -t {app}:latest .
docker run -p 8080:8080 \
  -e DB_URL=jdbc:postgresql://host/db \
  -e DB_USER=appuser \
  -e DB_PASSWORD=secret \
  {app}:latest
```

Secrets are injected via environment variables — never baked into the image.

## Spring Profiles

Three profiles cover the standard environment lifecycle:

| Profile | Database | Seed data | Log level | Cookie flags |
|---------|----------|-----------|-----------|--------------|
| `dev` | H2 in-memory | Yes | DEBUG | Off |
| `staging` | PostgreSQL | No | INFO | Off |
| `prod` | PostgreSQL | No | WARN | On (Secure, SameSite) |

```properties
# application-dev.properties
spring.datasource.url=jdbc:h2:mem:devdb;MODE=PostgreSQL
spring.jpa.hibernate.ddl-auto=validate
spring.profiles.include=seed
logging.level.root=DEBUG
logging.level.org.hibernate.SQL=DEBUG
logging.level.org.hibernate.orm.jdbc.bind=TRACE
spring.jpa.properties.hibernate.format_sql=true
spring.jpa.properties.hibernate.generate_statistics=true

# application-prod.properties
spring.datasource.url=${DB_URL}
spring.datasource.username=${DB_USER}
spring.datasource.password=${DB_PASSWORD}
spring.jpa.hibernate.ddl-auto=validate
logging.level.root=WARN
server.servlet.session.cookie.secure=true
server.servlet.session.cookie.same-site=strict
server.servlet.session.cookie.http-only=true
```

Switch profiles: `--spring.profiles.active=prod`. No code changes required.

## Database Migration — Flyway

Schema migrations are managed by Flyway, applied automatically on application startup
before requests are accepted.

### Versioned, Immutable Scripts

```
src/main/resources/db/migration/
    V1__initial_schema.sql
    V2__add_employee_photo.sql
    V3__add_department_manager.sql

src/main/resources/db/seed/          (dev/test only)
    S1__default_admin_user.sql
    S2__example_departments.sql
```

Applied scripts are **never modified**. A checksum mismatch causes startup failure.

### Idempotent Scripts

Write migrations defensively:

```sql
CREATE TABLE IF NOT EXISTS employees (
    employee_key  BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    active        BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_employees_email ON employees (email);
```

Re-running a migration on an already-migrated database produces no errors.

### DDL Strategy by Profile

```properties
# Never use create-drop in production — always validate
spring.jpa.hibernate.ddl-auto=validate
```

- `validate` in all profiles after the initial migration script exists
- Hibernate validates the schema against entity mappings at startup
- A schema mismatch causes startup failure with a descriptive error

### Seed Data — Dev/Test Only

Seed scripts live in `db/seed/` and are applied only when the active profile includes
`dev` or `test`. Production startup never executes seed scripts.

Seed data provides:
- An initial admin user with a documented initial password (which that admin can change later via the application's own password-change flow)
- Enough example data to demonstrate all views without manual data entry

### Rollback Documentation

Each migration script includes a comment block describing how to reverse the change
manually if needed:

```sql
-- V3__add_department_manager.sql
-- Rollback: ALTER TABLE departments DROP COLUMN manager_key;
--           DROP INDEX IF EXISTS idx_departments_manager;

ALTER TABLE departments ADD COLUMN manager_key BIGINT REFERENCES employees(employee_key);
CREATE INDEX idx_departments_manager ON departments (manager_key);
```

Full automated rollback (Flyway undo) is not required, but the manual procedure must be
documented for every migration.

## Test Data with @Transactional Rollback

Integration tests manage their own data using `@Transactional` rollback. The database
is in a known state at the start of each test. Tests use H2 in PostgreSQL compatibility
mode:

```properties
# application-test.properties
spring.datasource.url=jdbc:h2:mem:testdb;MODE=PostgreSQL;DB_CLOSE_DELAY=-1
spring.jpa.hibernate.ddl-auto=validate
```

See `docs/patterns/testing/patterns.md` for test data patterns.

## Health Check

Spring Boot Actuator exposes a health check endpoint:

```properties
management.endpoints.web.exposure.include=health
management.endpoint.health.show-details=never   # hide details in production
```

`/actuator/health` returns HTTP 200 when the application and database are healthy.
All other Actuator endpoints are restricted or disabled in production.

## Logging

Application code logs through **SLF4J** — the facade defined by `org.slf4j.Logger` and
`org.slf4j.LoggerFactory`, declared in classes via the `@Slf4j` annotation (see
`docs/patterns/conventions/lombok.md` → "for Logging"). The implementation that actually
formats and writes log records is **Logback**, pulled in transitively by
`spring-boot-starter-logging` (which every other `spring-boot-starter-*` depends on, so
nothing else needs to be declared). The same starter includes bridges for
`java.util.logging` and Log4j, so any third-party library using those APIs routes
through the same pipeline.

Practical consequences of the facade/implementation split:

- **Code imports** `org.slf4j.Logger`. Application code never imports from
  `ch.qos.logback.*` or `java.util.logging.*`.
- **Configuration uses Logback's native formats** — `logback-spring.xml` for structured
  setup, `application.properties` `logging.level.*` for level tuning. The XML below and
  the `LogstashEncoder` reference are Logback-specific.
- **Swapping implementations is possible** — a project could replace Logback with Log4j2
  or another SLF4J binding without changing a single line of application code. This kit
  does not do that; the separation is noted so the relationship is explicit.

### Structured Logging

Use structured logging (JSON or key-value pairs) compatible with log aggregation tools
(CloudWatch, Datadog, etc.). Each field is a distinct key in the log entry — not embedded
in the message string.

```xml
<!-- logback-spring.xml — JSON output in non-dev profiles -->
<springProfile name="staging,prod">
    <appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
        <encoder class="net.logstash.logback.encoder.LogstashEncoder"/>
    </appender>
    <root level="WARN">
        <appender-ref ref="JSON"/>
    </root>
</springProfile>

<springProfile name="dev">
    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%d{HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>
    <root level="DEBUG">
        <appender-ref ref="CONSOLE"/>
    </root>
</springProfile>
```

### Startup Logging

The application logs resolved configuration at INFO level on startup:

```java
@EventListener(ApplicationReadyEvent.class)
public void logStartupConfig() {
    log.info("Application started",
        kv("profile", activeProfile),
        kv("dbUrl", maskedDbUrl),
        kv("vaadinVersion", VaadinVersion.getFullVersion()));
}
```

Log enough to diagnose configuration issues without exposing secrets. Mask or omit
credential values from startup logs.

## HikariCP Connection Pool

Configure the connection pool for production load:

```properties
spring.datasource.hikari.maximum-pool-size=20
spring.datasource.hikari.minimum-idle=5
spring.datasource.hikari.connection-timeout=30000
spring.datasource.hikari.keepalive-time=600000
spring.datasource.hikari.max-lifetime=1800000
```

`keepalive-time` and `max-lifetime` detect and replace stale connections automatically.
A database connection that drops and recovers does not require an application restart.

## Java Runtime

The application runs on any Java-spec-compliant JVM — no vendor-specific APIs. The
recommended runtime is Eclipse Temurin (OpenJDK). Do not use Oracle-only or IBM J9-only APIs.
