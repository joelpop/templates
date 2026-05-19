# Recipe: Passkey (WebAuthn) Authentication

Spring Security 7's WebAuthn DSL handles the cryptographic
ceremonies and HTTP endpoints; this recipe covers the integration
surface — the user-entity adapter, credential persistence,
handle-lifecycle bookkeeping, and the Vaadin component that runs
the browser-side ceremony. Spring Security's own WebAuthn docs
are thin and assume an in-memory or JDBC repository; this recipe
is for projects that need their own JPA-backed credential store
and want their auth to participate in the kit's
audited-principal contract.

> **Prerequisite reminder.** Passkey is *supplementary* to form
> login per the conditional-auth combinability rule — users sign
> in with a password first, then register a passkey from their
> account preferences. Don't enable passkey without form login.

## What this produces

- A `WebAuthnSettings` interface exposing `rpId` (registrable
  domain) and `origin` (full URL the browser sees).
- A `webauthn_user_handle` column on the user entity (nullable;
  16 bytes when present).
- A `UserCredentialEntity` + `PasskeyCredentialRepository` for
  per-credential records (mirrors Spring Security's reference
  JDBC schema).
- An `AuditedPasskeyPrincipal` implementing both
  `PublicKeyCredentialUserEntity` (Spring Security WebAuthn) and
  `AuditedPrincipal` (the kit's audit contract).
- A `WebAuthnUserAdapter` implementing
  `PublicKeyCredentialUserEntityRepository` — bridges Spring
  Security's user-entity API to your domain user model.
- A triple-role `JpaPasskeyService` implementing
  `PasskeyService` (UI list/revoke), `PasskeyHandleManager`
  (handle lifecycle), and `UserCredentialRepository` (Spring
  Security persistence) — three concerns, one transactional
  home.
- A `WebAuthnConfig` `@Configuration` exposing the
  `WebAuthnRelyingPartyOperations` bean.
- The `SecurityConfig` branch that engages the WebAuthn DSL and
  coordinates CSRF with Vaadin.
- A `PasskeyButton` Vaadin Flow component wrapping a Lit web
  component that runs `navigator.credentials.create()` /
  `navigator.credentials.get()` browser-side.

## Dependencies

- Spring Boot 3+, Spring Security **7+** (the WebAuthn DSL —
  not present in 6.x).
- Vaadin 24+ with `vaadin-spring-boot-starter`.
- The [audited-principal recipe](audited-principal.md) — the
  passkey principal implements `AuditedPrincipal`.
- The [conditional-auth recipe](conditional-auth.md) — supplies
  `AuthMethods.isPasskeyEnabled()` and the
  `@ConditionalOnProperty` flags.
- The [form-login recipe](form-login.md) — passkey requires it.
- A custom Lit web component (`{app}-passkey-button.ts`) that
  performs the browser-side WebAuthn ceremony. Sketch only in
  this recipe; the full TypeScript implementation lives in
  `docs/solutions/` or your client-side codebase.

## Step 1 — Define the `WebAuthnSettings` interface

The two values WebAuthn requires from your environment. Read from
`AuthProperties.WebAuthnConfig` (set up in conditional-auth);
expose via this interface so consumers don't reach into the
typed-config record.

```java
package {base_package}.security;

public interface WebAuthnSettings {

    /** Registrable domain — e.g., {@code localhost} dev,
     *  {@code {app}.example.com} prod. */
    String rpId();

    /** Full URL the browser sees — scheme + host + port. Must
     *  match exactly; the browser refuses on any mismatch. */
    String origin();
}
```

## Step 2 — Add the user-handle column to the user entity

WebAuthn identifies a user via a stable opaque handle, not the
username. The handle is a 16-byte SecureRandom blob, generated
on first passkey-related interaction and cleared when the user
has no credentials (the invariant in Step 6).

```java
@Entity
@Table(name = "usr")
public class UserEntity extends AuditedEntity<Long> {

    @Column(name = "webauthn_user_handle")
    private byte[] webauthnUserHandle;

    public byte[] getWebauthnUserHandle() { return webauthnUserHandle; }
    public void setWebauthnUserHandle(byte[] handle) { this.webauthnUserHandle = handle; }
}
```

The column is nullable. Its nullability *is* the
"has-no-credentials" signal that prevents the device-side
`InvalidStateError` problem (Step 6).

## Step 3 — Add the credential entity and repository

`UserCredentialEntity` mirrors Spring Security's reference
`JdbcUserCredentialRepository` schema. Storing fields directly
avoids round-tripping through `WebauthnJacksonModule` — which
has no deserializer for `ImmutableCredentialRecord` and would
otherwise force a workaround.

```java
@Entity
@Table(name = "user_credential")
public class UserCredentialEntity extends AuditedEntity<Long> {

    @Column(name = "credential_id", nullable = false, unique = true)
    private byte[] credentialId;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "user_key", nullable = false)
    private UserEntity user;

    @Column(name = "user_entity_user_id", nullable = false)
    private byte[] userEntityUserId;  // mirrors WebAuthn's user-entity ID

    @Column(name = "public_key", nullable = false) private byte[] publicKey;
    @Column(name = "signature_count")              private long signatureCount;
    @Column(name = "uv_initialized")               private boolean uvInitialized;
    @Column(name = "backup_eligible")              private boolean backupEligible;
    @Column(name = "backup_state")                 private boolean backupState;
    @Column(name = "credential_type")              private String credentialType;
    @Column(name = "transports")                   private String transports;     // CSV
    @Column(name = "attestation_object")           private byte[] attestationObject;
    @Column(name = "attestation_client_data_json") private byte[] attestationClientDataJson;
    @Column(name = "registered_at", nullable = false) private Instant registeredAt;
    @Column(name = "last_used_at")                 private Instant lastUsedAt;
    @Column(name = "label")                        private String label;

    // getters / setters omitted
}
```

```java
public interface PasskeyCredentialRepository extends JpaRepository<UserCredentialEntity, Long> {
    Optional<UserCredentialEntity> findByCredentialId(byte[] credentialId);
    List<UserCredentialEntity> findByUser_Key(Long userKey);
    List<UserCredentialEntity> findByUserEntityUserId(byte[] userEntityUserId);
    boolean existsByUser_Key(Long userKey);
    long deleteByCredentialId(byte[] credentialId);
}
```

## Step 4 — Define the `AuditedPasskeyPrincipal`

Spring Security 7's `WebAuthnAuthentication.getPrincipal()`
returns a `PublicKeyCredentialUserEntity`. The kit's principal
also implements `AuditedPrincipal` so audit pipelines and
`CurrentUser.get()` see the same shape regardless of auth flow.

```java
public final class AuditedPasskeyPrincipal implements PublicKeyCredentialUserEntity, AuditedPrincipal {

    private final PublicKeyCredentialUserEntity delegate;
    private final Long key;
    private final String username;

    private AuditedPasskeyPrincipal(PublicKeyCredentialUserEntity delegate,
                                    Long key, String username) {
        this.delegate = delegate;
        this.key = key;
        this.username = username;
    }

    public static AuditedPasskeyPrincipal of(AuthenticatedUser user) {
        var delegate = ImmutablePublicKeyCredentialUserEntity.builder()
                .id(new Bytes(user.webauthnUserHandle()))
                .name(user.username())            // matches Authentication.getName()
                .displayName(user.displayName())  // for the device's passkey UI
                .build();
        return new AuditedPasskeyPrincipal(delegate, user.key(), user.username());
    }

    @Override public Long getKey()       { return key; }
    @Override public String getUsername() { return username; }

    // delegated PublicKeyCredentialUserEntity methods
    @Override public Bytes getId()             { return delegate.getId(); }
    @Override public String getName()          { return delegate.getName(); }
    @Override public String getDisplayName()   { return delegate.getDisplayName(); }
}
```

> The `AuthenticatedUser` record is whatever your `UserLookup`
> returns — at minimum `key`, `username`, `displayName`, and
> `webauthnUserHandle`.

## Step 5 — Implement the `WebAuthnUserAdapter`

Bridges Spring Security WebAuthn's user-entity API to your domain
user model. Spring calls into this adapter during registration
(`save`, `findByUsername`) and during assertion
(`findById` — by handle).

```java
@Component
@ConditionalOnProperty(name = "{app}.auth.passkey.enabled", havingValue = "true")
public class WebAuthnUserAdapter implements PublicKeyCredentialUserEntityRepository {

    private final UserLookup userLookup;
    private final PasskeyHandleManager handleManager;

    public WebAuthnUserAdapter(UserLookup userLookup, PasskeyHandleManager handleManager) {
        this.userLookup = userLookup;
        this.handleManager = handleManager;
    }

    @Override
    public PublicKeyCredentialUserEntity findById(Bytes id) {
        return userLookup.findByPasskeyHandle(id.getBytes())
                .map(AuditedPasskeyPrincipal::of)
                .orElse(null);
    }

    @Override
    public PublicKeyCredentialUserEntity findByUsername(String username) {
        return userLookup.findActiveUser(username)
                .map(u -> { handleManager.ensureHandleFor(u.key()); return u; })
                .map(AuditedPasskeyPrincipal::of)
                .orElse(null);
    }

    @Override
    public void save(PublicKeyCredentialUserEntity userEntity) {
        if (userEntity == null) return;
        userLookup.findActiveUser(userEntity.getName())
                .ifPresent(u -> handleManager.setHandle(u.key(), userEntity.getId().getBytes()));
    }

    @Override
    public void delete(Bytes id) {
        userLookup.findByPasskeyHandle(id.getBytes())
                .ifPresent(u -> handleManager.clearHandleByKey(u.key()));
    }
}
```

`UserLookup` (extended from form-login) gains two passkey-only
methods:

```java
Optional<AuthenticatedUser> findByPasskeyHandle(byte[] handle);
Optional<AuthenticatedUser> findActiveUser(String username);
```

## Step 6 — Implement the triple-role `JpaPasskeyService`

The class's three responsibilities map to three interfaces, but
they all share the same transactional context, the same
`SecureRandom`, and the same handle-lifecycle invariant. Splitting
across multiple classes would force coordination across
transactions; collapsing into one class keeps the invariant
local.

```java
@Service
public class JpaPasskeyService implements PasskeyService,
                                          PasskeyHandleManager,
                                          UserCredentialRepository {

    private static final int USER_HANDLE_BYTES = 16;

    private final UserLookup userLookup;
    private final UserRepository userRepository;
    private final PasskeyCredentialRepository credentialRepository;
    private final EntityManager entityManager;
    private final SecureRandom secureRandom = new SecureRandom();

    // ... constructor, then methods grouped by interface ...
}
```

The three interfaces:

```java
// User-facing — for the UI's "manage my passkeys" page.
public interface PasskeyService {
    boolean hasRegisteredPasskey(String username);
    boolean anyRegisteredPasskey();   // system-wide; used by login view
    List<RegisteredPasskey> listRegisteredPasskeys(String username);
    boolean deleteRegisteredPasskey(String username, long passkeyId);
}

public record RegisteredPasskey(long id, String label,
                                Instant registeredAt, Instant lastUsedAt) { }
```

```java
// Handle lifecycle — owns the webauthn_user_handle column.
public interface PasskeyHandleManager {
    byte[] ensureHandleFor(Long userKey);     // generate-if-null
    void   setHandle(Long userKey, byte[] handle);
    void   clearHandleByKey(Long userKey);
    void   clearHandleIfOrphaned(Long userKey);  // clears iff no credentials
}
```

`UserCredentialRepository` is Spring Security's interface — Spring
Security calls `save`, `findByCredentialId`, `findByUserId`,
`delete` on it.

### The handle invariant

> **`webauthn_user_handle` is null ⇔ user has no registered
> credentials.**

Why it matters: when a user has registered a passkey, the
device's passkey manager (Apple Keychain, iCloud, Windows Hello,
hardware key) caches the (rpId, userHandle) tuple. If the
server-side credential is deleted but the handle persists,
`navigator.credentials.create()` on the next register attempt
throws `InvalidStateError` because the device thinks a credential
already exists for that user — even though the server has none.

