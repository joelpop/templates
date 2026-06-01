# Rollback Documentation in Migration Scripts

When writing a Flyway migration, include a comment block describing how to reverse
the change manually — Flyway Community does not support automated undo, and a
procedure documented at write-time is far more reliable than reconstructing it under
pressure during an incident.

```sql
-- V3__add_department_manager.sql
-- Rollback: ALTER TABLE departments DROP COLUMN manager_key;
--           DROP INDEX IF EXISTS idx_departments_manager;

ALTER TABLE departments ADD COLUMN manager_key BIGINT REFERENCES employees(employee_key);
CREATE INDEX idx_departments_manager ON departments (manager_key);
```
