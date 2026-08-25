# @Menu Annotation for Routed Views

When registering navigation entries on Vaadin ≥24.4, annotate each `@Route` class
with `@Menu` so navigation entries are derived declaratively from routes rather than
registered manually in `MainLayout`.

```java
@Route("employees")
@Menu(order = 1, title = "Employees")
@ViewIcon(AppIcon.EMPLOYEES)
@RolesAllowed(UserRole.ROLE_USER)
public class EmployeeView extends Composite<VerticalLayout> { /* ... */ }
```

`@Menu` attributes:
- `order` — sort position in the navigation menu
- `title` — display label (defaults to the class name if omitted)
- `icon` — do not use; declare the icon via `@ViewIcon` instead

In `MainLayout`, use `MenuConfiguration.getMenuEntries()` to populate the `SideNav`,
reading the icon from `@ViewIcon` on the view class:

```java
var nav = new SideNav();
MenuConfiguration.getMenuEntries().forEach(entry -> {
    var item = new SideNavItem(entry.title(), entry.path());
    var viewIcon = entry.navigationTarget().getAnnotation(ViewIcon.class);
    if (viewIcon != null) {
        item.setPrefixComponent(viewIcon.value().create());
    }
    nav.addItem(item);
});
addToDrawer(nav);
```

## Vaadin <24.4

`@Menu` is not available. Register `SideNavItem` instances manually — see `conditional-nav.md`.

**Related:** `conditional-nav.md` — role-conditional filtering of navigation items;
`docs/patterns/ui/vaadin/recipes/view-icon.md` — `@ViewIcon` annotation and reading pattern;
`docs/patterns/ui/vaadin/recipes/app-icon.md` — `AppIcon` enum and `IconFactory` adapters.
