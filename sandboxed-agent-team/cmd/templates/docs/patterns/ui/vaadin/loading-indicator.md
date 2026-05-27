# Loading Indicators

When triggering an operation that takes longer than 200 ms, disable the triggering
control and show a loading cue before the call completes so the user knows the
action is in progress.

## Button Loading State

Disable the button and update its label before the service call; restore in
`finally`:

```java
saveButton.setEnabled(false);
saveButton.setText("Saving…");
try {
    itemService.save(binder.getBean());
    Notification.show("Saved successfully").addThemeVariants(NotificationVariant.LUMO_SUCCESS);
    close();
} catch (ValidationException e) {
    e.getErrors().forEach(msg ->
            Notification.show(msg).addThemeVariants(NotificationVariant.LUMO_ERROR));
} finally {
    saveButton.setEnabled(true);
    saveButton.setText("Save");
}
```

## Grid Loading State

Grid data loads via `CallbackDataProvider` show Vaadin's built-in loading indicator
automatically:

```java
grid.setItems(query -> {
    // Vaadin shows built-in loading indicator automatically for CallbackDataProvider
    ...
});
```

For `ListDataProvider`, display a visual cue before the service call and remove it
after:

```java
progressBar.setVisible(true);
UI.getCurrent().access(() -> {
    grid.setItems(itemService.listAll());
    progressBar.setVisible(false);
});
```
