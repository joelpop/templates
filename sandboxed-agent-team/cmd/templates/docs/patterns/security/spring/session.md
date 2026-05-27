# Session Configuration

When configuring Spring Boot session management for a Vaadin application, set
the inactivity timeout and cookie security flags so idle sessions expire and the
session cookie is protected from theft and cross-site requests.

## Session Timeout

```properties
server.servlet.session.timeout=30m
```

When a session expires during an active Vaadin session, Vaadin's built-in
handling displays a "Session expired" notification, offers a page reload, and
redirects the user to login.

## Session Cookie Flags

```properties
server.servlet.session.cookie.http-only=true
server.servlet.session.cookie.secure=true
server.servlet.session.cookie.same-site=strict
```

`HttpOnly` prevents JavaScript access. `Secure` prevents transmission over
plain HTTP. `SameSite=Strict` prevents CSRF via cross-site requests. Apply in
production profiles.
