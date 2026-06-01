# @Enumerated Always STRING

Always store enum constants by name, never by ordinal — adding or reordering constants silently maps existing data to the wrong value when ordinal is used.

```java
@Enumerated(EnumType.STRING)
private EmploymentStatusCode status;
```
