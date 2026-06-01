# Lambdas vs. Method References

Prefer method references over inline lambdas for multi-line handlers. Use inline lambdas
when the handler needs to capture a local variable:

```java
// Preferred: method reference for multi-line logic
saveButton.addClickListener(this::onSaveButtonClick);

// Acceptable: lambda captures a local variable
var confirmDialog = new Dialog();
deleteButton.addClickListener(_ -> {
    confirmDialog.close();
    onDeleteConfirmed();
});
```
