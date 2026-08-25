# Route-Path Grouping Convention

When naming routes for views that belong to the same section, use a shared
path prefix so nav entries group by section automatically and `@RolesAllowed`
coverage is easy to audit across the section.

| Route             | View             | Access            |
|-------------------|------------------|-------------------|
| `/admin/user`     | User management  | Admin only        |
| `/admin/settings` | System settings  | Admin only        |
| `/item`           | Item list        | All authenticated |
| `/item/:key`      | Item detail      | All authenticated |

**Related:** `conditional-nav.md` — `firstPathSegment` grouping in `SideNav`;
`navigation-guards.md` — per-record access checks within a route;
`security/authorization/view-access.md` — `@RolesAllowed` on route classes.
