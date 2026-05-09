# Security Patterns

Authentication, authorization, session management, and hardening patterns for
Vaadin 24+ with Spring Security (Spring Boot 3+). Version-sensitive notes inline;
see `docs/patterns/README.md` → "Version Compatibility" for the summary matrix.

## Authentication

### Password Authentication with BCrypt

Passwords must be stored as BCrypt hashes — never in plaintext or with reversible encryption.

```java
@Bean
public PasswordEncoder passwordEncoder() {
    return new BCryptPasswordEncoder(10); // work factor >= 10
}
```

The `users` table contains no plaintext password column. The hash column stores BCrypt
format (`$2a$...`).

### Entropy-Based Password Validation

Reject passwords by entropy, not by character-class rules (uppercase +
number + symbol rules are gamed easily and frustrate users):

- Minimum entropy: 50 bits
- Minimum length: 8 characters
- Maximum length: 128 characters
- Common/breached passwords rejected via blocklist check
- Display a visual strength indicator (Weak / Fair / Good / Strong / Very Strong) during entry

```java
// Example validation at service layer
public void validatePasswordStrength(String password) {
    if (password.length() < 8 || password.length() > 128) {
        throw new ValidationException("Password must be 8–128 characters.");
    }
    if (entropyBits(password) < 50) {
        throw new ValidationException("Password is too weak. Try adding more variety.");
    }
    if (blocklist.contains(password.toLowerCase())) {
        throw new ValidationException("This password is too common. Please choose another.");
    }
}
```

### Passkey / WebAuthn Authentication

Support passkey-based authentication as an alternative to passwords:

- Users register a passkey from their profile/preferences view
- The login view displays a "Sign in with Passkey" button alongside the password form
- Clicking invokes the browser's WebAuthn API
- Supported authenticators: platform (TouchID, FaceID, Windows Hello, Android biometrics)
  and roaming (YubiKey, FIDO2 security keys)
- Successful passkey authentication creates a valid session identical to password auth
- Failure falls back gracefully to password login with a user-friendly error

### Authentication Error Messages

Do not reveal whether an email address exists in the system:

- "Email not found" → display: "Incorrect email or password"
- "Wrong password" → display: "Incorrect email or password"

Both conditions produce the same generic message.

### Account Lockout Messaging

Communicate lockout status without revealing security details:

```
"Account is locked. Contact an administrator."
```

Server logs record the lockout reason and actor. The user never sees threshold values,
attempt counts, or countdown timers.

### Session Fixation Protection

Session fixation protection is handled by **`VaadinWebSecurity`**
(`com.vaadin.flow.spring.security`), the Vaadin-supplied base class for
Spring Security configuration in Vaadin 24+ Flow applications. Extend it
for your security configuration and you get session-ID regeneration on
login — along with the rest of the Vaadin-Spring-Security integration
tuned for a server-side SPA — without writing the
`http.sessionManagement(...)` DSL yourself.

```java
@EnableWebSecurity
@Configuration
public class SecurityConfig extends VaadinWebSecurity {

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        // Add your app-specific overrides (login view, custom authorize rules, etc.)
        super.configure(http);
        setLoginView(http, LoginView.class);
    }
}
```

What `VaadinWebSecurity` configures for you (do not re-declare these manually):
- **Session fixation protection** — regenerates the session ID on successful authentication.
- **CSRF protection** — enabled for non-Vaadin endpoints; Vaadin's own internal requests
  are carved out so the framework's built-in CSRF handling works correctly (see "CSRF
  Protection" below).
- **Vaadin internal request matchers** — Vaadin's servlet paths, Push endpoints, static
  resources, and development-mode URLs are permitted without authentication so the
  framework can function.
- **Access checker wiring** — `AnnotatedViewAccessChecker` (or the equivalent on your
  Vaadin version) is registered to enforce `@AnonymousAllowed` / `@PermitAll` /
  `@RolesAllowed` / `@DenyAll` on routes.
- **Login / logout endpoints** — `setLoginView(...)` wires form-login and logout redirects
  to your login view.

**Do not** mirror these with raw Spring Security DSL — re-declaring them
can override Vaadin's defaults in subtle, breakage-prone ways.

