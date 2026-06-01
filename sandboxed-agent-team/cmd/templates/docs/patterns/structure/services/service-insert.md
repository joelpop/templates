# Insert Pattern — save() for New Entities Only

When implementing an insert, create a new entity instance, apply the UI model via
MapStruct, and call `repository.save()` only on the transient (null-key) entity.

```java
@Override
@Transactional
public EmployeeDetail create(EmployeeDetail detail) {
    var entity = new EmployeeEntity();
    mapper.toEntity(detail, entity);
    var saved = repository.save(entity);   // null key triggers persist()
    return mapper.toDetail(saved);
}
```
