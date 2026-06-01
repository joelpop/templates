# Avoid @PostConstruct for Constructor-Eligible Work

When startup logic can run before other beans are needed, put it in the constructor
rather than `@PostConstruct` — a constructor that throws fails immediately with a
clear stack trace; `@PostConstruct` defers the failure to a later lifecycle phase
with no benefit.

```java
// Avoid — @PostConstruct for work the constructor can do
@SpringComponent
public class AuthMethodValidator {
    @Autowired private AuthProperties props;

    @PostConstruct
    void validate() { /* ... */ }
}
```

```java
// Preferred — constructor handles it directly
@SpringComponent
public class AuthMethodValidator {
    public AuthMethodValidator(AuthProperties props) {
        if (!props.formLogin().enabled()
                && !props.passkey().enabled()
                && !props.sso().enabled()) {
            throw new IllegalStateException(
                    "No authentication methods enabled — at least one of " +
                    "auth.{form-login,passkey,sso}.enabled must be true.");
        }
    }
}
```

Reserve `@PostConstruct` for work that genuinely requires all beans' lifecycles to
be complete before it can begin (e.g., cross-bean warm-up caches).