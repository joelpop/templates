# Dialog Delegation

When building a Vaadin dialog class, wrap `Dialog` via delegation — never
extend it — so the class exposes a focused public API rather than `Dialog`'s
50+ public methods.

```java
// Avoid
public class EditItemDialog extends Dialog { }
```

```java
// Preferred
public class EditItemDialog {
    private final Dialog dialog;

    public EditItemDialog(/* dependencies */) {
        dialog = new Dialog();
        // configure dialog contents
    }

    public void open()  { dialog.open(); }
    public void close() { dialog.close(); }
}
```

**Draggable dialogs — attachment required:** A draggable `Dialog` retains its
screen position across close/reopen only when it is explicitly added to a
parent component before opening — otherwise Vaadin auto-attaches it to the UI
on `open()` and removes it from the DOM on `close()`, resetting the position.
A delegating wrapper handles this by providing an `attachTo(HasComponents parent)`
method that adds the inner `dialog` to the given parent without exposing it:

```java
public void attachTo(HasComponents parent) { parent.add(dialog); }
```

The view calls `attachTo(this)` once during construction, and position is
preserved across open/close cycles.

For event publishing from a delegating dialog, see `non-component-events.md`.
