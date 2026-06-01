# @SpringComponent in Vaadin Modules

When registering a Spring bean in a Vaadin-facing module (`app` or `ui`),
use `@SpringComponent` instead of `@Component` so the bean participates in
Vaadin's UI lifecycle and does not silently fail when it later grows a Vaadin
scope.

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
// Avoid — plain @Component in a Vaadin module
@Component
public class AuthMethodCombinabilityValidator { ... }
```

```java
// Preferred — in the app or ui module
import com.vaadin.flow.spring.annotation.SpringComponent;

@SpringComponent
public class AuthMethodCombinabilityValidator { ... }
```

`@Configuration` stays `@Configuration` — it doesn't need swapping. `@Route`-annotated
views are already scope-aware and don't need `@SpringComponent` on top.
