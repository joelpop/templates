# Confirmation Dialog

When a consequential or irreversible action (delete, deactivate) is triggered, show a `ConfirmDialog`
that states the consequence in plain language before executing so the action cannot
be taken accidentally.

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

**Related:** `item-button.md` — carrying the item on the button; `item-confirmation-dialog.md` — carrying the item on the dialog.

Cancel leaves the record unchanged. The confirmation dialog always states what
will happen in plain language in the body text.
