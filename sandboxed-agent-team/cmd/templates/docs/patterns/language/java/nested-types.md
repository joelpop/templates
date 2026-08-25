# Nested Type Placement

When a class contains nested types (static inner classes, enums), place them after all
fields, constructors, and methods so the primary class body is readable before its
nested types. Give nested types the most restrictive visibility that satisfies their
callers — `private` when used only within the enclosing class.

```java
public class EditDialog {
    // Fields
    // Constructor
    // Public API methods
    // Private helper methods

    public static class SaveEvent extends ComponentEvent<EditDialog> {
        public SaveEvent(EditDialog source, boolean fromClient) { super(source, fromClient); }
    }
    public static class CancelEvent extends ComponentEvent<EditDialog> {
        public CancelEvent(EditDialog source, boolean fromClient) { super(source, fromClient); }
    }

    private enum Mode { CREATE, EDIT }
}
```
