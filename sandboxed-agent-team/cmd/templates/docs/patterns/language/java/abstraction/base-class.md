# Shared Base Class for Structural Templates

When several classes share structure — layout, chrome, lifecycle, fixture wiring — rather than state, extract a base class so each subclass says only what differs.

Recognition signal: every implementer rebuilds the same scaffolding before doing
the part that's actually different.

```java
// Avoid — every view re-implements the same chrome
public class CustomersView extends Composite<VerticalLayout> {
    public CustomersView() {
        var title = new H2("Customers");
        title.addClassNames(/* standard view-title styling */);
        var header = new HorizontalLayout(title, /* actions slot */);
        // ... shared scaffolding repeated in every view ...
        getContent().add(header, body);
    }
}
public class OrdersView extends Composite<VerticalLayout> { /* same shape, copied */ }
```

```java
// Preferred — chrome lives in a base; subclasses say only what differs
public class BaseView extends Composite<VerticalLayout> {
    private final HorizontalLayout headerActions = new HorizontalLayout();
    private final VerticalLayout body = new VerticalLayout();
    protected BaseView(String title) { /* builds header + body once */ }
    protected void setHeader(Component c)  { headerActions.removeAll(); headerActions.add(c); }
    protected void setContent(Component c) { body.removeAll(); body.add(c); }
}

public class CustomersView extends BaseView {
    public CustomersView() { super("Customers"); setContent(new CustomersGrid()); }
}
```

Wrong-size signs:

- The "base" defines a single hook and nothing else → callers are still doing the work.
- Subclasses override most methods → the shared shape isn't real.
- A new feature requires changing the base to accommodate one caller → the base is leaking caller-specific knowledge.
