# Picker Projections

When a `ComboBox` or `Select` needs to load a list of choices, define a picker projection that fetches only the key and the field(s) needed to render the display label — independent of the list-item and detail projections.

```java
public interface EmployeePickerItemProjection {
    Long getKey();
    String getFirstName();
    String getLastName();
}
```

The picker projection does not extend `EmployeeListItemProjection` or
`EmployeeDetailProjection` — it is independently sized for the picker's needs.

The corresponding UI model record implements `HasCaption` to carry the display
label. See `docs/patterns/ui/vaadin/uimodel.md` for the UI-side pattern.
