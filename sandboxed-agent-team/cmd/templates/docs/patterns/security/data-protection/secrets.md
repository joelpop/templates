# Secrets Out of Source Control

When configuring application secrets, inject them via environment variables
rather than committing plaintext values so credentials cannot be extracted from
the repository history.

```properties
spring.datasource.password=${DB_PASSWORD}
```

`application.properties` must contain no plaintext secrets. A scan of the
repository history should find no committed secrets.
