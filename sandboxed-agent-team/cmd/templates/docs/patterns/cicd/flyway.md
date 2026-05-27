# Flyway Database Migrations

When managing schema changes, use Flyway versioned migration scripts applied
automatically on startup so the schema evolves with the application and rollback
procedures are documented alongside the change.

## Versioned, Immutable Scripts

```
src/main/resources/db/migration/
    V1__initial_schema.sql
    V2__add_employee_photo.sql
    V3__add_department_manager.sql

src/main/resources/db/seed/          (dev/test only)
    S1__default_admin_user.sql
    S2__example_departments.sql
```

Applied scripts are **never modified**. A checksum mismatch causes startup
failure.

## Idempotent Scripts

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

## DDL Strategy

```properties
# Never use create-drop in production — always validate
spring.jpa.hibernate.ddl-auto=validate
```

Use `validate` in all profiles — Hibernate validates the schema against entity
mappings at startup; a mismatch causes startup failure with a descriptive error.

## Seed Data — Dev/Test Only

Seed scripts live in `db/seed/` and are applied only when the active profile
includes `dev` or `test`. Production startup never executes them.

Seed data provides:
- An initial admin user with a documented initial password (changeable via the app)
- Enough example data to demonstrate all views without manual entry

## Rollback Documentation

Each migration script includes a comment block describing how to reverse the
change manually if needed:

```sql
-- V3__add_department_manager.sql
-- Rollback: ALTER TABLE departments DROP COLUMN manager_key;
--           DROP INDEX IF EXISTS idx_departments_manager;

ALTER TABLE departments ADD COLUMN manager_key BIGINT REFERENCES employees(employee_key);
CREATE INDEX idx_departments_manager ON departments (manager_key);
```

Full automated rollback (Flyway undo) is not required; manual procedure must be
documented for every migration.
