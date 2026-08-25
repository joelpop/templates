# Form Layout

When laying out form fields, use `FormLayout` so column count adapts to available width automatically.

```java
var formLayout = new FormLayout();
formLayout.add(nameField);
formLayout.add(codeField);
formLayout.add(descriptionField);
formLayout.add(activeCheckbox);
formLayout.setResponsiveSteps(
    new FormLayout.ResponsiveStep(Breakpoints.XS.minWidth, 1),
    new FormLayout.ResponsiveStep(Breakpoints.SM.minWidth, 2)
);
```
