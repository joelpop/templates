# Vaadin Binder

When building a Vaadin form, use `Binder` for all field-to-model binding and validation — do not use manual `getValue()` / `setValue()` of fields.

```java
var itemDetailBinder = new Binder<>(ItemDetail.class);

itemDetailBinder.forField(nameField)
                .asRequired("Name is required")
                .withValidator(n -> n.length() <= 100, "Maximum 100 characters")
                .bind(ItemDetail::getName, ItemDetail::setName);

itemDetailBinder.forField(codeField)
                .asRequired("Code is required")
                .bind(ItemDetail::getCode, ItemDetail::setCode);

itemDetailBinder.setBean(item);
```

Validation errors appear inline, adjacent to the offending field, as Binder
field-level error messages.
