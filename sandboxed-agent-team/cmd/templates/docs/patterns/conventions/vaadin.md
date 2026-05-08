# Vaadin Coding Conventions

Standards for building Vaadin 24+ server-side UI components, views, layouts, and event
systems. Patterns are identical across 24.x and 25.x unless an inline **"Vaadin ≥X / <X"**
note calls out a version-specific difference. See `docs/agnostic/README.md` → "Version
Compatibility" for the summary matrix.

## `@SpringComponent` over `@Component` in Vaadin Modules

In modules that participate in Vaadin's UI lifecycle — views, layouts, listeners,
and anything that might grow a Vaadin scope (`@UIScope`, `@RouteScope`,
`@VaadinSessionScope`) — register Spring beans with Vaadin's `@SpringComponent`
rather than Spring's own `@Component`. `@SpringComponent` is `@Component` plus the
scope-aware plumbing that Vaadin's Flow-Spring integration expects; using plain
`@Component` on a UI bean works until the bean later grows a Vaadin scope and
silently fails to integrate with session or UI lifecycle.

**Project rule:** use `@SpringComponent` for every `@Component`-style bean in
the application's Vaadin-facing modules (typically the `app` and `ui` modules).
Non-Vaadin modules — service interfaces, service implementations, JPA model, JPA
client, common utilities — continue to use regular `@Component`.

```java
// Preferred — in the app or ui module
import com.vaadin.flow.spring.annotation.SpringComponent;

@SpringComponent
public class AuthMethodCombinabilityValidator { ... }

// Avoid — plain @Component in a Vaadin module, unless you know the bean will
// never need a Vaadin scope
@Component
public class AuthMethodCombinabilityValidator { ... }
```

`@Configuration` is unchanged — it is a specialization of `@Component` used for
`@Bean`-factory classes and does not need to be swapped to `@SpringComponent`. A
`@Configuration` class that hosts `@Bean` methods remains `@Configuration`; the
beans it returns are not subject to this convention unless they themselves are
scanned components in a Vaadin module.

`@Route`-annotated views are automatically scope-aware through Vaadin's router
infrastructure and do not need `@SpringComponent` on top of `@Route`.

## Views Must Extend Composite<T>

All custom view and component classes must extend `Composite<T>` with an appropriate root
component type. Do not extend layout classes directly.

```java
// Preferred
@Route("items")
@RolesAllowed(UserRole.ROLE_USER)
public class ItemView extends Composite<VerticalLayout> {
    // getContent() returns the root VerticalLayout
}

// Avoid
public class ItemView extends VerticalLayout { ... }
```

`Composite<T>` provides better encapsulation: only the explicitly exposed API is accessible
to callers, not the full component API of the root layout.

## Per-View Package Layout

Every `@Route`-annotated view lives in its own Java package named after the view's
terminal path segment. View packages are grouped under a path-prefix package — e.g.
`…ui.view.admin.organization` for a view at `@Route("admin/organization")`, or
`…ui.view.platform.tenant` for `@Route("platform/tenant")`. The view class itself
and any UI that serves only that view — per-view editors, grid cell renderers, form
helpers — live in the view's package. UI shared across views stays in a top-level
`component` (or equivalent) package.

```
ui/
├── layout/
│   └── BaseView.java           <- shared view base class
├── component/                  <- UI shared across views
└── view/
    ├── admin/
    │   ├── organization/
    │   │   └── OrganizationView.java        @Route("admin/organization")
    │   └── user/
    │       └── UserView.java                @Route("admin/user")
    └── platform/
        └── tenant/
            └── TenantView.java              @Route("platform/tenant")
```

This gives each view a clear home that grows naturally with its own helper classes,
mirrors the side-nav grouping on disk, and prevents the single-flat-package
anti-pattern flagged in architecture debt.

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
`NonComponentEventSupport<N>`) are defined in `docs/agnostic/ui/components.md` →
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

## vaadin.allowed-packages

`application.properties` must declare `vaadin.allowed-packages` to include all Vaadin
component package prefixes used by the application. Without this, Vaadin components in
add-ons may fail to render.

```properties
vaadin.allowed-packages=com.vaadin,org.vaadin
```

When adding a new Vaadin add-on, add its root package prefix to this property.

## vaadin-dev Dependency

The `vaadin-dev` artifact (or equivalent) must be included as a dev-mode dependency with
`<optional>true</optional>` so it is present in development but excluded from production JARs:

```xml
<dependency>
    <groupId>com.vaadin</groupId>
    <artifactId>vaadin-dev</artifactId>
    <optional>true</optional>
</dependency>
```

## Access Annotations on Layout Classes

