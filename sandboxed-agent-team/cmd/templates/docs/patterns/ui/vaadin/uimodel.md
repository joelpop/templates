# UI Model Conventions

Patterns for the `{app}-uimodel` module — the no-Vaadin-dependency layer that defines the data shapes the UI works with.

## Naming Conventions

UI model types are named for the UI context they serve — not for the entity or data source.

| Type          | Convention                                      | Example                                |
|---------------|-------------------------------------------------|----------------------------------------|
| Data POJO     | No suffix — named for its UI context            | `EmployeeListItem`, `EmployeeDetail`   |
| Enum          | No suffix — named for the concept it represents | `EmploymentStatus`, `PhoneType`        |
| Picker record | Suffix `PickerItem`                             | `EmployeePickerItem`                   |

Avoid generic suffixes like `Summary`, `Info`, or `Data` — `EmployeeListItem` is
self-explaining; `EmployeeSummary` is not.

### Pairing with the Implementation Layer

UI model POJOs sit at the service boundary — what they pair with on the implementation
side depends on the backing technology:

| Implementation   | Counterpart                                                        | Reference                                              |
|------------------|--------------------------------------------------------------------|--------------------------------------------------------|
| Spring Data JPA  | JPA interface projection (`EmployeeListItemProjection`)            | `docs/patterns/persistence/spring-data-jpa/naming.md` |
| REST client      | Response DTO or deserialized payload mapped by the service impl    | `docs/patterns/architecture/services.md`               |
| Internal/derived | Assembled in the service impl from other sources — no counterpart  | `docs/patterns/architecture/services.md`               |

The UI model name is chosen for the UI context it serves — not for the counterpart it
comes from. The name stays stable even if the backing technology changes.

## Enum Display — `HasCaption`

UI model enums that appear in selection components (`ComboBox`, `Select`, etc.) implement
`HasCaption`, a custom interface in `{app}-uimodel`:

```java
public interface HasCaption {
    String getCaption();
}
```

Any enum whose constants need a display label implements it:

```java
public enum ActiveFilterOption implements HasCaption {
    ALL("All"),
    ACTIVE("Active"),
    INACTIVE("Inactive");

    private final String caption;

    ActiveFilterOption(String caption) { this.caption = caption; }

    @Override public String getCaption() { return caption; }
}
```

This enables a uniform `setItemLabelGenerator` call across all display enums:

```java
statusSelect.setItemLabelGenerator(ActiveFilterOption::getCaption);
```

Without `HasCaption`, each enum has its own getter name (`getLabel()`, `getDisplayName()`,
`getName()`, …), forcing a per-call decision. With it, any `HasCaption` enum drops into
`setItemLabelGenerator` the same way.

Enums that are only stored or compared (not displayed in selection components) do not need
to implement `HasCaption`.

### Records also implement `HasCaption`

UI model records used as picker items in `ComboBox` or `Select` implement `HasCaption`
the same way — the caption method encodes how the record's fields combine into a display
string, keeping that logic in one place rather than repeated at every call site:

```java
public record EmployeePickerItem(long key, String firstName, String lastName)
        implements HasCaption {
    @Override
    public String getCaption() {
        return firstName + " " + lastName;
    }
}
```

```java
employeeSelect.setItems(employeeService.listPickerItems());
employeeSelect.setItemLabelGenerator(EmployeePickerItem::getCaption);
```

The picker projection that feeds this record (`EmployeePickerItemProjection`) fetches
only `key` + the display fields. See `docs/patterns/architecture/persistence.md` —
"Picker projections" for the JPA side.

## UI Model Capability Interfaces

UI model data records implement small single-method interfaces from `{app}-uimodel` to
advertise structural capabilities — the shapes that generic UI components bind against:

```java
public interface HasActive {
    boolean active();
}

public interface HasRole {
    UserRole role();
}
```

A list-item record declares what it carries:

```java
public record UserListItem(
        long id,
        String username,
        UserRole role,
        boolean active) implements HasActive, HasRole { }
```

UI components bind to the interface, not the concrete record type. This decouples
grid and filter logic from any particular record shape:

```java
// Mute inactive rows in any grid whose item type implements HasActive
browser.muteRowsWhen(item -> !item.active());

// FilterOption.matches() takes a method reference to the interface accessor —
// one static factory works for every HasRole item type
public static <T extends HasRole> SerializableBiPredicate<RoleFilterOption, T> matches() {
    return FilterOption.matches(HasRole::role);
}
```

### Naming convention

Interface names match the record component they expose, prefixed with `Has`:
`HasActive` exposes `active()`, `HasRole` exposes `role()`. The accessor name
matches the record component name — no `get` prefix (records use accessor syntax,
not getter syntax).

### When to define a new interface

Define a `Has*` interface when at least two record types share the same structural
property **and** at least one UI component binds to that property generically. A
property unique to one record type does not need an interface.

### Relationship to `HasCaption`

`HasCaption` is for *enum types* used in selection components — it normalises the
display label. Capability interfaces are for *data records* used in grids and
filters — they normalise the structural shape that generic components key off.
Both live in `{app}-uimodel`; neither imports Vaadin.

See `docs/patterns/recipes/item-browser.md` for the primary consumer: `muteRowsWhen`
and `FilterOption.matches()` both bind through capability interfaces.