Maintain the invariant by:

1. Generating a handle in `ensureHandleFor` only when the user
   has no credentials and no handle yet (Step 5's
   `findByUsername` triggers this on register start).
2. Calling `clearHandleIfOrphaned` after every credential delete:

```java
@Override
@Transactional
public boolean deleteRegisteredPasskey(String username, long passkeyId) {
    Optional<AuthenticatedUser> userOpt = userLookup.findActiveUser(username);
    if (userOpt.isEmpty()) return false;
    Long userKey = userOpt.get().key();
    return credentialRepository.findById(passkeyId)
            .filter(cred -> userKey.equals(cred.getUser().getKey()))
            .map(cred -> {
                credentialRepository.delete(cred);
                clearHandleIfOrphaned(userKey);
                return true;
            })
            .orElse(false);
}
```

### Resolving the current user during credential save

Spring Security calls `UserCredentialRepository.save(record)`
during registration; the record itself doesn't carry the
authenticated user (registration requires a prior login). Resolve
through `CurrentUser` — the audited-principal helper — and
`EntityManager.getReference()`, mirroring the `AuditorAware`
pattern:

```java
private UserEntity resolveCurrentUser() {
    AuditedPrincipal principal = CurrentUser.get().orElseThrow(() ->
            new IllegalStateException(
                    "No authenticated user — passkey registration requires a prior login"));
    return entityManager.getReference(UserEntity.class, principal.getKey());
}
```

