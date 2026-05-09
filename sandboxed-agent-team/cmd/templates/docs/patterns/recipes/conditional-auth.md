# Recipe: Conditional Authentication Methods via `application.properties`

A configuration-driven pattern for enabling any combination of
authentication methods — form login, passkey (WebAuthn), OIDC SSO —
with **invalid combinations failing fast at startup**. Build this
foundation once, then layer the per-method recipes
([passkey](passkey.md), [oidc-sso](oidc-sso.md)) on top.

## What this produces

- A typed config record (`AuthProperties`) bound to
  `application.properties` keys.
- A runtime API (`AuthMethods`) the rest of the app queries to ask
  "is this method enabled?" — no scattered `@Value` lookups.
- A startup validator (`AuthMethodCombinabilityValidator`) that
  rejects nonsensical combinations (e.g., passkey enabled without
  form login) before the app finishes booting.
- A `SecurityConfig` filter chain that branches on the typed
  `AuthMethods` API, adding only the filters and login flows the
  configuration calls for.

## Dependencies

- Spring Boot 3+ (Spring Security 6+ for the configurer DSL; 7+ for
  the WebAuthn DSL referenced by the passkey recipe)
- Vaadin 24+ with `vaadin-spring-boot-starter`
- `vaadin-sso-kit-starter` if SSO is in scope (or just Spring
  Security's OAuth2 client; SSO Kit adds Vaadin-aware redirect
  handling)

## Step 1 — Define the typed config record

Create a `@ConfigurationProperties`-bound record under your auth
package. Each method gets a nested `MethodConfig` with a single
`enabled` flag; method-specific settings (e.g., WebAuthn `rpId` /
`origin`) are nested records on the same root.

```java
package {base_package}.provider.security;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "{app}.auth")
public record AuthProperties(
        MethodConfig formLogin,
        MethodConfig passkey,
        MethodConfig sso,
        WebAuthnConfig webauthn) {

    public AuthProperties {
        if (formLogin == null) formLogin = new MethodConfig(false);
        if (passkey   == null) passkey   = new MethodConfig(false);
        if (sso       == null) sso       = new MethodConfig(false);
        if (webauthn  == null) webauthn  = new WebAuthnConfig(null, null);
    }

    public record MethodConfig(boolean enabled) {}

    /** {@code rpId}: registrable RP domain (e.g., {@code localhost} dev,
     *  {@code app.example.com} prod). {@code origin}: full URL the
     *  browser sees — must match exactly. */
    public record WebAuthnConfig(String rpId, String origin) {}
}
```

The compact constructor defaults each method to `disabled` when the
property is absent — no NPE on a missing block.

## Step 2 — Expose a runtime API

A small interface in `service` decouples consumers from the config
record. The rest of the app queries `AuthMethods`, never
`AuthProperties` directly.

```java
package {base_package}.service.security;

public interface AuthMethods {
    boolean isFormLoginEnabled();
    boolean isPasskeyEnabled();
    boolean isSsoEnabled();

    /** Route that initiates OIDC (e.g., {@code /oauth2/authorization/keycloak}).
     *  Only meaningful when SSO is enabled; null otherwise. */
    String getSsoLoginRoute();
}
```

## Step 3 — Wire `AuthMethods` to the typed properties

A `@Configuration` class binds the properties record, exposes the
`AuthMethods` bean, and (if you have method-specific settings like
WebAuthn) exposes those too as separate beans.

```java
package {base_package}.provider.security;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import {base_package}.service.security.AuthMethods;
import {base_package}.service.security.WebAuthnSettings;

@Configuration
@EnableConfigurationProperties(AuthProperties.class)
public class AuthMethodsConfig {

    private static final String DEFAULT_SSO_LOGIN_ROUTE =
            "/oauth2/authorization/keycloak";

    @Bean
    AuthMethods authMethods(AuthProperties properties, Environment environment) {
        return new PropertiesBackedAuthMethods(properties, environment);
    }

    @Bean
    WebAuthnSettings webAuthnSettings(AuthProperties properties) {
        return new PropertiesBackedWebAuthnSettings(properties);
    }

    private record PropertiesBackedAuthMethods(
            AuthProperties properties, Environment environment)
            implements AuthMethods {

        @Override public boolean isFormLoginEnabled() { return properties.formLogin().enabled(); }
        @Override public boolean isPasskeyEnabled()   { return properties.passkey().enabled(); }
        @Override public boolean isSsoEnabled()       { return properties.sso().enabled(); }

        @Override public String getSsoLoginRoute() {
            if (!isSsoEnabled()) return null;
            return environment.getProperty("vaadin.sso.login-route", DEFAULT_SSO_LOGIN_ROUTE);
        }
    }

    private record PropertiesBackedWebAuthnSettings(AuthProperties properties)
            implements WebAuthnSettings {
        @Override public String rpId()   { return properties.webauthn().rpId(); }
        @Override public String origin() { return properties.webauthn().origin(); }
    }
}
```

