# Grid Column Configuration

When configuring grid columns, declare each column explicitly with a header, sort
property, and appropriate width strategy (`setAutoWidth` or `setFlexGrow`) so
column widths are predictable and sortable columns are declared consistently.

```java
grid.addColumn(ItemListItem::getDisplayId)
    .setHeader("ID")
    .setSortable(true)
    .setAutoWidth(true);

grid.addColumn(ItemListItem::getName)
    .setHeader("Name")
    .setSortable(true)
    .setFlexGrow(1);

grid.addComponentColumn(item -> createStatusBadge(item.isActive()))
    .setHeader("Status")
    .setAutoWidth(true);
```

Use `setAutoWidth(true)` for columns where content width varies significantly.
Use `setFlexGrow(1)` on the primary text column to absorb remaining width.
