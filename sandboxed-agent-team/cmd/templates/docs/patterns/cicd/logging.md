# Application Logging

When logging application events, use SLF4J with Logback and structured output
in non-dev profiles so log entries are parseable by aggregation pipelines and
configuration mistakes are diagnosable from startup output.

## SLF4J + Logback

Application code logs through **SLF4J** — declared via `@Slf4j` (see
`docs/patterns/language/java/lombok.md` → "for Logging"). The implementation is
**Logback**, pulled in transitively by `spring-boot-starter-logging`. The same
starter bridges `java.util.logging` and Log4j, so third-party libraries route
through the same pipeline.

- **Code imports** `org.slf4j.Logger` — never `ch.qos.logback.*` or
  `java.util.logging.*`.
- **Configuration** uses `logback-spring.xml` for structured setup,
  `application.properties` `logging.level.*` for level tuning.

## Structured Logging

Use structured logging (JSON or key-value pairs) for log aggregation
(CloudWatch, Datadog, etc.). Each field is a distinct key — not embedded in the
message string.

```xml
<!-- logback-spring.xml — JSON output in non-dev profiles -->
<springProfile name="staging,prod">
    <appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
        <encoder class="net.logstash.logback.encoder.LogstashEncoder"/>
    </appender>
    <root level="WARN">
        <appender-ref ref="JSON"/>
    </root>
</springProfile>

<springProfile name="dev">
    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%d{HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>
    <root level="DEBUG">
        <appender-ref ref="CONSOLE"/>
    </root>
</springProfile>
```

## Startup Logging

The application logs resolved configuration at INFO level on startup:

```java
@EventListener(ApplicationReadyEvent.class)
public void logStartupConfig() {
    log.info("Application started",
        kv("profile", activeProfile),
        kv("dbUrl", maskedDbUrl),
        kv("vaadinVersion", VaadinVersion.getFullVersion()));
}
```

Log enough to diagnose configuration issues; mask or omit credential values.
