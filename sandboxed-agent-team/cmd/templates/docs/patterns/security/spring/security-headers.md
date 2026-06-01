# Security Response Headers

When adding HTTP security headers to a Vaadin application, configure them
through the `VaadinWebSecurity.configure(...)` override and respect Vaadin's
required CSP relaxations so the client runtime is not broken.

## What Vaadin Does and Does Not Set

- **`X-Frame-Options`** — **not set by default.** Per Vaadin's
  [Frequently Reported Issues](https://vaadin.com/docs/latest/flow/security/advanced-topics/frequent-issues)
  page: Vaadin doesn't automatically set the header because many applications
  need to run inside frames. If the application is not expected to be embedded
  in an iframe, set it explicitly for clickjacking protection.
- **Content-Security-Policy** — Vaadin's bootstrap requires specific relaxations:
  - `script-src 'unsafe-inline' 'unsafe-eval'`
  - `style-src 'unsafe-inline'`

  These are architectural limitations in Vaadin's client-side engine. Do not
  attempt to tighten these directives; doing so breaks Vaadin's client runtime.

## Adding a Header

Add security headers through Spring Security's `HttpSecurity.headers(...)` DSL
inside the `VaadinWebSecurity.configure(...)` override:

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
`Referrer-Policy` is not set by default; add it in the same `headers(...)`
block if needed.

