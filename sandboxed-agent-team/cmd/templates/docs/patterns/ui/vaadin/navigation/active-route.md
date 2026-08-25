# Active Route Highlighting

When adding navigation items, use the class-reference `SideNavItem` constructor
so active-route highlighting works automatically without manual highlighting code.

```java
// Avoid — string path, no automatic highlighting
new SideNavItem("Items", "items", "vaadin:list")
```

```java
// Preferred — class reference, automatic highlighting
new SideNavItem("Items", ItemView.class, "vaadin:list")
```
