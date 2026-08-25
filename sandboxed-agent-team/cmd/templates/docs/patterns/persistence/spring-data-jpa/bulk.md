# Bulk Operations

When all affected rows receive the same change, use a single-statement bulk
operation rather than loading entities so the database executes one SQL
statement regardless of how many rows are affected.

## Bulk inserts

Use `INSERT INTO ... SELECT ...` to copy or derive rows from existing data in
one statement without loading entities into the application:

```java
@Modifying
@Query("INSERT INTO EmployeeArchiveEntity (id, name, department) "
     + "SELECT e.id, e.name, e.department FROM EmployeeEntity e "
     + "WHERE e.status = :status")
int archiveByStatus(@Param("status") EmployeeStatus status);
```

## Bulk updates

```java
// Avoid — loads every entity into memory, then issues N individual UPDATEs
List<EmployeeEntity> employees = employeeRepository.findByDepartmentKey(deptKey);
employees.forEach(e -> e.setStatus(status));
// flushed on transaction close — one UPDATE per entity
```

```java
// Preferred — single SQL UPDATE, no entity loading
@Modifying(clearAutomatically = true)
@Query("UPDATE EmployeeEntity e SET e.status = :status WHERE e.department.key = :deptKey")
int updateStatusByDepartment(@Param("status") EmployeeStatus status,
                             @Param("deptKey") Long deptKey);
```

## Bulk deletes

```java
// Avoid — deleteAll(entities) loads every entity first, then issues N individual DELETEs
employeeRepository.deleteAll(employees);
```

```java
// Preferred — single DELETE WHERE key IN (...)
employeeRepository.deleteAllByIdInBatch(keys);
```

For condition-based deletes without a key list:

```java
@Modifying(clearAutomatically = true)
@Query("DELETE FROM EmployeeEntity e WHERE e.department.key = :deptKey")
int deleteByDepartment(@Param("deptKey") Long deptKey);
```

## Cache state after bulk operations

Bulk operations bypass the JPA first-level cache. Stale cached entities may
be returned if other JPA operations follow in the same transaction.
`clearAutomatically = true` on `@Modifying` clears the cache immediately
after the bulk statement so subsequent reads are fresh.

## Transactional requirement

`@Modifying` repository methods must be called from a `@Transactional` service
method. Calling them outside a transaction throws at runtime.

## Related

- `persistence/spring-data-jpa/batch.md` — batched operations when each row has distinct values
- `persistence/spring-data-jpa/n-plus-one.md` — N+1 queries on the read side; the fetch counterpart to this pattern
- `persistence/spring-data-jpa/osiv-disabled.md` — with OSIV disabled, bulk operations must complete within the service transaction before the view layer is reached