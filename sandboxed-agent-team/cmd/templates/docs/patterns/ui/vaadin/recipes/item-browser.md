# Recipe: Item Browser — List + Detail Component Family

When implementing a "list of items with optional CRUD" view (often used for admin
views), follow this recipe to produce the `ItemBrowser` / `EditableItemBrowser` /
`ItemEditor` component family — a list half (grid, toolbar, filters) plus a swappable
detail half (split pane, dialog, or separate view) that share the same wiring code
regardless of
layout.

## What this produces

| Class / Interface            | Role                                                                                                                                                                                             |
|------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `ItemBrowser<T>`             | List half: grid + toolbar (quick filter, filter popover, new-entity button) + row styling. Used standalone for read-only lists or as the inner list of any `EditableItemBrowser` implementation. |
| `EditableItemBrowser<T>`     | Interface: contract for CRUD-capable wrappers around `ItemBrowser`.                                                                                                                              |
| `ItemEditor<T>`              | Interface: contract for the detail-form component shown inside an `EditableItemBrowser`.                                                                                                         |
| `MasterDetailItemBrowser<T>` | Concrete `EditableItemBrowser`: split-pane layout (Vaadin `MasterDetailLayout`), responsive overlay below master min-width. The built implementation.                                            |
| `DialogItemBrowser<T>`       | Concrete `EditableItemBrowser` (project-built as needed): editor in a modal `Dialog`; list stays visible behind the dimmer.                                                                      |
| `ViewItemBrowser<T>`         | Concrete `EditableItemBrowser` (project-built as needed): editor is a separate routed view; `ItemEditor` implementation drives navigation.                                                       |
| `FilterOption<V>`            | Interface: closed-set filter options (enums) each carrying a match predicate.                                                                                                                    |

## Dependencies

- Vaadin 25+ — `ValueSignal`, `ListSignal`, `Signal.effect`,
  `bindVisible`, `bindText` are used throughout for reactive toolbar and
  detail state. On Vaadin 24 the signal API is absent; replace with
  imperative `setVisible` calls at each state-mutation site.
- `MasterDetailLayout` — a Vaadin preview feature at time of writing;
  enable it in `vaadin-featureflags.properties`.
- `HasCaption` — for filter-option enum display (see `ui/vaadin/uimodel/has-caption.md`).
- `AppIcon` — for toolbar and action icons (see `ui/vaadin/app-icon.md`).

## Component family overview

```
ItemBrowser<T>                     ← standalone read-only list
     ↑ owns (composition)
EditableItemBrowser<T>             ← CRUD interface
     ↑ implements
MasterDetailItemBrowser<T>         ← split-pane (built)
DialogItemBrowser<T>               ← modal dialog (project-built)
ViewItemBrowser<T>                 ← routed view (project-built)

ItemEditor<T>                      ← detail-form interface
     ↑ implements (per feature)
TenantEditor, UserEditor, …        ← feature-specific editors
```

`MasterDetailItemBrowser` and the two unbuilt variants all embed an
`ItemBrowser` by composition and delegate every list-configuration
method to it. Callers configure the list the same way regardless of
which layout they pick.

## Step 1 — `ItemBrowser<T>`

Place in `{base_package}.ui.component.itembrowser`. `ItemBrowser`
extends `Composite<VerticalLayout>` and implements `HasSize`.

### Toolbar

The toolbar is signal-driven: each segment (quick filter field, filter
button wrap, new-entity button, overflow menu) binds its own visibility
to a `ValueSignal<Boolean>`. The toolbar itself binds to the OR of all
segment signals, so it hides when nothing is configured and appears as
soon as any feature is enabled — no imperative `setVisible` calls at
each configuration site.

### Quick filter

```java
public final void setQuickFilter(SerializableFunction<T, String>... fieldExtractors) { /* ... */ }
```

Case-insensitive substring match across all extractors. Visibility binds
to "at least one extractor registered."

### Filter popover

The filter icon button and its active-count badge sit in a wrapper `Div`
(Vaadin `Button#add` is slot-private; the badge can't be a button
child). The wrapper is the popover target; clicks on either bubble up.
The badge binds its text and visibility to an `activeFilterCount`
`ValueSignal<Long>` maintained incrementally — ±1 on each filter
transition, not a full recount.

### `addCustomFilter`

```java
public <V, C extends Component & HasValue<?, V>> CustomFilter<V> addCustomFilter(
        String key, String label, C input, V defaultValue,
        SerializableBiPredicate<V, T> matches) { /* ... */ }
```

- `key` — stable identifier for save/restore.
- `defaultValue` — the "no filter" value; Clear All resets to it. Pass
  `null` when the component's natural empty value already means no
  filter.
