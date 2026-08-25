# Touch-Optimized Navigation on Mobile

When building navigation for phone and tablet users, replace the drawer `SideNav`
with a bottom tab bar in `AppLayout`'s touch-optimized navbar slot so primary
targets are within thumb reach — nested sidebar items are difficult to tap
accurately on small screens.

```java
// Avoid — drawer SideNav on all devices, including touch
public MainLayout(ClientDetailsService clientDetailsService) {
    var nav = new SideNav();
    // ... add items
    addToDrawer(nav);
    addToNavbar(new DrawerToggle());
}
```

```java
// Preferred — bottom tab bar on touch devices, drawer SideNav on desktop
public MainLayout(ClientDetailsService clientDetailsService) {
    if (clientDetailsService.isTouchDevice()) {
        buildMobileNav();
    } else {
        buildDesktopNav();
    }
}

private void buildMobileNav() {
    var dashboardTab = new Tab(AppIcon.DASHBOARD.create(), new Span("Dashboard"));
    var itemsTab = new Tab(AppIcon.ITEMS.create(), new Span("Items"));
    var tabs = new Tabs(dashboardTab, itemsTab);
    tabs.setOrientation(Tabs.Orientation.HORIZONTAL);
    tabs.addSelectedChangeListener(event -> {
        var selected = event.getSelectedTab();
        if (selected == dashboardTab) UI.getCurrent().navigate(DashboardView.class);
        else if (selected == itemsTab) UI.getCurrent().navigate(ItemView.class);
    });
    addToNavbar(true, tabs);
}

private void buildDesktopNav() {
    var nav = new SideNav();
    // ... add items — see conditional-nav.md
    addToDrawer(nav);
    addToNavbar(new DrawerToggle());
}
```

For secondary navigation within a mobile view, use an accordion rather than
nested `SideNavItem` entries.

**Related:** `after-navigation.md` — closing the drawer on navigation for touch devices;
`conditional-nav.md` — role-based filtering of navigation items;
`client-details-service.md` — `isTouchDevice()` via `ClientDetailsService`.