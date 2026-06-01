# Sensitive Information Protection in Error Views

When writing error views, never expose internal details in user-facing content
so attackers cannot probe the system by triggering errors.

No error view may surface:

- Stack traces
- Internal file paths
- Database error details or constraint violation messages
- Server version information
- Framework internals
- Security-sensitive information

Spring Boot's default `/error` endpoint must be suppressed:

```properties
# Vaadin's HasErrorParameter mechanism handles errors — suppress the default endpoint.
server.error.whitelabel.enabled=false
```

Database constraint violation messages caught by `DataIntegrityViolationException`
must be translated to user-friendly messages at the service layer before
propagating to the UI.