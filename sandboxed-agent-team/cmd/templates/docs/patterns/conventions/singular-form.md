# Singular Form for Named Things

Every named artifact — Maven module, Java type, package leaf, database table,
Vaadin route path — should use the **singular** form, not the plural.

| Artifact             | Right                                                       | Wrong                                                     |
|----------------------|-------------------------------------------------------------|-----------------------------------------------------------|
| Vaadin view class    | `ProductView`, `OrganizationView`, `UserView`               | `ProductsView`, `OrganizationsView`, `UsersView`          |
| Database table       | `product`, `organization`, `equipment`                      | `products`, `organizations`, `equipments`                 |
| Maven module         | `fleet-acuity-service`, `fleet-acuity-provider`             | `fleet-acuity-services`, `fleet-acuity-providers`         |
| Java package leaf    | `…ui.view.admin.product`                                    | `…ui.view.admin.products`                                 |
| Type / class name    | `OrganizationDetail`, `EmployeeListItem`                    | `OrganizationsDetail`, `EmployeesListItem`                |
| Vaadin `@Route` path | `@Route("admin/product")`                                   | `@Route("admin/products")`                                |

**Rationale:** the name describes the *kind* of thing the table / view / module deals in.
A `products` table doesn't hold one row called "products" — it holds rows of `product`.
Plural names also collide awkwardly when collections appear in code: `List<ProductsView>`
reads as "views of plural products"; double-pluralizing yields nonsense like `productss`.

The singular rule extends to Vaadin `@Route` paths — the REST convention that prefers plural
collection URIs (`/api/products`) does not apply here. Vaadin is a server-side UI framework;
route paths name the *view* of a kind of thing, not a *collection resource*. Keeping route
paths singular keeps the URL, the view class (`ProductView`), the package
(`…ui.view.admin.product`), and the table (`product`) all reading as the same noun.