# Error View Patterns

Pattern guidance for user-facing error views and sensitive-information
protection in Vaadin 24+ applications. The Vaadin APIs named below
(`RouteNotFoundError`, `HasErrorParameter`, `ErrorHandler`) are stable across
every supported Vaadin line.

This document describes the *shape* each error view must conform to. The
project applying this pattern picks its own copy, layout primitives, icons,
class names, and home-view target — the pattern doesn't dictate them.

## Error View Types

The application displays a user-friendly error view for each error condition:

| Error Type | HTTP Equivalent | Vaadin Mechanism |
|------------|-----------------|------------------|
| Not Found | 404 | `RouteNotFoundError` subclass |
| Access Denied | 403 (displayed as 404) | `HasErrorParameter<AccessDeniedException>` |
| System Error | 500 | `HasErrorParameter<Exception>` + `ErrorHandler` |
| Invalid Request | 400 | `HasErrorParameter<IllegalArgumentException>` |

No error condition results in a raw stack trace or framework error overlay
visible to the user.

## Common Base / Shared Chrome

The four error views should share most of their visible structure: page background,
container layout, the icon-heading-message stack, the action row, theme
integration. Push that commonality into a shared base so the four concrete
views are thin descendants that supply only their own type-specific copy
and behavior.

This buys two things:

- **Consistency by construction.** Visual or behavioral chrome changes (a
  background tweak, a new action) happen in one place and propagate to all
  four views automatically.
- **Disguise enforced at the type level.** The 403-as-404 view's content
  must be indistinguishable from the 404 view's. If both pull their chrome
  from the same base, they cannot drift.

Java single inheritance is a wrinkle: the 404 view must extend Vaadin's
`RouteNotFoundError`, while the other three implement `HasErrorParameter<X>`.
The cleanest path is to have the base class itself extend `RouteNotFoundError`;
the other three then subclass the base *and* implement their own
`HasErrorParameter<X>`. Vaadin routes errors by the most specific match, so
the parent's `RouteNotFoundError` role is invoked only when a class doesn't
also implement a more specific `HasErrorParameter`. An alternative is
composition (each view embeds a shared content component via `setContent`),
which is lighter on inheritance but loses the type-level guarantee.

## Action Row

Every error view renders the same action row, with **Go Back** to the left of
**Go Home**. The 500 view additionally appends **Retry** to the right of Go Home.

- **Go Back** invokes browser history back via
  `UI.getCurrent().getPage().getHistory().back()`. Returns the user to the page
  they came from — the most recoverable option when the error was incidental
  to navigation rather than to a destination the user actually wanted.
- **Go Home** navigates to the application's home view. The home view itself
  decides where to forward (dashboard for authenticated users, login page for
  unauthenticated ones), so a single target works in both states.
- **Retry** (500 only) reloads the current page via
  `UI.getCurrent().getPage().reload()`. A repeated failure renders the 500
  view again — intentional, since pretending the failure went away would be
  worse than honesty.

## Not Found View (404)

Subclass `RouteNotFoundError`, mark `@AnonymousAllowed`, populate `setContent`
with the friendly content and the standard action row.

Behavior contract:
- HTTP status `404` is returned (Vaadin's default for `RouteNotFoundError`).
- A heading and short message tell the user the page wasn't found.
- The standard action row is present.
- Navigation to a non-existent route renders this view, not a blank page.

## Access Denied View (403 Disguised as 404)

The disguise is application-security policy, not a Vaadin-mechanics quirk.
Same principle as why a failed login says "incorrect email or password"
rather than "incorrect password" — distinguishing a forbidden resource from
a non-existent one would let a probe enumerate which routes exist and which
only require a higher role. The 403-as-404 response must be indistinguishable
from a real 404 in everything the user (or an attacker) can observe: status
code, rendered content, headers, and timing characteristics.

Implement `HasErrorParameter<AccessDeniedException>`, mark `@AnonymousAllowed`.

Behavior contract:
- The rendered content matches the 404 view's content exactly. Enforce this
  via the shared base / chrome (see "Common Base / Shared Chrome" above) so
  the two views cannot drift.
- `setErrorParameter` returns `HttpStatusCode.NOT_FOUND.value()` (404), never
  `403`.
- The event is logged at WARN with the actor-context fields (see "Logging
  Standards" below) so support and security can investigate, but those
  details are never visible in the user-facing response.

## System Error View (500)

Implement `HasErrorParameter<Exception>` (the broad-catching parameter) plus
register an `ErrorHandler` so any uncaught exception thrown by a UI handler
or service routes here.

Behavior contract:
- A `UUID.randomUUID()` reference ID is generated per occurrence, displayed
  to the user, and included in the server-side ERROR log entry. Support
  staff correlate user-reported references to logs by this UUID.
- The user-facing content is a generic "something went wrong" message — no
  exception class, message, stack trace, file path, or database detail
  surfaces.
- The action row includes Retry (per the Action Row section above).
- `setErrorParameter` returns `HttpStatusCode.INTERNAL_SERVER_ERROR.value()`
  (500).
- Throwing any `RuntimeException` from a view event handler or service
  renders this view — no exception type escapes to a raw framework overlay.

## Invalid Request View (400)

Implement `HasErrorParameter<IllegalArgumentException>`.

Behavior contract:
- Surfaces the parameter's custom message when one is set
  (`parameter.hasCustomMessage()` / `parameter.getCustomMessage()`); falls
  back to a generic phrasing otherwise. Custom messages are the throwing
  code's responsibility to keep safe — never internal field names, database
  schema, or stack traces.
- `setErrorParameter` returns `HttpStatusCode.BAD_REQUEST.value()` (400).
- Only shown when a 403 error is *not* implicated. If a request could be
  both 400 and 403 (e.g., invalid ID format on a restricted resource), the
  403-as-404 view takes precedence — see the security rationale above.
  Otherwise a probe would learn the resource exists and only the role was
  missing.

## Sensitive Information Protection

No error view may expose:
- Stack traces
- Internal file paths
- Database error details or constraint violation messages
- Server version information
- Framework internals
- Security-sensitive information

Spring Boot's default `/error` endpoint must be suppressed:

```properties
# Disable the default /error endpoint — Vaadin's error mechanism handles errors.
server.error.whitelabel.enabled=false
```

Database constraint violation messages caught by `DataIntegrityViolationException`
must be translated to user-friendly messages at the service layer before
propagating to the UI.

## Logging Standards

Every error-view log entry carries the **actor-context fields** — the
who/where/when of the request — so logs from different error types are
correlatable and analyzable by the same set of dimensions:

- Actor identifier (the user)
- Actor role
- Operational scope (a tenant id, an organization, etc. — whatever scopes
  the action; project-defined)
- Request path
- Timestamp

How each project handles missing values (omit, empty string, sentinel) is a
project decision and should be specified in the project's requirements doc.

Per-error-type extras on top of the actor-context fields:

| Error type | Level | Additional fields |
|------------|-------|-------------------|
| 404 Not Found | WARN | — |
| 403-as-404 Access Denied | WARN | — |
| 500 System Error | ERROR | UUID correlation reference, exception class, message, stack trace |
| 400 Invalid Request | WARN | Safe exception message (no stack trace) |

Log message format follows the project's logging convention — typically
SLF4J's parameterized format with `key:value` pairs separated by commas:

```java
log.warn("Access denied - path:{}, email:{}, role:{}, tenant:\"{}\"",
        path, email, role, tenant);
```

Quote values that may contain whitespace. A structured-encoder upgrade
(Logstash JSON or equivalent) is a separable follow-up if and when log
aggregation is added.