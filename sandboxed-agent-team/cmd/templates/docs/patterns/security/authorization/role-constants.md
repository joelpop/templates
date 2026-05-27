# Role Name Constants

When defining application roles, declare them as a single `UserRole` enum that
owns both the canonical security name and the human-facing metadata so role
names are referenced by constant rather than string literal and renaming is a
single-file edit.

## The Enum

```java
package com.example.uimodel.type;

/**
 * UI representation of user roles.
 */
public enum UserRole {
    ADMIN("Admin", "Full system access", UserRole.ROLE_ADMIN),
    MANAGER("Manager", "Department-level access", UserRole.ROLE_MANAGER),
    STAFF("Staff", "Day-to-day operational access", UserRole.ROLE_STAFF);

    /** Full system access — manages users, settings, and all data. */
    public static final String ROLE_ADMIN = "ADMIN";

    /** Department-level access — manages staff and reads org-wide data. */
    public static final String ROLE_MANAGER = "MANAGER";

    /** Day-to-day operational access — read most data, mutate within scope. */
    public static final String ROLE_STAFF = "STAFF";

    private final String displayName;
    private final String description;
    private final String securityName;

    UserRole(String displayName, String description, String securityName) {
        this.displayName = displayName;
        this.description = description;
        this.securityName = securityName;
    }

    public String getDisplayName() {
        return displayName;
    }

    public String getDescription() {
        return description;
    }

    public String getSecurityName() {
        return securityName;
    }
}
```

## Usage in Annotations

```java
@Route("admin/user")
@RolesAllowed({UserRole.ROLE_ADMIN, UserRole.ROLE_MANAGER})
public class UserView extends BaseView { ... }
```

Not `@RolesAllowed({"ADMIN", "MANAGER"})`. String literals scatter role names
across the codebase; renaming becomes a project-wide grep, and typos are silent
(a misspelled `"ADMNI"` compiles but never matches anyone).

## Where the Enum Lives

Place it in a module accessible from both persistence (where the role is stored
on the user record) and UI (where the constants are referenced in
`@RolesAllowed`). In a typical multi-module Vaadin/Spring layout that is the
`uimodel` module — `service` already depends on it, and so does `ui`
transitively.

## Why `getSecurityName()`

Spring Security represents authorities as strings and some integrations need a
method-of-an-instance accessor — e.g., when mapping a `UserEntity` field to a
granted-authority list. The getter is the seam; the constants are the
compile-time form for annotations.
