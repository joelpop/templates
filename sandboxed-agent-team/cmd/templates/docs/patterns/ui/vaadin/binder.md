# Vaadin Binder

When building a Vaadin form, use `Binder` for all field-to-model binding and validation — manual `getValue()` / `setValue()` form handling is not permitted.

```java
var binder = new Binder<>(ItemDetail.class);

binder.forField(nameField)
      .asRequired("Name is required")
      .withValidator(n -> n.length() <= 100, "Maximum 100 characters")
      .bind(ItemDetail::getName, ItemDetail::setName);

binder.forField(codeField)
      .asRequired("Code is required")
      .bind(ItemDetail::getCode, ItemDetail::setCode);

binder.setBean(item);
```

Validation errors appear inline, adjacent to the offending field, as Binder
field-level error messages. No validation error may be shown only as a toast.
