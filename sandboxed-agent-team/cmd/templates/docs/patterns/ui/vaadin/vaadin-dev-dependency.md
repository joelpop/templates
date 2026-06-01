# Vaadin Dev Dependency

When configuring a Vaadin 25+ project, declare `vaadin-dev` as an optional
dependency so dev tools (Vite dev server, live reload, Copilot) are available
in development but excluded from production JARs.

`vaadin-dev` is not included transitively in Vaadin 25+ — it must be declared
explicitly.

```xml
<dependency>
    <groupId>com.vaadin</groupId>
    <artifactId>vaadin-dev</artifactId>
    <optional>true</optional>
</dependency>
```

In Vaadin 24, dev tools were included transitively; live reload came from
`spring-boot-devtools` instead.