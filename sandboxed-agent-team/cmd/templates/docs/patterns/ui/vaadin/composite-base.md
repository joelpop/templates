# Views Must Extend Composite<T>

When creating a routed view or UI component, extend `Composite<T>` with an appropriate root component type rather than a layout class directly so the root layout stays encapsulated.

```java
// Avoid
public class ItemView extends VerticalLayout { ... }
```

```java
// Preferred
@Route("items")
@RolesAllowed(UserRole.ROLE_USER)
public class ItemView extends Composite<VerticalLayout> {
    // getContent() returns the root VerticalLayout
}
```

`Composite<T>` encapsulates the root layout: callers see only what you explicitly expose.
