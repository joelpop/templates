# Singular Form for Named Things

Every named artifact — Maven module, Java type, package leaf, database entity,
Vaadin route path — should use the **singular** form, not the plural.

| Artifact             | Right                                                    | Wrong                                                        |
|----------------------|----------------------------------------------------------|--------------------------------------------------------------|
| Vaadin view class    | `ProductView`, `OrganizationView`, `UserView`            | `ProductsView`, `OrganizationsView`, `UsersView`             |
| Database entity       | `ProductEntity`, `OrganizationEntity`, `EquipmentEntity` | `ProductsEntity`, `OrganizationsEntity`, `EquipmentsEntity` |
| Maven module         | `my-app-service`, `my-app-provider`                      | `my-app-services`, `my-app-providers`                        |
| Java package leaf    | `…ui.view.admin.product`                                 | `…ui.view.admin.products`                                    |
| Type / class name    | `OrganizationDetail`, `EmployeeListItem`                 | `OrganizationsDetail`, `EmployeesListItem`                   |
| Vaadin `@Route` path | `@Route("admin/product")`                                | `@Route("admin/products")`                                   |

**Rationale:** the name describes the *kind* of thing the entity / view / module deals in.
A `ProductsEntity` class doesn't represent a collection — each instance is one `ProductEntity`.
Plural names also collide awkwardly when collections appear in code: `List<ProductsView>`
reads as "views of plural products"; double-pluralizing yields nonsense like `productss`.

The singular rule extends to Vaadin `@Route` paths — the REST convention that prefers plural
collection URIs (`/api/products`) does not apply here. Vaadin is a server-side UI framework;
route paths name the *view* of a kind of thing, not a *collection resource*. Keeping route
paths singular keeps the URL, the view class (`ProductView`), the package
(`…ui.view.admin.product`), and the entity (`ProductEntity`) all reading as the same noun.
