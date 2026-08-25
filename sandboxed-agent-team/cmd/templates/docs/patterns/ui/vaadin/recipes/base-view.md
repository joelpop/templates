# Recipe: `BaseView` — Consistent Per-View Chrome

When implementing a per-view chrome shared by every routed view (icon, title,
right-side action slot, body slot) without involving `MainLayout`, follow this
recipe to produce a `BaseView` base class that subclasses populate by supplying
only what makes them distinct.

## What this produces

- A `BaseView` extending `Composite<VerticalLayout>` with a header
  (icon + title left, actions right) separated from a body by a bottom
  border.
- Two protected mutators — `setHeader(Component)` and
  `setContent(Component)` — that replace the respective slot contents.
- Automatic `@ViewIcon` integration: if the subclass carries the
  annotation, its icon is rendered next to the title with no subclass
  code required.

## Dependencies

- Vaadin 24+ (`Composite`, `HorizontalLayout`, `VerticalLayout`,
  `LumoUtility`).
- The [view-icon recipe](view-icon.md) — `@ViewIcon` is read
  reflectively by `BaseView`. Omit the annotation block if the project
  does not use `@ViewIcon`.

## Step 1 — Extend `Composite`, not `VerticalLayout`

`Composite<VerticalLayout>` wraps a `VerticalLayout` as its root content
without exposing that layout's mutation API to subclasses. This is the
correct Vaadin pattern for a view base class: subclasses call
`setHeader()` and `setContent()`; they cannot bypass the intended
structure by calling `add()` or `remove()` directly on the root layout.
`HasSize` lets external callers (e.g., `MainLayout`) size the view
without needing to know its internal layout type.

```java
package {base_package}.ui.layout;

import com.vaadin.flow.component.Component;
import com.vaadin.flow.component.Composite;
import com.vaadin.flow.component.HasSize;
import com.vaadin.flow.component.html.H2;
import com.vaadin.flow.component.orderedlayout.FlexComponent;
import com.vaadin.flow.component.orderedlayout.HorizontalLayout;
import com.vaadin.flow.component.orderedlayout.VerticalLayout;
import com.vaadin.flow.theme.lumo.LumoUtility;
import {base_package}.ui.component.icon.ViewIcon;

public class BaseView extends Composite<VerticalLayout> implements HasSize {

    private final HorizontalLayout headerActions = new HorizontalLayout();
    private final VerticalLayout body = new VerticalLayout();

    protected BaseView(String title) {
        var viewTitle = new H2(title);
        viewTitle.addClassNames(
                LumoUtility.FontSize.XLARGE,
                LumoUtility.FontWeight.SEMIBOLD,
                LumoUtility.TextColor.HEADER,
                LumoUtility.Margin.NONE);

        var titleGroup = new HorizontalLayout();
        titleGroup.setAlignItems(FlexComponent.Alignment.CENTER);
        titleGroup.addClassNames(LumoUtility.Gap.SMALL);

        var viewIcon = getClass().getAnnotation(ViewIcon.class);
        if (viewIcon != null) {
            var icon = viewIcon.value().create();
            icon.addClassName(LumoUtility.IconSize.MEDIUM);
            titleGroup.add(icon);
        }
        titleGroup.add(viewTitle);

        headerActions.setAlignItems(FlexComponent.Alignment.CENTER);

        var header = new HorizontalLayout(titleGroup, headerActions);
        header.setWidthFull();
        header.setPadding(true);
        header.setAlignItems(FlexComponent.Alignment.CENTER);
        header.setJustifyContentMode(FlexComponent.JustifyContentMode.BETWEEN);
        header.addClassNames(
                LumoUtility.Border.BOTTOM,
                LumoUtility.BorderColor.CONTRAST_10);

        body.setSizeFull();
        body.setPadding(true);
        body.setSpacing(true);

        var content = getContent();
        content.setSizeFull();
        content.setPadding(false);
        content.setSpacing(false);
        content.add(header, body);
    }

    protected void setHeader(Component component) {
        headerActions.removeAll();
        headerActions.add(component);
    }

    protected void setContent(Component component) {
        body.removeAll();
        body.add(component);
    }
}
```

`getClass().getAnnotation(ViewIcon.class)` is called on the *subclass*
at construction time — the annotation on `BaseView` itself is never
present; the check is effectively "does this concrete view have an icon?"
See [view-icon.md](view-icon.md) for the annotation and its value type.

## Step 2 — Use in subclasses

A subclass calls `super(title)` and then `setContent()` with its primary
component. `setHeader()` is optional — omit it when the view has no
toolbar actions.

```java
@Route("users")
@Menu(order = 2, title = "Users")
@RolesAllowed(UserRole.ROLE_ADMIN)
@ViewIcon(AppIcon.USER_MANAGEMENT)
public class UserManagementView extends BaseView {

    private final Grid<User> grid;
    private final TextField quickFilter;

    public UserManagementView(UserService userService) {
        super("Users");

        grid = new Grid<>(User.class);
        grid.setItems(query -> userService.list(query, quickFilter.getValue()));

        quickFilter = new TextField();
        quickFilter.setPlaceholder("Search…");
        quickFilter.addValueChangeListener(e -> grid.getDataProvider().refreshAll());

        setHeader(quickFilter);
        setContent(grid);
    }
}
```

The header and content are wired together inside the view. The icon
appears next to "Users" in the subheader automatically — the same icon
used by the drawer nav entry.

## Decisions this recipe imposes

- **`Composite<VerticalLayout>` over `extends VerticalLayout`.** Hiding
  the root layout prevents subclasses from bypassing `setContent()` and
  adding directly to the outer layout. It also prevents the root layout's
  API from becoming part of the subclass's public interface.
- **Single-slot, replace-not-append semantics.** `setHeader()` and
  `setContent()` call `removeAll()` before `add()`. Subclasses cannot
  accidentally layer two grids or two toolbars by calling the setters
  twice; the last call wins.
- **Header chrome is not in `MainLayout`.** `MainLayout` handles the
  drawer and the overall page shell. Per-view title, icon, and action
  slot belong in the view — they change with navigation; `MainLayout`
  does not.
- **`@ViewIcon` is optional.** `BaseView` reads it when present and
  ignores it when absent. Views that have no meaningful icon simply omit
  the annotation; `BaseView` does not require it.
- **Title is a constructor parameter, not a setter.** The title is
  structural chrome, set once at construction and not expected to change.
  If a view needs a dynamic title, override the `H2` through a subclass
  field, not by adding a setter to `BaseView`.

## What to verify

- A view without `@ViewIcon` renders its title with no icon and no
  layout gap where the icon would have been.
- A view with `@ViewIcon` renders the icon at medium Lumo size to the
  left of the title, vertically centred.
- The header bottom border is visible; the body fills the remaining
  height (`setSizeFull()` on the outer layout propagates correctly when
  `MainLayout` calls `setSizeFull()` on the view).
- Calling `setContent()` twice from a subclass (e.g., to swap content on
  user action) removes the previous content and renders only the new
  component.
- A view with `setHeader()` aligns the actions flush right, vertically
  centred with the title.

## Related

- `ui/vaadin/recipes/view-icon.md` — `@ViewIcon` annotation read by
  `BaseView` to render the in-header icon.
- `ui/vaadin/app-icon.md` — `AppIcon` enum; the value type carried
  by `@ViewIcon`.
- `language/java/abstraction/base-class.md` — when a shared base class
  is the right abstraction and when it is not.
- `ui/vaadin/navigation/app-layout.md` — `MainLayout`, `@Menu`, and the
  drawer navigation that shares the same icon via `@ViewIcon`.