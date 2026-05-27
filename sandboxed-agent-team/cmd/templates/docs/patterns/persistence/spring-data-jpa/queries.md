# Query and Mutation Patterns

When writing Spring Data JPA queries and mutations, use named parameters, load
managed entities for updates rather than calling `save()` on detached ones, and
prefer `@EntityGraph` for eager-loading joins so the persistence layer does not
silently issue extra queries or overwrite concurrent changes.

## Insert vs. Update

**Insert (new entity):** Call `repository.save(newEntity)` on a transient entity
(null key). The null key triggers `persist()`.

**Update (managed entity):** Load the entity by key within a transaction, mutate
its fields, and let the transaction flush. Hibernate's dirty checking detects the
changes and issues a targeted UPDATE. Do not call `save()` on a managed entity —
it is redundant.

**Never call `save()` on a detached entity** (one loaded in a previous
transaction). This triggers `merge()`, which issues a SELECT + full-column UPDATE
that can silently overwrite concurrent changes. Instead, load a fresh managed
copy by key and apply changes selectively.

```java
// Insert
@Transactional
public EmployeeDetail create(EmployeeDetail detail) {
    var entity = new EmployeeEntity();
    mapper.toEntity(detail, entity);
    return mapper.toDetail(repository.save(entity));
}

// Update — dirty checking, no save() call
@Transactional
public void update(long key, EmployeeDetail detail) {
    var entity = repository.findById(key).orElseThrow(() -> new EntityNotFoundException(key));
    mapper.toEntity(detail, entity);
    // transaction flush performs the UPDATE automatically
}
```

## @DynamicUpdate for Wide Tables

For entities with many columns, large columns (TEXT, BLOB), or high write
volume, annotate with `@DynamicUpdate` to limit UPDATE statements to only
changed columns:

```java
@Entity
@Table(name = "employees")
@DynamicUpdate
public class EmployeeEntity extends BaseEntity<Long> { ... }
```

Without `@DynamicUpdate`, Hibernate generates a static full-column UPDATE that
overwrites all columns on every flush. `@DynamicUpdate` adds a small per-flush
overhead to compute the dirty field diff — use it only when profiling confirms
the full-column writes are a meaningful cost.

## N+1 Prevention

The N+1 problem: loading N entities then firing N additional queries to load a
lazy collection. Solutions:

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

## Named Parameters Only — SQL Injection Prevention

Every value placeholder in a `@Query` must use `:param` named-parameter syntax
bound via `@Param`. String concatenation into a query string is a SQL injection
vector and is forbidden — including in `nativeQuery = true` queries.

```java
// Correct — parameterized JPQL
@Query("SELECT e FROM EmployeeEntity e WHERE e.department.key = :deptKey AND e.active = true")
List<EmployeeListItemProjection> findActiveByDepartment(@Param("deptKey") Long deptKey);

// Never — string concatenation in queries
@Query("SELECT e FROM EmployeeEntity e WHERE e.name = '" + name + "'")  // NEVER
```

Derived query methods (`findByDepartmentKey(...)`, etc.) are inherently
parameterized and need no special handling. The rule applies whenever you reach
for `@Query`.

## Batch and Bulk Operations

For high-volume inserts, updates, or deletes, bypass entity loading:

```java
// Batch inserts — configure batch_size in application.properties
// spring.jpa.properties.hibernate.jdbc.batch_size=50
// spring.jpa.properties.hibernate.order_inserts=true
employeeRepository.saveAll(employees);   // grouped into batches

// Bulk update — single SQL UPDATE, no entity loading
@Modifying
@Query("UPDATE EmployeeEntity e SET e.status = :status WHERE e.department.key = :deptKey")
int updateStatusByDepartment(@Param("status") String status, @Param("deptKey") Long deptKey);

// Bulk delete by IDs — single DELETE WHERE key IN (...)
employeeRepository.deleteAllByIdInBatch(keys);
// AVOID: deleteAll() loads every entity first, then issues N individual DELETEs
```
