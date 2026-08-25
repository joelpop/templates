# Composition — Child Has No Existence Without Parent

When mapping a one-to-many relationship where the child belongs entirely to
the parent and has no meaning outside that relationship, use `CascadeType.ALL`
and `orphanRemoval = true` so the parent fully controls the child's lifecycle.

## Example

The example uses an employee → phone number relationship: a phone number
belongs exclusively to one employee; deleting the employee deletes their phone
numbers, and removing a phone from the collection deletes it from the database.

```java
// EmployeeEntity (parent) — non-owning side
@Getter(AccessLevel.NONE)
@Setter(AccessLevel.NONE)
@OneToMany(mappedBy = "employee",
           cascade = CascadeType.ALL,
           orphanRemoval = true)
private List<PhoneEntity> phones = new ArrayList<>();

// collection chaperone methods
public List<PhoneEntity> getPhones() {
    return Collections.unmodifiableList(phones);
}

public void clearPhones() {
    phones.clear();
}

public void setPhones(List<PhoneEntity> phones) {
    clearPhones();
    phones.forEach(this::addPhone);
}

public void addPhone(PhoneEntity... phones) {
    Stream.of(phones).forEach(p -> {
        this.phones.add(p);
        p.setEmployee(this);    // owning side — without this, FK is never written
    });
}

public void removePhone(PhoneEntity... phones) {
    Stream.of(phones).forEach(this.phones::remove);
}
```

```java
// PhoneEntity (child) — owning side; holds the FK column
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "employee_key", nullable = false)
private EmployeeEntity employee;
```

## Mapping

`PhoneEntity` owns the relationship — it holds the FK column. `EmployeeEntity`
is the non-owning side; `mappedBy = "employee"` tells JPA not to write the FK
from the collection.

## Cascade strategy

`CascadeType.ALL` propagates all operations to children, including `REMOVE`.
Deleting the employee deletes all their phone numbers. To persist changes,
save the employee — cascade handles the phone numbers. A `PhoneRepository`
is never needed.

## Orphan removal

`orphanRemoval = true`. Removing a phone from the employee's collection
deletes it from the database. Both `CascadeType.ALL` and `orphanRemoval`
are required — they cover different scenarios:

- `CascadeType.ALL` — employee deleted → all phones deleted
- `orphanRemoval` — phone removed from collection → that phone deleted

Omitting either one results in incorrect behavior.

## Collection chaperone methods

JPA writes the FK by reading `phone.employee` — it never reads
`employee.phones`. Adding to the collection alone leaves `phone.employee`
null; the FK column is never written. With `orphanRemoval = true`, Hibernate
deletes any phone no longer in the collection at flush — direct list mutation
bypasses `addPhone` and produces unpredictable deletes.

The unmodifiable list view and suppressed Lombok setter force all mutations
through the helpers so that `phone.employee` is always set to the owning
employee and Hibernate issues a DELETE for every phone absent from the
collection at flush.

## Related

- `persistence/spring-data-jpa/association.md` — when the child has independent existence
- `persistence/spring-data-jpa/many-to-one-lazy.md` — why `@ManyToOne` is always `FetchType.LAZY`
- `persistence/spring-data-jpa/entity-lombok.md` — `@Getter(AccessLevel.NONE)` / `@Setter(AccessLevel.NONE)` in the context of Lombok entity patterns
