# SLF4J + Logback

When logging from application code, use SLF4J via `@Slf4j` — never import Logback or `java.util.logging` classes directly.

Application code logs through **SLF4J** — declared via `@Slf4j` (see
`docs/patterns/language/java/lombok/logging.md`). The implementation is
**Logback**, pulled in transitively by `spring-boot-starter-logging`. The same
starter bridges `java.util.logging` and Log4j, so third-party libraries route
through the same pipeline.

- **Code imports** `org.slf4j.Logger` — never `ch.qos.logback.*` or
  `java.util.logging.*`.
- **Configuration** uses `logback-spring.xml` for structured setup,
  `application.properties` `logging.level.*` for level tuning
  (see `docs/patterns/cicd/structured-logging.md`).
