# Form Layout and Binding

When building a Vaadin form, use `Binder` for all field-to-model binding and
validation and `FormLayout` for the field container so form fields are bound
consistently, validation errors appear inline, and the layout adapts to screen
width automatically.

## Vaadin Binder

All forms use Vaadin `Binder` for field-to-model binding and validation. Manual
`getValue()` / `setValue()` form handling is not permitted.

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

## Form Layout

Use `FormLayout` for form fields. It adapts column count to available width
automatically:

```java
var form = new FormLayout();
form.add(nameField, codeField, descriptionField, activeCheckbox);
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),     // 1 column on small
    new FormLayout.ResponsiveStep("600px", 2)  // 2 columns on wider
);
```
