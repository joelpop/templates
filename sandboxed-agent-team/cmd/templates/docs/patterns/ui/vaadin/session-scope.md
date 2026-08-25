# Vaadin Session Scope

When a Spring bean holds per-session state, use
`@Scope(VaadinSessionScope.VAADIN_SESSION_SCOPE_NAME)` with
`proxyMode = ScopedProxyMode.INTERFACES` so the bean lifecycle is tied to
`VaadinSession` rather than `HttpSession` and is correct whether push is
enabled or not.

```java
@SpringComponent
@Scope(value = VaadinSessionScope.VAADIN_SESSION_SCOPE_NAME, proxyMode = ScopedProxyMode.INTERFACES)
public class MySessionContext implements MyContextInterface {
    // per-session state, e.g. a ValueSignal<...>
}
```

## Why not `@SessionScope`

Spring's `@SessionScope` binds to `HttpSession`. Vaadin's Push threads carry no
HTTP request context, so any attempt to resolve an `@SessionScope` bean from a Push
callback throws `No thread-bound request found`. `VaadinSessionScope.VAADIN_SESSION_SCOPE_NAME`
is backed by `VaadinSession`, which is accessible on both HTTP request and Push threads.

## Why `proxyMode = ScopedProxyMode.INTERFACES`

Session-scoped beans are typically injected into longer-lived beans (singletons,
UI-scoped components). Without a scoped proxy, Spring injects the raw instance at
construction time — before any session exists — and the same stale instance is
shared across all sessions. With `ScopedProxyMode.INTERFACES`, Spring injects a
proxy; each call is routed to the instance for the *current* session at invocation
time.

`ScopedProxyMode.INTERFACES` requires the bean to be referenced via an interface.
`ScopedProxyMode.TARGET_CLASS` works without an interface but requires CGLIB
subclassing; prefer the interface form.

## When to use

- Per-session reactive state that Push callbacks read or write (e.g., a
  `ValueSignal` carrying impersonation state, active filter presets).
- Context objects that cross the module boundary (declared as an interface in a
  non-Vaadin module; implemented with `@SpringComponent` + `@Scope(VaadinSessionScope.VAADIN_SESSION_SCOPE_NAME)`
  in the UI module).

## When not to use

- **Stateless beans** — `@SpringComponent` alone (singleton) is correct.
- **Per-UI-instance state** — use `@UIScope` (one instance per browser tab, not per
  session).
- **View-local state** — keep it as a field on the view class itself.
