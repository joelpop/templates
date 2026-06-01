# WAR Deployment

When deploying to an existing servlet container (Tomcat, Jetty, WildFly, or similar),
package as a WAR file so the container manages the JVM and multiple applications can
share it — requires three changes to a default Spring Boot fat-JAR project.

**1 — Set WAR packaging in the app module POM:**

```xml
<packaging>war</packaging>
```

**2 — Add `extends SpringBootServletInitializer` to the application entry point** (see
`docs/patterns/structure/modules/app-entry-point.md` for the base class shape):

```java
public class Application extends SpringBootServletInitializer implements AppShellConfigurator {

    @Override
    protected SpringApplicationBuilder configure(SpringApplicationBuilder builder) {
        return builder.sources(Application.class);
    }
    // ...
}
```

**3 — Mark the embedded Tomcat as `provided` so it is excluded from the WAR:**

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-tomcat</artifactId>
    <scope>provided</scope>
</dependency>
```

Build and deploy:

```bash
mvn clean package -Pproduction
# deploy {app}-app/target/{app}-app.war to the container
```

## Related

- `docs/patterns/cicd/deployment/recipes/fat-jar.md` — alternative: standalone executable JAR with no external container.