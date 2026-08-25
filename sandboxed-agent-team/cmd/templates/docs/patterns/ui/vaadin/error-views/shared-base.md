# Error View Shared Base

When building error views, extend a shared `BaseErrorView` so chrome changes
propagate to all four views automatically and the 403-as-404 disguise is
enforced at the type level.

`BaseErrorView extends Composite<VerticalLayout>` and constructs the full
shared chrome: page background, centered card, icon-heading-message stack,
and the standard action row with **Back** and **Home** buttons. Subclasses
extend `BaseErrorView`, implement `HasErrorParameter<X>` for their specific
exception type, and supply only their own copy and any extra content or actions.

```
BaseErrorView extends Composite<VerticalLayout>
    NotFoundView      implements HasErrorParameter<NotFoundException>
    AccessDeniedView  implements HasErrorParameter<AccessDeniedException>
    SystemErrorView   implements HasErrorParameter<Exception>
    InvalidRequestView implements HasErrorParameter<IllegalArgumentException>
```

Vaadin routes each exception to the most specific registered
`HasErrorParameter` implementation, so `AccessDeniedView` and
`InvalidRequestView` always win over `SystemErrorView`'s
`HasErrorParameter<Exception>` for their respective exception types.

## Action Row

Every error view renders the same action row built by the base class:

- **Back** — `UI.getCurrent().getPage().getHistory().back()`. Returns the user
  to the previous page.
- **Home** — navigates to the root view, which forwards appropriately for the
  user's authentication state.

Subclasses may append extra actions via `addExtraAction(component)` (e.g.,
`SystemErrorView` appends a **Retry** button). Extra content above the action
row (e.g., a UUID reference) is added via `addBelowMessage(component)`.

## Logging Helpers

`BaseErrorView` provides `logWarning` and `logError` helpers that format the
standard actor-context fields (path, email, role, tenant) and guard the log
call with `isWarnEnabled` / `isErrorEnabled`. Subclasses pass their own SLF4J
`Logger` so each log entry is attributed to the concrete subclass.
