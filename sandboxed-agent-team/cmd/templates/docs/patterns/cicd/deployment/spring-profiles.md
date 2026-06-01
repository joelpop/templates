# Spring Profiles per Environment

When configuring the application for different environments, put production-safe
defaults in `application.properties` and override only what differs in
`application-local.properties` — so a missing profile in production fails secure
rather than exposing development settings.

```properties
# application.properties — production-safe base
spring.datasource.url=${DB_URL}
spring.datasource.username=${DB_USER}
spring.datasource.password=${DB_PASSWORD}
spring.jpa.hibernate.ddl-auto=validate
spring.jpa.open-in-view=false
logging.level.root=WARN
server.servlet.session.cookie.secure=true
server.servlet.session.cookie.same-site=strict
server.servlet.session.cookie.http-only=true
server.servlet.session.timeout=30m
```

```properties
# application-local.properties — local development overrides
spring.datasource.url=jdbc:h2:mem:devdb;MODE=PostgreSQL;DATABASE_TO_LOWER=TRUE;CASE_INSENSITIVE_IDENTIFIERS=TRUE
spring.datasource.username=sa
spring.datasource.password=
logging.level.root=DEBUG
logging.level.org.hibernate.SQL=DEBUG
logging.level.org.hibernate.orm.jdbc.bind=TRACE
spring.jpa.properties.hibernate.format_sql=true
server.servlet.session.cookie.secure=false
spring.flyway.locations=classpath:db/migration,classpath:db/seed
```

Activate the local profile at startup:

```bash
SPRING_PROFILES_ACTIVE=local mvn spring-boot:run
# or
java -Dspring.profiles.active=local -jar {app}-app/target/{app}-app.jar
```

Set `SPRING_PROFILES_ACTIVE=local` in the IDE run configuration for everyday
development.

For staging, add `application-staging.properties` with staging-specific overrides
(database URL, log level) on top of the production base.

## Related

- `docs/patterns/security/spring/session-cookie-flags.md` — rationale for each cookie security flag.
- `docs/patterns/cicd/flyway/seed-data.md` — seed data activation via Flyway repeatable migrations.