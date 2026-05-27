# Recipe: Form Login — Username/Password with Vaadin `LoginForm`

When implementing username/password form login backed by a JPA user table and
participating in the audited-principal contract, follow this recipe to produce
the `UserLookup` seam, the `AuditedUserDetails` principal, the password-encoder
bean, and the `LoginView` that hosts Vaadin's `LoginForm`.

> **Naming note.** Spring Security calls the identity field
> *username*. Whether your project's username column actually
> holds an email address, an employee ID, or an opaque login
> handle is a project decision — this recipe uses *username*
> throughout and treats email-as-username as one valid choice
> alongside others. Don't bake "username = email" into the
> abstraction.

## What this produces

- A `UserLookup` interface — the project-specific seam your
  adapter delegates to (covers tenant scoping, role mapping,
  whatever-your-username-column-is).
- An `AuditedFormLoginUser` principal extending Spring's `User`
  and implementing `AuditedPrincipal`.
- A `UserDetailsAdapter` (`UserDetailsService` impl) gated by the
  conditional-auth `form-login.enabled` flag.
- A `BCryptPasswordEncoder` bean.
- A Vaadin `LoginForm`-driven login view that integrates with
  Spring Security's `/login` filter.

## Dependencies

- Spring Boot 3+ (Spring Security 6+).
- Vaadin 24+ with `vaadin-spring-boot-starter`.
- The [audited-principal recipe](audited-principal.md) — this
  recipe's principal implements that contract.
- The [conditional-auth recipe](conditional-auth.md) — supplies
  `AuthMethods.isFormLoginEnabled()` and the
  `@ConditionalOnProperty` flag the adapter uses.

## Step 1 — Define the project-specific `UserLookup` seam

Form-login needs to resolve a username string to "the user record
plus the password hash." Anything beyond that is project-specific:
which column the username matches, whether to scope by tenant /
role / status, what to do for soft-deleted users, etc. Keep that
logic out of the recipe and behind a generic interface.

```java
package {base_package}.security;

import java.util.Optional;

public interface UserLookup {

    /**
     * Resolve a form-login user by the username supplied at
     * sign-in. Implementations apply project-specific filters
     * (tenant scoping, account status, role gating). The
     * returned {@link FormLoginUser} carries the user's key and
     * password hash; the adapter wraps it for Spring Security.
     */
    Optional<FormLoginUser> findFormLoginUser(String username);
}
```

`FormLoginUser` is a record that crosses the seam:

```java
public record FormLoginUser(Long key, String username,
                            String passwordHash, String role) { }
```

> The username field on the projection is *whatever
> `loadUserByUsername` was called with* — i.e., the value the
> client typed into the form's username field. The project may
> match this against an `email` column, an `employee_id` column,
> a `login_handle` column, or anything else; that's the lookup
> implementation's concern, not the adapter's.

## Step 2 — Define the `AuditedFormLoginUser` principal

Extends Spring Security's `User` (so `DaoAuthenticationProvider`
can read the password hash and authorities) and implements
`AuditedPrincipal` (so the audit pipeline can read the key without
a DB lookup). See [audited-principal.md](audited-principal.md) for
the interface.

```java
package {base_package}.security;

import java.util.Set;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.User;

public final class AuditedFormLoginUser extends User implements AuditedPrincipal {

    private final Long key;

    private AuditedFormLoginUser(FormLoginUser projection) {
        super(projection.username(),
              projection.passwordHash(),
              Set.of(new SimpleGrantedAuthority("ROLE_" + projection.role())));
        this.key = projection.key();
    }

    public static AuditedFormLoginUser of(FormLoginUser projection) {
        return new AuditedFormLoginUser(projection);
    }

    @Override public Long getKey()      { return key; }
    @Override public String getUsername() { return super.getUsername(); }
}
```

The `getUsername()` override is just to make the
`AuditedPrincipal` contract explicit — Spring's `User` already
implements it; the override declares satisfaction of both
contracts at the same call site.

> **Password erasure.** Spring's `User` superclass invalidates
> the password field when `eraseCredentials()` runs after
> successful authentication, so the hash isn't carried in the
> session-stored `Authentication`. You don't need to do anything
> extra; just don't override `eraseCredentials()` to keep it.

