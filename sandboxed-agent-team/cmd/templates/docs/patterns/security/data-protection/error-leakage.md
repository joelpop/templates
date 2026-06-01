# Sensitive Information Leakage Prevention

When handling errors, suppress internal details from all responses so stack
traces, file paths, database error details, and server version information
cannot be extracted by an attacker.

## What Must Not Appear in Error Responses

- Stack traces
- Internal file paths
- Database error details or constraint names
- Server version information
- Framework version information

Stack traces appear only in server-side logs. Spring Boot's default `/error`
endpoint must be suppressed or overridden. See
`docs/patterns/ui/vaadin/error-views.md` for error view patterns.

