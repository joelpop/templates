# Service Error Contracts

When a service method encounters a predictable failure, throw `EntityNotFoundException`
or `ValidationException` rather than surfacing raw database errors or Vaadin UI
notifications so views control how errors are presented.

Define exceptions in the `{app}-service` module so both service implementations and views
can reference them without violating layer boundaries:

```java
// Entity not found (or belongs to a different context)
public class EntityNotFoundException extends RuntimeException {
    public EntityNotFoundException(long key) {
        super("Entity not found: " + key);
    }
}

// Business rule violation — may include multiple error messages
public class ValidationException extends RuntimeException {
    private final List<String> errors;

    public ValidationException(String message) {
        super(message);
        this.errors = List.of(message);
    }

    public ValidationException(List<String> errors) {
        super(String.join("; ", errors));
        this.errors = List.copyOf(errors);
    }

    public List<String> getErrors() { return errors; }
}
```

Only these two typed exceptions are defined. Uniqueness violations, required-field
omissions, cross-entity checks, and every other predictable business-rule failure surface
as a `ValidationException` with a field-specific message — not as a dedicated subclass
per field. A dedicated `DuplicateKeyException` / `DuplicateDisplayIdException` / etc.
proliferates exception types for no benefit; callers distinguish cases by the message,
not by the class.

Service implementations catch `DataIntegrityViolationException` from the database and
translate them to `ValidationException` before propagating to the UI. Raw database
error messages must never reach the user.

Services must not import or reference Vaadin — no `Notification.show(...)`, no
`UI.getCurrent()`, no component creation. An exception thrown from a service carries a
*signal*; the view decides whether that signal warrants a notification, an inline error,
a navigation reset, or something else.
