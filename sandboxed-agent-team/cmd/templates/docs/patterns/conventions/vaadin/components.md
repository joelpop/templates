# Vaadin Component Conventions

Component construction, dialogs, events, state management, and styling.

## UI Initialization in Constructors

Keep all UI initialization in the constructor rather than splitting it across helper methods
such as `createHeader()`, `createContent()`, etc. A non-trivial view constructor will grow
long — that is expected. Use **sectioning comments** inside the constructor instead of
extracting helpers; keep related setup together under a label so the constructor reads
top-to-bottom like a narrative:

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

Sectioning comments give readers the same scanning cues as factory method names, without
scattering initialization or creating arbitrary-extraction decisions.

## Code Organization Within Methods

Group code by operation type, not by component. Within a constructor or method, organize
in this order:

1. **Component initializations** — creating instances and configuring properties
2. **Signal definitions** — creating and configuring Vaadin Signals *(skip on Vaadin <25; see "Signals — When to Use Them" below)*
3. **Signal bindings** — connecting signals to components (reactive UI) *(skip on Vaadin <25)*
4. **Binder bindings** — connecting form fields to bean model (with validation)
5. **Value settings** — setting initial/default values on fields or bean on Binder

Use blank lines between groups for readability. On Vaadin <25, where Signals are not
available, steps 2–3 collapse into private state-management methods wired at the end of
the constructor instead; the overall "group by operation type" principle is unchanged.

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

## Custom Dialogs Use Delegation, Not Inheritance

No application class may extend `Dialog`. Wrapping `Dialog` via delegation exposes a
focused, intentional API rather than Dialog's 50+ public methods.

```java
// Preferred
public class EditItemDialog {
    private final Dialog dialog;

    public EditItemDialog(...) {
        dialog = new Dialog();
        // configure dialog contents
    }

    public void open() { dialog.open(); }
    public void close() { dialog.close(); }
}

// Avoid
public class EditItemDialog extends Dialog { ... }
```

## NonComponent Event System for Delegating Dialogs

Classes that do not extend `Component` — such as delegating dialogs — must use the
`NonComponent` event infrastructure for event publishing. This provides the same listener
registration and event firing semantics as Vaadin's `ComponentEvent` system.

The three infrastructure classes (`NonComponent`, `NonComponentEvent<N>`, and
`NonComponentEventSupport<N>`) are defined in `docs/patterns/ui/components.md` →
"NonComponent Event Infrastructure" — copy them into a shared event package in your UI
module. This section covers the *caller-side* pattern: defining typed event subclasses,
exposing convenience `add*Listener` methods, and firing through the support instance.

```java
public class EditItemDialog implements NonComponent {
    private final Dialog dialog;
    private final NonComponentEventSupport<EditItemDialog> eventSupport = new NonComponentEventSupport<>();

    // Event class — extend NonComponentEvent<SourceType>
    public static class SaveEvent extends NonComponentEvent<EditItemDialog> {
        private final ItemDetail item;

        public SaveEvent(EditItemDialog source, ItemDetail item) {
            super(source);
            this.item = item;
        }

        public ItemDetail getItem() { return item; }
    }

    public static class CancelEvent extends NonComponentEvent<EditItemDialog> {
        public CancelEvent(EditItemDialog source) { super(source); }
    }

    // Implement NonComponent
    @Override
    public <E extends NonComponentEvent<?>> Registration addListener(Class<E> eventType,
                                                                      Consumer<E> listener) {
        return eventSupport.addListener((Class) eventType, listener);
    }

    // Convenience listener registration methods
    public Registration addSaveListener(Consumer<SaveEvent> listener) {
        return eventSupport.addListener(SaveEvent.class, listener);
    }

    public Registration addCancelListener(Consumer<CancelEvent> listener) {
        return eventSupport.addListener(CancelEvent.class, listener);
    }

    protected void fireEvent(NonComponentEvent<EditItemDialog> event) {
        eventSupport.fireEvent(event);
    }
}
```

The caller attaches listeners and handles events:

```java
var dialog = new EditItemDialog(...);
dialog.addSaveListener(e -> handleSave(e.getItem()));
dialog.addCancelListener(e -> dialog.close());
```

## Signals — When to Use Them

