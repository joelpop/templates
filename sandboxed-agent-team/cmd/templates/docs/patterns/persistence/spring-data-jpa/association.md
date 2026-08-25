# Association — Child Has Independent Existence

When mapping a one-to-many relationship where the child has a lifecycle
independent of any particular parent, use `{CascadeType.PERSIST, CascadeType.MERGE}`
so saves propagate but deletes never cascade to children.

## Example

The example uses a department → employee relationship: an employee can be
added, reassigned, or removed without affecting the department, and deleting
a department should not cascade-delete its employees.

```java
// DepartmentEntity (parent) — non-owning side
@Getter(AccessLevel.NONE)
@OneToMany(mappedBy = "department",
           cascade = {CascadeType.PERSIST, CascadeType.MERGE})
private List<EmployeeEntity> employees = new ArrayList<>();

// read-only collection chaperone method
public List<EmployeeEntity> getEmployees() {
    return Collections.unmodifiableList(employees);
}
```

```java
// EmployeeEntity (child) — owning side; holds the FK column
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "department_key")
private DepartmentEntity department;
```

## Mapping

`EmployeeEntity` owns the relationship — it holds the FK column and is the
side JPA reads when writing the foreign key. `DepartmentEntity` is the
non-owning side; `mappedBy = "department"` tells JPA not to write the FK
from the collection.

## Cascade strategy

`{PERSIST, MERGE}` propagates saves: persisting or merging a department
that has unsaved employees in its collection saves them too. `REMOVE` is
not included — deleting a department does not cascade-delete its employees.

`CascadeType.ALL` is wrong here: it includes `REMOVE`, so
`departmentRepository.delete(dept)` would cascade-delete every employee in
that department. With `{PERSIST, MERGE}`, attempting to delete a department
that still has employees fails at the DB level (FK constraint) — this is the
correct behavior; it forces the caller to reassign the employees first.

## Orphan removal

`orphanRemoval = false` (the default). Removing an employee from the
department's collection does not delete the employee — the employee continues
to exist and can be reassigned to another department.

## Collection chaperone methods

JPA writes the FK by reading `employee.department` — it never reads
`department.employees`. To assign or reassign an employee, set
`employee.setDepartment(dept)` on the owning side. Service methods should provide
the ergonomic API (e.g., `assignToDepartment(employee, department)`).

The unmodifiable list view prevents accidental in-memory mutation of the
collection, which is a read-only view populated by JPA on load.

## Related

- `persistence/spring-data-jpa/composition.md` — when the child cannot exist without the parent
- `persistence/spring-data-jpa/many-to-one-lazy.md` — why `@ManyToOne` is always `FetchType.LAZY`
- `persistence/spring-data-jpa/entity-lombok.md` — Lombok patterns for entities
