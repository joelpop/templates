# Relationship Cascade Strategies

When mapping a one-to-many relationship, choose cascade strategy based on whether the child has independent existence — reach for `{PERSIST, MERGE}` first and escalate to `ALL` + `orphanRemoval` only when the child truly cannot exist without the parent.

## Association — child has independent existence

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

## Composition — child has no meaning without parent

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
