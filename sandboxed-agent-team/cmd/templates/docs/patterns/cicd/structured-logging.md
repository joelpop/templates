# Structured Logging

When configuring logging for production, use JSON output via `LogstashEncoder` so
each field is a distinct key parseable by log aggregation pipelines.

```xml
<dependency>
    <groupId>net.logstash.logback</groupId>
    <artifactId>logstash-logback-encoder</artifactId>
    <version>8.0</version>
</dependency>
```

```xml
<!-- logback-spring.xml — JSON output in production, readable output locally -->
<springProfile name="prod">
    <appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
        <encoder class="net.logstash.logback.encoder.LogstashEncoder"/>
    </appender>
    <root level="WARN">
        <appender-ref ref="JSON"/>
    </root>
</springProfile>

<springProfile name="local">
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

## Related

- `docs/patterns/cicd/logging-framework.md` — SLF4J + Logback setup.
- `docs/patterns/cicd/deployment/spring-profiles.md` — profile activation and `logging.level.*` overrides.