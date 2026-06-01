# @Menu Annotation for Routed Views

When registering navigation entries, use `@Menu` on `@Route` classes (Vaadin ≥24.4) rather than manually adding `SideNavItem` instances so navigation entries are derived declaratively from routes.

> **Vaadin ≥24.4:** declare navigation entries with `@Menu` on `@Route` classes rather than
> manually adding `SideNavItem` instances for every view. `MenuConfiguration.getMenuEntries()`
> provides the discovered entries for the `SideNav`.
>
> **Vaadin <24.4:** `@Menu` is not available. Register `SideNavItem` instances manually in
> `MainLayout` and use the role-conditional rendering pattern (see conditional-nav.md).

```java
// Vaadin ≥24.4
@Route("employees")
@Menu(order = 1, icon = "vaadin:group")
@RolesAllowed(UserRole.ROLE_USER)
public class EmployeeView extends Composite<VerticalLayout> { ... }
```

`@Menu` provides:
- `order` — sort position in the navigation menu
- `icon` — Vaadin icon or Lumo icon name
- `title` — display label (defaults to the class name if omitted)
