# AuditProjection for Shared Audit Slice

When audit columns recur on multiple entity detail projections, extract a shared `AuditProjection` slice so each detail projection inherits audit fields rather than re-declaring them.

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
