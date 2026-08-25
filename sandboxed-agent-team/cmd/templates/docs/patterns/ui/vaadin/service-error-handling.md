# Service Error Handling in Views

When a view calls a service method that can throw, branch on the service's
specific exception types — such as `ValidationException` and
`EntityNotFoundException` — before a generic `Exception` catch so user-authored
error messages surface directly while internal details are never exposed.

```java
try {
    itemService.save(binder.getBean());
    close();
} catch (ValidationException e) {
    // ValidationException messages are authored by the service for user consumption —
    // show them directly. Each error in the list is a discrete user-facing message.
    e.getErrors().forEach(msg ->
            Notification.show(msg, 5000, Notification.Position.BOTTOM_END)
                    .addThemeVariants(NotificationVariant.LUMO_ERROR));
} catch (EntityNotFoundException e) {
    // Record disappeared between load and save. Navigate away — leaving the user on a
    // detail form for a nonexistent record is worse than a navigation reset.
    Notification.show("Record no longer exists.", 5000, Notification.Position.BOTTOM_END)
            .addThemeVariants(NotificationVariant.LUMO_ERROR);
    getUI().ifPresent(ui -> ui.navigate(ItemListView.class));
} catch (Exception e) {
    // Unexpected — log server-side, show nothing that reveals internals.
    log.error("Unexpected error saving item", e);
    Notification.show("An error occurred. Please try again.", 5000, Notification.Position.BOTTOM_END)
            .addThemeVariants(NotificationVariant.LUMO_ERROR);
}
```

`ValidationException.getErrors()` messages are the only exception messages safe to
surface directly — they are authored for that purpose. `EntityNotFoundException` and
any unexpected exception always show a generic message. Never pass an arbitrary
`e.getMessage()` to `Notification.show()`; the exception message is internal state.

For the service-layer error contract that produces these exception types, see
`docs/patterns/structure/services.md`.