## Step 4 — Validate combinations at startup

A `@Component` whose constructor throws when configuration is
invalid — the right place for "at least one method must be enabled"
and "passkey requires form login." Invalid config becomes a clear
startup error, not a runtime mystery.

```java
package {base_package}.provider.security;

import org.springframework.stereotype.Component;

@Component
public class AuthMethodCombinabilityValidator {

    public AuthMethodCombinabilityValidator(AuthProperties properties) {
        boolean form    = properties.formLogin().enabled();
        boolean passkey = properties.passkey().enabled();
        boolean sso     = properties.sso().enabled();

        if (!form && !passkey && !sso) {
            throw new IllegalStateException(
                    "No authentication methods enabled. At least one of " +
                    "{app}.auth.{form-login,passkey,sso}.enabled must be true.");
        }
        if (passkey && !form) {
            throw new IllegalStateException(
                    "{app}.auth.passkey.enabled=true requires " +
                    "{app}.auth.form-login.enabled=true. Passkey is a " +
                    "supplementary authentication method (users register a " +
                    "passkey after logging in with a password) and cannot be " +
                    "configured without form login.");
        }
    }
}
```

Project-specific rules go in the same validator — keep all
combinability constraints in one place.

## Step 5 — Branch the SecurityConfig on `AuthMethods`

The `SecurityFilterChain` injects `AuthMethods` and adds only
configured branches:

```java
@Bean
SecurityFilterChain securityFilterChain(
        HttpSecurity http,
        AuthMethods authMethods,
        ObjectProvider<OidcUserService> oidcUserServiceProvider,
        ObjectProvider<ClientRegistrationRepository> clientRegistrationRepositoryProvider)
        throws Exception {

    http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
        configurer.loginView(LoginView.class);
        configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
    });

    if (authMethods.isSsoEnabled()) {
        // OIDC branch — see oidc-sso.md
    }

    if (authMethods.isPasskeyEnabled()) {
        // WebAuthn branch — see passkey.md
    }

    return http.build();
}
```

`ObjectProvider<>` is correct for OIDC beans that only exist when
SSO is enabled — it returns `null` when no bean is registered.

## Step 6 — `application.properties` defaults

Set conservative defaults in `application.properties` (form login
on, passkey/SSO off) with inline documentation. Per-env overrides
go in `application-{profile}.properties` or
`application-local.properties` (gitignored).

```properties
# Authentication methods. At least one must be enabled.
# Passkey requires form-login. Invalid combinations fail startup.
{app}.auth.form-login.enabled=true
{app}.auth.passkey.enabled=false
{app}.auth.sso.enabled=false

# WebAuthn / Passkey — RP id is the registrable domain (localhost
# dev, app.example.com prod). Origin is the full URL the browser
# sees (scheme + host + port). Browser rejects on any mismatch.
#{app}.auth.webauthn.rp-id=localhost
#{app}.auth.webauthn.origin=http://localhost:8080
```

## Decisions this recipe imposes

- **Typed config over string lookups.** A
  `@ConfigurationProperties` record gives compile-time refactoring
  safety and a single place to default missing blocks. Avoid
  `@Value("${...}")` lookups scattered across the codebase.
- **Fail fast on invalid combinations.** The validator catches bad
  configs before users hit them; the cost is a few lines of code.
- **`AuthMethods` decouples consumers from the config record.**
  Future config changes (renaming a key, splitting a method into
  sub-flags) don't ripple through callers.
- **Default to *off* for non-essential methods.** Form login on,
  passkey/SSO off. Operators opt in per environment.
- **Method-specific config beans are conditional.** `WebAuthnConfig`
  uses `@ConditionalOnProperty` so its beans don't load when passkey
  is off — see [passkey.md](passkey.md).

## What to verify

- App starts with **only `form-login.enabled=true`** → login view
  shows, no OAuth2 / WebAuthn endpoints registered.
- App starts with **`form-login=true`, `passkey=true`** → login +
  passkey ceremonies both work; `/webauthn/register/options` is
  reachable for an authenticated session.
- App starts with **`form-login=true`, `sso=true`** → login view
  shows, OIDC redirect from "Login with SSO" works, end_session
  logout closes both local + IdP sessions.
- App **fails to start** with all three disabled (validator throws
  with a clear message).
- App **fails to start** with `passkey=true` and `form-login=false`
  (validator throws with the supplementary-method explanation).

## Related

- [passkey.md](passkey.md) — the passkey-specific filter chain,
  WebAuthn bean wiring, and Vaadin/CSRF coordination.
- [oidc-sso.md](oidc-sso.md) — OIDC client setup, SSO Kit
  integration, RP-initiated logout via Vaadin's UIDL channel.
- `docs/patterns/architecture/security.md` — the surrounding
  security patterns this recipe operates within (BCrypt for form
  login, session config, RBAC).
