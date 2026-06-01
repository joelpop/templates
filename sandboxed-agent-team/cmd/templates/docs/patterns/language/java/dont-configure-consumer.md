# Don't Configure Your Consumer

When designing a helper, provider, or service that another component will use,
expose the interface and let the consumer wire itself — do not accept the
consumer as a parameter and configure it internally — so the helper stays
ignorant of its consumers and can be reused anywhere.

**Avoid** — provider accepts its consumer and configures it:

```java
public GridLazyDataView<T> setItems(Grid<T> grid) {
    return grid.setItems(this::fetch, this::count);
}

var dataView = provider.setItems(grid);
```

**Preferred** — consumer wires itself from the provider's interface:

```java
public Stream<T> fetch(Query<T, Void> query) { ... }
public int count(Query<T, Void> query) { ... }

var dataView = grid.setItems(provider::fetch, provider::count);
```

The provider has no dependency on `Grid` or any other consumer. Its interface
is a set of methods; any consumer that needs those methods can wire itself
without the provider knowing it exists.