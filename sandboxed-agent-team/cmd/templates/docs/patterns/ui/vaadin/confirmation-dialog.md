# Confirmation Dialog

When a destructive action (delete, deactivate) is triggered, show a `ConfirmDialog`
that states the consequence in plain language before executing so the action cannot
be taken accidentally.

```java
private void confirmDeactivate(ItemListItem item) {
    var dialog = new ConfirmDialog();
    dialog.setHeader("Deactivate " + item.getName() + "?");
    dialog.setText("The item will no longer appear in active lists.");
    dialog.setCancelable(true);
    dialog.setConfirmText("Deactivate");
    dialog.setConfirmButtonTheme("error primary");
    dialog.addConfirmListener(_ -> deactivate(item));
    dialog.open();
}
```

Cancel leaves the record unchanged. The confirmation dialog always states what
will happen in plain language in the body text.
