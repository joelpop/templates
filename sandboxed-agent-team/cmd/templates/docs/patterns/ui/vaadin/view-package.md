# Per-View Package Layout

When creating a `@Route`-annotated view, place it in its own Java package named after the view's terminal path segment so per-view companions (dialogs, editors) have a stable home that mirrors the side-nav grouping.

Every `@Route`-annotated view lives in its own Java package named after the view's
terminal path segment. View packages are grouped under a path-prefix package — e.g.
`…ui.view.admin.organization` for a view at `@Route("admin/organization")`, or
`…ui.view.platform.tenant` for `@Route("platform/tenant")`. The view class itself
and any UI that serves only that view — per-view editors, grid cell renderers, form
helpers — live in the view's package. UI shared across views stays in a top-level
`component` (or equivalent) package.

```
ui/
├── layout/
│   └── BaseView.java           <- shared view base class
├── component/                  <- UI shared across views
└── view/
    ├── admin/
    │   ├── organization/
    │   │   └── OrganizationView.java        @Route("admin/organization")
    │   └── user/
    │       └── UserView.java                @Route("admin/user")
    └── platform/
        └── tenant/
            └── TenantView.java              @Route("platform/tenant")
```