- Returns the `CustomFilter` so `activeOnlyWhen` can be chained for
  compound cases (e.g., a range filter where any non-null bound counts
  as active regardless of the other).

The filter input is any component that also implements `HasValue` —
`Select`, `ComboBox`, `Checkbox`, `DatePicker`, a `CustomField`
wrapping multiple sub-fields, or anything else.

### Row styling

Five built-in tones via CSS `part` names on the grid row:

```java
browser.muteRowsWhen(item -> !item.isActive());          // secondary color, italic
browser.markRowsAsPrimaryWhen(item -> item.isCurrent()); // primary tint, semibold
browser.markRowsAsSuccessWhen(item -> item.isHealthy());
browser.markRowsAsWarningWhen(item -> item.isExpiring());
browser.markRowsAsErrorWhen(item -> item.isFailed());
```

Tones compose: multiple matching predicates space-join their part names,
so an item can be both muted and warning simultaneously. CSS rules live
in `frontend/itembrowser/item-browser.css`, bundled with `ItemBrowser`
via `@CssImport`. To add a reusable tone, add a
`vaadin-grid::part(item-browser-<tone>)` rule and a convenience method
that calls `markRowsWithCustomPartWhen`.

### Record count

```java
browser.setRecordCountVisible(true);
```

Opt-in footer showing the count of rows currently passing all active
filters. Updates reactively via a `dataProvider` listener.

### Saved filters

```java
browser.enableSavedFilters();
```

Adds "Save current filter as…" and "Recall filter…" items to the
overflow menu. Storage is in-session, in-memory. Snapshot/restore
methods (`filterSnapshot()`, `restoreFilterSnapshot(Map)`) are public
for callers that want to persist snapshots externally (localStorage,
server-side per-user presets).

## Step 2 — `FilterOption<V>` and filter option enums

`FilterOption<V>` is a thin interface for enum-based filter options:

```java
public interface FilterOption<V> {
    SerializablePredicate<V> match();

    static <V, T, F extends FilterOption<V>> SerializableBiPredicate<F, T> matches(
            SerializableFunction<T, V> accessor) {
        return (option, item) -> option.match().test(accessor.apply(item));
    }
}
```

The static `matches(accessor)` helper bridges from "option has a
predicate on a value of `V`" to the `SerializableBiPredicate<F, T>`
shape `addCustomFilter` requires — the caller supplies the item-to-value
accessor; the helper composes.

Implement it as an enum that also implements `HasCaption` (for
`setItemLabelGenerator`):

```java
public enum ActiveFilterOption implements HasCaption, FilterOption<Boolean> {
    ANY("Any", _ -> true),
    ACTIVE("Active",   active -> active),
    INACTIVE("Inactive", active -> !active);

    private final String caption;
    private final SerializablePredicate<Boolean> match;

    ActiveFilterOption(String caption, SerializablePredicate<Boolean> match) {
        this.caption = caption;
        this.match = match;
    }

    @Override public String getCaption() { return caption; }
    @Override public SerializablePredicate<Boolean> match() { return match; }

    // Convenience — wraps the accessor for HasActive items
    public static <T extends HasActive> SerializableBiPredicate<ActiveFilterOption, T> matches() {
        return FilterOption.matches(HasActive::active);
    }
}
```

Wiring in the view:

```java
var statusSelect = new Select<>(ActiveFilterOption.values());
statusSelect.setItemLabelGenerator(ActiveFilterOption::getCaption);
statusSelect.setValue(ActiveFilterOption.ANY);
browser.addCustomFilter("status", "Status", statusSelect,
        ActiveFilterOption.ANY, ActiveFilterOption.matches());
```

The `ANY` constant is the `defaultValue` — Clear All resets the select
to it, and the default `activeOnlyWhen` predicate treats any deviation
from `ANY` as an active filter without needing `activeOnlyWhen`.

## Step 3 — `ItemEditor<T>`

`ItemEditor<T>` is the contract for the detail-form component hosted
inside any `EditableItemBrowser`:

```java
public interface ItemEditor<T> {
    Component getContent();               // called once; host caches and reuses
    String getTitle(T item);              // shown in detail chrome (VIEW + EDIT)
    default String getCreateTitle() { return "New entry"; } // override with entity name
    void show(T item, EditableItemBrowser.Mode mode); // bind + apply read-only/editable
    default boolean hasUnsavedChanges() { return false; }
    default void revert() { }
    default boolean trySave() { return true; }
    default Component leadingFooter() { return null; }
}
```

### Critical contract: reuse, don't rebuild

`getContent()` is called **once** by the host when the first item is
opened; the same component instance is reused for every subsequent open.
`show(item, mode)` is called on every open and on every mode flip — it
rebinds the binder's bean, not the form's component tree. Never
construct new form fields in `show()`.

