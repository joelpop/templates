# Unused Lambda Parameters

When a lambda or catch clause has a parameter whose value is not used, name it
explicitly so readers know the omission is deliberate rather than an oversight.
The correct form depends on the Java version.

```java
// Avoid — undeclared intent; readers can't tell if the parameter was forgotten
saveButton.addClickListener(e -> save());
```

```java
// Preferred — Java 21+: unnamed variable makes intent unambiguous
saveButton.addClickListener(_ -> save());

try {
    // ...
} catch (ValidationException _) {
    Notification.show("Please fix the validation errors");
}

// Preferred — Java 17–20: name the parameter `unused`
saveButton.addClickListener(unused -> save());

try {
    // ...
} catch (ValidationException unused) {
    Notification.show("Please fix the validation errors");
}
```
