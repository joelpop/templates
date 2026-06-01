# Local Variable Declaration

Declare local variables close to their first use, not at the top of the method:

```java
// Avoid
var content = new VerticalLayout();   // declared too early
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
