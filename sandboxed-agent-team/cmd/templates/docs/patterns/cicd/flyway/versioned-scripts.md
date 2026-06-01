# Versioned, Immutable Migration Scripts

When managing database schema changes, use Flyway versioned scripts in
`db/migration/` — Flyway is the schema authority; applied scripts are never
modified, and a checksum mismatch causes startup failure.

```
src/main/resources/db/migration/
    V1__initial_schema.sql
    V2__add_employee_photo.sql
    V3__add_department_manager.sql

src/main/resources/db/seed/          (local profile only)
    R__local_seed.sql
```

Seed files use the `R__` (repeatable migration) prefix and live in a separate
`db/seed/` directory — see `docs/patterns/cicd/flyway/seed-data.md` for activation
and idempotency rules.

## Related

- `docs/patterns/cicd/flyway/ddl-strategy.md` — set `ddl-auto=validate` so Hibernate defers to Flyway rather than managing the schema itself.
- `docs/patterns/cicd/flyway/seed-data.md` — seed data activation and idempotency rules.