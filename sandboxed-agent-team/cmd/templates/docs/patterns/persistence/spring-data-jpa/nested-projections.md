# Nested Projections

When a projection getter returns another projection interface, Spring Data generates a JOIN and populates the nested object in the same query — no separate SELECT per row.

```java
public interface UserAuditProjection {
    String getFirstName();
    String getLastName();
}

// Spring Data resolves getCreatedBy() via a JOIN — no N+1 per row
public interface EmployeeDetailProjection extends EmployeeListItemProjection, AuditProjection {
    String getEmail();
    Long getDepartmentKey();
}
```

The same technique applies to any `@ManyToOne` relationship where the caller
only needs a narrow slice of the related entity.
