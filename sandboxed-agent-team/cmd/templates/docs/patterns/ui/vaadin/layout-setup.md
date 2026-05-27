# Main Layout Setup

When creating the main layout class, annotate it with `@Layout` on Vaadin 24.1+
(or set `layout = MainLayout.class` on each `@Route` on Vaadin 24.0) and never
add `implements RouterLayout` so Vaadin's router finds the layout without a
redundant interface declaration.

## @Layout Annotation Version Compatibility

The `@Layout` annotation was introduced in Vaadin 24.1. It designates a class
as the application's main layout without requiring each `@Route` to name it:

```java
// Vaadin 24.1+ / 25+
@Layout
@PermitAll
public class MainLayout extends AppLayout {
    // DO NOT add "implements RouterLayout" — already inherited from AppLayout
}
```

On Vaadin 24.0, `@Layout` is not available. Set the layout on each route
explicitly:

```java
// Vaadin 24.0 — no @Layout; set layout per @Route
@Route(value = "items", layout = MainLayout.class)
public class ItemView extends Composite<VerticalLayout> { ... }
```

For the access annotation choice (`@PermitAll` vs `@AnonymousAllowed`) and the
Vaadin ≥25 `anyRequest` override that makes `@PermitAll` work, see
`docs/patterns/security/authorization/layout-access.md`.

## MainLayout Must Not Implement RouterLayout

`AppLayout` already implements `RouterLayout`. Adding `implements RouterLayout`
explicitly is redundant and can cause unexpected behavior (especially in Vaadin
≥25 where the access-checker treats layouts differently).
