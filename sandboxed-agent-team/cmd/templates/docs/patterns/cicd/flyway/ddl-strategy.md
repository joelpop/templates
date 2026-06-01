# DDL Strategy — Flyway Owns the Schema

When configuring Hibernate DDL behavior, set `ddl-auto=validate` and let Flyway
own all schema changes — so Hibernate catches entity/schema mismatches at startup
and no DDL is silently applied.

```properties
spring.jpa.hibernate.ddl-auto=validate
```

Avoid `update`: it silently alters the schema to match entities, bypassing Flyway
and accumulating drift that is invisible in version control. Avoid `create-drop`:
it destroys and recreates the schema on every restart, including staging if the
wrong profile is activated.

Avoid `spring.sql.init` (`data.sql` / `schema.sql`) when Flyway is present —
ordering between Spring's SQL initializer and Flyway is unreliable and splits
schema authority across two mechanisms. Route all initialization through Flyway
instead: versioned migrations for schema, repeatable migrations for seed data.

## Related

- `docs/patterns/cicd/flyway/versioned-scripts.md` — Flyway migration naming and location conventions.
- `docs/patterns/cicd/flyway/seed-data.md` — seed data via repeatable migrations instead of `data.sql`.