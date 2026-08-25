# AppLayout and SideNav

When creating the application layout, use `AppLayout` as the navigation shell so
views share a consistent drawer, navbar, and content slot.

```java
@Layout
@PermitAll
public class MainLayout extends AppLayout {

    public MainLayout(AuthenticatedUser currentUser) {
        var nav = new SideNav();
        addToDrawer(nav);
        addToNavbar(new DrawerToggle());
    }
}
```

**Related:** `layout-annotation.md` — `@Layout` version compatibility;
`router-layout-interface.md` — why `implements RouterLayout` is not needed;
`conditional-nav.md` — filtering nav items by role;
`docs/patterns/security/authorization/layout-access.md` — access annotation choice and inheritance.