The full `save(CredentialRecord)` implementation handles two
distinct contexts in one method: (1) registration (record is new,
linked to current user) and (2) authentication update (Spring
Security re-saves the record post-assertion to bump the signature
counter and `lastUsedAt`).

```java
@Override
@Transactional
public void save(CredentialRecord credentialRecord) {
    byte[] credentialIdBytes = credentialRecord.getCredentialId().getBytes();
    UserCredentialEntity entity = credentialRepository.findByCredentialId(credentialIdBytes)
            .orElseGet(() -> {
                UserCredentialEntity fresh = new UserCredentialEntity();
                fresh.setCredentialId(credentialIdBytes);
                fresh.setUser(resolveCurrentUser());
                return fresh;
            });
    copyInto(entity, credentialRecord);          // field-by-field map
    credentialRepository.save(entity);
}
```

The `copyInto` helper translates between Spring Security's
`CredentialRecord` (with `Bytes` wrappers and an
`AuthenticatorTransport` enum set) and the JPA columns. The
inverse helper `toCredentialRecord` is needed for `findBy*` reads.

## Step 7 — Configure `WebAuthnRelyingPartyOperations`

The single bean Spring Security 7's WebAuthn DSL needs.
Conditional on passkey enabled — when off, neither this config
nor the conditional adapters load, so Spring Security's default
in-memory credential repository can't be wired by accident.

