# Executable Fat JAR

When deploying a Spring Boot application as a standalone process without an external
servlet container, use the Spring Boot Maven plugin to produce a self-contained
executable JAR so all dependencies are bundled and the application runs with `java -jar`.

The `spring-boot-maven-plugin` must be declared in the app module's POM:

```xml
<plugin>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-maven-plugin</artifactId>
</plugin>
```

```bash
mvn clean package -Pproduction
java -jar {app}-app/target/{app}-app.jar
```

## Related

- `docs/patterns/cicd/deployment/recipes/war.md` — alternative: WAR packaging for deployment to an existing servlet container.
