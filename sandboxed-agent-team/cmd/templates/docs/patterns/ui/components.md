# Component Patterns

Reusable UI component patterns for Vaadin 24+ applications: Quick Filter, Avatar, grids,
forms, dialogs, notifications, and loading indicators. The `_` unused-parameter syntax in
code examples requires Java 21+ (see `docs/patterns/conventions/java.md` → "Unused Lambda
Parameters" for the Java 17–20 alternative).

## Quick Filter

The Quick Filter is a cross-cutting pattern applied to every grid view. A single text input
field filters results across multiple columns simultaneously using case-insensitive partial
match.

```java
var quickFilter = new TextField();
quickFilter.setPlaceholder("Search…");
quickFilter.setPrefixComponent(VaadinIcon.SEARCH.create());
quickFilter.setClearButtonVisible(true);
quickFilter.setValueChangeMode(ValueChangeMode.LAZY);

// Wire to ListDataProvider filter
quickFilter.addValueChangeListener(e -> {
    dataProvider.setFilter(item ->
        matchesQuickFilter(item, e.getValue())
    );
});

private boolean matchesQuickFilter(ItemListItem item, String query) {
    if (query == null || query.isBlank()) return true;
    var q = query.strip().toLowerCase();
    return item.getName().toLowerCase().contains(q)
        || item.getCode().toLowerCase().contains(q);
}
```

The columns that participate in the Quick Filter are defined per view in the feature
requirements. The Quick Filter may coexist with additional filter controls (status toggle,
role filter, date range) — it does not replace view-specific filters.

## Avatar Component

Use the Vaadin `Avatar` component for all user profile photos and entity logos throughout
the application. Do not use `<img>` tags or custom image components for this purpose.

```java
// User avatar
var avatar = new Avatar(user.getFullName());
avatar.setImage(user.getPhotoUrl());  // displays initials fallback when null

// Entity logo
var logo = new Avatar(entity.getName());
logo.setImage(entity.getLogoUrl());
```

When no photo or logo is set, `Avatar` displays its built-in initials fallback automatically.

Application-specific avatar ring patterns (e.g., role indicator rings) are defined in the
application-specific UI requirements — see `docs/reqs/external-interfaces/user-interfaces.md`.

## Grid Standards

> **Note:** The admin-list mechanics that previously lived in this section
> (full-dataset loading, record count, row-click → edit) were
> project-specific requirements, not agnostic patterns. They have been
> moved to `docs/reqs/functional/cross-cutting/admin-list.md` (the
> canonical cross-cutting requirement) and `docs/architecture/admin-grid.md`
> (the impl companion). Only the genuinely agnostic Column Configuration
> guidance remains below.
>
> Other sections of this file (loading indicators, confirmation dialogs)
> similarly mix req-flavored content with agnostic guidance — that
> contamination is pending cleanup in a separate branch and is not
> addressed here.

### Column Configuration

Define columns explicitly with appropriate widths, sort properties, and headers. Use
`setAutoWidth(true)` for columns where content width varies significantly:

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

## Form Standards

### Vaadin Binder

All forms use Vaadin `Binder` for field-to-model binding and validation. No manual
`getValue()` / `setValue()` form handling.

```java
var binder = new Binder<>(ItemDetail.class);

binder.forField(nameField)
      .asRequired("Name is required")
      .withValidator(n -> n.length() <= 100, "Maximum 100 characters")
      .bind(ItemDetail::getName, ItemDetail::setName);

binder.forField(codeField)
      .asRequired("Code is required")
      .bind(ItemDetail::getCode, ItemDetail::setCode);

binder.setBean(item);
```

Validation errors appear inline, adjacent to the offending field, as Binder field-level
error messages. No validation error is shown only as a toast.

### Form Layout

Use `FormLayout` for form fields. It adapts column count to available width automatically:

```java
var form = new FormLayout();
form.add(nameField, codeField, descriptionField, activeCheckbox);
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),     // 1 column on small
    new FormLayout.ResponsiveStep("600px", 2)  // 2 columns on wider
);
```

## Dialogs — Delegation Pattern

See `docs/patterns/conventions/vaadin/components.md` for the full dialog delegation pattern. Summary:

- Never extend `Dialog` — wrap it via delegation
- Implement `NonComponent` for event publishing from dialog classes
- Expose only the methods needed by callers (`open()`, `close()`, `addSaveListener()`, etc.)

### NonComponent Event Infrastructure

Dialogs assembled by delegation aren't `Component` instances, so they can't use Vaadin's
built-in `ComponentEventBus`. The `NonComponent*` trio below provides a deliberate
parallel to `Component` / `ComponentEvent` / `ComponentEventBus`, so anyone familiar
with Vaadin's component-event pattern recognizes the shape immediately:

| Component-event equivalent | Non-component counterpart |
|----------------------------|---------------------------|
| `Component` (fires events) | `NonComponent` (marker interface) |
| `ComponentEvent<S>` | `NonComponentEvent<N extends NonComponent>` |
| `ComponentEventBus` (held by the `Component`) | `NonComponentEventSupport<N>` (held by composition) |

Copy the three classes below into a shared event package in your UI module (for
example, `{base_package}.ui.event`). They are small and have no dependencies beyond
`com.vaadin.flow.shared.Registration`.

```java
/**
 * Interface for classes that can fire events but don't extend
 * {@link com.vaadin.flow.component.Component}.
 * Analogous to how {@link com.vaadin.flow.component.Component} provides event capabilities.
 */
public interface NonComponent {

    /**
     * Adds a listener for events of the given type.
     *
     * @param eventType the type of event to listen for
     * @param listener  the listener to add
     * @param <E>       the event type
     * @return a registration that can be used to remove the listener
     */
    <E extends NonComponentEvent<?>> Registration addListener(Class<E> eventType, Consumer<E> listener);
}
```

```java
/**
 * Base class for events fired by non-component sources.
 * Analogous to {@link com.vaadin.flow.component.ComponentEvent} but for classes
 * that don't extend {@link com.vaadin.flow.component.Component}.
 *
 * @param <N> the type of the event source, must implement {@link NonComponent}
 */
public abstract class NonComponentEvent<N extends NonComponent> {
    private final N source;

    /**
     * Creates a new event with the given source.
     *
     * @param source the source of the event
     */
    protected NonComponentEvent(N source) {
        this.source = source;
    }

    /**
     * Returns the source of the event.
     *
     * @return the source of the event
     */
    public N getSource() {
        return source;
    }
}
```

```java
/**
 * Helper class that provides event listener management for {@link NonComponent} implementations.
 * Use via composition to add event support to classes that don't extend
 * {@link com.vaadin.flow.component.Component}.
 *
 * <p>Designed for UI-thread use, matching {@link com.vaadin.flow.component.ComponentEventBus}
 * conventions. Callers dispatching or registering from a Push thread should wrap the call
 * in {@code UI.access(...)}.</p>
 *
 * @param <N> the type of the event source
 */
public class NonComponentEventSupport<N extends NonComponent> {
    private final Map<Class<?>, List<Consumer<?>>> listeners = new HashMap<>();

    /**
     * Adds a listener for events of the given type.
     *
     * @param eventType the type of event to listen for
     * @param listener  the listener to add
     * @param <E>       the event type
     * @return a registration that can be used to remove the listener
     */
    public <E extends NonComponentEvent<N>> Registration addListener(Class<E> eventType, Consumer<E> listener) {
        listeners.computeIfAbsent(eventType, k -> new ArrayList<>()).add(listener);
        return () -> {
            var list = listeners.get(eventType);
            if (list != null) {
                list.remove(listener);
            }
        };
    }

    /**
     * Fires an event to all registered listeners of the event's type.
     *
     * @param event the event to fire
     * @param <E>   the event type
     */
    @SuppressWarnings("unchecked")
    public <E extends NonComponentEvent<N>> void fireEvent(E event) {
        var list = listeners.get(event.getClass());
        if (list != null) {
            // Snapshot before iterating in case a listener adds or removes listeners during dispatch.
            for (var listener : new ArrayList<>(list)) {
                ((Consumer<E>) listener).accept(event);
            }
        }
    }
}
```

Dispatch is keyed by the event's runtime class, so a listener registered on
`SaveEvent.class` only fires for exact-class `SaveEvent` instances — matching
`ComponentEventBus` behavior. The `Registration` returned from `addListener` removes the
listener on `remove()`.

See `docs/patterns/conventions/vaadin/components.md` → "NonComponent Event System for Delegating
Dialogs" for the caller-side pattern (defining `SaveEvent` / `CancelEvent` subclasses,
exposing typed `addSaveListener(Consumer<SaveEvent>)` convenience methods, and firing
through the support instance).

