# Nested Type Placement

Place nested types (inner classes, enums) at the end of the class, after all methods:

```java
public class EditDialog {
    // Fields
    // Constructor
    // Public API methods
    // Private helper methods

    // ========== Event Classes ==========
    public static class SaveEvent extends NonComponentEvent<EditDialog> { ... }
    public static class CancelEvent extends NonComponentEvent<EditDialog> { ... }

    // Private enums last
    private enum Mode { CREATE, EDIT }
}
```

Use `private` visibility for enums and inner classes used only internally.
