# Notification Toasts

When showing an outcome that cannot be displayed inline, use a Vaadin `Notification`
toast — 3 s for success, 5 s for error — positioned at `BOTTOM_END` so feedback
is non-blocking and temporally appropriate.

```java
// Success — auto-dismiss after 3 seconds
var notification = Notification.show("Saved successfully");
notification.setDuration(3000);
notification.setPosition(Notification.Position.BOTTOM_END);
notification.addThemeVariants(NotificationVariant.LUMO_SUCCESS);

// Error
var notification = Notification.show("An error occurred. Please try again.");
notification.setDuration(5000);
notification.setPosition(Notification.Position.BOTTOM_END);
notification.addThemeVariants(NotificationVariant.LUMO_ERROR);
```

Reserve toasts for outcomes that cannot be shown inline (service errors, async
completions). Do not use toasts as the sole means of reporting form validation
errors.
