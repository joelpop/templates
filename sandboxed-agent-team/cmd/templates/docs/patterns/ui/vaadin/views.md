# Vaadin View Conventions

When creating a routed Vaadin view, extend `Composite<T>` rather than a layout class,
place the view class in its own sub-package under `*.ui.view.<viewname>`, and declare
navigation via `@Route` plus `@Menu` so layout primitives stay encapsulated, per-view
companions (dialogs, editors) have a stable home, and the navigation tree is
discoverable without manually registering each entry.

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

`Composite<T>` encapsulates the root layout: callers see only what you explicitly expose.

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

Each view has a clear home that grows naturally with its helper classes and mirrors the
side-nav grouping on disk.

## @Menu Annotation for Navigation

> **Vaadin ≥24.4:** use the `@Menu` annotation on `@Route` classes for views that appear
> in the main navigation menu. Do not manually add menu items to the `SideNav` for routed
> views — let `MenuConfiguration.getMenuEntries()` discover them.
>
> **Vaadin <24.4:** `@Menu` is not available. Register `SideNavItem` instances manually in
> `MainLayout` and use the role-conditional rendering pattern (see
> `docs/patterns/ui/navigation.md` → Conditional Navigation Rendering).

```java
@Route("items")
@Menu(order = 1, icon = "vaadin:list")      // Vaadin ≥24.4
@RolesAllowed(UserRole.ROLE_USER)
public class ItemView extends Composite<VerticalLayout> { ... }
```
