# Quick Filter

When adding a text search control to a grid view, implement it as a Quick Filter —
a single `TextField` with `LAZY` value-change mode and case-insensitive partial
match across the specified columns — so every grid view has consistent search
behavior.

```java
var quickFilter = new TextField();
quickFilter.setPlaceholder("Search…");
quickFilter.setPrefixComponent(VaadinIcon.SEARCH.create());
quickFilter.setClearButtonVisible(true);
quickFilter.setValueChangeMode(ValueChangeMode.LAZY);

// Wire to ListDataProvider filter
quickFilter.addValueChangeListener(e -> {
    dataProvider.setFilter(item ->
        matchesQuickFilter(item, e.getValue())
    );
});

private boolean matchesQuickFilter(ItemListItem item, String query) {
    if (query == null || query.isBlank()) return true;
    var q = query.strip().toLowerCase();
    return item.getName().toLowerCase().contains(q)
        || item.getCode().toLowerCase().contains(q);
}
```

The columns that participate in the Quick Filter are defined per view in the feature
requirements. The Quick Filter may coexist with additional filter controls (status
toggle, role filter, date range) — it does not replace view-specific filters.