## Step 3 — Implement `UserDetailsService` via the adapter

The adapter is thin: delegate to `UserLookup`, wrap in
`AuditedFormLoginUser`, throw `UsernameNotFoundException` on any
miss. Spring's `DaoAuthenticationProvider` converts that to
`BadCredentialsException`, so the client can't tell whether the
miss was username, password, or status — by design.

```java
package {base_package}.security;

import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Component;

@Component
@ConditionalOnProperty(name = "{app}.auth.form-login.enabled", havingValue = "true")
public class UserDetailsAdapter implements UserDetailsService {

    private final UserLookup userLookup;

    public UserDetailsAdapter(UserLookup userLookup) {
        this.userLookup = userLookup;
    }

    @Override
    public UserDetails loadUserByUsername(String username) {
        return userLookup.findFormLoginUser(username)
                .filter(u -> u.passwordHash() != null && !u.passwordHash().isBlank())
                .map(AuditedFormLoginUser::of)
                .orElseThrow(() -> new UsernameNotFoundException("Invalid credentials"));
    }
}
```

The `@ConditionalOnProperty` gate matches the conditional-auth
foundation: when form login is disabled the adapter doesn't load,
so no rogue `UserDetailsService` competes with OIDC- or
passkey-only deployments.

## Step 4 — Expose the `PasswordEncoder` bean

BCrypt at strength 10 is the kit default. Strength is the
log₂(rounds), so 10 = 2¹⁰ = 1024 rounds; raise it (12, 14) if
your threat model demands more, with the cost being a
proportional CPU hit per login.

```java
@Bean
public PasswordEncoder passwordEncoder() {
    return new BCryptPasswordEncoder(10);
}
```

This bean is required regardless of whether form-login is the
*only* auth method — passkey registration and OIDC user
provisioning may still need to encode-or-verify a password
elsewhere in the system.

## Step 5 — Wire form-login into `SecurityConfig`

`VaadinSecurityConfigurer` does the heavy lifting: it binds
Spring Security's form-login filter to your Vaadin `LoginView`
route, configures CSRF with Vaadin's UIDL-aware token strategy,
and sets up the appropriate session-management defaults. You
don't call `http.formLogin(...)` explicitly.

```java
@Bean
SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {

    http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
        configurer.loginView(LoginView.class);
        // Override the configurer's denyAll default so unknown URLs
        // reach Vaadin's NotFoundView instead of returning a bare
        // HTTP 403 at the filter-chain level. View-based
        // annotations (@AnonymousAllowed, @PermitAll,
        // @RolesAllowed) remain the actual gate.
        configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
    });

    // Vaadin's default config permits only a couple of static-resource
    // paths. Opt into Spring Boot's common static locations so the
    // login background, theme assets, etc. resolve for anonymous
    // visitors. This must come AFTER the Vaadin configurer.
    http.authorizeHttpRequests(auth ->
            auth.requestMatchers(PathRequest.toStaticResources().atCommonLocations()).permitAll());

    return http.build();
}
```

If you also enable passkey or OIDC, add their branches after
this — see [passkey.md](passkey.md) and [oidc-sso.md](oidc-sso.md).

## Step 6 — Build the Vaadin `LoginForm` view

The view is a `VerticalLayout` (or your `BaseView` analogue)
hosting a `LoginForm`. The form's `setAction("login")` posts
straight to Spring Security's filter — **no custom controller
needed**. Spring's filter handles the success redirect (back to
the originally requested URL or `/`) and the failure redirect
(`/login?error`); the view reads the query param to render the
error state.

```java
@Route(value = "login", autoLayout = false)
@PageTitle("Sign In")
@AnonymousAllowed
public class LoginView extends VerticalLayout implements BeforeEnterObserver {

    public static final String ROUTE = "login";

    private final LoginForm loginForm;

    public LoginView(AuthMethods authMethods) {
        // ... layout, branding, sso/passkey buttons (see other recipes) ...

        if (authMethods.isFormLoginEnabled()) {
            loginForm = new LoginForm();

            // i18n: customise the label for the username field.
            // The default is "Username"; override only if your project
            // uses a different identity column AND you want the form
            // to advertise it (e.g., "Email" if your username column
            // is email; "Employee ID" if it's that).
            LoginI18n i18n = LoginI18n.createDefault();
            i18n.getForm().setTitle("Sign in with password");
            // i18n.getForm().setUsername("Email");   // example override
            loginForm.setI18n(i18n);

            loginForm.setAction("login");
            loginForm.setForgotPasswordButtonVisible(false);
            card.add(loginForm);
        } else {
            loginForm = null;
        }
    }

    @Override
    public void beforeEnter(BeforeEnterEvent event) {
        if (loginForm != null
                && event.getLocation().getQueryParameters().getParameters().containsKey("error")) {
            loginForm.setError(true);
        }
    }
}
```

