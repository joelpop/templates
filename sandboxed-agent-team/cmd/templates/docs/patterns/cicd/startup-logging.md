# Startup Configuration Logging

When the application starts, log resolved configuration at INFO level so deployment
configuration issues are diagnosable without reading the raw properties — mask or
omit credential values.

```java
import static net.logstash.logback.argument.StructuredArguments.kv;

@Slf4j
@Component
public class StartupLogger {

    @Value("${spring.datasource.url}")
    private String dbUrl;

    @Autowired
    private Environment environment;

    @EventListener(ApplicationReadyEvent.class)
    public void logStartupConfig() {
        String activeProfile = Arrays.toString(environment.getActiveProfiles());
        String maskedDbUrl = dbUrl.replaceAll(":[^:@]+@", ":***@");

        log.info("Application started",
                kv("profile", activeProfile),
                kv("dbUrl", maskedDbUrl),
                kv("vaadinVersion", VaadinVersion.getFullVersion()));
    }
}
```

`kv()` is from `net.logstash.logback.encoder.LogstashEncoder` — it emits each pair
as a distinct JSON field in production and as `key=value` in local console output.

## Related

- `docs/patterns/cicd/structured-logging.md` — `LogstashEncoder` setup required for `kv()` JSON output.