# Lambdas vs. Method References

When writing an event handler or callback, use a method reference instead of an inline
lambda so the handler body stays in a named method where it can be read and tested
independently. Use an inline lambda only when the body must capture a local variable
that should not become a field.

```java
// Avoid — multi-line logic buried in the lambda
saveButton.addClickListener(e -> {
    if (binder.validate().isOk()) {
        service.save(binder.getBean());
        close();
    }
});
```

```java
// Preferred — delegate to a named method
saveButton.addClickListener(this::onSaveButtonClick);

// Acceptable — lambda needed to capture a local variable
var confirmDialog = new Dialog();
deleteButton.addClickListener(_ -> {
    confirmDialog.close();
    onDeleteConfirmed();
});
```
