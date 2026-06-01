# Type Inference

Use `var` for local variable type inference whenever the type is obvious from the right-hand side:

```java
// Avoid — type is redundant
FormLayout form = new FormLayout();
List<ItemListItem> items = itemService.listAll();
```

```java
// Preferred
var form = new FormLayout();
var items = itemService.listAll();
```
