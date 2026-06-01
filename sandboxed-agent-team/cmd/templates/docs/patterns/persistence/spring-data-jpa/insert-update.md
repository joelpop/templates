# Insert and Update Patterns

When persisting an entity, call `save()` only on transient entities (null key); for updates, load a managed entity within a transaction and mutate its fields — never call `save()` on a detached entity.

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
