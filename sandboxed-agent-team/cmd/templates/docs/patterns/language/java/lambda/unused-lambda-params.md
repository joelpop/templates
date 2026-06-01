# Unused Lambda Parameters

> **Java 21+:** use `_` (the unnamed variable) for unused lambda parameters and unused
> catch clause variables. This is the preferred form on any Java-21-or-newer project
> (including Vaadin 24.1+ on Java 21 and every Vaadin 25+ project).
>
> **Java 17–20:** `_` is not available. Name the parameter `unused` to signal intent at
> the call site — the name itself tells any reader the value is deliberately ignored.

```java
// Java 21+: preferred
saveButton.addClickListener(_ -> save());
cancelButton.addClickListener(_ -> close());

try {
    // ...
} catch (ValidationException _) {
    Notification.show("Please fix the validation errors");
}

// Java 17–20: name the parameter `unused`
saveButton.addClickListener(unused -> save());

try {
    // ...
} catch (ValidationException unused) {
    Notification.show("Please fix the validation errors");
}
```