## Step 7 — Encode passwords on user create

Whenever you create or update a user's password (admin reset,
self-service change, registration), encode through the
`PasswordEncoder` bean — never store the raw password and never
write a custom `MessageDigest` call. The encoded form is what the
adapter reads back via `FormLoginUser.passwordHash()`.

```java
@Service
public class UserService {

    private final PasswordEncoder passwordEncoder;
    private final UserRepository userRepository;

    public UserService(PasswordEncoder passwordEncoder, UserRepository userRepository) {
        this.passwordEncoder = passwordEncoder;
        this.userRepository = userRepository;
    }

    public void setPassword(Long userKey, String rawPassword) {
        validateStrength(rawPassword);  // see docs/patterns/architecture/security.md
        String hash = passwordEncoder.encode(rawPassword);
        userRepository.updatePasswordHash(userKey, hash);
    }
}
```

Password-strength validation (entropy, length, pwned-list check)
is project policy — see
`docs/patterns/architecture/security.md` for the kit's stance and
defer to `docs/solutions/` for project-specific rules.

## Decisions this recipe imposes

- **Username is whatever your project uses.** Spring's "username"
  is a *concept*, not "email." The recipe never assumes email; the
  adapter accepts whatever string comes through `loadUserByUsername`,
  and the i18n label on the form is the project's choice.
- **`UserLookup` is the project seam.** The adapter doesn't query
  the database directly; project-specific lookup rules (tenant,
  status, role) stay behind the interface.
- **`AuditedFormLoginUser` extends Spring's `User`, not a custom
  base.** Spring's `User` already handles credential erasure,
  serialization, equality. Don't reinvent.
- **`@ConditionalOnProperty` gates the adapter.** Disabling form
  login removes the `UserDetailsService` from the context; OIDC- or
  passkey-only deployments don't compete with a stub.
- **`VaadinSecurityConfigurer` is the form-login wiring.** Don't
  call `http.formLogin(...)` explicitly. The configurer wires the
  filter, the login view, the CSRF strategy, and the session
  defaults coherently with Vaadin.
- **BCrypt at strength 10 by default.** Adjustable; document any
  deviation in `docs/solutions/security.md`.

## What to verify

- POST to `/login` with valid `username` + `password` form fields
  → redirect to `/` (or originally-requested URL); session cookie
  present on the response.
- POST to `/login` with bad credentials → redirect to
  `/login?error`; the view renders the error state.
- The authenticated user's principal class is
  `AuditedFormLoginUser` (cast and verify) and
  `principal.getKey()` returns the user's `usr.key`.
- An audited write performed under a form-login session populates
  `created_by_key` with the same key.
- `Authentication.getName()` and
  `principal.getUsername()` both return the value the client typed
  — they don't silently transform it.
- A login attempt against a soft-deleted / inactive user fails
  the same as a bad-password attempt (no information leak).

## Related

- [conditional-auth.md](conditional-auth.md) — the
  `form-login.enabled` flag and `AuthMethods.isFormLoginEnabled()`.
- [audited-principal.md](audited-principal.md) — the principal
  contract `AuditedFormLoginUser` satisfies.
- [passkey.md](passkey.md) — passkey is a *supplementary*
  method that requires form login (see the conditional-auth
  combinability rule).
- [oidc-sso.md](oidc-sso.md) — alternative auth method;
  configurable independently.
- `docs/patterns/architecture/security.md` — password strength
  policy, BCrypt strength rationale, session config, RBAC.
- `docs/solutions/` — your project's username column choice,
  user-lookup logic, role mapping, password policy.
