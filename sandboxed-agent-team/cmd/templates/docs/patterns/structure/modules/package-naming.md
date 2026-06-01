# Package Naming

When naming Java packages in each module, use this standard sub-package structure
under `{base_package}` so code is consistently located across projects.

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
{base_package}.service                 — service interfaces
{base_package}.jpaservice              — service implementations
{base_package}.jpaservice.mapper       — MapStruct mappers
{base_package}.ui                      — shared UI components and MainLayout
{base_package}.ui.component            — shared reusable UI components
{base_package}.ui.view                 — views (each view in its own sub-package)
{base_package}.app.config              — application configuration
{base_package}.app.config.security     — security configuration
```

Package leaf names use the **singular** form — `…ui.view.admin.product`, not
`…ui.view.admin.products`. See `docs/patterns/conventions/singular-form.md`.
