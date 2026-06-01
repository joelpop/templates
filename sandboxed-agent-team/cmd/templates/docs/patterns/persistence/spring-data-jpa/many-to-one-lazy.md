# @ManyToOne Always Lazy

Always declare `@ManyToOne` with `fetch = FetchType.LAZY` — JPA's default of `EAGER` causes surprise queries everywhere.

```java
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "department_key")
private DepartmentEntity department;
```

No `@ManyToOne` annotation in the codebase should omit `fetch = FetchType.LAZY`.
