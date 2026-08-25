# Recipe: Application Icon Catalog — `AppIcon` and `IconFactory` Adapters

When implementing an intent-named icon catalog that decouples icon *meaning*
from icon *appearance* and lets any icon library plug in through Vaadin's
`IconFactory` interface, follow this recipe to produce an `AppIcon` enum plus
per-library adapter enums callable as `AppIcon.EDIT.create()` everywhere.

## What this produces

- An `AppIcon` enum implementing `IconFactory` — the single catalog of
  icons named by application intent (not by visual identity or library
  source).
- An adapter enum (e.g., `UntitledUiIcon`) implementing `IconFactory`
  for each third-party icon library beyond Vaadin's built-in set,
  shipped as a Vaadin custom iconset JS module.
- Uniform call sites: `AppIcon.EDIT.create()` everywhere, regardless of
  which backing library any given constant uses.

## Dependencies

- Vaadin 24+ (`com.vaadin.flow.component.icon.IconFactory`,
  `com.vaadin.flow.component.icon.Icon`)
- One JS iconset file per third-party library (e.g.,
  `untitled-ui-iconset.js`) placed under `src/main/frontend/icons/` —
  see Vaadin's "Custom Icon Sets" documentation for the
  `<vaadin-iconset>` file format.

## Step 1 — Define `AppIcon`

The catalog enum lives in the UI module (`{app}-ui`). Each constant
names an *intent*, not a visual. Two constants can delegate to the same
icon today and diverge independently later — there is no cost to having
`IMPERSONATE` and `USER_MANAGEMENT` both render the same person icon
until Phase 2 makes them distinct.

The enum implements `IconFactory` itself, so call sites are uniform
regardless of which library a constant delegates to:

```java
package {base_package}.ui.component.icon;

import com.vaadin.flow.component.icon.Icon;
import com.vaadin.flow.component.icon.IconFactory;
import com.vaadin.flow.component.icon.VaadinIcon;

public enum AppIcon implements IconFactory {

    // ---- Navigation (drawer) -------------------------------------------
    DASHBOARD(VaadinIcon.DASHBOARD),
    USER_MANAGEMENT(VaadinIcon.USERS),
    SYSTEM_SETTINGS(VaadinIcon.COG),

    // ---- Row / item actions --------------------------------------------
    EDIT(UntitledUiIcon.EDIT_05),
    REMOVE(UntitledUiIcon.DELETE),
    SEARCH(UntitledUiIcon.SEARCH_SM),
    CLOSE(UntitledUiIcon.X_CLOSE),
    ACTIVATE(VaadinIcon.CHECK_CIRCLE),
    DEACTIVATE(VaadinIcon.BAN),

    // ---- Error views ---------------------------------------------------
    ERROR_NOT_FOUND(VaadinIcon.SEARCH),
    ERROR_SYSTEM(VaadinIcon.WARNING);

    private final IconFactory delegate;

    AppIcon(IconFactory delegate) {
        this.delegate = delegate;
    }

    @Override
    public Icon create() {
        return delegate.create();
    }
}
```

Group constants by usage area with comments — the grouping communicates
where each icon appears without encoding location in the name.

`VaadinIcon` implements `IconFactory` natively, so built-in Vaadin icons
are passed directly as delegates. Third-party icons are passed as adapter
enum constants (Step 2). Both satisfy the same `IconFactory` parameter
type — the heterogeneous mix requires no common base class or wrapper.

## Step 2 — Define adapter enums for third-party libraries

For each icon library beyond Vaadin's built-in set, create an enum
implementing `IconFactory`. Constants match the source library's icon
naming (typically the SVG filename, minus extension) for traceability.
The `create()` method constructs an `Icon` using Vaadin's
`new Icon("collection", "icon-name")` API.

The `@JsModule` annotation causes Flow to include the iconset JS file in
the bundle whenever this enum class is referenced — no manual import
configuration needed.

