# Main Layout Access Annotation

When annotating the main layout class, use `@PermitAll` so anonymous users
hitting protected routes are cleanly redirected to login by
`ExceptionHandlingConfigurer` instead of seeing layout chrome before the
redirect.

## Vaadin 24.x

`@PermitAll` works without additional configuration.

## Vaadin ≥ 25

`VaadinSecurityConfigurer` defaults `anyRequest` to `denyAll` at the URL filter
level. Without overriding that default, `@PermitAll` on the main layout silently
fails — anonymous users get a bare HTTP 403 from Spring Security before Vaadin's
exception handler runs. Override the default inside the configurer block:

```java
http.with(VaadinSecurityConfigurer.vaadin(), configurer -> {
    configurer.loginView(LoginView.class);
    configurer.anyRequest(AuthorizeHttpRequestsConfigurer.AuthorizedUrl::permitAll);
});
```

Without this override, `@AnonymousAllowed` on the layout is the workaround that
keeps things working — it puts the layout into Vaadin's `defaultPermitMatcher`,
sidestepping the URL-level `denyAll`. It is a workaround, not the intended
design: it exposes layout chrome to anonymous users momentarily before the login
redirect, and the Vaadin docs explicitly recommend `@PermitAll`. Use the
`anyRequest=permitAll` override to make `@PermitAll` work as documented; fall
back to `@AnonymousAllowed` only when touching the configurer is not an option.

## Each View Still Controls Its Own Access

The layout annotation only determines whether the layout itself blocks
navigation. Each view's own annotation still controls which users can reach that
view.

See `docs/patterns/ui/vaadin/layout-setup.md` for `@Layout` annotation
mechanics and Vaadin 24.0 vs 24.1+ version compatibility.