```java
@Configuration
@ConditionalOnProperty(name = "{app}.auth.passkey.enabled", havingValue = "true")
public class WebAuthnConfig {

    @Bean
    WebAuthnRelyingPartyOperations webAuthnRelyingPartyOperations(
            PublicKeyCredentialUserEntityRepository userEntityRepository,
            UserCredentialRepository credentialRepository,
            WebAuthnSettings settings) {
        String rpId   = requireProp(settings.rpId(),   "{app}.auth.webauthn.rp-id");
        String origin = requireProp(settings.origin(), "{app}.auth.webauthn.origin");
        return new Webauthn4JRelyingPartyOperations(
                userEntityRepository,
                credentialRepository,
                PublicKeyCredentialRpEntity.builder()
                        .id(rpId)
                        .name("{App Display Name}")
                        .build(),
                Set.of(origin));
    }

    private static String requireProp(String value, String key) {
        if (value == null || value.isBlank()) {
            throw new IllegalStateException(
                    "{app}.auth.passkey.enabled=true requires " + key);
        }
        return value;
    }
}
```

## Step 8 — Wire passkey into `SecurityConfig`

The WebAuthn DSL contributes filters that register the registration
and assertion endpoints (`/webauthn/register/options`,
`/webauthn/register`, `/webauthn/authenticate/options`,
`/login/webauthn`). Three coordination points to handle:

