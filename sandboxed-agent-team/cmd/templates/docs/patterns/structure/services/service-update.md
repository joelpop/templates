# Update Pattern — Dirty Checking

When implementing an update, load the managed entity by key, apply changes via
MapStruct's `@MappingTarget`, and let the transaction flush automatically — do not call
`save()` on a managed entity.

```java
@Override
@Transactional
public EmployeeDetail update(EmployeeDetail detail) {
    var entity = repository.findById(detail.getKey())
        .orElseThrow(() -> new EntityNotFoundException(detail.getKey()));
    mapper.toEntity(detail, entity);
    // transaction flush performs the UPDATE automatically — no save() call needed
    return mapper.toDetail(entity);
}
```

The method returns the refreshed UI model so the caller has the latest state without a
re-query. The entity itself implements the detail projection interface, so
`mapper.toDetail(entity)` is valid without a separate projection query.

Do not call `save()` on a managed (already-loaded) entity — it is redundant.
Do not call `save()` on a detached entity — it triggers a full-column overwrite.
