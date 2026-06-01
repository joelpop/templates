# System Error View (500)

When implementing the system error view, extend `BaseErrorView` (see
`shared-base.md`) and implement `HasErrorParameter<Exception>` so all
uncaught exceptions are caught here. Vaadin specificity-based routing ensures
more-specific `HasErrorParameter` implementations (`AccessDeniedView`,
`InvalidRequestView`) win for their own exception types; this view handles
everything else.

Behavior contract:

- A `UUID.randomUUID()` reference ID is generated per occurrence, displayed
  to the user as a monospace "Reference: …" paragraph (added via
  `addBelowMessage`), and included in the server-side ERROR log entry.
  Support staff correlate user-reported references to logs by this UUID.
- User-facing content is a generic "something went wrong" message — no
  exception class, message, stack trace, file path, or database detail.
- A **Retry** button is appended via `addExtraAction`; it reloads the page.
  A repeated failure renders the 500 view again — intentional.
- `setErrorParameter` returns `HttpStatus.INTERNAL_SERVER_ERROR.value()` (500).
- Logs at ERROR with actor-context fields, the UUID, exception class, and
  exception message via the `logError` helper.

```java
@PermitAll
@Slf4j
public class SystemErrorView extends BaseErrorView
        implements HasErrorParameter<Exception> {

    private final Paragraph reference;

    public SystemErrorView(CurrentUserTenantService currentUser) {
        super(AppIcon.ERROR_SYSTEM,
                "Something Went Wrong",
                "An unexpected error occurred. Use the reference below if you contact support.",
                currentUser);

        reference = new Paragraph();
        addBelowMessage(reference);
        addExtraAction(new Button("Retry", _ -> UI.getCurrent().getPage().reload()));
    }

    @Override
    public int setErrorParameter(BeforeEnterEvent event,
                                 ErrorParameter<Exception> parameter) {
        var errorRef = UUID.randomUUID().toString();
        reference.setText("Reference: " + errorRef);

        var caught = parameter.getCaughtException();
        logError(log, "Unhandled exception", event,
                ", ref:%s, exceptionClass:%s, exceptionMessage:\"%s\"".formatted(
                        errorRef,
                        caught != null ? caught.getClass().getName() : "",
                        caught != null && caught.getMessage() != null ? caught.getMessage() : ""),
                caught);

        return HttpStatus.INTERNAL_SERVER_ERROR.value();
    }
}
```