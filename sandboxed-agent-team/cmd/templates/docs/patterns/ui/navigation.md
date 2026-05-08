# Navigation Patterns

AppLayout, SideNav, route-path grouping, responsive navigation, and navigation guards
for Vaadin 24+ applications. Version-sensitive APIs (`@Menu`, `@Layout`, touch-optimized
bottom tab bar) carry inline "Vaadin ≥X / <X" notes.

## AppLayout + SideNav

Use `AppLayout` as the shell for all authenticated views. `AppLayout` provides the drawer,
navbar, and content slot. `SideNav` populates the drawer with navigation items.

```java
@Layout
@AnonymousAllowed  // See docs/agnostic/conventions/vaadin.md — do NOT use @PermitAll
public class MainLayout extends AppLayout {
    // DO NOT add "implements RouterLayout" — already inherited from AppLayout

    public MainLayout(AuthenticatedUser currentUser) {
        var nav = new SideNav();
        // Add nav items based on role — see Conditional Navigation below
        addToDrawer(nav);
        addToNavbar(new DrawerToggle());
    }
}
```

## @Menu Annotation for Routed Views

> **Vaadin ≥24.4:** declare navigation entries with `@Menu` on `@Route` classes rather than
> manually adding `SideNavItem` instances for every view. `MenuConfiguration.getMenuEntries()`
> provides the discovered entries for the `SideNav`.
>
> **Vaadin <24.4:** `@Menu` is not available. Register `SideNavItem` instances manually in
> `MainLayout` and use the role-conditional rendering pattern below.

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

## Route-Path Grouping Convention

Group related views under a common path prefix to signal their relationship:

```
/admin/user           → User management (admin only)
/admin/settings       → System settings (admin only)
/item                 → Item list (all authenticated)
/item/:key            → Item detail (all authenticated)
```

This makes access control patterns easier to express and audit.

## DrawerToggle and Drawer Behavior

Always include `DrawerToggle` in the navbar so users can open and close the drawer:

```java
addToNavbar(new DrawerToggle());
```

On tablet-width viewports (600px–1023px), the drawer overlays the content rather than
pushing it aside. This is Vaadin AppLayout's default responsive behavior.

## Conditional Navigation Rendering

Navigation items for routes the current user cannot access must not be generated — not
styled hidden, not disabled. They never enter the DOM:

```java
private void buildNav(SideNav nav, AuthenticatedUser user) {
    // All authenticated users
    nav.addItem(new SideNavItem("Dashboard", DashboardView.class, "vaadin:dashboard"));
    nav.addItem(new SideNavItem("Items", ItemView.class, "vaadin:list"));

    // Admin only
    if (user.hasRole(UserRole.ROLE_ADMIN)) {
        nav.addItem(new SideNavItem("User Management", UserManagementView.class, "vaadin:users"));
        nav.addItem(new SideNavItem("Settings", SettingsView.class, "vaadin:cog"));
    }
}
```

## Active Route Highlighting

Vaadin's `SideNavItem` automatically highlights the item corresponding to the current route.
No manual highlighting code is required when using `SideNavItem(String, Class<?>)` or
`SideNavItem(String, String)` constructors with the view class.

## Navigation Guards

Use `BeforeEnterObserver` on views to perform programmatic access checks beyond what
`@RolesAllowed` provides (e.g., entity-level ownership checks):

```java
@Route("items/:key/edit")
@RolesAllowed(UserRole.ROLE_ADMIN)
public class ItemEditView extends Composite<VerticalLayout> implements BeforeEnterObserver {

    @Override
    public void beforeEnter(BeforeEnterEvent event) {
        var key = event.getRouteParameters().get("key")
            .map(Long::parseLong)
            .orElse(null);
        if (key == null || !itemService.exists(key)) {
            event.rerouteToError(NotFoundException.class);
        }
    }
}
```

## Touch-Optimized Navigation on Mobile

On mobile viewports (< 600px), consider replacing or supplementing the sidebar drawer with
a bottom tab bar for the most frequently accessed views. Vaadin's `Tabs` component with
`HORIZONTAL` orientation and bottom positioning serves this purpose.

For secondary navigation within a mobile view, use an accordion pattern rather than nested
sidebar items, which are difficult to tap accurately on small screens.

## Body Scrolling

In standard `AppLayout`, body scrolling is managed by the content area, not the `<body>`
element. Ensure that views within the layout use `setSizeFull()` or explicit height
constraints appropriate for their content — otherwise the content area may not scroll
correctly on mobile.

## AfterNavigationObserver for Per-Navigation Actions

Use `AfterNavigationObserver` on `MainLayout` for actions that must run after every
successful navigation (e.g., closing the mobile drawer, updating a breadcrumb):

```java
public class MainLayout extends AppLayout implements AfterNavigationObserver {

    @Override
    public void afterNavigation(AfterNavigationEvent event) {
        // Close drawer on mobile after navigation
        setDrawerOpened(false);
    }
}
```