### `trySave()` — editor-owned persistence

When an editor owns its own service call (the common case):

```java
@Override
public boolean trySave() {
    if (!binder.validate().isOk()) {
        return false;          // validation failed; host stays in EDIT/CREATE
    }
    try {
        tenantService.save(binder.getBean());
        return true;           // host proceeds: close (CREATE) or flip to VIEW (EDIT)
    } catch (ValidationException e) {
        e.getErrors().forEach(msg ->
                Notification.show(msg).addThemeVariants(NotificationVariant.ERROR));  // v25.1+; use LUMO_ERROR on v24/v25.0
        return false;
    }
}
```

Returning `false` tells the host to stay open. The editor is responsible
for surfacing the error; the host does nothing further.

### `leadingFooter()` — state-toggle actions

State-change actions (Activate, Deactivate) that should remain reachable
in both VIEW and EDIT modes belong in `leadingFooter()`:

```java
@Override
public Component leadingFooter() {
    return activateButton; // kept up-to-date in show()
}
```

The host hides the leading footer slot in CREATE mode (no persisted
entity yet).

### `hasUnsavedChanges()` and `revert()`

Implement both when the binder tracks edits, so the "discard changes?"
guard fires correctly:

```java
@Override public boolean hasUnsavedChanges() { return binder.hasChanges(); }
@Override public void revert()               { binder.readBean(binder.getBean()); }
```

## Step 4 — `EditableItemBrowser<T>`

The interface all CRUD layouts implement:

```java
public interface EditableItemBrowser<T> {

    enum Mode { VIEW, EDIT, CREATE }

    void setEditor(ItemEditor<T> editor);
    void setEditable(SerializablePredicate<T> canEdit);   // default: all forbidden
    void setOnSave(SerializableConsumer<T> onSave);       // called after trySave() returns true
    void openDetail(T item, Mode mode);
    void closeDetail();
    void closeDetailIfOpenOn(SerializablePredicate<T> predicate);
    void setNewEntityButton(String label, SerializableSupplier<T> newEntityFactory);
}
```

`setEditable` defaults to "nothing is editable." Call it explicitly;
forgetting it silently removes all edit affordances with no error.

`closeDetailIfOpenOn` is for state-change actions — pass a key-match
predicate so a mutation on row B doesn't disturb a detail open on row A:

```java
browser.closeDetailIfOpenOn(open -> open.key() == mutatedItem.key());
browser.setItems(refreshedList);
```

## Step 5 — `MasterDetailItemBrowser<T>`

The built implementation. Extends
`Composite<MasterDetailLayout>`, composes an `ItemBrowser`, and
delegates every list-configuration method to it:

```java
var browser = new MasterDetailItemBrowser<Tenant>();
browser.addColumn(Tenant::name).setHeader("Name");
browser.addColumn(Tenant::status).setHeader("Status");
browser.setQuickFilter(Tenant::name, Tenant::code);
browser.addCustomFilter("status", "Status", statusSelect,
        ActiveFilterOption.ANY, ActiveFilterOption.matches());
browser.muteRowsWhen(t -> !t.isActive());
browser.setEditor(new TenantEditor(tenantService));
browser.setEditable(t -> currentUser.isAdmin());
browser.setNewEntityButton("New Tenant", Tenant::new);
browser.setItems(tenantService.listAll());
```

### Mode transitions

| User action | Transition |
|---|---|
| Row click | → VIEW (safe default; avoids inadvertent edits) |
| Edit pencil in chrome | VIEW → EDIT in place (no grid round-trip) |
| Edit action in action column | → EDIT directly |
| New entity button | → CREATE (factory invoked on each click; fresh bean) |
| Save (CREATE) | persist → **close detail** (new entity has no display state yet) |
| Save (EDIT) | persist → **VIEW** (stay on the same row) |
| Cancel (CREATE) | close detail |
| Cancel (EDIT) | revert → VIEW |

### "Discard changes?" guard

Every close path runs through `tryProceed`:

- Close button, backdrop click, Escape → prompt if dirty
- Row-switch click → prompt if dirty; on "Keep editing", restore prior
  row selection
- New entity button → prompt if dirty before opening CREATE

The confirm dialog intentionally puts "Keep editing" on the primary
(`Enter`) button and "Discard" on the cancel slot with an error theme —
losing edits shouldn't be the default keyboard action.

## Step 6 — Choosing a layout