**Verify:** the session cookie value changes between pre-login and
post-login responses.

### Session Timeout

Sessions expire after a configured period of inactivity:

```properties
server.servlet.session.timeout=30m
```

When a session expires during an active Vaadin session, Vaadin's
built-in handling displays a "Session expired" notification, offers a
page reload, and redirects the user to login.

### Session Cookie Flags

```properties
server.servlet.session.cookie.http-only=true
server.servlet.session.cookie.secure=true
server.servlet.session.cookie.same-site=strict
```

`HttpOnly` prevents JavaScript access. `Secure` prevents transmission over HTTP.
`SameSite=Strict` prevents CSRF via cross-site requests. Apply in production profiles.

## Authorization — RBAC with Jakarta Security Annotations

### View-Level Access Control

Every `@Route`-annotated class must carry exactly one of:

| Annotation | Meaning |
|------------|---------|
| `@AnonymousAllowed` | Public — no authentication required |
| `@PermitAll` | Any authenticated user |
| `@RolesAllowed(...)` | Only users with one of the listed roles |
| `@DenyAll` | No one (used to explicitly block a view) |

A view without an access annotation is a security defect.

### Role Name Constants

Define the project's roles as a single enum owning both the canonical
name (the exact string Spring Security uses for the granted authority)
and the human-facing metadata (display label, description). The enum
exposes a parallel set of `public static final String ROLE_*` constants
whose values match each enum's canonical name. `@RolesAllowed` and
`hasRole(...)` reference the constants; UI code uses the enum entries
for display and runtime checks.

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

**Use the constants in `@RolesAllowed`:**

```java
@Route("admin/user")
@RolesAllowed({UserRole.ROLE_ADMIN, UserRole.ROLE_MANAGER})
public class UserView extends BaseView { ... }
```

Not `@RolesAllowed({"ADMIN", "MANAGER"})`. String literals scatter role
names across the codebase; renaming becomes a project-wide grep instead
of a single-file edit, and typos are silent (a misspelled `"ADMNI"`
compiles but never matches anyone).

**Where the enum lives.** Place it in a module accessible from both
persistence (where the role is stored on the user record) and UI (where
the constants are referenced in `@RolesAllowed`). In a typical
multi-module Vaadin/Spring layout, that's the `uimodel` module —
`service` already depends on it, and so does `ui` transitively. Putting
it in the JPA model module would force UI to reach across an
architectural boundary just to gate routes.

**Why a `getSecurityName()` getter when the value matches the enum
name.** Spring Security represents authorities as strings (typically
`ROLE_<NAME>`), and some integrations need a method-of-an-instance
accessor — e.g., when mapping a `UserEntity` field to a
granted-authority list. The getter is the seam; the constants are the
compile-time form for annotations.

### Access Annotation on the Main Layout

The main layout class should carry `@PermitAll` — that's Vaadin's documented
recommendation, and it lets anonymous users hitting protected routes get cleanly
redirected to login by `ExceptionHandlingConfigurer` instead of seeing layout
chrome before the redirect.

> **Precondition (Vaadin ≥25):** `VaadinSecurityConfigurer` defaults
> `anyRequest` to `denyAll` at the URL filter level. Without overriding that
> default, `@PermitAll` on the main layout silently fails — anonymous users
> hitting protected routes get a bare HTTP 403 from Spring Security before
> Vaadin's exception handler runs, instead of being redirected to login.
> Override the default inside the configurer block:
>
> ```java
> http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
>     configurer.loginView(LoginView.class);
>     configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
> });
> ```
>
> Without this override, `@AnonymousAllowed` on the layout is the workaround
> that keeps things working — it puts the layout into Vaadin's
> `defaultPermitMatcher` ("anonymous routes"), sidestepping the URL-level
> `denyAll`. It's a workaround, not the intended design: it exposes layout
> chrome to anonymous users momentarily before the login redirect, and the
> Vaadin docs explicitly recommend `@PermitAll`. Use the `anyRequest=permitAll`
> override to make `@PermitAll` work as documented; fall back to
> `@AnonymousAllowed` only when touching the configurer isn't an option.

On **Vaadin 24.x**, the access-checker issue does not apply and `@PermitAll`
works without the `anyRequest=permitAll` override.

