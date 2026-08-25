# Local Variable Declaration

When declaring a local variable, place it as close as possible to its first use —
not at the top of the method — so the declaration, initialization, and usage are
readable together without scrolling.

```java
// Avoid
var content = new VerticalLayout();   // declared before needed - too early
var form = new FormLayout();
// ... configure form ...
content.add(form);
return content;
```

```java
// Preferred
var form = new FormLayout();
// ... configure form ...

var content = new VerticalLayout();   // declared just before it's needed
content.add(form);
return content;
```
