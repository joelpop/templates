# Error View Logging Standards

When logging from an error view, use `BaseErrorView`'s `logWarning` /
`logError` helpers so every entry carries the standard actor-context fields
and the `isEnabled` guard is consistent.

Every log entry includes the **actor-context fields**, such as:

- Path (request path that triggered the error)
- Email (actor identifier)
- Role
- Tenant (if applicable, operational scope; project-defined)

Per-error-type level and extras on top of the actor-context fields:

| Error type               | Level | Additional fields                                           |
|:-------------------------|:------|:------------------------------------------------------------|
| 404 Not Found            | WARN  | —                                                           |
| 403-as-404 Access Denied | WARN  | —                                                           |
| 500 System Error         | ERROR | UUID correlation reference, exception class, message, stack |
| 400 Invalid Request      | WARN  | Safe exception message (no stack trace)                     |

Log format uses SLF4J parameterized format with `key:value` pairs; quote
values that may contain whitespace:

```java
log.warn("Access denied - path:{}, email:{}, role:{}, tenant:\"{}\"",
        path, email, role, tenant);
```