See `docs/patterns/conventions/vaadin.md` → "Access Annotations on Layout
Classes" for additional context. Each view's own annotation still controls
access — the layout annotation only determines whether the layout itself
blocks navigation.

### No @PreAuthorize on Views

Don't use `@PreAuthorize` (Spring Security method security) on Vaadin
views. Access control is enforced at the view level via Jakarta Security
annotations and Vaadin's `AnnotatedViewAccessChecker` (or its older
equivalent). `@PreAuthorize` adds confusion and may not behave as
expected with Vaadin's navigation model.

### Role-Based Rendering — Security Application of the UI Rubrics

Authorization gating at the UI layer must use the "do not generate" mode:
a component or navigation item the current user's role doesn't authorize
must never be constructed, never added to the layout, and therefore
never present in the DOM — leaving no artifact for the user to discover
or re-enable via the browser inspector. `setVisible(false)` and CSS
concealment are **not** acceptable: the underlying component is still
present and (if the server considers it enabled) interactive.

For the full three-mode rubric (do not generate / hide / disable),
layout-preservation guidance, and the Vaadin-server-state authority
rule, see:

- `docs/patterns/ui/components.md` → "Conditional Component Rendering —
  Do Not Generate vs. Hide vs. Disable" and its "Layout Preservation —
  When a Placeholder Is Needed" subsection.
- `docs/patterns/ui/navigation.md` → "Conditional Navigation Rendering"
  for `SideNavItem` / menu entries.

This section stops at the security *rule*; the canonical UI patterns
live with the component and navigation docs.

### Self-Editing Restrictions

Implement service-layer guards for self-editing actions that would leave
the system inconsistent:

```java
// Service implementation
@Transactional
public void deactivate(long key) {
    if (key == currentUser.getKey()) {
        throw new ValidationException("You cannot deactivate your own account.");
    }
    // ...
}
```

Mirror these in the UI (disabled button with tooltip) but enforce them
in the service layer so they hold against programmatic callers.

## Security Headers

Vaadin's documented behavior for HTTP security headers is narrower than
generic Spring Security guidance. Start from Vaadin's rules; add more
only when you have a concrete reason.

### What Vaadin Does and Does Not Set

