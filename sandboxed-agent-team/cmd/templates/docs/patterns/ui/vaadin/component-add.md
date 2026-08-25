# Adding Components to a Layout

When adding components to a `HasComponents`, use one `add` call per component so
each child is on its own line — easier to read, reorder, and diff than multi-argument
calls or constructor arguments.

```java
// Avoid — components passed through the constructor
var layout = new VerticalLayout(titleLabel, nameField, saveButton);

// Avoid — multiple components in one add call
layout.add(titleLabel, nameField, saveButton);
```

```java
// Preferred
var layout = new VerticalLayout();
layout.add(titleLabel);
layout.add(nameField);
layout.add(saveButton);
```