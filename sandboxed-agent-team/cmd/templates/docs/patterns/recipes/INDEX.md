# Recipes — Index

Step-by-step implementation recipes for capabilities that recur
across projects on this stack (Vaadin + Spring Boot + Spring Data
JPA). Recipes are different from the rest of `docs/patterns/`:

- **Conventions, architecture, ui, testing, deployment** entries
  are *rules and idioms* — "do this, not that," "here's the
  convention." They guide ongoing implementation work.
- **Recipes** are *how-to guides for a specific capability* — "here
  is the sequence of steps to integrate passkey auth," "here is the
  configuration to switch between OIDC SSO and uid/pwd via
  `application.properties`." They guide the initial implementation
  of a feature.

A recipe belongs here when the steps are fixed by the framework
stack (Vaadin's passkey API, Spring Security's OIDC client) rather
than project-specific decisions. Project-specific configuration
(which IdP, tenancy model, entity layout) lives in
`docs/architecture/`, referencing the recipe.

| Path | Description |
|------|-------------|
| [conditional-auth.md](conditional-auth.md) | Configuration-driven authentication: typed `@ConfigurationProperties`, runtime `AuthMethods` API, startup combinability validator, `SecurityConfig` filter-chain branching. Foundation for the per-method recipes below. |
| _(more recipes pending: passkey integration, OIDC/SSO with Vaadin SSO Kit)_ | |

## Recipe-writing conventions

- Open with **what the recipe produces** (one paragraph) and
  **what it depends on** (stack versions, MCP servers, libraries).
- Walk the steps in order. Each step has a heading, a brief
  description of *why* the step is needed, and the concrete code or
  configuration.
- Surface the **decisions** the recipe imposes (e.g., "this recipe
  assumes user identities are stored in the application DB; if you
  use an external IdP for identities, see [other-recipe]").
- Close with **what to verify** — the smoke test or manual check
  that confirms the recipe was applied correctly.
- Cross-link to related conventions, architecture entries, or other
  recipes.

When extracting a recipe from a project-specific feature, keep only
framework-fixed steps. Move project-specific UX, error handling,
and business rules to a `docs/architecture/` entry that references
the recipe.