# UI Component Field Naming

When declaring a Vaadin UI component as a class field, always include the
component type as a name suffix so field types are unambiguous.

```java
private final Span totalValueSpan;
private final Button saveButton;
private final TextField displayIdField;
private final Grid<ItemListItem> itemGrid;
```
