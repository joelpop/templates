# Code Organization Within Constructors

When writing a Vaadin view or component constructor, organize code by
operation type and use sectioning comments instead of extracting helper or
initialization methods so the constructor reads as a single top-to-bottom
narrative.

For field declaration and initialization, see
`language/java/variable/initialization.md`.

## Sectioning Comments Instead of Helper Methods

Use sectioning comments to give readers scanning cues without scattering
initialization across `createHeader()`, `init()`, `setup()`, or similar
extracted methods:

```java
// Avoid
public MyView(...) {
    getContent().add(createToolbar(), createGrid());
}

private HorizontalLayout createToolbar() { /* ... */ }
private Grid<Item> createGrid() { /* ... */ }
```

```java
// Preferred
public MyView(...) {

    // ---------- Toolbar ----------
    var title = new H2("Items");
    title.addClassNames(LumoUtility.Margin.NONE);

    var toolbar = new HorizontalLayout(title, addButton);
    toolbar.setAlignItems(Alignment.CENTER);

    // ---------- Grid ----------
    grid.addColumn(Item::getDisplayId).setHeader("ID").setSortable(true);
    grid.addColumn(Item::getName).setHeader("Name").setFlexGrow(1);
    grid.addItemClickListener(this::onGridItemClick);

    // ---------- Assembly ----------
    var content = getContent();
    content.add(toolbar);
    content.add(grid);
}
```

A non-trivial view constructor will grow long — that is expected.

## Operation-Type Ordering

Within each section, group code by operation type in this order:

1. **Component initializations** — creating instances and configuring properties
2. **Signal definitions** — creating and configuring Vaadin Signals *(Vaadin 25+ only; see `signals.md`)*
3. **Signal bindings** — connecting signals to components *(Vaadin 25+ only)*
4. **Binder bindings** — connecting form fields to the bean model with validation
5. **Value settings** — setting initial or default values on fields or via `binder.setBean(...)`

Use blank lines between groups. On Vaadin 24, where Signals are not available,
steps 2–3 collapse into private state-management methods wired at the end of
the constructor; the overall "group by operation type" principle is unchanged.

```java
var nameField = new TextField("Name");
nameField.setRequired(true);
nameField.setMaxLength(100);

var departmentComboBox = new ComboBox<DepartmentSummary>("Department");
departmentComboBox.setItemLabelGenerator(DepartmentSummary::getCaption);
departmentComboBox.setRequired(true);
departmentComboBox.setItems(departments);   // items last — after configuration is complete

var activeCheckbox = new Checkbox("Active");

// signal definitions
var nameSignal = new ValueSignal<>("");

// signal bindings
nameField.bindValue(nameSignal);

// binder bindings
binder.forField(nameField)
      .asRequired("Name is required")
      .bind(Item::getName, Item::setName);
binder.forField(departmentComboBox)
      .asRequired("Department is required")
      .bind(Item::getDepartment, Item::setDepartment);
binder.forField(activeCheckbox)
      .bind(Item::isActive, Item::setActive);

// value settings
binder.setBean(item);
```

**`ComboBox` (and other selection components):** `setItems(...)` belongs in
step 1 — it configures the available options, not the selection. The selected
value comes from `binder.setBean(...)` in step 5.
