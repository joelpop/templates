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

## Suppress Lombok on Managed Collection Fields

For `@OneToMany` and `@ManyToMany` collections that require bidirectional
synchronization, suppress Lombok's getter and setter and provide manual
implementations:

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