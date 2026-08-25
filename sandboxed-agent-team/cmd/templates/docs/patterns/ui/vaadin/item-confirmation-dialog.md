# ItemConfirmationDialog<T>

When a dialog needs to act on a specific item, give it an item field and a typed confirm event so handlers retrieve the item via `event.getItem()` rather than casting or external state.

For example:

```java
public class ItemConfirmationDialog<T> extends Composite<ConfirmDialog> {

    private final ConfirmDialog dialog;
    private final T item;

    public ItemConfirmationDialog(T item) {
        this.item = item;
        dialog = getContent();
        dialog.addConfirmListener(_ -> fireEvent(new ConfirmEvent<>(this, false)));
    }

    public T getItem() {
        return item;
    }

    public void setHeader(String header)            { dialog.setHeader(header); }
    public void setText(String text)                { dialog.setText(text); }
    public void setCancelable(boolean cancelable)   { dialog.setCancelable(cancelable); }
    public void setConfirmText(String confirmText)  { dialog.setConfirmText(confirmText); }
    public void setConfirmButtonTheme(String theme) { dialog.setConfirmButtonTheme(theme); }
    public void open()                              { dialog.open(); }

    @SuppressWarnings("unchecked")
    public Registration addConfirmListener(ComponentEventListener<ConfirmEvent<T>> listener) {
        return addListener(ConfirmEvent.class, (ComponentEventListener) listener);
    }

    public static class ConfirmEvent<T> extends ComponentEvent<ItemConfirmationDialog<T>> {

        public ConfirmEvent(ItemConfirmationDialog<T> source, boolean fromClient) {
            super(source, fromClient);
        }

        public T getItem() {
            return getSource().getItem();
        }
    }
}
```

Configure the dialog and register a typed confirm listener:

```java
private void onDeactivateButtonClick(ItemButton.ItemClickEvent<ItemListItem> event) {
    var item = event.getItem();
    var dialog = new ItemConfirmationDialog<>(item);
    dialog.setHeader("Deactivate " + item.getName() + "?");
    dialog.setText("The item will no longer appear in active lists.");
    dialog.setCancelable(true);
    dialog.setConfirmText("Deactivate");
    dialog.setConfirmButtonTheme("error primary");
    dialog.addConfirmListener(this::onDeactivateDialogConfirm);
    dialog.open();
}

private void onDeactivateDialogConfirm(ItemConfirmationDialog.ConfirmEvent<ItemListItem> event) {
    deactivate(event.getItem());
}
```

**Related:** `item-button.md` — the same item-carrying pattern applied to `Button`;
`dialog-delegation.md` — the `Composite<Dialog>` delegation pattern.