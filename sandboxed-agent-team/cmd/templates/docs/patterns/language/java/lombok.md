# Lombok Guidelines

Use Lombok only on JPA entities and enums that carry properties; most pitfalls
here cause runtime failures, not compile errors — easy to miss in development.

## for JPA Entities

### Safe Annotations

| Annotation | Notes |
|------------|-------|
| `@Getter` | Generates accessors without touching `equals`/`hashCode` |
| `@Setter` | Safe on entities; use selectively on immutable fields |
| `@NoArgsConstructor` | JPA requires a no-arg constructor |
| `@RequiredArgsConstructor` | Safe when used carefully |
| `@Builder` | Requires explicit `@NoArgsConstructor` and `@AllArgsConstructor` alongside it — see "@Builder" below |

### Unsafe Annotations — Never Use

| Annotation | Why it breaks |
|------------|---------------|
| `@Data` | Bundles `@EqualsAndHashCode` and `@ToString`, so it brings every problem below at once. Also implies `@Setter` on every field, which is rarely what you want on an audited entity. |
| `@EqualsAndHashCode` | Generates `equals`/`hashCode` from field values. Hibernate may return a proxy rather than the actual entity; uninitialized proxy fields compare as `null`, making two references to the same row appear unequal. On bidirectional relationships, field-based `equals`/`hashCode` also recurses into related entities. |
| `@ToString` | Traverses all fields including lazy-loaded relationships — triggers `LazyInitializationException` outside a transaction and `StackOverflowError` on bidirectional relationships (A → B → A → …). |

### Do Not Use @ToString on Entities with Relationships

`@ToString` traverses all fields, including lazy-loaded relationships. This triggers:
- `LazyInitializationException` outside a transaction
- `StackOverflowError` on bidirectional relationships (A → B → A → ...)

If you need a `toString` on an entity, use:
```java
@ToString(onlyExplicitlyIncluded = true)
```
and mark only simple scalar fields with `@ToString.Include`.

### Suppress Lombok on Managed Collection Fields

For `@OneToMany` and `@ManyToMany` collections that require bidirectional synchronization,
suppress Lombok's getter and setter and provide manual implementations:

```java
@Getter(AccessLevel.NONE)
@Setter(AccessLevel.NONE)
@OneToMany(mappedBy = "employee", cascade = CascadeType.ALL, orphanRemoval = true)
private List<PhoneEntity> phones = new ArrayList<>();

// Manual getter — returns unmodifiable view so callers cannot bypass sync helpers
public List<PhoneEntity> getPhones() {
    return Collections.unmodifiableList(phones);
}

// Manual setter — routes through addPhone to maintain both sides of the relationship
public void setPhones(List<PhoneEntity> phones) {
    this.phones.clear();
    phones.forEach(this::addPhone);
}

// Varargs add helper — maintains the back-reference
public void addPhone(PhoneEntity... phones) {
    Stream.of(phones).forEach(p -> {
        this.phones.add(p);
        p.setEmployee(this);
    });
}

// Varargs remove helper
public void removePhone(PhoneEntity... phones) {
    Stream.of(phones).forEach(this.phones::remove);
}
```

This pattern ensures:
- Callers cannot mutate the collection and bypass back-reference synchronization
- Orphan removal works correctly (items removed from the collection are deleted)
- The bidirectional relationship is always consistent

### @Builder

When using `@Builder` on a JPA entity, you must also declare `@NoArgsConstructor` and
`@AllArgsConstructor` explicitly, because `@Builder` alone replaces the no-arg constructor
that JPA requires:

```java
@Entity
@NoArgsConstructor          // required by JPA
@AllArgsConstructor         // required by @Builder
@Builder
@Getter
@Setter
public class SomeEntity extends BaseEntity<Long> {
    // fields
    // manual equals/hashCode inherited from BaseEntity — do not override
}
```

### Recommended Annotation Set

For most principal entities:

```java
@Entity
@Table(name = "employees")
@NoArgsConstructor
@Getter
@Setter
public class EmployeeEntity extends BaseEntity<Long> {
    // fields with @Column, @ManyToOne, @OneToMany, etc.
}
```

No `@Data`, no `@ToString`, no `equals`/`hashCode` override — all inherited from `BaseEntity`.

## for Enumerations

Enums that carry additional per-constant properties (labels, sort order, symbols,
categories, etc.) are a good fit for `@Getter` + `@RequiredArgsConstructor` — Lombok
generates the constructor Java requires for enum constants and the accessors for the
fields, removing boilerplate without touching any identity semantics.

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

## for Logging

Use Lombok's `@Slf4j` annotation on any class that needs a logger. It generates a
ready-to-use `log` field with no boilerplate and no copy-paste hazard of declaring the
logger against the wrong class:

```java
@Service
@Slf4j
public class JpaEmployeeService implements EmployeeService {

    @Override
    @Transactional
    public EmployeeDetail create(EmployeeDetail detail) {
        log.info("Creating employee with id {}", detail.getKey());
        // ...
    }
}
```

Use SLF4J parameterized logging (`log.debug("loaded {} records", count)`), not string
concatenation. Parameters are stringified only when the level is enabled, which matters
for DEBUG/TRACE paths.

See `docs/patterns/architecture/security.md` → "PII Not in Logs" for the rule against
logging user-identifying information at INFO level and below.
