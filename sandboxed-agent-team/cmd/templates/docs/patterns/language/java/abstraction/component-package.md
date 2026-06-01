# Component-Family Package for Multi-Class Capabilities

When N classes together realize one capability the requirements name as a unit, the abstraction is a package — not a class. Each class is a partial answer; the capability is the package.

Recognition signal: the package boundary maps to a requirement boundary, and the
public face is one or two types — the rest are package-private.

```
// Avoid — three list-view features each rebuild the same

// toolbar + quick-filter + filter-popover + grid skeleton
ui/customers/CustomersGrid.java       + toolbar, filter, popover code
ui/orders/OrdersGrid.java             + toolbar, filter, popover code
ui/products/ProductsGrid.java         + toolbar, filter, popover code
```

```
// Preferred — the capability is one package, parameterized over the

// row type; features use the package, not the parts inside
ui/component/itembrowser/
    ItemBrowser.java<T>           // public face: parameterized list view
    Toolbar.java                  // package-private; an ItemBrowser part
    FilterPopover.java            // package-private; an ItemBrowser part
    filter/
        CustomFilter.java         // public; features pass these in

ui/customers/CustomersView.java extends BaseView {
    setContent(new ItemBrowser<Customer>(/* configured for customers */));
}
```

Wrong-size signs:

- A package whose name is just the class it contains → the package adds no organizational meaning; collapse it.
  - Exception: a **category package** like `*.ui.component` is fine with one entry if the name reflects intent.
  - Exception: a **convention-driven per-item package** like `*.ui.view.useradmin` is correct even with one class today, because consistency requires every view to own its sub-package.
- Internals reused independently outside the package → either the capability isn't cohesive or you built a feature-specific subsystem where a generic mechanism belonged.
- Extending the package requires deep cross-knowledge of its internals → it's a directory, not a package.
