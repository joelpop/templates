# Session Timeout

When configuring Spring Boot session management for a Vaadin application, set
the inactivity timeout per environment so idle sessions expire in production,
development sessions are uninterrupted, and timeout behavior can be exercised
in tests.

```properties
# application-production.properties — expire idle sessions
server.servlet.session.timeout=30m

# application-development.properties — no expiry during active development
server.servlet.session.timeout=0
```

To test session-expiry behavior in isolation, override the timeout on the
specific test class rather than in a shared properties file — a global short
timeout would expire sessions across all tests:

```java
@SpringBootTest
@TestPropertySource(properties = "server.servlet.session.timeout=10s")
class SessionExpiryTest { /* ... */ }
```

When a session expires during an active Vaadin session, Vaadin's built-in
handling displays a "Session expired" notification, offers a page reload, and
redirects the user to login.