> **Vaadin ≥25.1:** Vaadin Signals are the **preferred mechanism for non-`Binder` component
> state management** — UI state that is not backed by a JPA-style bean (visibility toggles,
> selection state, reactive counts, cross-session shared data, computed derivations). Use
> `ValueSignal`, `ListSignal`, `SharedNumberSignal`, `SharedMapSignal`, `Signal.computed(...)`
> as the default state-management primitives.
>
> **Vaadin 25.0:** Signals exist but are not yet positioned as the universal preference. Use
> them for reactive/shared state where they add clear value (cross-session updates, computed
> values); private state-management methods remain acceptable for simple cases.
>
> **Vaadin <25 (24.x):** Signals are **not available**. Manage non-`Binder` state through
> private fields, manual listener wiring, and explicit rebind/refresh methods on the view or
> composite. Do not try to backport Signals.
>
> **`Binder` is always preferred for bean-backed forms**, regardless of Vaadin version —
> field-to-property binding with validation is the job `Binder` does best. Signals
> complement `Binder`; they do not replace it.

## Signal Field Naming

*(Applies only when Signals are in use — see "Signals — When to Use Them" above.)*

Suffix Signal fields with their signal type for clarity:

- `ListSignal` fields: `itemsListSignal`
- `ValueSignal` fields: `selectedItemSignal`, `editingSignal`
- Computed signals: `totalValueSignal`, `filteredCountSignal`

```java
private final ListSignal<ItemListItem> itemsListSignal;
private final ValueSignal<Boolean> editingSignal;

// computed signal — local variable in constructor is fine
Signal<Integer> visibleCountSignal = Signal.computed(() -> ...);
```

## UI Component Field Naming

When a field holds a UI component and the component type is not obvious from the name alone,
include the component type as a suffix:

```java
private final Span totalValueSpan;          // not: totalValue
private final Button saveButton;            // clear enough
private final TextField displayIdField;     // "Field" distinguishes from model property
private final Grid<ItemListItem> itemGrid;  // "Grid" clarifies
```

## Lumo Theme and LumoUtility

### Styling priority

Reach for each level in order before moving to the next:

1. **Component API** — `setWidth()`, `setMinWidth()`, `setMaxWidth()`, `setHeight()`,
   `setPadding()`, `setSpacing()`, `setMargin()`, etc. If the component exposes a Java method for the
   style concern, use it — it is the most refactor-safe and type-safe option.
2. **Component theme variants** — `addThemeVariants(...)`. Zero CSS, zero class names;
   the component expresses its own intent.
3. **`LumoUtility` class constants** — `addClassNames(LumoUtility.Padding.MEDIUM, ...)`.
   Covers padding, margin, gap, colour, typography, flexbox, sizing, and more without
   writing any CSS.
4. **`getStyle().set(...)`** — inline style on a specific element. Use when the value
   is dynamic or not covered by a `LumoUtility` constant.
5. **Custom CSS** — a `.css` file imported via `@StyleSheet`. Last resort, for
   structural rules that cannot be expressed any other way.

```java
// Level 1 — component API
layout.setPadding(true);
layout.setSpacing(true);
layout.setWidth("320px");

// Level 2 — component variant
button.addThemeVariants(ButtonVariant.LUMO_PRIMARY);
grid.addThemeVariants(GridVariant.LUMO_ROW_STRIPES);

// Level 3 — LumoUtility constants
content.addClassNames(
    LumoUtility.Padding.Horizontal.MEDIUM,
    LumoUtility.Gap.SMALL,
    LumoUtility.Display.FLEX,
    LumoUtility.BoxSizing.BORDER,
    LumoUtility.FlexDirection.COLUMN
);

// Level 4 — inline style (dynamic value)
badge.getStyle().set("--badge-color", item.getColor());

// Level 5 — custom CSS (last resort)
// @StyleSheet("styles.css") on the AppShell class
```

### Applying the Lumo theme (`@Theme` deprecated since Vaadin 25)

`@Theme` is deprecated as of Vaadin 25. Apply Lumo, the utility stylesheet, and
project CSS via `@StyleSheet` annotations on the application's `AppShellConfigurator`
class in this order:

```java
@StyleSheet(Lumo.STYLESHEET)          // Lumo base theme
@StyleSheet(Lumo.UTILITY_STYLESHEET)  // LumoUtility classes
@StyleSheet("styles.css")             // project overrides
public class AppShell implements AppShellConfigurator { }
```

Order matters: project styles load after Lumo so they can override it. On
Vaadin 24.x, continue using `@Theme(Lumo.class)` — the `@StyleSheet` approach
is a Vaadin 25+ API.

## Binder for Forms

All forms must use Vaadin `Binder` for field-to-model binding and validation. Manual
`getValue()` / `setValue()` form handling is not permitted.

```java
var binder = new Binder<>(Item.class);

binder.forField(nameField)
      .asRequired("Name is required")
      .withValidator(n -> n.length() <= 100, "Name must be 100 characters or fewer")
      .bind(Item::getName, Item::setName);

binder.setBean(item);
```

Validation errors must appear inline, adjacent to the offending field. No validation error
may be displayed only as a toast.