```java
@Bean
SecurityFilterChain securityFilterChain(
        HttpSecurity http,
        SecurityContextHolderStrategy securityContextHolderStrategy,
        AuthMethods authMethods)
        throws Exception {

    // 1. SecurityContextHolderStrategy gotcha — must be set BEFORE building
    //    the filter chain. Spring Security's WebAuthn filters capture the
    //    default ThreadLocal strategy at construction time; Vaadin's
    //    SmartInitializingSingleton swaps the global later, so the filters
    //    end up reading from a different ThreadLocal than
    //    SecurityContextHolderFilter writes to — silent 400s on
    //    /webauthn/register/options (empty auth). Injecting the bean and
    //    setting the global here guarantees every filter constructed during
    //    http.build() captures the Vaadin strategy.
    SecurityContextHolder.setContextHolderStrategy(securityContextHolderStrategy);

    http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
        configurer.loginView(LoginView.class);
        configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
    });

    if (authMethods.isPasskeyEnabled()) {
        // 2. Engage the WebAuthn DSL.
        http.webAuthn(Customizer.withDefaults());

        // 3. Permit the WebAuthn endpoints. View-based annotations remain
        //    the actual authz gate; this only releases the filter-chain
        //    block. Registration still requires a prior authenticated
        //    session — the filters check the SecurityContext.
        http.authorizeHttpRequests(auth -> auth.requestMatchers(
                "/webauthn/**", "/login/webauthn").permitAll());

        // 4. CSRF coordination with Vaadin: VaadinSecurityConfigurer above
        //    set up CSRF with Vaadin's UIDL-aware token strategy. Calling
        //    http.csrf(...) here *adds* ignore rules; it doesn't replace
        //    the configuration. WebAuthn endpoints are JS-driven (not
        //    form-based) and must be reachable without the UIDL CSRF token.
        http.csrf(csrf -> csrf.ignoringRequestMatchers(
                "/webauthn/**", "/login/webauthn"));
    }

    return http.build();
}
```

## Step 9 — Build the `PasskeyButton` Vaadin component

The browser-side WebAuthn ceremony runs in a Lit web component
(`{app}-passkey-button.ts`) so all `navigator.credentials.*`
logic stays client-side. The Java side is a thin Flow wrapper
that exposes a typed event API to views. Two factory methods:
`forAuthentication` (login view) and `forRegistration`
(account-preferences view).

```java
@Tag("{app}-passkey-button")
@JsModule("./{app}-passkey-button.ts")
public class PasskeyButton extends Component {

    private PasskeyButton() { }

    public static PasskeyButton forAuthentication(String label) {
        PasskeyButton b = new PasskeyButton();
        b.getElement().setAttribute("mode", "authenticate");
        b.getElement().setAttribute("label", label);
        return b;
    }

    public static PasskeyButton forRegistration(String label, String credentialLabel) {
        PasskeyButton b = new PasskeyButton();
        b.getElement().setAttribute("mode", "register");
        b.getElement().setAttribute("label", label);
        b.getElement().setAttribute("credential-label", credentialLabel);
        return b;
    }

    public Registration addSuccessListener(ComponentEventListener<PasskeySuccessEvent> listener) {
        return addListener(PasskeySuccessEvent.class, listener);
    }

    public Registration addErrorListener(ComponentEventListener<PasskeyErrorEvent> listener) {
        return addListener(PasskeyErrorEvent.class, listener);
    }

    @DomEvent("passkey-success")
    public static class PasskeySuccessEvent extends ComponentEvent<PasskeyButton> {
        private final String redirectUrl;
        public PasskeySuccessEvent(PasskeyButton source, boolean fromClient,
                                   @EventData("event.detail.redirectUrl") String redirectUrl) {
            super(source, fromClient);
            this.redirectUrl = redirectUrl;
        }
        public Optional<String> getRedirectUrl() {
            return Optional.ofNullable(redirectUrl).filter(s -> !s.isEmpty());
        }
    }

    @DomEvent("passkey-error")
    public static class PasskeyErrorEvent extends ComponentEvent<PasskeyButton> {
        private final String message;
        public PasskeyErrorEvent(PasskeyButton source, boolean fromClient,
                                 @EventData("event.detail.message") String message) {
            super(source, fromClient);
            this.message = message;
        }
        public String getMessage() {
            return message != null ? message : "Passkey operation failed";
        }
    }
}
```

The Lit component contract (sketch — the full implementation is a
separate concern):

- Authenticate mode: POST to `/webauthn/authenticate/options` →
  `navigator.credentials.get()` → POST credential to
  `/login/webauthn` → fire `passkey-success` with the redirect
  URL on the response, or `passkey-error` on any throw.
- Register mode: POST to `/webauthn/register/options` →
  `navigator.credentials.create()` → POST attestation to
  `/webauthn/register` with the `credential-label` attribute as
  the human-readable name → fire `passkey-success` (no redirect)
  or `passkey-error`.

## Step 10 — Place the button on login and account-preferences views

