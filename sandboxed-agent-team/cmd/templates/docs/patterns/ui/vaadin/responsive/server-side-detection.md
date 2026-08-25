# Server-Side Breakpoint Detection

When a view must adjust its component structure based on viewport width — not
just its visual appearance — call `clientDetailsService.getWindowInnerWidth()`
and compare against the standard breakpoints so the server applies the correct
layout before the first render.

Prefer CSS-based responsive rules via `LumoUtility` or `@media` queries for
purely visual changes. Use server-side detection only when components must be
added, removed, or structurally reconfigured.

```java
var width = clientDetailsService.getWindowInnerWidth();
var isMediumOrWider = width >= Breakpoints.MD.minWidthPx;
```

Use `clientDetailsService.isTouchDevice()` for interaction decisions such as
rendering a touch-optimized navigation layout — orthogonal to viewport width.

**Related:** `responsive/breakpoints.md` — breakpoint pixel values;
`client-details-service.md` — the `ClientDetailsService` interface;
`client-details-impl.md` — how details are retrieved per Vaadin version.