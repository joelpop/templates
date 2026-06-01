# Application Entry Point and Route Scanning

When placing `Application.java`, put it in the root package `{base_package}` so that
Vaadin automatically scans `{base_package}.ui.*` subpackages for `@Route`-annotated
views. Avoid `@EnableVaadin` — Vaadin's Spring Boot starter auto-configures Vaadin automatically;
no explicit opt-in is needed. The annotation exists for Spring MVC package scanning only.
Adding it to a Spring Boot application is a common mistake when copying from Spring MVC
examples.

```java
@SpringBootApplication
@StyleSheet(Lumo.STYLESHEET)
@StyleSheet(Lumo.UTILITY_STYLESHEET)
public class Application implements AppShellConfigurator {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
```

Exactly one class in the application may implement `AppShellConfigurator`. In a
multi-module project it must be in the deployable module (`{app}-app`), not a
dependency JAR.

## Related

- `docs/patterns/cicd/deployment/war.md` — WAR packaging adds `extends SpringBootServletInitializer` to this class.
