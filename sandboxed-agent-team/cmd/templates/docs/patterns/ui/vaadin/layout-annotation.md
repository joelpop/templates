# @Layout Annotation Version Compatibility

When creating the main layout class on Vaadin 24.1+, annotate it with `@Layout` so views are discovered without each `@Route` naming the layout explicitly.

```java
// Vaadin 24.1+ / 25+
@Layout
@PermitAll
public class MainLayout extends AppLayout { /* ... */ }
```

On Vaadin 24.0, `@Layout` is not available. Set the layout on each route explicitly:

```java
// Vaadin 24.0 — no @Layout; set layout per @Route
@Route(value = "items", layout = MainLayout.class)
public class ItemView extends Composite<VerticalLayout> { /* ... */ }
```

**Related:** `docs/patterns/security/authorization/layout-access.md` — access annotation
choice (`@PermitAll` vs `@AnonymousAllowed`) and the Vaadin ≥25 `anyRequest` override;
`router-layout-interface.md` — why `implements RouterLayout` need not be added.
