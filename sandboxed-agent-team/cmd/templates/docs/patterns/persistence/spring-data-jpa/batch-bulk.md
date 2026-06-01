# Batch and Bulk Operations

When performing high-volume inserts, updates, or deletes, bypass entity loading using batch insert configuration or `@Modifying` bulk queries.

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
