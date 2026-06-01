# Event Handler Naming

Use the pattern `on{ComponentName}{EventType}` for event handler methods:

```java
saveButton.addClickListener(this::onSaveButtonClick);
nameField.addValueChangeListener(this::onNameFieldValueChanged);
upload.addSucceededListener(this::onPhotoUploadSucceeded);
```