The main layout class should carry `@PermitAll` per Vaadin's
[security doc](https://vaadin.com/docs/latest/flow/security/enabling-security#error-messages-for-unauthorized-views).
This is the doc-recommended choice: anonymous users hitting protected routes
get redirected to login by Vaadin's `ExceptionHandlingConfigurer` instead of
seeing the layout chrome before the redirect.

> **Precondition:** `VaadinSecurityConfigurer` defaults `anyRequest` to
> `denyAll` at the URL filter level. Without overriding that default,
> `@PermitAll` on the layout silently fails — anonymous users get a bare
> HTTP 403 from Spring Security before Vaadin's exception handler runs,
> instead of being redirected to login. Override the default inside the
> configurer block:
>
> ```java
> http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
>     configurer.loginView(LoginView.class);
>     configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
> });
> ```
>
> Without this override, `@AnonymousAllowed` on the layout is the workaround
> that keeps things working — it puts the layout into Vaadin's
> `defaultPermitMatcher` ("anonymous routes"), sidestepping the URL-level
> `denyAll`. It's a workaround, not the intended design: it exposes layout
> chrome to anonymous users momentarily before the login redirect, and the
> Vaadin docs explicitly recommend `@PermitAll`. Use the `anyRequest=permitAll`
> override to make `@PermitAll` work as documented; fall back to
> `@AnonymousAllowed` only when touching the configurer isn't an option.

The `@Layout` annotation itself was introduced in Vaadin 24.1 — on Vaadin
24.0, set the layout on each route via `@Route(value = "...", layout =
MainLayout.class)` and have `MainLayout` extend `AppLayout` (which already
implements `RouterLayout`).

```java
// Vaadin 24.1+ / 25+
@Layout
@PermitAll
public class MainLayout extends AppLayout {
    // DO NOT add "implements RouterLayout" — already inherited from AppLayout
}
```

```java
// Vaadin 24.0 — no @Layout; set layout per @Route
@Route(value = "items", layout = MainLayout.class)
public class ItemView extends Composite<VerticalLayout> { ... }
```

## MainLayout Must Not Implement RouterLayout

`AppLayout` already implements `RouterLayout`. Adding `implements RouterLayout` explicitly
is redundant and can cause unexpected behavior (especially in Vaadin ≥25 where the
access-checker treats layouts differently).

## Lumo Theme and LumoUtility

Use `LumoUtility` class constants for padding, margin, color, flexbox, and sizing. Do not
write custom CSS for things Lumo provides. Use component theme variants before resorting
to custom styling.

```java
// Preferred — LumoUtility constants
content.addClassNames(
    LumoUtility.Padding.MEDIUM,
    LumoUtility.Gap.SMALL,
    LumoUtility.Display.FLEX,
    LumoUtility.FlexDirection.COLUMN
);

// Preferred — component variants
button.addThemeVariants(ButtonVariant.LUMO_PRIMARY);
grid.addThemeVariants(GridVariant.LUMO_ROW_STRIPES);
```

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

## @Menu Annotation for Navigation

> **Vaadin ≥24.4:** use the `@Menu` annotation on `@Route` classes for views that appear
> in the main navigation menu. Do not manually add menu items to the `SideNav` for routed
> views — let `MenuConfiguration.getMenuEntries()` discover them.
>
> **Vaadin <24.4:** `@Menu` is not available. Register `SideNavItem` instances manually in
> `MainLayout` and use the role-conditional rendering pattern (see
> `docs/agnostic/ui/navigation.md` → Conditional Navigation Rendering).

```java
@Route("items")
@Menu(order = 1, icon = "vaadin:list")      // Vaadin ≥24.4
@RolesAllowed(UserRole.ROLE_USER)
public class ItemView extends Composite<VerticalLayout> { ... }
```

## Browser Client Details — Bridging the SoC Wall

Browser-level details (timezone, screen dimensions, locale) are available on the UI
thread but not in the service or persistence layers, which have no access to Vaadin
APIs (enforced by module boundaries). When service-layer code needs browser context —
most commonly timezone for `InstantMapper` (see `persistence.md`) — a bridging
service is required.

### The Pattern

1. **Service interface** (in the service module, no Vaadin dependency):

```java
public interface ClientDetailsService {
    ZoneId getBrowserTimezone();
    // add other browser details as needed
}
```

2. **UI-module implementation** (has Vaadin access):

The implementation obtains browser details via Vaadin's client-detail APIs and
caches them for the duration of the session. Cache in `VaadinSession` attributes —
not via `@SessionScope`, which is not compatible with Push threads.

On the UI thread, browser details may be directly accessible (no bridging needed).
The service exists primarily so that **non-UI code** (service implementations,
MapStruct mappers) can access browser context without a Vaadin dependency.

### Version-Specific Notes

The mechanism for obtaining extended browser/client details varies by Vaadin version.
Always consult the `vaadin` MCP server or current Vaadin documentation for the
correct API before implementing.

> **Vaadin 24.x:** `UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> ...)`
> — asynchronous callback pattern. The callback-based API makes caching in `VaadinSession`
> essential, since the details are not synchronously available.
>
> **Vaadin ≥25:** Consult the current documentation. The mechanism may differ — for
> instance, if browser details are synchronously available on the UI thread, the
> `VaadinSession` cache is primarily useful for the persistence side of the module
> boundary (where `UI.getCurrent()` is not available).

### Why Not @SessionScope?

`@SessionScope` beans are tied to the HTTP request's `HttpSession`. Vaadin's Push
threads (used for Signals and server push) do not run within an HTTP request context,
so `@SessionScope` beans are not accessible from Push threads. `VaadinSession`
attributes are accessible from both HTTP request threads and Push threads, making
them the correct storage mechanism for session-scoped data in a Vaadin application.

### Fallback Behavior

If browser details are not yet available (e.g., first request before client
round-trip), the service should return a sensible default — typically the system
default timezone (`ZoneId.systemDefault()`). The fallback is transparent to callers.
