# Route-Path Grouping Convention

When naming view routes, group related views under a common path prefix so access control patterns are easier to express and audit.

```
/admin/user           → User management (admin only)
/admin/settings       → System settings (admin only)
/item                 → Item list (all authenticated)
/item/:key            → Item detail (all authenticated)
```
