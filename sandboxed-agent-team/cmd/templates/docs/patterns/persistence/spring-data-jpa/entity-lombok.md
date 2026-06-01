# Lombok on JPA Entities

When a JPA entity has accessor or constructor boilerplate, apply only
`@NoArgsConstructor`, `@Getter`, and `@Setter`, and suppress Lombok on managed
collection fields so Hibernate proxies, lazy loading, and bidirectional
relationships do not cause silent runtime failures or performance issues.

## Safe Annotations

| Annotation           | Notes                                                    |
|----------------------|----------------------------------------------------------|
| `@Getter`            | Generates accessors without touching `equals`/`hashCode` |
| `@Setter`            | Safe on entities; use selectively on immutable fields    |
| `@NoArgsConstructor` | JPA requires a no-arg constructor                        |

## Unsafe Annotations — Never Use on Entities

| Annotation           | Why it breaks                                                                                                                                                                                                                                              |
|----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `@Data`              | Bundles `@EqualsAndHashCode` and `@ToString`, bringing every problem below at once. Also implies `@Setter` on every field, which is rarely correct on an audited entity.                                                                                   |
| `@EqualsAndHashCode` | Generates `equals`/`hashCode` from field values. Hibernate may return a proxy rather than the actual entity; uninitialized proxy fields compare as `null`, making two references to the same row appear unequal. Bidirectional relationships also recurse. |
| `@ToString`          | Traverses all fields including lazy-loaded relationships — triggers `LazyInitializationException` outside a transaction and `StackOverflowError` on bidirectional relationships (A → B → A → …).                                                           |

If you need a `toString` on an entity, include only scalar fields explicitly:

```java
@ToString(onlyExplicitlyIncluded = true)
public class EmployeeEntity extends BaseEntity<Long> {
    @ToString.Include
    private String lastName;
}
```

## Recommended Annotation Set

For most entities:

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

No `@Data`, no `@ToString`, no `equals`/`hashCode` override — all inherited from
`BaseEntity`.