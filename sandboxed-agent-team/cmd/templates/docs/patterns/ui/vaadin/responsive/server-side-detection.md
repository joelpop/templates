# Server-Side Breakpoint Detection

When server-side layout decisions depend on viewport width, use Vaadin's `Page` API — but prefer CSS-based responsive rules for visual-only changes.

```java
UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> {
    var isMobile = details.getWindowInnerWidth() < 600;
    // adjust layout server-side
});
```

Use sparingly — prefer CSS-based responsive rules via `LumoUtility` or `@media` queries
for visual-only changes.

The `retrieveExtendedClientDetails` API differs between Vaadin 24 and 25 — see
`docs/patterns/README.md` → "Version Compatibility".
