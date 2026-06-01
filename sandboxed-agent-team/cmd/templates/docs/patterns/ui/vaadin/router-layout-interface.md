# MainLayout Must Not Implement RouterLayout

When creating the `MainLayout` class, do not add `implements RouterLayout` — `AppLayout` already implements it, and re-declaring it can cause unexpected behavior especially in Vaadin ≥25.

```java
// Avoid — redundant and potentially harmful in Vaadin ≥25
public class MainLayout extends AppLayout implements RouterLayout { ... }
```

```java
// Preferred
public class MainLayout extends AppLayout { ... }
```
