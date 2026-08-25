# Dialog Delegation

When building a Vaadin dialog class, extend `Composite<Dialog>` or
`Composite<ConfirmationDialog>` rather than `Dialog` itself so the class exposes
a focused public API rather than `Dialog`'s  50+ public methods, and the dialog stays
in the DOM when added to a parent for draggable positioning.

```java
// Avoid
public class EditItemDialog extends Dialog { /* ... */ }
```

```java
// Preferred
public class EditItemDialog extends Composite<Dialog> {
    private final Dialog dialog;

    public EditItemDialog(/* dependencies */) {
        dialog = getContent();
        // configure dialog contents
        var cancelButton = new Button("Cancel");
        cancelButton.addClickListener(this::onCancelButtonClick);
        var saveButton = new Button("Save");
        saveButton.addClickListener(this::onSaveButtonClick);
    }

    private void onCancelButtonClick(ClickEvent<Button> _) {
        dialog.close();
    }

    private void onSaveButtonClick(ClickEvent<Button> event) {
        fireEvent(new SaveEvent(this, event.isFromClient()));
        dialog.close();
    }

    public void open()  {
        dialog.open();
    }

    public Registration addSaveListener(ComponentEventListener<SaveEvent> listener) {
        return addListener(SaveEvent.class, listener);
    }

    public static class SaveEvent extends ComponentEvent<EditItemDialog> {
        public SaveEvent(EditItemDialog source, boolean fromClient) {
            super(source, fromClient);
        }
    }
}
```
