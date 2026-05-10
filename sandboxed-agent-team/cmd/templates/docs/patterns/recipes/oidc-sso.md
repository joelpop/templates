# Recipe: OIDC / SSO Authentication

OpenID Connect single sign-on for Vaadin Flow apps. Spring
Security's OAuth2 client handles the authorization code flow,
token exchange, and userinfo call; this recipe covers the Vaadin
integration — wiring the IdP redirect to the login view, handling
RP-initiated logout through Vaadin's UIDL channel, and mapping
IdP-supplied identities to your domain user model. It assumes
**existing-user-only authentication** by default — the IdP's
authentication is trusted, but a user must already exist in the
application database to be granted access. Auto-provisioning is a
project decision; see "Decisions" below.

## What this produces

- OAuth2 client configuration in `application.properties` (one
  registration per IdP).
- An `AuditedOidcUser` principal extending Spring's
  `DefaultOidcUser` and implementing `AuditedPrincipal`.
- An `OidcUserAdapter` extending `OidcUserService`: extracts the
  email claim, validates it, resolves the active user via
  `UserLookup`, wraps in `AuditedOidcUser`. Refuses with a
  descriptive `OAuth2Error` on any miss.
- An RP-initiated logout handler (Spring Security's
  `OidcClientInitiatedLogoutSuccessHandler`) wrapped with
  Vaadin's `UidlRedirectStrategy` so the logout redirect is
  emitted as a UIDL navigation rather than a plain 302 — Vaadin's
  client otherwise rejects it.
- `SecurityConfig` branches that engage `oauth2Login` and the
  custom logout handler, gated by the conditional-auth foundation.
- A "Sign in with SSO" button on the login view that redirects to
  the OAuth2 authorization endpoint.

## Dependencies

- Spring Boot 3+ with `spring-boot-starter-oauth2-client`.
- Vaadin 24+ with `vaadin-spring-boot-starter`.
- `vaadin-sso-kit-starter` for the `UidlRedirectStrategy` class
  — but **exclude its auto-configuration**; this recipe wires
  the pieces explicitly to coexist with the multi-method auth
  layer from [conditional-auth](conditional-auth.md). See Step
  5's gotcha.
- The [audited-principal recipe](audited-principal.md) — the
  OIDC principal implements `AuditedPrincipal`.
- The [conditional-auth recipe](conditional-auth.md) — supplies
  `AuthMethods.isSsoEnabled()`, `getSsoLoginRoute()`, and the
  `@ConditionalOnProperty` flag the adapter uses.
- An OIDC-compliant IdP (Keycloak, Okta, Auth0, Entra ID,
  Cognito, etc.).

## Step 1 — Configure the OAuth2 client

Spring Boot reads two property trees: `registration` (client
credentials) and `provider` (IdP metadata). The provider's
`issuer-uri` triggers OIDC discovery — Spring fetches
`/.well-known/openid-configuration` to populate authorization,
token, userinfo, and end-session endpoints automatically.

```properties
# Client registration — credentials your IdP issued for this app.
spring.security.oauth2.client.registration.{provider-name}.client-id=...
spring.security.oauth2.client.registration.{provider-name}.client-secret=...
spring.security.oauth2.client.registration.{provider-name}.scope=openid,profile,email
spring.security.oauth2.client.registration.{provider-name}.authorization-grant-type=authorization_code
spring.security.oauth2.client.registration.{provider-name}.redirect-uri={baseUrl}/login/oauth2/code/{registrationId}

# Provider metadata — the issuer URL discovers the rest.
spring.security.oauth2.client.provider.{provider-name}.issuer-uri=https://idp.example.com/realms/{realm}
```

`{provider-name}` is your label for this IdP — e.g., `keycloak`,
`okta`, `entra`. It's the segment that appears in
`/oauth2/authorization/{provider-name}` (the URL the SSO button
redirects to).

Conditional-auth's default `getSsoLoginRoute()` value is
`/oauth2/authorization/keycloak`; if your provider name differs,
set:

```properties
vaadin.sso.login-route=/oauth2/authorization/{provider-name}
```

> The IdP must be configured with the exact redirect URI Spring
> Security advertises (`{baseUrl}/login/oauth2/code/{provider-name}`,
> resolved at request time). A mismatch surfaces as a generic
> "redirect_uri_mismatch" error from the IdP — easy to miss
> because the application's logs only show the failed callback,
> not the underlying claim.

