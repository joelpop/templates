# Projection Inheritance — Detail Extends List-Item

When defining a detail projection, extend the list-item projection rather than re-declaring its fields so the detail query returns everything the grid already shows plus the detail-only fields.

```java
public interface EmployeeListItemProjection {
    Long getKey();
    String getFirstName();
    String getLastName();
    String getStatus();
}

// Inherits all list-item getters; adds only what the edit form needs beyond the grid
public interface EmployeeDetailProjection extends EmployeeListItemProjection {
    String getEmail();
    Long getDepartmentKey();
}
```

The entity implements both — satisfying `EmployeeDetailProjection` implies
satisfying `EmployeeListItemProjection` by inheritance:

```java
public class EmployeeEntity extends AuditedEntity<Long>
        implements EmployeeDetailProjection { /* ... */ }
```
