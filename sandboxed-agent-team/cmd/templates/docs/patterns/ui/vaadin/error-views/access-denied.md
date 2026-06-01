# Access Denied View (403 Disguised as 404)

When implementing the access-denied error view, extend `BaseErrorView` (see
`shared-base.md`), implement `HasErrorParameter<AccessDeniedException>`, and
return status 404 — never 403 — so the response is indistinguishable from a
real 404.

The disguise is security policy: distinguishing a forbidden resource from a
non-existent one lets a probe enumerate which routes require higher roles.
The 403-as-404 response must be indistinguishable in everything observable:
status code, content, headers, and timing.

Behavior contract:

- Rendered content is identical to `NotFoundView` — same icon, heading, and
  message. Extending the same `BaseErrorView` and using the same constructor
  arguments enforces this; the two views cannot drift.
- `setErrorParameter` returns `HttpStatus.NOT_FOUND.value()` (404), never 403.
- Marked `@PermitAll` (not `@AnonymousAllowed`) — the view is for authenticated
  users who lack the required role.
- Logs at WARN with full actor-context fields; the log entry is the only place
  where the access-denied nature is visible.

```java
@PermitAll
@Slf4j
public class AccessDeniedView extends BaseErrorView
        implements HasErrorParameter<AccessDeniedException> {

    public AccessDeniedView(CurrentUserTenantService currentUser) {
        super(AppIcon.ERROR_NOT_FOUND,
                "Page Not Found",
                "The page you're looking for doesn't exist or has been moved.",
                currentUser);
    }

    @Override
    public int setErrorParameter(BeforeEnterEvent event,
                                 ErrorParameter<AccessDeniedException> parameter) {
        logWarning(log, "Access denied", event);
        return HttpStatus.NOT_FOUND.value();
    }
}
```