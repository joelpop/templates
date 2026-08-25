# Per-View Package Layout

When creating a `@Route`-annotated view, place it in its own Java package named after the view's terminal path segment so per-view companions (dialogs, editors) have a stable home that mirrors the side-nav grouping.

View packages are grouped under a path-prefix package — e.g.
`…ui.view.admin.organization` for a view at `@Route("admin/organization")`. The view
class and any UI that serves only that view — editors, grid cell renderers, form helpers
— live in the view's package. UI shared across views stays in a top-level `component`
package.

```
ui/
├── layout/
│   └── BaseView.java               <- shared view base class
├── component/                      <- UI shared across views
└── view/
    ├── admin/
    │   ├── organization/
    │   │   ├── OrganizationView.java        @Route("admin/organization")
    │   │   └── OrganizationEditor.java      <- per-view companion
    │   └── user/
    │       └── UserView.java                @Route("admin/user")
    └── platform/
        └── tenant/
            └── TenantView.java              @Route("platform/tenant")
```
