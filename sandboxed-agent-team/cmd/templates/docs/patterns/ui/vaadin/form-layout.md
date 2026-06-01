# Form Layout

When laying out form fields, use `FormLayout` so column count adapts to available width automatically.

```java
var form = new FormLayout();
form.add(nameField, codeField, descriptionField, activeCheckbox);
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),     // 1 column on small
    new FormLayout.ResponsiveStep("600px", 2)  // 2 columns on wider
);
```
