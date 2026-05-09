# Java Coding Conventions

Standards for Java code style, organization, and idiomatic usage in Vaadin + Spring Boot projects.

## Type Inference

Use `var` for local variable type inference whenever the type is obvious from the right-hand side:

```java
// Preferred
var form = new FormLayout();
var items = itemService.listAll();

// Avoid — type is redundant
FormLayout form = new FormLayout();
List<ItemListItem> items = itemService.listAll();
```

## Member Variable Initialization

Initialize member variables in constructors, not at the declaration site. All initialization
is visible in one place; components needed only during construction stay as locals.

```java
// Preferred
public class EditDialog {
    private final Dialog dialog;
    private final TextField nameField;

    public EditDialog(...) {
        dialog = new Dialog();
        nameField = new TextField("Name");
        // ...
    }
}

// Avoid
public class EditDialog {
    private final Dialog dialog = new Dialog();       // initialization scattered
    private final TextField nameField = new TextField("Name");
}
```

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
2. **Signal definitions** — creating and configuring Vaadin Signals *(skip on Vaadin <25; see `docs/patterns/conventions/vaadin.md` → "Signals — When to Use Them")*
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

## Local Variable Declaration

Declare local variables close to their first use, not at the top of the method:

```java
// Preferred
var form = new FormLayout();
// ... configure form ...

var content = new VerticalLayout();   // declared just before it's needed
content.add(form);
return content;

// Avoid
var content = new VerticalLayout();   // declared too early
var form = new FormLayout();
// ... configure form ...
content.add(form);
return content;
```

## Nested Types Placement

Place nested types (inner classes, enums) at the **end** of the class, after all methods:

```java
public class EditDialog {
    // Fields
    // Constructor
    // Public API methods
    // Private helper methods

    // ========== Event Classes ==========
    public static class SaveEvent extends NonComponentEvent<EditDialog> { ... }
    public static class CancelEvent extends NonComponentEvent<EditDialog> { ... }

    // Private enums last
    private enum Mode { CREATE, EDIT }
}
```

Use `private` visibility for enums and inner classes used only internally.

## Event Handler Naming

Use the pattern `on{ComponentName}{EventType}` for event handler methods:

```java
saveButton.addClickListener(this::onSaveButtonClick);
nameField.addValueChangeListener(this::onNameFieldValueChanged);
upload.addSucceededListener(this::onPhotoUploadSucceeded);
```

## Lambdas vs. Method References

Prefer method references over inline lambdas for multi-line handlers. Use inline lambdas
when the handler needs to capture a local variable:

```java
// Preferred: method reference for multi-line logic
saveButton.addClickListener(this::onSaveButtonClick);

// Acceptable: lambda captures a local variable
var confirmDialog = new Dialog();
deleteButton.addClickListener(_ -> {
    confirmDialog.close();
    onDeleteConfirmed();
});
```

## Unused Lambda Parameters

> **Java 21+:** use `_` (the unnamed variable) for unused lambda parameters and unused
> catch clause variables. This is the preferred form on any Java-21-or-newer project
> (including Vaadin 24.1+ on Java 21 and every Vaadin 25+ project).
>
> **Java 17–20:** `_` is not available. Name the parameter `unused` to signal intent at
> the call site — the name itself tells any reader the value is deliberately ignored.

```java
// Java 21+: preferred
saveButton.addClickListener(_ -> save());
cancelButton.addClickListener(_ -> close());

try {
    // ...
} catch (ValidationException _) {
    Notification.show("Please fix the validation errors");
}

// Java 17–20: name the parameter `unused`
saveButton.addClickListener(unused -> save());

try {
    // ...
} catch (ValidationException unused) {
    Notification.show("Please fix the validation errors");
}
```

## Null Check Policy

Do not null-check where the framework guarantees non-null results:

- `TextField.getValue()` returns `""`, never `null`
- `MultiSelectComboBox.getValue()` returns an empty `Set`, never `null`
- `Signal` fields initialized with a non-null value remain non-null

Do not filter Streams for nulls when the data source guarantees non-null elements.

## Dependency Injection

