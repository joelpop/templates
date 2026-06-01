# UserRole Enum Module Placement

When deciding where to place the `UserRole` enum, put it in a module accessible from both persistence (where the role is stored on the user record) and UI (where the constants appear in `@RolesAllowed`).

In a typical multi-module Vaadin/Spring layout that is the `uimodel` module — `service` already depends on it, and so does `ui` transitively.
