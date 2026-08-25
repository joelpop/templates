# @SpringComponent in Vaadin Modules

When registering a Spring bean in the `ui` module, use `@SpringComponent`
instead of `@Component` so imports do not conflict with
`com.vaadin.flow.component.Component`.

```java
// Avoid — plain @Component in a Vaadin module
@Component
public class AuthMethodCombinabilityValidator { /* ... */ }
```

```java
// Preferred — in the ui module
@SpringComponent
public class AuthMethodCombinabilityValidator { /* ... */ }
```

`@Configuration` stays `@Configuration` — it doesn't need swapping. `@Route`-annotated
views are already scope-aware and don't need `@SpringComponent` on top.
