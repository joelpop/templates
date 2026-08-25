# Invalid Request View (400)

When implementing the invalid-request error view, extend `BaseErrorView` (see
`shared-base.md`) and implement `HasErrorParameter<IllegalArgumentException>`
so invalid-argument exceptions thrown during navigation render a 400 view
rather than the generic 500.

Behavior contract:

- When `parameter.hasCustomMessage()` is true, the custom message is surfaced
  in a secondary details paragraph (added via `addBelowMessage`). Throwing
  code is responsible for keeping custom messages safe — never internal field
  names, database schema, or stack traces.
- When no custom message is set, a generic fallback ("The request contained
  invalid parameters.") is shown and the details paragraph stays hidden.
- `setErrorParameter` returns `HttpStatus.BAD_REQUEST.value()` (400).
- If a request would also raise `AccessDeniedException`, Vaadin's
  specificity-based routing sends it to `AccessDeniedView` first — the 400
  view is reached only when access control is not implicated.
- Logs at WARN with actor-context fields and the safe message string.

```java
@PermitAll
@Slf4j
public class InvalidRequestView extends BaseErrorView
        implements HasErrorParameter<IllegalArgumentException> {

    private static final String FALLBACK_MESSAGE =
            "The request contained invalid parameters.";

    private final Paragraph details;

    public InvalidRequestView(CurrentUserTenantService currentUser) {
        super(AppIcon.ERROR_INVALID_REQUEST, "Invalid Request", FALLBACK_MESSAGE, currentUser);

        details = new Paragraph();
        details.setVisible(false);
        addBelowMessage(details);
    }

    @Override
    public int setErrorParameter(BeforeEnterEvent event,
                                 ErrorParameter<IllegalArgumentException> parameter) {
        String safeMessage;
        if (parameter.hasCustomMessage()) {
            safeMessage = parameter.getCustomMessage();
            details.setText(safeMessage);
            details.setVisible(true);
        } else {
            safeMessage = FALLBACK_MESSAGE;
        }

        logWarning(log, "Invalid request", event, ", message:\"%s\"".formatted(safeMessage));
        return HttpStatus.BAD_REQUEST.value();
    }
}
```