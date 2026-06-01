# Compile-Time Layer Separation

When assembling the application, declare `{app}-jpaservice` as a runtime-only
dependency of `{app}-app` so the UI module cannot import JPA entities, repositories,
or service implementations at compile time — violations produce compiler errors
immediately, not runtime surprises.

```xml
<!-- In {app}-app/pom.xml -->
<dependency>
    <groupId>com.example</groupId>
    <artifactId>{app}-jpaservice</artifactId>
    <scope>runtime</scope>   <!-- enforces UI-to-persistence separation -->
</dependency>
```

What this enforces:
- UI code can call service interfaces (`{app}-service`)
- UI code can work with UI models (`{app}-uimodel`)
- UI code **cannot** import JPA entities, repositories, or mappers at compile time
- Service implementations **cannot** import Vaadin classes
