# Lombok on Enumerations

When an enum carries per-constant properties (labels, sort order, symbols), use `@Getter` + `@RequiredArgsConstructor` to remove boilerplate without touching any identity semantics.

```java
@Getter
@RequiredArgsConstructor
public enum PriorityCode {
    LOW("Low", 1),
    MEDIUM("Medium", 2),
    HIGH("High", 3),
    URGENT("Urgent", 4);

    private final String label;
    private final int sortOrder;
}
```

Enum constants are singletons, so `equals`/`hashCode`/`toString` identity is already
correct by default — **do not** apply `@Data`, `@EqualsAndHashCode`, or `@ToString` to an
enum. Override `toString()` by hand only if you specifically need a display form
different from the constant name (most code should call the dedicated getter — e.g.,
`PriorityCode.HIGH.getLabel()` — rather than rely on `toString()`).

Fields on an enum should be `private final` so each constant is immutable; `@Setter`
has no place here.
