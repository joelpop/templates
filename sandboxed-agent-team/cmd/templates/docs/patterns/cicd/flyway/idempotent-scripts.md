# Idempotent Repeatable Migrations

When writing a Flyway repeatable migration (`R__`), use `IF NOT EXISTS` / `ON CONFLICT`
guards so the script can re-run safely whenever its content changes — unlike versioned
migrations, repeatable migrations execute again each time their checksum changes.

```sql
-- R__local_seed.sql
INSERT INTO roles (name)
VALUES ('ROLE_ADMIN'), ('ROLE_USER')
ON CONFLICT (name) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_employees_email ON employees (email);
```

Versioned migrations (`V__`) are applied exactly once; do not add `IF NOT EXISTS`
guards to them as a safety net — Flyway's checksum enforcement is the protection,
and masking errors with guards defeats it.

## Related

- `docs/patterns/cicd/flyway/versioned-scripts.md` — versioned vs. repeatable migration conventions.
- `docs/patterns/cicd/flyway/seed-data.md` — where repeatable seed migrations live and how they are activated.