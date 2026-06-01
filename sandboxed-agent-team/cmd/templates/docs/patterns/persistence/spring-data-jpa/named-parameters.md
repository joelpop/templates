# Named Parameters Only

When writing a `@Query`, use `:param` named-parameter syntax bound via `@Param` — string concatenation into query strings is a SQL injection vector and is forbidden.

```java
// Preferred — parameterized JPQL
@Query("SELECT e FROM EmployeeEntity e WHERE e.department.key = :deptKey AND e.active = true")
List<EmployeeListItemProjection> findActiveByDepartment(@Param("deptKey") Long deptKey);

// Never — string concatenation in queries
@Query("SELECT e FROM EmployeeEntity e WHERE e.name = '" + name + "'")  // NEVER
```

This rule applies to both JPQL and `nativeQuery = true` queries. Derived query
methods (`findByDepartmentKey(...)`, etc.) are inherently parameterized and need
no special handling.
