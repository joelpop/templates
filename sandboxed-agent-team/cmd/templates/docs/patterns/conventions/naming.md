# Naming Conventions

Consistent naming across all layers of a Vaadin + Spring Boot + Spring Data JPA project.
These conventions make every class's role immediately clear from its name alone.

## Java Class Naming by Layer

| Type | Convention | Example |
|------|------------|---------|
| JPA entity | Suffix `Entity` | `EmployeeEntity`, `DepartmentEntity` |
| JPA entity enum | Suffix `Code` | `EmploymentStatusCode`, `PhoneTypeCode` |
| Spring Data repository | Suffix `Repository` | `EmployeeRepository` |
| JPA interface projection | Suffix `Projection`; prefix matches the UI model it supplies | `EmployeeListItemProjection`, `EmployeeDetailProjection` |
| Service interface | Suffix `Service` | `EmployeeService` |
| Service implementation | Prefix with technology, suffix `Service` | `JpaEmployeeService` |
| UI model POJO | No suffix — named for its UI context | `EmployeeListItem`, `EmployeeDetail` |
| UI model enum | No suffix | `EmploymentStatus`, `PhoneType` |

**Rationale:**
- `Entity` suffix signals a JPA-managed object with lifecycle implications (proxies, lazy loading, dirty checking).
- `Code` suffix on JPA enums distinguishes stored values from their UI counterparts. The suffix signals "this value is persisted as a code string."
- `Projection` suffix signals a read-only, partial view of data — no persistence lifecycle.
- `Service` on both interface and implementation makes the contract name stable for callers. The implementation prefix (`Jpa`) names the backing technology, making it easy to add alternatives (`Mock`, `Cache`).
- UI models are named for the screen or operation they serve, not their data source. They carry no persistence semantics.

## Singular Form for Named Things

Every named thing — Maven module, Vaadin view class, Java type, package
leaf, database table — uses the **singular** form, not the plural.

| What | Right | Wrong |
|------|-------|-------|
| Vaadin view class | `ProductView`, `OrganizationView`, `UserView` | `ProductsView`, `OrganizationsView`, `UsersView` |
| Database table | `product`, `organization`, `equipment` | `products`, `organizations`, `equipments` |
| Maven module | `fleet-acuity-service`, `fleet-acuity-provider` | `fleet-acuity-services`, `fleet-acuity-providers` |
| Java package leaf | `…ui.view.admin.product` | `…ui.view.admin.products` |
| Type / class name | `OrganizationDetail`, `EmployeeListItem` | `OrganizationsDetail`, `EmployeesListItem` |

**Rationale:** the name describes the *kind* of thing the table / view /
module deals in. A `products` table doesn't hold one row called
"products" — it holds rows of `product`. Plural names also collide
awkwardly when collections of the kind appear in code: `List<ProductsView>`
reads as "views of plural products"; double-pluralizing yields nonsense
like `productss`.

The singular rule extends to Vaadin `@Route` paths — `@Route("admin/product")`,
not `@Route("admin/products")`. The REST collection-URI tradition that
prefers plural (`/api/products`) does not apply: Vaadin is a server-side
UI framework, not a REST API; route paths name the *view* of a kind of
thing, not a *collection resource*. Keeping route paths singular keeps
the URL, the view class (`ProductView`), the package
(`…ui.view.admin.product`), and the table (`product`) all reading as
the same noun.

## Database Column Naming

| Convention | Meaning | Examples |
|------------|---------|---------|
| `_key` suffix | Surrogate primary and foreign keys — system-generated, no business meaning | `employee_key`, `department_key`, `manager_key` |
| `_id` suffix | Business identifiers — human-meaningful values | `display_id`, `employee_id` (externally assigned) |

**This is a critical distinction.** Every entity has a `key` (surrogate PK). Some entities also
have an `_id` field that humans use to identify records. Never use `_id` for a surrogate
primary key — that suffix is reserved for business identifiers.

## Package Naming

All packages follow the base package `{base_package}.*`. Sub-packages by concern:

```
{base_package}                         — root; Application.java lives here
{base_package}.common.util             — shared utilities
{base_package}.jpamodel.entity         — JPA entity classes
{base_package}.jpamodel.code           — JPA entity enums (*Code suffix)
{base_package}.jpamodel.projection     — JPA interface projections (*Projection suffix)
{base_package}.jpaclient.repository    — Spring Data repositories
{base_package}.jpaclient.config        — JpaConfig (@EntityScan, @EnableJpaRepositories)
{base_package}.uimodel.data            — UI model POJOs (no suffix)
{base_package}.uimodel.type            — UI model enums (no suffix)
{base_package}.service                 — Service interfaces
{base_package}.jpaservice              — Service implementations
{base_package}.jpaservice.mapper       — MapStruct mappers
{base_package}.ui                      — Shared UI components and MainLayout
{base_package}.ui.component            — Shared reusable UI components
{base_package}.ui.view                 — Views (each view in its own sub-package)
{base_package}.app.config              — Application configuration
{base_package}.app.config.security     — Security configuration
```

## Projection and UI Model Naming

JPA interface projections are named for the UI context they serve — not generic
suffixes like `Summary` or `Info`. Each projection pairs one-to-one with a
similarly-named UI model POJO: the projection carries the fields; the UI model is
the shape the service layer returns. Keeping the names symmetric makes it obvious
which projection feeds which model.

| UI context | JPA interface projection | UI model POJO |
|------------|--------------------------|----------------|
| Grid row | `EmployeeListItemProjection` | `EmployeeListItem` |
| Edit form / detail view | `EmployeeDetailProjection` | `EmployeeDetail` |
| (other contexts as needed — e.g., a picker dropdown) | `EmployeePickerItemProjection` | `EmployeePickerItem` |

For new UI contexts, choose a UI-contextual name for the pair — not a generic
`Summary`, `Info`, `View`, or `Data` suffix. `EmployeeListItem` /
`EmployeeListItemProjection` are self-explaining; `EmployeeSummary` is not.

## Method Naming Conventions

| Pattern | Used For | Example |
|---------|----------|---------|
| `findBy*` | Query — returns single result or Optional | `findByKey(long key)` |
| `listAll*` | Query — returns complete list | `listAll()`, `listAllActive()` |
| `count*` | Query — returns count | `countByStatus(status)` |
| `is*Available` | Query — uniqueness check | `isDisplayIdAvailable(id, excludeKey)` |
| `create*` | Mutation — insert | `createEmployee(detail)` |
| `update*` | Mutation — update | `updateEmployee(detail)` |
| `deactivate*` | Mutation — logical delete | `deactivateEmployee(key)` |
| `on{Component}{Event}` | UI event handler | `onSaveButtonClick`, `onNameFieldValueChanged` |

## Enum Value Style

Enum constants use `SCREAMING_SNAKE_CASE`. JPA (`Code`) enums store the constant name
via `@Enumerated(EnumType.STRING)`. Never use `ORDINAL` — ordinal values break if the
enum declaration order changes.
