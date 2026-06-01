# Session Cookie Security Flags

When configuring the session cookie for a production Vaadin application, set
all three security flags so the cookie is protected from theft and cross-site
requests regardless of the servlet container.

```properties
# application-production.properties
server.servlet.session.cookie.http-only=true
server.servlet.session.cookie.secure=true
server.servlet.session.cookie.same-site=strict
```

```properties
# application-development.properties — override Secure so the browser sends
# the cookie over plain HTTP; all other flags inherit from production
server.servlet.session.cookie.secure=false
```

**`HttpOnly`** prevents JavaScript from reading the session cookie, closing the
main XSS session-hijacking vector. Embedded Tomcat sets this by default
(`StandardContext.useHttpOnly = true`), but other servlet containers (Jetty,
Undertow, external Tomcat deployments) may not. Setting it explicitly
guarantees the flag regardless of container.

**`Secure`** prevents transmission over plain HTTP. No Spring Boot or Vaadin
default sets this — it must be configured explicitly for production. Without it
the session cookie travels over unencrypted connections even when HTTPS is
available.

**`SameSite=Strict`** prevents the browser from sending the session cookie on
any cross-site request, eliminating CSRF for session-authenticated endpoints.
Neither Spring Boot, Spring Security, nor Vaadin sets a SameSite value for the
session cookie by default — it requires explicit configuration. `Strict` is
appropriate for applications that do not need cross-site link navigation to
preserve session state; use `Lax` if direct navigation from external links must
carry the session.