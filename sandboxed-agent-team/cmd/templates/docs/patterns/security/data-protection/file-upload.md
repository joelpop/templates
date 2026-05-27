# File Upload Validation

When accepting file uploads, validate type and size before storage so malicious
or oversized files are rejected at the boundary rather than written to the
filesystem or object store.

```java
if (file.size() > 2 * 1024 * 1024) {
    throw new ValidationException("File must be 2 MB or smaller.");
}
if (!Set.of("image/jpeg", "image/png").contains(file.contentType())) {
    throw new ValidationException("Only JPEG and PNG files are accepted.");
}
```

Return a user-facing error on invalid uploads — not a stack trace.
