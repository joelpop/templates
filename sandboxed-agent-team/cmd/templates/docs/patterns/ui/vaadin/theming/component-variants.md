# Component Theme Variants

When styling a Vaadin component, use theme variants before resorting to custom CSS — check available variants in the Vaadin documentation first.

```java
// Button variants
button.addThemeVariants(ButtonVariant.LUMO_PRIMARY);
button.addThemeVariants(ButtonVariant.LUMO_ERROR);
button.addThemeVariants(ButtonVariant.LUMO_TERTIARY);

// Grid variants
grid.addThemeVariants(GridVariant.LUMO_ROW_STRIPES);
grid.addThemeVariants(GridVariant.LUMO_NO_BORDER);

// TextField variants
field.addThemeVariants(TextFieldVariant.LUMO_SMALL);

// Badge variants (for status indicators)
badge.getElement().getThemeList().add("badge success");
badge.getElement().getThemeList().add("badge error");
```
