# Lombok @Slf4j for Logging

When a class needs a logger, use `@Slf4j` to generate a ready-to-use `log` field — no boilerplate and no risk of declaring the logger against the wrong class.

```java
@Service
@Slf4j
public class JpaEmployeeService implements EmployeeService {

    @Override
    @Transactional
    public EmployeeDetail create(EmployeeDetail detail) {
        log.info("Creating employee with id {}", detail.getKey());
        // ...
    }
}
```

Use SLF4J parameterized logging (`log.debug("loaded {} records", count)`), not string
concatenation. Parameters are stringified only when the level is enabled, which matters
for DEBUG/TRACE paths.

See `docs/patterns/security/data-protection/pii-logging.md` for the rule against
logging user-identifying information at INFO level and below.
