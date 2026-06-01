# AppLayout and SideNav

When building the navigation shell, use `AppLayout` as the shell for all authenticated views — it provides the drawer, navbar, and content slot. `SideNav` populates the drawer.

```java
@Layout
@AnonymousAllowed  // See docs/patterns/conventions/vaadin/spring.md — do NOT use @PermitAll
public class MainLayout extends AppLayout {
    // DO NOT add "implements RouterLayout" — already inherited from AppLayout

    public MainLayout(AuthenticatedUser currentUser) {
        var nav = new SideNav();
        // Add nav items based on role — see conditional-nav.md
        addToDrawer(nav);
        addToNavbar(new DrawerToggle());
    }
}
```
