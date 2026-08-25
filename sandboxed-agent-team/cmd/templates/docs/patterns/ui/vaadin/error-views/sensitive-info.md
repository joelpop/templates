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