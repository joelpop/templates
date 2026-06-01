# Grid Data Loading Pattern

When loading data for a grid, return the full dataset as a `List<T>` from the service
via `CachingDataProvider` in the view so the Grid shows its own loading indicator
on the first page request while sorting and filtering run in-memory with no further
backend calls.

```java
@Transactional(readOnly = true)
public List<EmployeeListItem> listAll() {
    return mapper.toListItems(repository.findAll(...));
}
```

```java
var provider = new CachingDataProvider<>(employeeService::listAll);
var dataView = grid.setItems(provider::fetch, provider::count);
```

Quick filter and column sorting are applied by `CachingDataProvider`'s in-memory cache,
not by database queries. No `Pageable`, `Page<T>`, or offset/limit parameters are used
for grid display.

## Related

- `docs/patterns/ui/vaadin/grid-loading-state.md` — when and why to use `CachingDataProvider`.
- `docs/patterns/ui/vaadin/recipes/caching-data-provider.md` — full implementation including filter, sort, and reload.