## Notification Toasts

Success and error notifications that cannot be shown inline appear as non-blocking Vaadin
`Notification` toasts:

```java
// Success — auto-dismiss after 3 seconds
var notification = Notification.show("Saved successfully");
notification.setDuration(3000);
notification.setPosition(Notification.Position.BOTTOM_END);
notification.addThemeVariants(NotificationVariant.LUMO_SUCCESS);

// Error
var notification = Notification.show("An error occurred. Please try again.");
notification.setDuration(5000);
notification.setPosition(Notification.Position.BOTTOM_END);
notification.addThemeVariants(NotificationVariant.LUMO_ERROR);
```

Reserve toasts for outcomes that cannot be shown inline (service errors, async completions).
Do not use toasts as the sole means of reporting form validation errors.

## Service Error Handling

Services throw typed exceptions; views catch them and decide what to do. The exception
type drives the response — that is why exception types must be specific enough for the
view to branch on (see `docs/patterns/architecture/services.md` — "Error Contracts").

Three branches cover all cases:

```java
try {
    itemService.save(binder.getBean());
    Notification.show("Saved successfully").addThemeVariants(NotificationVariant.LUMO_SUCCESS);
    close();
} catch (ValidationException e) {
    // ValidationException messages are authored by the service for user consumption —
    // show them directly. Each error in the list is a discrete user-facing message.
    e.getErrors().forEach(msg ->
            Notification.show(msg)
                    .addThemeVariants(NotificationVariant.LUMO_ERROR));
} catch (EntityNotFoundException e) {
    // Record disappeared between load and save. Navigate away — leaving the user on a
    // detail form for a nonexistent record is worse than a navigation reset.
    Notification.show("Record no longer exists.")
            .addThemeVariants(NotificationVariant.LUMO_ERROR);
    getUI().ifPresent(ui -> ui.navigate(ItemListView.class));
} catch (Exception e) {
    // Unexpected — log server-side, show nothing that reveals internals.
    log.error("Unexpected error saving item", e);
    Notification.show("An error occurred. Please try again.")
            .addThemeVariants(NotificationVariant.LUMO_ERROR);
}
```

