# Vaadin Spring Boot Configuration

When configuring Spring Boot for a Vaadin application, declare
`vaadin.allowed-packages` in `application.properties` and include `vaadin-dev`
as an optional dependency so Vaadin add-on components resolve at runtime and
dev tools are excluded from production JARs.

## vaadin.allowed-packages

`application.properties` must declare `vaadin.allowed-packages` to include all
Vaadin component package prefixes used by the application. Without this, Vaadin
components in add-ons may fail to render.

```properties
vaadin.allowed-packages=com.vaadin,org.vaadin
```

When adding a new Vaadin add-on, add its root package prefix to this property.

## vaadin-dev Dependency

The `vaadin-dev` artifact (or equivalent) must be included as a dev-mode dependency
with `<optional>true</optional>` so it is present in development but excluded from
production JARs:

```xml
<dependency>
    <groupId>com.vaadin</groupId>
    <artifactId>vaadin-dev</artifactId>
    <optional>true</optional>
</dependency>
```
