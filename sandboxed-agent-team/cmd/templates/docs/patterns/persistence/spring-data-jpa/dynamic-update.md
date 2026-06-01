# @DynamicUpdate for Wide Tables

When an entity table is wide, has large columns (TEXT, BLOB), or has high write volume, annotate with `@DynamicUpdate` so UPDATE statements are limited to changed columns rather than overwriting all columns on every flush.

```java
@Entity
@Table(name = "employees")
@DynamicUpdate
public class EmployeeEntity extends BaseEntity<Long> { ... }
```

`@DynamicUpdate` adds a small per-flush overhead to compute the dirty field diff — use it only when profiling confirms the full-column writes are a meaningful cost.
