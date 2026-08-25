# Event Handler Naming

When naming a private method registered as an event handler, use
`on{ComponentName}{EventType}` so the method name identifies both the source
component and the event without reading the registration line:

```java
saveButton.addClickListener(this::onSaveButtonClick);
nameField.addValueChangeListener(this::onNameFieldValueChanged);
photoUpload.addSucceededListener(this::onPhotoUploadSucceeded);
```
