# Seed Data — Local Profile Only

When providing example data for local development, place seed scripts in `db/seed/`
as Flyway repeatable migrations rather than `data.sql` / `spring.sql.init` — Flyway
is the single schema and data authority, and Spring's SQL initializer has unreliable
ordering relative to Flyway (see `docs/patterns/cicd/flyway/ddl-strategy.md`).
Activate seed scripts only via `application-local.properties` — production startup
must never execute them.

```
src/main/resources/db/seed/
    R__local_seed.sql
```

The `R__` prefix marks the script as a repeatable migration: it runs after all
versioned migrations and re-runs whenever its content changes. Write every statement
to be idempotent (see `docs/patterns/cicd/flyway/idempotent-scripts.md`).

Activate by extending Flyway's locations in `application-local.properties` only:

```properties
# application-local.properties
spring.flyway.locations=classpath:db/migration,classpath:db/seed
```

The base `application.properties` omits `db/seed` from the locations, so seed data
is never applied in production or staging.

Seed data should provide:
- An initial admin user with a documented default password (changeable via the app)
- Enough example data to demonstrate all views without manual entry

## Related

- `docs/patterns/cicd/flyway/ddl-strategy.md` — why `data.sql` / `spring.sql.init` must not be used alongside Flyway.
- `docs/patterns/cicd/flyway/versioned-scripts.md` — versioned vs. repeatable migration conventions.
- `docs/patterns/cicd/flyway/idempotent-scripts.md` — `ON CONFLICT` / `IF NOT EXISTS` guards for repeatable scripts.
- `docs/patterns/cicd/deployment/spring-profiles.md` — `application-local.properties` pattern.