## Step 2 — Define the `AuditedOidcUser` principal

Extends `DefaultOidcUser` (so Spring Security's OAuth2 plumbing
treats it as a normal OIDC user) and implements `AuditedPrincipal`
(so the audit pipeline and `CurrentUser.get()` see the same shape
across all auth flows).

```java
public final class AuditedOidcUser extends DefaultOidcUser implements AuditedPrincipal {

    private final Long key;
    private final String username;

    private AuditedOidcUser(AuthenticatedUser user, OidcUser delegate, String nameAttributeKey) {
        super(authoritiesFor(user), delegate.getIdToken(), delegate.getUserInfo(), nameAttributeKey);
        this.key = user.key();
        this.username = user.username();
    }

    public static AuditedOidcUser of(AuthenticatedUser user, OidcUser delegate) {
        // The name-attribute key drives Authentication.getName(). Pick
        // whichever claim your project uses as the username — "email" is
        // common; "preferred_username" or "sub" are equally valid choices.
        // Match the column your UserLookup queries against.
        return new AuditedOidcUser(user, delegate, "email");
    }

    private static Set<GrantedAuthority> authoritiesFor(AuthenticatedUser user) {
        return Set.of(new SimpleGrantedAuthority("ROLE_" + user.role()));
    }

    @Override public Long getKey()      { return key; }
    @Override public String getUsername() { return username; }
}
```

> **Name-attribute key choice.** `"email"` is what the example
> uses; a project that identifies users by `preferred_username`
> or `sub` should pass that key instead. The IdP must include the
> chosen claim — verify by inspecting the ID token at
> `https://jwt.io` during integration.

## Step 3 — Implement the `OidcUserAdapter`

Extends `OidcUserService`, overrides `loadUser` to:

1. Call `super.loadUser` — the network round-trip to the userinfo
   endpoint plus ID-token validation.
2. Extract the username claim (email in the example; whatever
   your project chose in Step 2).
3. Validate the claim is present and well-formed.
4. Resolve the active user via `UserLookup.findActiveUser`.
5. Wrap in `AuditedOidcUser` — or refuse with a descriptive
   `OAuth2Error`.

```java
@Component
@ConditionalOnProperty(name = "{app}.auth.sso.enabled", havingValue = "true")
public class OidcUserAdapter extends OidcUserService {

    private final UserLookup userLookup;

    public OidcUserAdapter(UserLookup userLookup) {
        this.userLookup = userLookup;
    }

    @Override
    public OidcUser loadUser(OidcUserRequest userRequest) throws OAuth2AuthenticationException {
        return mapToDomainUser(super.loadUser(userRequest));
    }

    /** Package-private so unit tests can exercise mapping without the
     *  network call inside super.loadUser(). */
    OidcUser mapToDomainUser(OidcUser oidcUser) {
        String email = oidcUser.getEmail();
        if (email == null || email.isBlank() || !email.contains("@")) {
            throw new OAuth2AuthenticationException(
                    new OAuth2Error("missing_or_malformed_email_claim",
                            "The ID token did not include a usable email claim.", null));
        }
        return userLookup.findActiveUser(email)
                .map(user -> (OidcUser) AuditedOidcUser.of(user, oidcUser))
                .orElseThrow(() -> new OAuth2AuthenticationException(
                        new OAuth2Error("user_not_provisioned",
                                "No active user provisioned for email " + email, null)));
    }
}
```

The two `OAuth2Error` codes are project-defined. The error
surfaces to the user as a redirect to `/login?error` (Spring's
default OAuth2 failure handler); the message goes to the server
log. If you want a more descriptive UI, register a
`AuthenticationFailureHandler` — see "Decisions" for why we don't
do that by default.

## Step 4 — Configure `OidcClientInitiatedLogoutSuccessHandler`

When the user logs out, Vaadin's
`AuthenticationContext.logout()` POSTs to `/logout`. Spring
Security invalidates the local session and hands off to a
`LogoutSuccessHandler`. For OIDC, the handler must redirect to
the IdP's `end_session_endpoint` with an `id_token_hint` so the
IdP terminates its session too — otherwise the next click on
"Sign in with SSO" silently re-authenticates the previous user.