| Layout | When to use |
|---|---|
| `MasterDetailItemBrowser` | Editor fits in ~400 px width alongside the list. The standard choice for most admin entities. |
| `DialogItemBrowser` | Editor needs more horizontal space than MDL allows, but the form is still self-contained. The list remains visible behind the dimmer. |
| `ViewItemBrowser` | The edited object is large and complex enough to deserve its own view — multi-section forms, nested grids, embedded tabs. The editor is a full routed view; `ItemEditor.show()` drives navigation rather than rendering inline content. |

## Step 7 — Building `DialogItemBrowser` and `ViewItemBrowser`

Neither is built yet. Both implement `EditableItemBrowser<T>` and compose
an `ItemBrowser<T>`, delegating list-configuration methods to it — same
pattern as `MasterDetailItemBrowser`.

**`DialogItemBrowser`** wraps the `ItemEditor`'s content in a Vaadin
`Dialog`. `openDetail` opens the dialog and calls `editor.show(item,
mode)`; `closeDetail` closes it. The "discard changes?" guard and Cancel
/ Save footer follow the same logic as `MasterDetailItemBrowser`.

**`ViewItemBrowser`** navigates to the editor view rather than rendering
inline content. `openDetail` calls `ui.navigate(EditorView.class,
item.key())`. The editor view is a separate `@Route` that fetches the
item from the service and calls `editor.show()` in its `setParameter`
method. `closeDetail` navigates back. Because `ItemEditor.getContent()`
is meaningless in this layout, the interface contract relaxes: the
`ItemEditor` implementation is optional; the routed view itself is the
"editor."

## Decisions this recipe imposes

- **`ItemBrowser` is the list; layout implementations are wrappers.**
  `MasterDetailItemBrowser` doesn't subclass `ItemBrowser` — it
  composes one. This lets `DialogItemBrowser` and `ViewItemBrowser` use
  the same list without inheriting the MDL-specific detail wiring.
- **`ItemEditor.getContent()` is called once; `show()` rebinds.**
  Rebuilding the form on every open is expensive and breaks focus
  management. Constructing once and rebinding via `Binder.setBean()` in
  `show()` is the only correct approach.
- **`trySave()` returning false keeps the host open.** The editor
  surfaces the error (notification, inline field error); the host stays
  in EDIT/CREATE. Returning true surrenders control back to the host.
- **`setEditable` defaults to no editing.** A forgotten `setEditable`
  call silently removes all edit affordances. This is safer than
  defaulting to "all editable" — a list view that inadvertently allows
  edits is a worse bug than one that inadvertently hides them.
- **Signals drive toolbar and detail visibility.** Adding a feature
  (quick filter, new-entity button) changes a signal; the component
  binds to that signal. No imperative `setVisible` calls scattered at
  each configuration site.
- **Filter enums implement both `HasCaption` and `FilterOption<V>`.**
  `HasCaption` gives `setItemLabelGenerator` a uniform method reference.
  `FilterOption<V>.match()` gives `addCustomFilter` the predicate.
  `FilterOption.matches(accessor)` composes both without caller-side
  lambdas.

## What to verify

- A standalone `ItemBrowser` (no `EditableItemBrowser`) with quick
  filter and two custom filters: toolbar appears; active-filter badge
  increments / decrements on filter change; Clear All resets badge to 0.
- Saved-filter round-trip: save a named filter, recall it, verify the
  select/checkbox values restore correctly.
- Row tones: a muted row and a warning row on the same grid; both parts
  appear in the row's `part` attribute.
- `MasterDetailItemBrowser` CREATE → save: detail closes, list refreshes
  with the new entity.
- `MasterDetailItemBrowser` EDIT → dirty → row click: "Discard changes?"
  dialog fires; "Keep editing" restores the prior selection; "Discard"
  opens the new row in VIEW.
- `trySave()` returning false: detail stays open in EDIT; the error
  notification fires; save button remains active.
- `closeDetailIfOpenOn` with a matching predicate closes; with a
  non-matching predicate leaves the detail open.

## Related

- `ui/vaadin/recipes/base-view.md` — `BaseView` chrome that hosts
  `ItemBrowser` via `setContent(browser)`.
- `ui/vaadin/app-icon.md` — `AppIcon` used for toolbar and action
  icons (`FILTER`, `SEARCH`, `OVERFLOW_MENU`, `EDIT`, `REMOVE`).
- `ui/vaadin/uimodel/has-caption.md` — `HasCaption` for filter-option
  enum display; `ComboBox` / `Select` item-label-generator convention;
  UI model capability interfaces (`HasActive`, `HasRole`, …)
  that `muteRowsWhen` and `FilterOption.matches()` bind against.
- `structure/services/service-errors.md` — service error contracts
  (`ValidationException`, `EntityNotFoundException`); `trySave()`
  catch pattern follows `ui/vaadin/service-error-handling.md`.