- **`X-Frame-Options`** — **not set by default.** Per Vaadin's
  [Frequently Reported Issues](https://vaadin.com/docs/latest/flow/security/advanced-topics/frequent-issues)
  page: *"Vaadin doesn't automatically set the X-Frame-Options HTTP header, because many
  times applications need to run inside frames."* If the application is not expected to be
  embedded in an iframe, set it explicitly for clickjacking protection (example below).
- **Content-Security-Policy (CSP)** — Vaadin's bootstrap requires specific relaxations to
  function. The Vaadin docs state these are *"architectural limitations in Vaadin, so that
  the framework can start its client-side engine in the browser"*:
  - `script-src 'unsafe-inline' 'unsafe-eval'`
  - `style-src 'unsafe-inline'`

  Do not attempt to tighten these directives as part of the "normal" CSP configuration;
  doing so will break Vaadin's client runtime.

### Adding a Header via `VaadinWebSecurity`

Add security headers through Spring Security's `HttpSecurity.headers(...)`
DSL within the `VaadinWebSecurity.configure(...)` override — `http` comes
from there. Nothing Vaadin-specific wraps the headers API; the idiom is
just to keep everything inside the `VaadinWebSecurity` subclass.

```java
@EnableWebSecurity
@Configuration
public class SecurityConfig extends VaadinWebSecurity {

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        super.configure(http);
        setLoginView(http, LoginView.class);

        // Vaadin does not set X-Frame-Options by default. Add it here if the
        // application should not be embedded in iframes.
        http.headers(headers -> headers
            .frameOptions(frame -> frame.sameOrigin())
        );
    }
}
```

Spring Security's defaults cover `X-Content-Type-Options: nosniff`,
cache-control headers, and (over HTTPS) `Strict-Transport-Security`.
`Referrer-Policy` is not set by default; add it in the same
`headers(...)` block if needed. These are generic Spring Security
concerns — Vaadin has no documented guidance — so follow standard Spring
Security / OWASP practice.

### Strict CSP (Nonce-Based)

If a strict CSP beyond Vaadin's default relaxations is needed, do **not**
use Spring Security's `headers().contentSecurityPolicy(...)` — use
Vaadin's [nonce-based strict CSP](https://vaadin.com/docs/latest/flow/security/advanced-topics/strict-csp)
mechanism. It requires:

- An `IndexHtmlRequestListener` that generates a random nonce per request and sets the
  `Content-Security-Policy` header on the response with `script-src 'nonce-...'`.
- Nonce injection into script tags in the index file.
- JavaScript overrides that replace `Function` / `eval()` with CSP-compliant versions.

Production-mode only with prerequisites — follow the linked Vaadin docs
rather than reproducing the recipe here.

## CSRF Protection

`VaadinWebSecurity` configures CSRF protection automatically for Vaadin
endpoints. No manual CSRF token handling is required in views.

Endpoints outside Vaadin's filter chain (actuator, custom REST) must
implement CSRF protection independently. In Phase 1, no such endpoints
should exist.

## Rate Limiting

Limit login attempts per IP to prevent brute-force attacks:

- Threshold: more than 10 failed attempts within 5 minutes from a single IP
- Response: temporary block or CAPTCHA challenge
- Log the rate-limit event server-side

Return HTTP 429 with a generic "Too many requests" message. Do not
reveal thresholds or countdowns in any header or body.

The implementation (in-memory, Redis, database) is a per-project
architectural decision.

## Input Validation and SQL Injection Prevention

Use Spring Data JPA repository methods and JPQL with named parameters exclusively:

```java
// Correct — parameterized JPQL
@Query("SELECT e FROM EmployeeEntity e WHERE e.department.key = :deptKey AND e.active = true")
List<EmployeeListItemProjection> findActiveByDepartment(@Param("deptKey") Long deptKey);

// Never — string concatenation in queries
@Query("SELECT e FROM EmployeeEntity e WHERE e.name = '" + name + "'")  // NEVER
```

No `nativeQuery = true` with string concatenation. All custom JPQL uses
`:param` named parameters.

Vaadin's component layer escapes output by default, preventing XSS from
user-supplied data rendered in the UI.

## Sensitive Information Leakage Prevention

Error responses must not expose:
- Stack traces
- Internal file paths
- Database error details or constraint names
- Server version information
- Framework version information

Stack traces appear only in server-side logs. Spring Boot's default `/error` endpoint must
be suppressed or overridden. See `docs/patterns/ui/error-views.md` for error view patterns.

HTTP response headers must not reveal the server technology stack:
- `Server` header: absent or generic
- `X-Powered-By` header: absent

## Dependency Security Scanning

Run OWASP Dependency-Check or equivalent as part of CI:
- Fail the build on any dependency with CVSS score >= 9.0 (critical)
- Review and remediate or accept high-severity CVEs (7.0–8.9) within 30 days

## Data Protection

### Secrets Not in Source Control

`application.properties` must contain no plaintext secrets. Inject secrets via environment
variables using `${ENV_VAR_NAME}` syntax:

```properties
spring.datasource.password=${DB_PASSWORD}
```

A scan of the repository history should find no committed secrets.

### TLS in Production

All client-server communication in production uses TLS (HTTPS). HTTP requests are
redirected to HTTPS. The application never serves content over plain HTTP in production
or staging.

### PII Not in Logs

Personally identifiable information — user names, email addresses, contact details — must
not appear in application logs at INFO level or below. Error logs may include user
identifiers (surrogate keys) for correlation but not display names or contact details.

### File Upload Validation

Validate uploaded files for type and size before storage:

```java
if (file.size() > 2 * 1024 * 1024) {
    throw new ValidationException("File must be 2 MB or smaller.");
}
if (!Set.of("image/jpeg", "image/png").contains(file.contentType())) {
    throw new ValidationException("Only JPEG and PNG files are accepted.");
}
```

Return a user-facing error on invalid uploads — not a stack trace.

### BCrypt Password Storage

```java
// Store
String hash = passwordEncoder.encode(rawPassword);

// Verify
boolean valid = passwordEncoder.matches(rawPassword, storedHash);
```

Work factor must be >= 10. Never store raw passwords or reversible hashes.
