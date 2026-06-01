# Quick Filter

When adding a text search control to a grid view, implement it as a Quick Filter —
a single `TextField` with `LAZY` value-change mode and case-insensitive partial
match across the specified columns — so every grid view has consistent search
behavior.

```java
var provider = new CachingDataProvider<>(itemService::listAll);
var dataView = grid.setItems(provider::fetch, provider::count);

var quickFilter = new TextField();
quickFilter.setPlaceholder("Search…");
quickFilter.setPrefixComponent(VaadinIcon.SEARCH.create());
quickFilter.setClearButtonVisible(true);
quickFilter.setValueChangeMode(ValueChangeMode.LAZY);

quickFilter.addValueChangeListener(e -> {
    var term = e.getValue().strip().toLowerCase();
    if (term.isEmpty()) {
        provider.clearFilter(dataView);
    } else {
        provider.setFilter(item -> matchesQuickFilter(item, term), dataView);
    }
});

private boolean matchesQuickFilter(ItemListItem item, String term) {
    return item.getName().toLowerCase().contains(term)
        || item.getCode().toLowerCase().contains(term);
}
```

The columns that participate in the Quick Filter are defined per view in the feature
requirements. The Quick Filter may coexist with additional filter controls (status
toggle, role filter, date range) — it does not replace view-specific filters.

## Related

- `docs/patterns/ui/vaadin/grid-loading-state.md` — when and why to use `CachingDataProvider`.
- `docs/patterns/ui/vaadin/recipes/caching-data-provider.md` — full `CachingDataProvider` implementation including sort and reload.
