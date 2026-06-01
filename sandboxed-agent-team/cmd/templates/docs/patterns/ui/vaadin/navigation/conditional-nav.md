# Conditional Navigation Rendering

When populating the navigation drawer, generate items only for routes the current user can access — never add items and then hide or disable them; unauthorized entries must never enter the DOM.

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
