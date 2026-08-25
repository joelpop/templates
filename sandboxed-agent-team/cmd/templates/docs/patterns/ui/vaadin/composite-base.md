# Extend Composite<T> for Views and Components

When creating a routed view or a reusable UI component, extend `Composite<T>` with an appropriate root component type rather than a layout class directly so the root layout stays encapsulated and callers see only what you explicitly expose.

```java
// Avoid
public class ItemView extends VerticalLayout { /* ... */ }
public class ItemCard extends HorizontalLayout { /* ... */ }
```

```java
// Preferred
@Route("items")
@RolesAllowed(UserRole.ROLE_USER)
public class ItemView extends Composite<VerticalLayout> { /* ... */ }

public class ItemCard extends Composite<HorizontalLayout> { /* ... */ }
```

The exception is when deliberately extending a specific component type — such as a `SsnField` that extends `PasswordField` and adds SSN validation.

Because `Composite<T>` encapsulates its root, the component does not automatically expose sizing or styling APIs. Implement `HasSize`, `HasStyle`, or other capability interfaces explicitly when callers need them:

```java
public class ItemCard extends Composite<HorizontalLayout> implements HasSize, HasStyle { /* ... */ }
```