`ValidationException.getErrors()` messages are the only exception messages safe to
surface directly — they are authored for that purpose. `EntityNotFoundException` and
any unexpected exception always show a generic message. Never pass an arbitrary
`e.getMessage()` to `Notification.show()`; the exception message is internal state.

## Confirmation Dialogs

Destructive actions (deactivate, delete) require explicit confirmation before execution:

```java
private void confirmDeactivate(ItemListItem item) {
    var dialog = new ConfirmDialog();
    dialog.setHeader("Deactivate " + item.getName() + "?");
    dialog.setText("The item will no longer appear in active lists.");
    dialog.setCancelable(true);
    dialog.setConfirmText("Deactivate");
    dialog.setConfirmButtonTheme("error primary");
    dialog.addConfirmListener(_ -> deactivate(item));
    dialog.open();
}
```

Cancel leaves the record unchanged. The confirmation dialog always states what will happen
in plain language.

## Loading Indicators

Long-running operations (>200ms) show a loading indicator:

```java
saveButton.setEnabled(false);
saveButton.setText("Saving…");
try {
    itemService.save(binder.getBean());
    Notification.show("Saved successfully").addThemeVariants(NotificationVariant.LUMO_SUCCESS);
    close();
} catch (ValidationException e) {
    e.getErrors().forEach(msg ->
            Notification.show(msg).addThemeVariants(NotificationVariant.LUMO_ERROR));
} finally {
    saveButton.setEnabled(true);
    saveButton.setText("Save");
}
```

Grid data loads show a Vaadin progress bar while loading:

