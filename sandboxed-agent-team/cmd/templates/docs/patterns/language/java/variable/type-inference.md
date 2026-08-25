# Type Inference

When declaring a local variable, use `var` so the type declaration is not
repeated when it is already expressed on the right-hand side:

```java
// Avoid — type is redundant
FormLayout userFormLayout = new FormLayout();
List<ItemListItem> itemListItems = itemService.listAll();
```

```java
// Preferred
var userFormLayout = new FormLayout();
var itemListItems = itemService.listAll();
```