Use **constructor injection** for all mandatory dependencies. Declare fields `private
final` and let Spring call the single public constructor — no `@Autowired` annotation
needed (Spring 4.3+ auto-detects the single constructor).

```java
// Preferred
@SpringComponent
public class FleetAcuityOidcUserService {
    private final TenantRepository tenantRepository;
    private final UserRepository userRepository;

    public FleetAcuityOidcUserService(TenantRepository tenantRepository,
                                      UserRepository userRepository) {
        this.tenantRepository = tenantRepository;
        this.userRepository = userRepository;
    }
}

// Avoid — field injection
@SpringComponent
public class FleetAcuityOidcUserService {
    @Autowired private TenantRepository tenantRepository;
    @Autowired private UserRepository userRepository;
}
```

Constructor injection makes dependencies explicit, allows `final` fields, works in
plain unit tests without reflection, and fails fast at instantiation.

**Avoid `@PostConstruct` for validation or simple setup** — throw from the constructor
instead. A constructor that throws fails Spring Boot startup with a clear stack trace;
there's no reason to defer validation to a later lifecycle phase.

```java
// Preferred — fail from the constructor
@SpringComponent
public class AuthMethodCombinabilityValidator {
    public AuthMethodCombinabilityValidator(FleetAcuityAuthProperties props) {
        if (!props.formLogin().enabled()
                && !props.passkey().enabled()
                && !props.sso().enabled()) {
            throw new IllegalStateException(
                    "No authentication methods enabled — at least one of " +
                    "fleet-acuity.auth.{form-login,passkey,sso}.enabled must be true.");
        }
        // ... additional combinability checks
    }
}

// Avoid — @PostConstruct for validation
@SpringComponent
public class AuthMethodCombinabilityValidator {
    @Autowired private FleetAcuityAuthProperties props;

    @PostConstruct
    void validate() { /* ... */ }
}
```

Reserve `@PostConstruct` for the rare case requiring all beans' lifecycles to be
complete before work can begin (e.g., cross-bean warm-up). Setter / field injection
is acceptable only for truly optional or reconfigurable dependencies — not needed here.

## SOLID Principles

- **Single Responsibility:** Each class has one reason to change. Views display; services
  enforce business rules; repositories store.
- **Open/Closed:** Extend behavior through interfaces and composition, not by modifying
  existing classes.
- **Liskov Substitution:** Subtypes must be substitutable for their base types.
- **Interface Segregation:** Prefer small, focused interfaces over large omnibus ones.
- **Dependency Inversion:** Depend on abstractions (`*Service` interfaces), not
  implementations (`Jpa*Service`).

## Access Modifiers

Use the most restrictive access modifier that satisfies the requirement:
- `private` for internal implementation details
- package-private for intra-package collaborators
- `protected` only when designed for inheritance
- `public` only for intentional public API

## Suppressing Warnings

Every `@SuppressWarnings` annotation must carry an inline `//` comment naming the
specific warning being suppressed and the reason. A reader should not have to leave
the file to understand why a warning was silenced.

```java
// Preferred
@SuppressWarnings("java:S2160") // key-based equality is inherited from RootEntity
public class EmployeeEntity extends BaseEntity<Long> { ... }

@SuppressWarnings("unchecked") // raw type comes from legacy third-party API
var bean = (Map<String, Object>) legacyApi.getProperties();

// Avoid — no context
@SuppressWarnings("java:S2160")
public class EmployeeEntity extends BaseEntity<Long> { ... }
```

Without a comment, a suppression is indistinguishable from a cargo-culted one — future
readers can't tell whether it's still load-bearing or the underlying issue has been fixed.

## JavaDoc

Write JavaDoc on all public service interface methods. Include:
- `@param` for every parameter
- `@return` describing the return value (omit for `void`)
- `@throws` for every checked and documented unchecked exception

```java
/**
 * Finds the entity with the given key.
 *
 * @param key the surrogate primary key
 * @return the matching entity data
 * @throws EntityNotFoundException if no entity with that key exists in the current context
 */
ItemDetail findByKey(long key);
```
