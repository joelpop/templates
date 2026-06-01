# Active Route Highlighting

Vaadin's `SideNavItem` automatically highlights the current route — no manual highlighting code is needed when using the class-reference constructor form.

```java
// Automatic highlighting — use class reference, not string path
new SideNavItem("Items", ItemView.class, "vaadin:list")
```
