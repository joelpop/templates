# Not Found View (404)

When implementing the 404 error view, extend `BaseErrorView` (see
`shared-base.md`), implement `HasErrorParameter<NotFoundException>`, and mark
`@AnonymousAllowed` so unauthenticated users see a friendly page instead of
a redirect to login.

Behavior contract:

- `setErrorParameter` returns `HttpStatus.NOT_FOUND.value()` (404).
- Heading and message tell the user the page wasn't found.
- The standard action row is present.
- Navigation to a non-existent route renders this view, not a blank page.
- Logs at WARN with actor-context fields via the `logWarning` helper.

```java
@AnonymousAllowed
@Slf4j
public class NotFoundView extends BaseErrorView
        implements HasErrorParameter<NotFoundException> {

    public NotFoundView(CurrentUserTenantService currentUser) {
        super(AppIcon.ERROR_NOT_FOUND,
                "Page Not Found",
                "The page you're looking for doesn't exist or has been moved.",
                currentUser);
    }

    @Override
    public int setErrorParameter(BeforeEnterEvent event,
                                 ErrorParameter<NotFoundException> parameter) {
        logWarning(log, "Not found", event);
        return HttpStatus.NOT_FOUND.value();
    }
}
```