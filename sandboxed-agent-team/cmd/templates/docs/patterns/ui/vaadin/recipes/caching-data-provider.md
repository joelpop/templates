# Recipe: `CachingDataProvider` — Grid with Lazy-Loaded In-Memory Cache

Follow this recipe to produce a `CachingDataProvider<T>` that loads the full
result once on the first page request, caches it in memory, and serves all
subsequent pages, filter changes, and sort changes from that cache with no
further backend calls.

## What this produces

- A `CachingDataProvider<T>` whose `fetch` and `count` methods are passed to
  `grid.setItems(provider::fetch, provider::count)`, returning a
  `GridLazyDataView<T>` — the Grid's own built-in loading indicator fires
  automatically on the first page request.
- `setFilter` / `clearFilter` for Quick Filter: applies a predicate to the
  cached list and triggers a Grid refresh with no backend call.
- Optional `sortResolver` at construction: maps sort-property strings from
  column headers to in-memory `Comparator<T>` instances — no separate sort
  listener needed.
- `refresh(dataView)` to clear the cache and force a full reload, e.g. after a
  save or delete.
- `cache` is `transient`: on session failover the cache is cleared and reloads
  transparently on next Grid access.

## Dependencies

- Vaadin 24+ (`Grid`, `GridLazyDataView`, `Query`, `QuerySortOrder`,
  `SortDirection`, `SerializableFunction`, `SerializablePredicate`,
  `SerializableSupplier`).

## Step 1 — The `CachingDataProvider<T>` class

The count callback is called first by Vaadin; whichever callback fires first
triggers `getOrLoad()` and populates the cache. The other callback reuses it.
Sort is applied to the filtered stream on every fetch — it is not cached
separately because sort order can change without a full reload.

```java
package {base_package}.ui.data;

import com.vaadin.flow.component.grid.dataview.GridLazyDataView;
import com.vaadin.flow.data.provider.Query;
import com.vaadin.flow.data.provider.QuerySortOrder;
import com.vaadin.flow.data.provider.SortDirection;
import com.vaadin.flow.function.SerializableFunction;
import com.vaadin.flow.function.SerializablePredicate;
import com.vaadin.flow.function.SerializableSupplier;

import java.io.Serializable;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Stream;

public class CachingDataProvider<T> implements Serializable {

    private final SerializableSupplier<List<T>> loader;
    private final SerializableFunction<QuerySortOrder, Comparator<T>> sortResolver;
    private transient List<T> cache;
    private SerializablePredicate<T> filter;

    public CachingDataProvider(SerializableSupplier<List<T>> loader) {
        this(loader, order -> (a, b) -> 0);
    }

    public CachingDataProvider(
            SerializableSupplier<List<T>> loader,
            SerializableFunction<QuerySortOrder, Comparator<T>> sortResolver) {
        this.loader = loader;
        this.sortResolver = sortResolver;
    }

    public Stream<T> fetch(Query<T, Void> query) {
        return sorted(query).skip(query.getOffset()).limit(query.getLimit());
    }

    public int count(Query<T, Void> query) {
        return (int) filtered().count();
    }

    public void setFilter(SerializablePredicate<T> predicate, GridLazyDataView<T> dataView) {
        filter = predicate;
        dataView.refreshAll();
    }

    public void clearFilter(GridLazyDataView<T> dataView) {
        filter = null;
        dataView.refreshAll();
    }

    public void refresh(GridLazyDataView<T> dataView) {
        cache = null;
        dataView.refreshAll();
    }

    private Stream<T> sorted(Query<T, Void> query) {
        var stream = filtered();
        var comparator = query.getSortOrders().stream()
            .map(order -> {
                var c = sortResolver.apply(order);
                return order.getDirection() == SortDirection.DESCENDING ? c.reversed() : c;
            })
            .reduce(Comparator::thenComparing)
            .orElse(null);
        if (comparator != null) {
            stream = stream.sorted(comparator);
        }
        return stream;
    }

    private Stream<T> filtered() {
        var stream = getOrLoad().stream();
        return filter != null ? stream.filter(filter) : stream;
    }

    private List<T> getOrLoad() {
        if (cache == null) {
            cache = loader.get();
        }
        return cache;
    }
}
```

## Step 2 — Basic use (no filter, no sort)

```java
var provider = new CachingDataProvider<>(itemService::listAll);
var dataView = grid.setItems(provider::fetch, provider::count);
```

## Step 3 — Add Quick Filter

```java
var provider = new CachingDataProvider<>(itemService::listAll);
var dataView = grid.setItems(provider::fetch, provider::count);

quickFilterField.addValueChangeListener(event -> {
    var term = event.getValue().trim().toLowerCase();
    if (term.isEmpty()) {
        provider.clearFilter(dataView);
    } else {
        provider.setFilter(
            item -> item.getName().toLowerCase().contains(term),
            dataView
        );
    }
});
```

## Step 4 — Add column sort

Pass a `sortResolver` at construction. Each entry maps the sort-property string
(set via `column.setSortProperty("key")`) to a `Comparator<T>`. The provider
handles direction reversal and multi-column chaining automatically.

```java
var provider = new CachingDataProvider<ItemListProjection>(
    itemService::listAll,
    order -> switch (order.getSorted()) {
        case "name"   -> Comparator.comparing(ItemListProjection::getName);
        case "status" -> Comparator.comparing(p -> p.getStatus().name());
        default       -> (a, b) -> 0;
    }
);
var dataView = grid.setItems(provider::fetch, provider::count);
```

Column declarations must set `setSortProperty` with the matching key:

```java
grid.addColumn(ItemListProjection::getName)
    .setHeader("Name")
    .setSortProperty("name");

grid.addColumn(p -> p.getStatus().getCaption())
    .setHeader("Status")
    .setSortProperty("status");
```

## Step 5 — Reload after mutation

Call `refresh(dataView)` after a save or delete to clear the cache and let the
Grid re-fetch from the service on its next render cycle:

```java
private void onSave(ItemUiModel item) {
    itemService.save(item);
    provider.refresh(dataView);
}
```

## Decisions this recipe imposes

- **Full query on first page request.** The Grid's own loading indicator fires
  while the service call runs; the rest of the view is available immediately.
- **One query per `refresh()`, zero for filter and sort changes.** Filter and
  sort operate on the in-memory cache without touching the backend.
- **Sort applied on every fetch, not cached.** The sorted stream is rebuilt per
  page request because sort order can change without a cache reset. For typical
  data sizes this is negligible.
- **`cache` is `transient`.** On session failover or deserialization, `cache`
  clears to `null`. The active `filter` predicate survives; the next Grid access
  reloads the full list and re-applies it transparently.
- **No `setSort` / `clearSort` methods.** Sort order is owned by the Grid and
  passed through `query.getSortOrders()`; the provider reads it directly rather
  than requiring a separate sort-change listener.

## What to verify

- On first navigation to the view, the Grid shows its built-in loading indicator
  while the first page loads, then renders.
- Typing in the Quick Filter field filters the Grid without a backend round-trip.
- Clicking a sortable column header re-sorts without a backend round-trip.
- Clearing the Quick Filter restores the full unfiltered list.
- After a save, `refresh(dataView)` reloads from the backend and the updated
  row appears in the Grid.

## Related

- `docs/patterns/ui/vaadin/grid-loading-state.md` — when to use
  `CachingDataProvider`.
- `docs/patterns/structure/services/service-grid-loading.md` — service-layer
  convention: return `List<T>` from the service; `CachingDataProvider` is the
  UI-layer complement.
