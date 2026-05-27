# Entity Design Rules

When writing Spring Data JPA entity classes, declare every `@ManyToOne` as
`LAZY`, store enums as `STRING`, choose cascade strategy based on child
lifecycle, and map custom or multi-field types via `@AttributeConverter` or
`@Embeddable` so Hibernate does not silently issue extra queries or corrupt data.

## @ManyToOne Always Lazy

JPA's default for `@ManyToOne` is `FetchType.EAGER` — a spec decision that
causes surprise queries everywhere. Always override it:

```java
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "department_key")
private DepartmentEntity department;
```

No `@ManyToOne` annotation in the codebase should omit `fetch = FetchType.LAZY`.

## @Enumerated Always STRING

Store enum constants by name, never by ordinal:

```java
@Enumerated(EnumType.STRING)
private EmploymentStatusCode status;
```

`ORDINAL` stores the position of the constant in the declaration. Adding or
reordering constants silently maps existing data to the wrong value.

## Relationship Cascade Strategies

Choose cascade strategy based on the nature of the relationship.

### Association — child has independent existence

Child can exist without the parent; parent deletion should not cascade-delete
children.

```java
// DepartmentEntity — one-to-many, association
@OneToMany(mappedBy = "department",
           cascade = {CascadeType.PERSIST, CascadeType.MERGE})
private List<EmployeeEntity> employees = new ArrayList<>();

// EmployeeEntity — many-to-one, always LAZY
@ManyToOne(fetch = FetchType.LAZY,
           cascade = {CascadeType.PERSIST, CascadeType.MERGE})
@JoinColumn(name = "department_key")
private DepartmentEntity department;
```

`orphanRemoval = false` (the default). Removing from the collection does not
delete the child. Attempting to delete a parent that still has children fails at
the DB level (FK constraint) — this is the correct behavior.

### Composition — child has no meaning without parent

Child lifecycle belongs entirely to the parent. Deleting the parent deletes the
children. Removing from the collection deletes from the database.

```java
// EmployeeEntity — one-to-many, composition
@Getter(AccessLevel.NONE)
@Setter(AccessLevel.NONE)
@OneToMany(mappedBy = "employee",
           cascade = CascadeType.ALL,
           orphanRemoval = true)
private List<PhoneEntity> phones = new ArrayList<>();

// Manual collection management — see docs/patterns/persistence/spring-data-jpa/entity-lombok.md
public List<PhoneEntity> getPhones() { return Collections.unmodifiableList(phones); }
public void addPhone(PhoneEntity... phones) { ... }
public void removePhone(PhoneEntity... phones) { ... }
```

Reach for `{PERSIST, MERGE}` first. Escalate to `ALL` + `orphanRemoval = true`
only when the child truly cannot exist without the parent.

## @AttributeConverter for Custom Types

For any Java type that maps to a single DB column and is not a simple enum:

```java
@Converter(autoApply = true)
public class MoneyConverter implements AttributeConverter<Money, BigDecimal> {
    @Override
    public BigDecimal convertToDatabaseColumn(Money m) { return m == null ? null : m.amount(); }
    @Override
    public Money convertToEntityAttribute(BigDecimal v) { return v == null ? null : new Money(v); }
}
```

`autoApply = true` applies the converter to every field of that Java type across
all entities.

## @Embeddable for Multi-Field Value Objects

When a value object has multiple fields that belong to the same table row (no
separate table, no join, no independent identity):

```java
@Embeddable
@NoArgsConstructor @Getter @Setter
public class Address {
    private String street;
    private String city;
    private String state;
    private String postalCode;
}

@Entity
public class EmployeeEntity extends BaseEntity<Long> {
    @Embedded
    private Address homeAddress;
}
```

When the same `@Embeddable` is embedded twice, use `@AttributeOverrides` to
disambiguate column names.