```java
// Login view — authenticate mode. Shown only when passkey is
// enabled AND at least one passkey exists in the system (the
// login view has no identified user, so we gate on system-wide
// presence to avoid showing a button no one can use).
if (authMethods.isPasskeyEnabled() && passkeyService.anyRegisteredPasskey()) {
    PasskeyButton signIn = PasskeyButton.forAuthentication("Sign in with passkey");
    signIn.addSuccessListener(e ->
            getUI().ifPresent(ui ->
                    ui.getPage().setLocation(e.getRedirectUrl().orElse("/"))));
    signIn.addErrorListener(e ->
            Notification.show("Passkey sign-in failed: " + e.getMessage()));
    card.add(signIn);
}
```

```java
// Account preferences — register mode. Shown only when the user
// has fewer than the project-imposed maximum (no kit limit).
PasskeyButton register = PasskeyButton.forRegistration(
        "Register a passkey", suggestedLabel);
register.addSuccessListener(_ -> {
    Notification.show("Passkey registered. Use it to sign in next time.");
    refreshPasskeyList();
});
register.addErrorListener(e ->
        Notification.show("Passkey registration failed: " + e.getMessage()));
add(register);
```

## Decisions this recipe imposes

- **One transactional service for three interfaces.** The handle
  invariant spans credential delete + handle clear; splitting
  across services would force cross-transaction coordination.
  `JpaPasskeyService` is the canonical home.
- **`webauthn_user_handle` is null ⇔ no credentials.** Stale
  handles cause device-side `InvalidStateError`. Always pair a
  credential delete with `clearHandleIfOrphaned`.
- **Direct-field credential persistence (no JSON round-trip).**
  `WebauthnJacksonModule` has no deserializer for
  `ImmutableCredentialRecord`; mapping fields explicitly is the
  cleanest path.
- **`SecurityContextHolderStrategy` set BEFORE filter-chain
  build.** Skipping this causes silent 400s on
  `/webauthn/register/options` because the filters capture the
  default ThreadLocal strategy at construction time, while
  Vaadin's `SmartInitializingSingleton` swaps it later.
- **WebAuthn endpoints permit-all at the filter chain.** The
  filters themselves enforce the auth requirement; the
  filter-chain `permitAll` only releases the URL-level block.
- **CSRF: ignore for WebAuthn endpoints.** They're JS-driven and
  must be reachable without Vaadin's UIDL CSRF token. Add the
  ignore rule alongside Vaadin's CSRF setup, don't replace it.
- **Lit component, not raw JS in Java.** WebAuthn ceremonies are
  complex async flows; a Lit web component is the natural unit.
  The Flow wrapper exposes typed events; the Java side stays
  pure Java.
- **Login button gates on `anyRegisteredPasskey()`.** No
  identified user at login time — gate system-wide. Don't show a
  button no one can use.

## What to verify

- `{app}.auth.passkey.enabled=true` requires
  `{app}.auth.form-login.enabled=true` — the conditional-auth
  validator throws on startup otherwise.
- With passkey enabled, the WebAuthn endpoints are reachable for
  an authenticated session: `GET /webauthn/register/options`
  returns 200 with a JSON challenge.
- Register-then-revoke cycle: register a passkey from account
  preferences, then revoke it — the user's
  `webauthn_user_handle` column is `NULL` again (run the SQL).
  Re-registering on the same device works without
  `InvalidStateError`.
- After a successful passkey sign-in, the principal is
  `AuditedPasskeyPrincipal` (cast and verify),
  `principal.getKey()` returns the user's `usr.key`, and
  `Authentication.getName()` returns the username (matching what
  was stored on the user-entity at register time).
- An audited write performed under a passkey-authenticated
  session populates `created_by_key` correctly — same audit
  pipeline as form login.
- `{app}.auth.passkey.enabled=false` removes the WebAuthn
  endpoints (404 / no filter), the conditional adapters from the
  context, and the `WebAuthnConfig` bean.

## Related

- [conditional-auth.md](conditional-auth.md) — the
  `passkey.enabled` flag, combinability rules, and the
  `requires form-login` constraint.
- [audited-principal.md](audited-principal.md) — the principal
  contract `AuditedPasskeyPrincipal` satisfies.
- [form-login.md](form-login.md) — passkey is supplementary;
  form login is the primary method.
- [oidc-sso.md](oidc-sso.md) — alternative auth method.
- `docs/patterns/architecture/security.md` — surrounding
  security architecture.
- `docs/solutions/` — your project's Lit component
  implementation, credential-label policy, registration-attempt
  limits, device-binding policy.
