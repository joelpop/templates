# AuditProjection for Shared Audit Slice

When projections for audited entities need to expose audit fields, extract a
shared `AuditProjection` interface and mix it into each detail projection so
audit accessors are declared once rather than re-declared per entity.

```java
public interface AuditProjection {
    Instant getCreatedAt();
    Instant getUpdatedAt();
    UserAuditProjection getCreatedBy();
    UserAuditProjection getUpdatedBy();
}
```

Any detail projection for an audited entity mixes it in:

```java
public interface EmployeeDetailProjection
        extends EmployeeListItemProjection, AuditProjection {
    String getEmail();
    Long getDepartmentKey();
}
```

`UserAuditProjection` is the canonical nested projection:

```java
public interface UserAuditProjection {
    String getFirstName();
    String getLastName();
}
```

Spring Data resolves `getCreatedBy()` → `UserAuditProjection` via a JOIN on the
audit FK column. The audit-by references can be `null` for system-originated
writes; null-check before rendering.

## Placement

`AuditProjection` and `UserAuditProjection` live in the persistence module
(`{app}-data`) alongside the other shared projection interfaces, not in any
single entity's package.

## Related

- `persistence/spring-data-jpa/entity-hierarchy/audited-entity.md` — the entity base class that carries the audit columns these projections expose
- `persistence/spring-data-jpa/projection-inheritance.md` — how Spring Data resolves projection interface hierarchies
- `persistence/spring-data-jpa/projection-naming.md` — naming conventions for projection interfaces