```java
private static LogoutSuccessHandler createRpInitiatedLogoutHandler(
        ClientRegistrationRepository clientRegistrationRepository) {
    if (clientRegistrationRepository == null) return null;

    var handler = new OidcClientInitiatedLogoutSuccessHandler(clientRegistrationRepository);

    // {baseUrl} resolves at request time to scheme+host+port+context
    // (e.g., http://localhost:8080). The IdP client must list this
    // URL in its post.logout.redirect.uris attribute.
    handler.setPostLogoutRedirectUri("{baseUrl}/");

    // Vaadin's AuthenticationContext.logout() triggers /logout over
    // the UIDL (AJAX) channel. A plain HTTP 302 causes Vaadin's
    // client to throw "Invalid JSON response from server".
    // UidlRedirectStrategy detects UIDL requests and emits a JSON
    // navigation instruction the client understands. Applies to both
    // the OIDC end_session redirect and the form-login fallback
    // redirect (super.onLogoutSuccess in the parent class).
    handler.setRedirectStrategy(new UidlRedirectStrategy());
    return handler;
}
```

`UidlRedirectStrategy` comes from `com.vaadin.sso.starter` — the
SSO Kit. The kit's auto-configuration is excluded (Step 5); only
this class is borrowed.

## Step 5 — Wire OIDC into `SecurityConfig`

Two coordination points: the OAuth2 login DSL (with custom user
service) and the logout handler (wired through
`VaadinSecurityConfigurer`).

```java
@Bean
SecurityFilterChain securityFilterChain(
        HttpSecurity http,
        SecurityContextHolderStrategy securityContextHolderStrategy,
        AuthMethods authMethods,
        ObjectProvider<OidcUserService> oidcUserServiceProvider,
        ObjectProvider<ClientRegistrationRepository> clientRegistrationRepositoryProvider)
        throws Exception {

    SecurityContextHolder.setContextHolderStrategy(securityContextHolderStrategy);

    OidcUserService oidcUserService = authMethods.isSsoEnabled()
            ? oidcUserServiceProvider.getIfAvailable()
            : null;
    LogoutSuccessHandler oidcLogoutHandler = authMethods.isSsoEnabled()
            ? createRpInitiatedLogoutHandler(clientRegistrationRepositoryProvider.getIfAvailable())
            : null;

    http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
        configurer.loginView(LoginView.class);
        configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
        if (oidcLogoutHandler != null) {
            // RP-initiated logout: AuthenticationContext.logout() POSTs
            // /logout → Spring Security invalidates the local session and
            // hands off to this handler, which redirects to the IdP's
            // end_session_endpoint with an id_token_hint. The IdP ends
            // its session and redirects back to post_logout_redirect_uri.
            // For non-OIDC (form-login) users the handler falls back to
            // a plain redirect — no id token to hint with.
            configurer.logoutSuccessHandler(oidcLogoutHandler);
        }
    });

    if (authMethods.isSsoEnabled()) {
        http.oauth2Login(oauth2 -> {
            // Route the login UI back to our Vaadin LoginView. Without
            // this, oauth2Login() registers DefaultLoginPageGeneratingFilter,
            // which intercepts /login before Vaadin's router and renders
            // the generic "Login with OAuth 2.0" page (including on
            // /login?error after a failed OIDC round-trip).
            oauth2.loginPage("/" + LoginView.ROUTE);
            if (oidcUserService != null) {
                oauth2.userInfoEndpoint(userInfo -> userInfo.oidcUserService(oidcUserService));
            }
        });
    }

    return http.build();
}
```

> **SSO Kit auto-config exclusion.** Add this to your application
> properties to disable the SSO Kit's `@AutoConfiguration`
> classes — they wire a single-method OIDC-only flow that
> conflicts with the multi-method config built here:
>
> ```properties
> spring.autoconfigure.exclude=\
>   com.vaadin.sso.starter.SingleSignOnConfiguration,\
>   com.vaadin.sso.starter.SingleSignOnDefaultBeansConfiguration
> ```

## Step 6 — Add the SSO button to the login view

```java
if (authMethods.isSsoEnabled()) {
    Button ssoButton = new Button("Sign in with single sign-on",
            _ -> redirectToSsoLogin());
    ssoButton.addThemeVariants(ButtonVariant.LUMO_PRIMARY);
    card.add(ssoButton);
}

// Helper.
private void redirectToSsoLogin() {
    String loginRoute = authMethods.getSsoLoginRoute();
    if (loginRoute == null) return;
    getUI().ifPresent(ui -> ui.getPage().setLocation(loginRoute));
}
```

