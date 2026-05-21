# Vaadin–Spring Integration

Spring bean registration, scopes, security annotations, and build configuration for Vaadin modules.

## `@SpringComponent` over `@Component` in Vaadin Modules

In modules that participate in Vaadin's UI lifecycle — views, layouts, listeners,
and anything that might grow a Vaadin scope (`@UIScope`, `@RouteScope`,
`@VaadinSessionScope`) — register Spring beans with Vaadin's `@SpringComponent`
rather than Spring's own `@Component`. `@SpringComponent` is `@Component` plus the
scope-aware plumbing that Vaadin's Flow-Spring integration expects; using plain
`@Component` on a UI bean works until the bean later grows a Vaadin scope and
silently fails to integrate with session or UI lifecycle.

**Project rule:** use `@SpringComponent` for every `@Component`-style bean in
the application's Vaadin-facing modules (typically `app` and `ui`). Non-Vaadin
modules — services, JPA model, repositories, utilities — continue to use `@Component`.

```java
// Preferred — in the app or ui module
import com.vaadin.flow.spring.annotation.SpringComponent;

@SpringComponent
public class AuthMethodCombinabilityValidator { ... }

// Avoid — plain @Component in a Vaadin module
@Component
public class AuthMethodCombinabilityValidator { ... }
```

`@Configuration` stays `@Configuration` — it doesn't need swapping. `@Route`-annotated
views are already scope-aware and don't need `@SpringComponent` on top.

## vaadin.allowed-packages

`application.properties` must declare `vaadin.allowed-packages` to include all Vaadin
component package prefixes used by the application. Without this, Vaadin components in
add-ons may fail to render.

```properties
vaadin.allowed-packages=com.vaadin,org.vaadin
```

When adding a new Vaadin add-on, add its root package prefix to this property.

## vaadin-dev Dependency

The `vaadin-dev` artifact (or equivalent) must be included as a dev-mode dependency with
`<optional>true</optional>` so it is present in development but excluded from production JARs:

```xml
<dependency>
    <groupId>com.vaadin</groupId>
    <artifactId>vaadin-dev</artifactId>
    <optional>true</optional>
</dependency>
```

## Access Annotations on Layout Classes

The main layout class should carry `@PermitAll` per Vaadin's
[security doc](https://vaadin.com/docs/latest/flow/security/enabling-security#error-messages-for-unauthorized-views).
This is the doc-recommended choice: anonymous users hitting protected routes
get redirected to login by Vaadin's `ExceptionHandlingConfigurer` instead of
seeing the layout chrome before the redirect.

> **Precondition:** `VaadinSecurityConfigurer` defaults `anyRequest` to
> `denyAll` at the URL filter level. Without overriding that default,
> `@PermitAll` on the layout silently fails — anonymous users get a bare
> HTTP 403 from Spring Security before Vaadin's exception handler runs,
> instead of being redirected to login. Override the default inside the
> configurer block:
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

The `@Layout` annotation itself was introduced in Vaadin 24.1 — on Vaadin
24.0, set the layout on each route via `@Route(value = "...", layout =
MainLayout.class)` and have `MainLayout` extend `AppLayout` (which already
implements `RouterLayout`).

```java
// Vaadin 24.1+ / 25+
@Layout
@PermitAll
public class MainLayout extends AppLayout {
    // DO NOT add "implements RouterLayout" — already inherited from AppLayout
}
```

```java
// Vaadin 24.0 — no @Layout; set layout per @Route
@Route(value = "items", layout = MainLayout.class)
public class ItemView extends Composite<VerticalLayout> { ... }
```

## MainLayout Must Not Implement RouterLayout

`AppLayout` already implements `RouterLayout`. Adding `implements RouterLayout` explicitly
is redundant and can cause unexpected behavior (especially in Vaadin ≥25 where the
access-checker treats layouts differently).

## Vaadin Session Scope

Use `@Scope("vaadin-session")` for beans that hold per-session state and must be
reachable from Push threads:

```java
@SpringComponent
@Scope(value = "vaadin-session", proxyMode = ScopedProxyMode.INTERFACES)
public class MySessionContext implements MyContextInterface {
    // per-session state, e.g. a ValueSignal<...>
}
```

### Why not `@SessionScope`

Spring's `@SessionScope` binds to `HttpSession`. Vaadin's Push threads carry no
HTTP request context, so any attempt to resolve an `@SessionScope` bean from a Push
callback throws `No thread-bound request found`. `vaadin-session` is backed by
`VaadinSession`, which is accessible on both HTTP request and Push threads.

### Why `proxyMode = ScopedProxyMode.INTERFACES`

Session-scoped beans are typically injected into longer-lived beans (singletons,
UI-scoped components). Without a scoped proxy, Spring injects the raw instance at
construction time — before any session exists — and the same stale instance is
shared across all sessions. With `ScopedProxyMode.INTERFACES`, Spring injects a
proxy; each call is routed to the instance for the *current* session at invocation
time.

`ScopedProxyMode.INTERFACES` requires the bean to be referenced via an interface.
`ScopedProxyMode.TARGET_CLASS` works without an interface but requires CGLIB
subclassing; prefer the interface form.

### When to use

- Per-session reactive state that Push callbacks read or write (e.g., a
  `ValueSignal` carrying impersonation state, active filter presets).
- Context objects that cross the module boundary (declared as an interface in a
  non-Vaadin module; implemented with `@SpringComponent` + `@Scope("vaadin-session")`
  in the UI module).

### When not to use

- **Stateless beans** — `@SpringComponent` alone (singleton) is correct.
- **Per-UI-instance state** — use `@UIScope` (one instance per browser tab, not per
  session).
- **View-local state** — keep it as a field on the view class itself.