```java
package {base_package}.ui.component.icon;

import com.vaadin.flow.component.dependency.JsModule;
import com.vaadin.flow.component.icon.Icon;
import com.vaadin.flow.component.icon.IconFactory;
import java.util.Locale;

@JsModule("./icons/untitled-ui-iconset.js")
public enum UntitledUiIcon implements IconFactory {

    DELETE,
    DOTS_VERTICAL,
    EDIT_05,
    FILTER_LINES,
    SEARCH_SM,
    X_CLOSE;
    // add constants as icons are needed — keep the list to what AppIcon uses

    @Override
    public Icon create() {
        return new Icon("untitled-ui",
                name().toLowerCase(Locale.ENGLISH).replace('_', '-'));
    }
}
```

The convention for `create()`:
- The collection name (`"untitled-ui"`) matches the name registered in
  the JS iconset file.
- The icon name is the constant name lowercased and underscores replaced
  with hyphens — `EDIT_05` → `"edit-05"`. This matches the convention
  used by Vaadin's own iconsets and makes the mapping predictable without
  a lookup table.

The JS file (`src/main/frontend/icons/untitled-ui-iconset.js`) registers
a `<vaadin-iconset>` element under the name `"untitled-ui"` containing
inline SVG sprites. See Vaadin's "Custom Icon Sets" documentation for the
exact format. Once the file is in place, `new Icon("untitled-ui",
"edit-05")` resolves the same way `new Icon("vaadin", "edit")` does.

## Step 3 — Call sites

Every place in the application that needs an icon calls
`AppIcon.<INTENT>.create()`:

```java
// Toolbar search prefix
searchField.setPrefixComponent(AppIcon.SEARCH.create());

// Grid row action
editButton.setIcon(AppIcon.EDIT.create());

// Error view
errorIcon.add(AppIcon.ERROR_NOT_FOUND.create());
```

No call site needs to know which library backs the icon or what its
collection / icon-name strings are. Changing `AppIcon.EDIT`'s delegate
from `UntitledUiIcon.EDIT_05` to some other factory updates every call
site automatically.

## Decisions this recipe imposes

- **Named by intent, not appearance.** `AppIcon.IMPERSONATE` and
  `AppIcon.USER_MANAGEMENT` can share the same visual today and diverge
  later — one line change each. A catalog named by visual (`PERSON_ICON`,
  `USER_ICON`) obscures which constant to change.
- **Adapters are also named by source, not by intent.** `UntitledUiIcon`
  constant names match the upstream library (e.g., `EDIT_05`) for
  traceability to the asset. Intent naming lives only on `AppIcon`.
- **`AppIcon` is in the UI module, not `{app}-uimodel`.** It creates
  `Icon` components (a Vaadin type); it belongs where Vaadin dependencies
  live. `{app}-uimodel` must not import Vaadin.
- **One adapter enum per library.** Do not mix icons from multiple
  libraries into a single enum — the collection name in `create()` would
  vary per constant, breaking the uniform construction pattern.
- **Keep the adapter enum list to what `AppIcon` uses.** Don't add every
  icon in the source library — add on demand. The adapter list is a
  subset.

## What to verify

- `AppIcon.EDIT.create()` returns a non-null `Icon`; inspecting it in
  the browser shows the correct visual.
- `AppIcon.DASHBOARD.create()` (a `VaadinIcon` delegate) renders
  correctly alongside `AppIcon.EDIT.create()` (an `UntitledUiIcon`
  delegate) — heterogeneous delegates coexist.
- Changing one `AppIcon` constant's delegate to a different `IconFactory`
  updates that intent's visual everywhere without touching other
  constants.
- Removing an unused `UntitledUiIcon` constant causes a compile error
  on the referencing `AppIcon` constant — not a silent runtime failure.
- The JS iconset file is included in the bundle (check browser dev tools
  → network or `vaadin-dev-server` logs); `new Icon("untitled-ui",
  "edit-05")` resolves without a console error.

## Related

- `ui/vaadin/recipes/view-icon.md` — `@ViewIcon` annotation: type-safe icon
  declaration on view classes, backed by `AppIcon`. Eliminates
  the string literal in `@Menu(icon = "vaadin:xxx")`.
- `ui/vaadin/navigation/menu-annotation.md` — `@Menu` and sidebar navigation;
  explains the string-literal icon problem that `@ViewIcon` addresses.
- `ui/vaadin/uimodel/has-caption.md` — `HasCaption` for display enums;
  a related pattern for enum display in selection components.