For SSO-only deployments (form login + passkey both off), the
`LoginView.beforeEnter` should detect the configuration and redirect
immediately rather than rendering an empty card:

```java
@Override
public void beforeEnter(BeforeEnterEvent event) {
    if (isSsoOnly()) {
        redirectToSsoLogin();
        return;
    }
    // ...
}

private boolean isSsoOnly() {
    return authMethods.isSsoEnabled()
            && !authMethods.isFormLoginEnabled()
            && !authMethods.isPasskeyEnabled();
}
```

## Decisions this recipe imposes

- **No auto-provisioning by default.** A user must exist in the
  application database before the IdP can authenticate them.
  Auto-provisioning (creating a user record on first SSO login)
  couples user lifecycle to the IdP and is a project policy
  decision — typically wrong for B2B / regulated apps, sometimes
  right for B2C. The recipe's `OidcUserAdapter` refuses unknown
  emails with `user_not_provisioned`; opt into auto-provisioning
  in `docs/architecture/` and add the JIT-creation step there.
- **Email is one valid claim choice; the recipe doesn't bake it
  in.** The `OidcUserAdapter` matches the example to fleet-acuity's
  email-as-username choice; substitute `preferred_username` or
  `sub` if your project uses a different identity claim. Match it
  in `AuditedOidcUser`'s `nameAttributeKey` argument and in the
  `UserLookup` query.
- **RP-initiated logout, not local-only.** Skipping the IdP
  end-session leaves a half-logged-out state — local app session
  gone, IdP session kept. The next "Sign in with SSO" click
  silently re-authenticates the previous user, which feels
  broken to operators of shared workstations.
- **`UidlRedirectStrategy` for the logout handler.** Plain 302s
  don't work over Vaadin's UIDL channel; the SSO Kit's helper
  emits a JSON navigation that Vaadin's client understands.
- **SSO Kit class, not SSO Kit auto-config.** The auto-config
  wires an OIDC-only flow that conflicts with the multi-method
  layout. Exclude the auto-config; depend on
  `vaadin-sso-kit-starter` only for the `UidlRedirectStrategy`
  class.
- **`OAuth2Error` codes are descriptive, the UI is generic.**
  The error code goes to the server log (operator visibility);
  the user sees `/login?error` (no information leak).
- **Refuse, don't degrade, on missing email claim.** If the IdP
  doesn't include the username claim you configured, the auth
  fails — don't fall back to `sub` or generate a synthetic value.
  Misconfigured IdPs surface as fix-able errors; silent fallbacks
  surface as inconsistent audit logs months later.

## What to verify

- Click "Sign in with SSO" on the login view → redirect to the
  IdP → log in there → return to the application
  authenticated. The principal class is `AuditedOidcUser`,
  `principal.getKey()` returns the user's `usr.key`,
  `Authentication.getName()` returns the configured name claim.
- Logout: click sign-out (Vaadin
  `AuthenticationContext.logout()`) → end_session redirect to the
  IdP → IdP terminates its session → return to `{baseUrl}/`. The
  next "Sign in with SSO" click prompts for credentials again.
- An IdP-authenticated user not in the application DB → fails
  with `user_not_provisioned` in the server log; user redirected
  to `/login?error`.
- An ID token without an email claim → fails with
  `missing_or_malformed_email_claim`; same UI path.
- An audited write performed under an SSO session populates
  `created_by_key` with the resolved user's key (via
  `AuditorAware` reading `AuditedOidcUser.getKey()`).
- `{app}.auth.sso.enabled=false` removes the OIDC adapter, the
  logout handler customisation falls back to the plain redirect,
  and the SSO button disappears from the login view.

## Related

- [conditional-auth.md](conditional-auth.md) — the
  `sso.enabled` flag, `getSsoLoginRoute()`, and the SSO Kit
  auto-config exclusion rationale.
- [audited-principal.md](audited-principal.md) — the principal
  contract `AuditedOidcUser` satisfies.
- [form-login.md](form-login.md) — alternative auth method;
  configurable independently.
- [passkey.md](passkey.md) — alternative auth method.
- `docs/patterns/architecture/security.md` — surrounding
  security architecture.
- `docs/architecture/` — your project's IdP choice, role mapping
  from IdP claims to application roles, auto-provisioning policy
  if any, name-attribute claim choice, and SSO-only-mode UX.
