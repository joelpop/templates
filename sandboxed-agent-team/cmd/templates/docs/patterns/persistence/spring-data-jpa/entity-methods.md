# Entity Methods — Chaperones Only

When direct access to a JPA entity field or collection would allow callers
to corrupt the mapped model, add a method that makes the incorrect path
impossible — and no others — so the entity remains a faithful model
representation and domain logic stays in the service layer.

## What entity methods belong

JPA entities are data-model classes. A chaperone method belongs in an entity
only when it guards a JPA concern that callers could otherwise violate:

- **Unmodifiable collection getters** — prevent callers from mutating the
  backing list directly and bypassing JPA mechanics
- **Back-reference setters** — required in composition so the FK column is
  written (`addPhone` calls `p.setEmployee(this)`)
- **Collection mutators that trigger cascade or orphanRemoval** — `addPhone`,
  `removePhone`, `clearPhones`, `setPhones` on a composition parent
- **Owning-side collection management** — `addEmployee`, `removeEmployee`,
  `setEmployees`, `clearEmployees` on the owning side of a many-to-many

## What does not belong

Domain actions do not belong as entity chaperone methods. Beyond inverting
responsibility — a non-owning entity reaching into an owning entity to set
its state — they introduce hidden reads: syncing the other side's in-memory
collection initializes that collection if it has not been loaded, triggering
a lazy load of every related entity just to append one entry. Ergonomic APIs
("assign employee to department") belong in service methods where the cost
is explicit and the transaction boundary is controlled.

```java
// Avoid — domain ergonomic on the non-owning side
// DepartmentEntity
public void addEmployee(EmployeeEntity employee) {
    employees.add(employee);
    employee.setDepartment(this);
}
```

```java
// Preferred — JPA operation on the owning side; domain API in the service
// EmployeeEntity
public void setDepartment(DepartmentEntity department) {
    this.department = department;
}

// EmployeeService
public void assignToDepartment(EmployeeEntity employee, DepartmentEntity department) {
    employee.setDepartment(department);
    employeeRepository.save(employee);
}
```

The tell: if a method reads as a domain action ("assign", "transfer",
"promote") or reaches from a non-owning entity into an owning entity to set
its state, it belongs in a service method.

## Related

- `persistence/spring-data-jpa/association.md` — non-owning side has only a read-only getter; assignment is done via the owning field
- `persistence/spring-data-jpa/composition.md` — owning-side mechanics that do belong in the entity
- `persistence/spring-data-jpa/many-to-many.md` — owning side manages the join table; non-owning side exposes only a getter
- `persistence/spring-data-jpa/entity-lombok.md` — `equals()` and `hashCode()` for JPA identity semantics