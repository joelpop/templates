# View Constructor Initialization

When building a Vaadin view or component constructor, keep all UI assembly inline
with sectioning comments and organize code by operation type (component inits →
signal definitions → signal bindings → binder bindings → value settings) so the
constructor reads as a single top-to-bottom narrative without arbitrary
factory-method extractions.

## UI Initialization in Constructors

Keep all UI initialization in the constructor rather than splitting it across helper
methods such as `createHeader()`, `createContent()`, etc. A non-trivial view
constructor will grow long — that is expected. Use **sectioning comments** inside the
constructor instead of extracting helpers; keep related setup together under a label
so the constructor reads top-to-bottom like a narrative:

```java
// Preferred
public MyView(...) {

    // ---------- Toolbar ----------
    var title = new H2("Items");              // local — only used here
    title.addClassNames(LumoUtility.Margin.NONE);

    var toolbar = new HorizontalLayout(title, addButton);
    toolbar.setAlignItems(Alignment.CENTER);

    // ---------- Grid ----------
    grid.addColumn(Item::getDisplayId).setHeader("ID").setSortable(true);
    grid.addColumn(Item::getName).setHeader("Name").setFlexGrow(1);
    grid.addItemClickListener(this::onGridItemClick);

    // ---------- Binder / Signals / Value settings ----------
    // (follow the order documented in "Code Organization Within Methods" below)

    // ---------- Assembly ----------
    getContent().add(toolbar, grid);
}

// Avoid
public MyView(...) {
    getContent().add(createToolbar(), createGrid());  // what does createToolbar add?
}

private HorizontalLayout createToolbar() { ... }
private Grid<Item> createGrid() { ... }
```

Sectioning comments give readers the same scanning cues as factory method names,
without scattering initialization or creating arbitrary-extraction decisions.
Inline assembly also keeps field count down: a component used only during
construction stays a local variable; extraction forces it into a field so the
helper method can reach it. Fields that do need to persist can be declared
`final` — Java requires `final` fields to be assigned directly in the
constructor body, not via helper methods.

## Code Organization Within Methods

Group code by operation type, not by component. Within a constructor or method,
organize in this order:

1. **Component initializations** — creating instances and configuring properties
2. **Signal definitions** — creating and configuring Vaadin Signals *(skip on Vaadin <25; see `docs/patterns/ui/vaadin/signals.md`)*
3. **Signal bindings** — connecting signals to components (reactive UI) *(skip on Vaadin <25)*
4. **Binder bindings** — connecting form fields to bean model (with validation)
5. **Value settings** — setting initial/default values on fields or bean on Binder

Use blank lines between groups for readability. On Vaadin <25, where Signals are not
available, steps 2–3 collapse into private state-management methods wired at the end
of the constructor instead; the overall "group by operation type" principle is unchanged.

```java
// Preferred: grouped by operation type
var nameField = new TextField("Name");
nameField.setRequired(true);
nameField.setMaxLength(100);

var departmentComboBox = new ComboBox<DepartmentSummary>("Department");
departmentComboBox.setItemLabelGenerator(DepartmentSummary::getName);
departmentComboBox.setRequired(true);
departmentComboBox.setItems(departments);                   // items last — after configuration is complete

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

**`ComboBox` (and other selection components):** `setItems(...)` belongs in step 1
(component initialization) — it configures the available options, not the selection.
The selected value comes from `binder.setBean(...)` in step 5. Keep `setItems` next to
other `setXxx` calls, not mixed with bindings.
