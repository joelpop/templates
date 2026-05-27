# Deployment Packaging

When packaging the application for deployment, produce a Spring Boot fat JAR,
run it inside a Docker image keyed by environment-specific Spring profile, and
expose `/actuator/health` for orchestrator probes so every environment deploys
from one reproducible artifact.

## Executable Fat JAR

The application artifact is a self-contained executable JAR produced by the
Spring Boot Maven plugin:

```bash
mvn clean package
java -jar {app}-app/target/{app}-app.jar
```

All dependencies are bundled inside the JAR. This is the deployment unit for all
environments.

## Docker

A `Dockerfile` in the project root builds a runnable image from the fat JAR.
Pick the Temurin tag that matches your Java target — Spring Boot 3 supports Java
17+. The example uses Java 25; substitute `21` or `17` for older lines:

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

| Profile   | Database        | Seed data | Log level | Cookie flags |
|-----------|-----------------|-----------|-----------|--------------|
| `dev`     | H2 in-memory    | Yes       | DEBUG     | Off          |
| `staging` | PostgreSQL      | No        | INFO      | Off          |
| `prod`    | PostgreSQL      | No        | WARN      | On (Secure, SameSite) |

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

## Health Check

Spring Boot Actuator exposes a health check endpoint:

```properties
management.endpoints.web.exposure.include=health
management.endpoint.health.show-details=never   # hide details in production
```

`/actuator/health` returns HTTP 200 when the application and database are
healthy. All other Actuator endpoints are restricted or disabled in production.

## HikariCP Connection Pool

Configure the connection pool for production load:

```properties
spring.datasource.hikari.maximum-pool-size=20
spring.datasource.hikari.minimum-idle=5
spring.datasource.hikari.connection-timeout=30000
spring.datasource.hikari.keepalive-time=600000
spring.datasource.hikari.max-lifetime=1800000
```

`keepalive-time` and `max-lifetime` detect and replace stale connections; a
connection that drops and recovers does not require an application restart.

## Java Runtime

The application runs on any Java-spec-compliant JVM. Recommended runtime:
Eclipse Temurin (OpenJDK). Do not use Oracle-only or IBM J9-only APIs.