```java
grid.setItems(query -> {
    // Vaadin shows built-in loading indicator automatically for CallbackDataProvider
    ...
});
```

For `ListDataProvider`, display a visual cue before the service call and remove it after:

```java
progressBar.setVisible(true);
UI.getCurrent().access(() -> {
    grid.setItems(itemService.listAll());
    progressBar.setVisible(false);
});
```

## Conditional Component Rendering — Do Not Generate vs. Hide vs. Disable

Three modes, distinguished by *why*:

- **Do not generate** — **authorization** (the current user's role does not permit the
  action) or any "this user should never see this at all" condition. Never constructed,
  never in the DOM. Nothing to discover via dev tools, no attribute to re-enable by
  tampering.
- **Hide** (`setVisible(false)`) — **not applicable to the current situation** but could
  apply in another state the same user encounters (e.g., a Reactivate button hidden on
  an active record, shown on an inactive one). Component is constructed and lives in the
  layout so it can be revealed without a rebuild.
- **Disable** (`setEnabled(false)` + tooltip) — **applicable and authorized, but not
  possible right now** (e.g., cannot deactivate the last remaining admin, cannot submit
  while a save is in flight). Tooltip must explain *why*.

```java
// Do not generate — e.g., security gating by role. The button is never constructed for
// users who cannot perform the action; there is no attribute to tamper with in the DOM.
if (currentUser.hasRole(UserRole.ROLE_ADMIN)) {
    var deactivateButton = new Button("Deactivate");
    deactivateButton.addClickListener(_ -> confirmDeactivate());
    toolbar.add(deactivateButton);
}

// Hide — contextually not applicable. The button exists but is not visible for this
// record's current state; it re-appears when the state changes.
reactivateButton.setVisible(!record.isActive());

// Disable — applicable and authorized, temporarily blocked. Show a tooltip explaining why.
deactivateButton.setEnabled(!isLastAdmin);
deactivateButton.setTooltipText(isLastAdmin
    ? "Cannot deactivate the only admin account" : null);
```

**Do not** use `setVisible(false)` for authorization gating — that is the "do not
generate" case. `setEnabled(false)` is not a substitute for the other two: a
permanently-disabled control communicates "try again later" and invites futile
interaction.

### Layout Preservation — When a Placeholder Is Needed

"Do not generate" and `setVisible(false)` both remove the component from layout flow;
surrounding elements collapse to fill the space. "Disable" preserves layout at its
normal size and position. When absence would cause jarring layout shift or misalignment
(e.g., a toolbar with positionally-dependent buttons), the absent component needs a
**placeholder** that occupies the same space without being interactive.

> **Vaadin state lives on the server.** Interactivity is governed by the server-side
> component state (`setEnabled(false)`, `setVisible(false)`), not by client-side CSS or
> HTML attributes. Any purely-CSS concealment (`visibility: hidden`, `display: none` via
> `getStyle().set(...)`, a client-side `disabled` attribute toggled in the browser
> inspector) can be reverted by the user in dev tools; if the server still considers the
> component enabled and visible, the server will happily process clicks on the
> "re-revealed" element. **Use server-side state, not CSS, to prevent interaction.**

Options, in preference order:

1. **Rethink the mode.** If a missing control would disrupt layout, "disable with
   tooltip" preserves the slot by design, blocks interaction on the server, and
   communicates intent more clearly than an empty placeholder.
2. **An inert placeholder** the same size as the control — a `Div`, empty `Span`, or
   subdued "—" label. A separate component, not the real control styled invisible — no
   enabled element hiding in the DOM.
3. **A neutral affordance** ("No action" label, subdued status pill) when the absence
   itself is informative.

Do not construct the real interactive component and hide it with CSS alone
(`getStyle().set("visibility", "hidden")` on an enabled button). That leaves a
fully-interactive server-side component reachable from the browser inspector. If the real
component must occupy the slot, `setEnabled(false)` on the server brings you back to
disable mode — the slot is already preserved.

The three-mode rubric answers *whether* the component exists or functions. Placeholder
decisions answer *what occupies the space* — they compose with the mode choice and must
never weaken server-enforced interactivity state.
