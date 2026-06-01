# Caching at the Service Layer

When caching frequently read, rarely changing data, apply Spring `@Cacheable` at the
service layer and `@CacheEvict` on mutations so cache invalidation is co-located with
the data it guards.

```java
@Cacheable("departments")
@Transactional(readOnly = true)
public List<DepartmentListItem> listAllDepartments() { ... }

@CacheEvict(value = "departments", allEntries = true)
@Transactional
public void createDepartment(DepartmentDetail detail) { ... }
```

Enable caching with `@EnableCaching` on a configuration class. This approach is simpler
and more predictable than Hibernate's L2 cache, and works with projections and UI models
(not just entities).
