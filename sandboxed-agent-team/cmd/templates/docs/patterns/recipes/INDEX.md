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

**Recipe dependency order.** `conditional-auth` is the
configuration foundation; `audited-principal` is the
identity/audit foundation. The three auth-method recipes layer on
both. `passkey` additionally requires `form-login` (passkey is
supplementary, not a primary method). `oidc-sso` is independent
of `form-login` and `passkey`. Read in this order:
`conditional-auth` → `audited-principal` →
`form-login` / `oidc-sso` → `passkey` (if applicable).

| Path | Description |
|------|-------------|
| [conditional-auth.md](conditional-auth.md) | Configuration-driven authentication: typed `@ConfigurationProperties`, runtime `AuthMethods` API, startup combinability validator, `SecurityConfig` filter-chain branching. Configuration foundation for the per-method recipes below. |
| [audited-principal.md](audited-principal.md) | Cross-flow auth/audit foundation: `AuditedPrincipal` interface (`getKey()` + `getUsername()`), `CurrentUser` helper returning the principal, `AuditedEntity<K>` mapped superclass, `AuditorAware` with the `EntityManager.getReference()` optimisation (no SELECT per audited write). Every per-method recipe's principal implements this contract. |
| [form-login.md](form-login.md) | Username/password authentication via `VaadinSecurityConfigurer` and Spring Security's `DaoAuthenticationProvider`: `UserLookup` seam, `AuditedFormLoginUser` principal, `UserDetailsAdapter`, `BCryptPasswordEncoder`, Vaadin `LoginForm` integration. Username is whatever your project uses — not necessarily email. |
| [passkey.md](passkey.md) | WebAuthn / FIDO2 passkey authentication on Spring Security 7+: `WebAuthnUserAdapter`, triple-role `JpaPasskeyService` (`PasskeyService` + `PasskeyHandleManager` + `UserCredentialRepository`), the `webauthn_user_handle` invariant, Vaadin Flow ↔ Lit `PasskeyButton` bridge, `SecurityContextHolderStrategy` ordering gotcha, CSRF coordination with Vaadin's UIDL token. Requires `form-login`. |
| [oidc-sso.md](oidc-sso.md) | OpenID Connect SSO with Spring Boot's OAuth2 client: `OidcUserAdapter`, `AuditedOidcUser`, RP-initiated logout via `OidcClientInitiatedLogoutSuccessHandler` wrapped with Vaadin's `UidlRedirectStrategy` (so `/logout` over UIDL doesn't return "Invalid JSON response"), SSO Kit auto-config exclusion. Existing-user-only by default; auto-provisioning is a project decision. |
| [app-icon.md](app-icon.md) | Application icon catalog: `AppIcon` enum implementing `IconFactory`, named by intent not appearance; `UntitledUiIcon` (or any third-party library) adapter enum with `@JsModule` wiring. Heterogeneous delegates (`VaadinIcon` + custom sets) in one catalog; uniform `AppIcon.INTENT.create()` call sites. |
| [view-icon.md](view-icon.md) | `@ViewIcon` annotation for type-safe per-view icon declaration: eliminates the fragile `@Menu(icon = "vaadin:xxx")` string literal; value type is any `IconFactory`-implementing enum (`AppIcon` for projects with a full catalog, `VaadinIcon` for simpler projects). Read reflectively at runtime by navigation and layout components. Depends on app-icon. |
| [base-view.md](base-view.md) | `BaseView` shared view chrome: `Composite<VerticalLayout>` base class with icon + title header (left) and action slot (right), separated from a body slot by a bottom border. `@ViewIcon` integration is automatic. `setHeader()` / `setContent()` use replace-not-append semantics. Depends on view-icon. |
| [item-browser.md](item-browser.md) | List + detail component family: `ItemBrowser<T>` (grid + toolbar with quick filter, filter popover + badge, row tones, saved filters); `EditableItemBrowser<T>` interface; `ItemEditor<T>` contract (reuse not rebuild; `trySave()` returns false to keep host open); `MasterDetailItemBrowser<T>` (built); `DialogItemBrowser` / `ViewItemBrowser` (project-built variants). `FilterOption<V>` + `HasCaption` enum pattern for filter options. |

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