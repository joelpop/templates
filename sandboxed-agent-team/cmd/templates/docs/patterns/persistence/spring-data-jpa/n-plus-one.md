# N+1 Prevention

When a query needs associated collections loaded, use `@EntityGraph` or `JOIN FETCH` to load them in a single query — not N additional lazy-load queries.

- `@EntityGraph` — declarative, specified on the repository method:

```java
@EntityGraph(attributePaths = {"phones"})
List<EmployeeEntity> findAllWithPhones();
```

- `JOIN FETCH` in `@Query` — explicit, when the query is already in JPQL:

```java
@Query("SELECT e FROM EmployeeEntity e JOIN FETCH e.department WHERE e.status = :status")
List<EmployeeEntity> findActiveWithDepartment(@Param("status") String status);
```

Prefer `@EntityGraph` as the first choice — it keeps fetch semantics separate
from the query string. Use `JOIN FETCH` when a `@Query` is already required.
