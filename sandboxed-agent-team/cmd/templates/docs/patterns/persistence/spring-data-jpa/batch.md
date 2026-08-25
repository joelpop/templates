# Batch Operations

When performing high-volume inserts, updates, or deletes where each row has
distinct values, configure JDBC batching so Hibernate groups statements into
JDBC batches rather than issuing them one at a time.

JDBC batching is disabled by default (`batch_size=0`). Enable it in
`application.properties`:

```properties
spring.jpa.properties.hibernate.jdbc.batch_size=50
spring.jpa.properties.hibernate.order_inserts=true
spring.jpa.properties.hibernate.order_updates=true
```

- `batch_size` — maximum statements per JDBC batch
- `order_inserts` — groups INSERT statements by entity type to maximize batching
- `order_updates` — groups UPDATE statements by entity type to maximize batching

## Batch inserts

```java
// Avoid — save() in a loop issues N individual INSERTs regardless of batch config
for (EmployeeEntity e : employees) {
    employeeRepository.save(e);
}
```

```java
// Preferred — saveAll() with batch config groups INSERTs into batches of 50
employeeRepository.saveAll(employees);
```

## Batch updates

Load the entities, modify them within the transaction, and let Hibernate flush.
Dirty checking detects the changes and flushes on transaction commit. With
`order_updates` and `batch_size` configured the flush issues N / batch_size
JDBC batches; without them it issues N individual UPDATEs.

```java
List<EmployeeEntity> employees = employeeRepository.findByDepartmentKey(deptKey);
employees.forEach(e -> e.applyPromotion(newGrade));
// no save() needed — dirty checking flushes on transaction commit;
// order_updates + batch_size batch the UPDATEs; without them, N individual UPDATEs
```

## Batch deletes

```java
// Avoid — delete() in a loop issues N individual DELETEs regardless of batch config
for (EmployeeEntity e : employees) {
    employeeRepository.delete(e);
}
```

```java
// Preferred — deleteAll(entities) with batch config groups DELETEs into batches of 50
employeeRepository.deleteAll(employees);
```

## Related

- `persistence/spring-data-jpa/bulk.md` — single-statement operations when all affected rows receive the same change
- `persistence/spring-data-jpa/n-plus-one.md` — N+1 queries on the read side; the fetch counterpart to this pattern