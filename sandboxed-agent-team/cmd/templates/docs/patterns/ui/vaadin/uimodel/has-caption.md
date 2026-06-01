# HasCaption — Enum and Record Display Labels

When a UI model enum or picker record appears in a selection component
(`ComboBox`, `Select`, etc.), implement `HasCaption` so every display enum
and picker drops into `setItemLabelGenerator` the same way.

```java
public interface HasCaption {
    String getCaption();
}
```

Without `HasCaption`, each enum has its own getter name (`getLabel()`,
`getDisplayName()`, `getName()`, …), forcing a per-call decision. With it,
any `HasCaption` type plugs into `setItemLabelGenerator` uniformly.

## Enums

Any enum whose constants need a display label implements `HasCaption`:

```java
@AllArgsConstructor
@Getter
public enum ActiveFilterOption implements HasCaption {
    ALL("Any"),
    ACTIVE("Active"),
    INACTIVE("Inactive");

    private final String caption;
}
```

```java
statusSelect.setItemLabelGenerator(ActiveFilterOption::getCaption);
```

Enums that are only stored or compared — not displayed in selection components
— do not need to implement `HasCaption`.

## Picker Records

UI model records used as picker items implement `HasCaption` and carry a
`HasNames` component. `getCaption()` delegates to `fullName()`, keeping the
display-string logic in one place:

```java
public record EmployeePickerItem(long key, HasNames hasNames)
        implements HasCaption {
    @Override
    public String getCaption() { return hasNames.fullName(); }
}
```

```java
employeeSelect.setItems(employeeService.listPickerItems());
employeeSelect.setItemLabelGenerator(EmployeePickerItem::getCaption);
```

The picker projection that feeds this record (`EmployeePickerItemProjection`)
fetches only `key` plus the display fields.
