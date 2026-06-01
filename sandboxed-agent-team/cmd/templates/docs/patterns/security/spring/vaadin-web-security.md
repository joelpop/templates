# Extending VaadinWebSecurity

When configuring Spring Security for a Vaadin application, extend
`VaadinWebSecurity` rather than `WebSecurityConfigurerAdapter` so session
fixation, CSRF handling, Vaadin internal request matchers, and the
`AnnotatedViewAccessChecker` integration are all configured correctly without
re-declaring framework defaults.

Version-sensitive notes are inline; see `docs/patterns/README.md` →
"Version Compatibility" for the summary matrix.

## Configuration Class

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

## What VaadinWebSecurity Provides

Do not re-declare these manually — re-declaring them overrides Vaadin's defaults
in subtle, breakage-prone ways:

- **Session fixation protection** — regenerates the session ID on successful authentication.
- **CSRF protection** — enabled for non-Vaadin endpoints; Vaadin's own internal requests
  are carved out so the framework's built-in CSRF handling works correctly.
- **Vaadin internal request matchers** — Vaadin's servlet paths, Push endpoints, static
  resources, and development-mode URLs are permitted without authentication.
- **Access checker wiring** — `AnnotatedViewAccessChecker` is registered to enforce
  `@AnonymousAllowed` / `@PermitAll` / `@RolesAllowed` / `@DenyAll` on routes.
- **Login / logout endpoints** — `setLoginView(...)` wires form-login and logout redirects
  to your login view.

**Verify session fixation works:** the session cookie value changes between
pre-login and post-login responses.

