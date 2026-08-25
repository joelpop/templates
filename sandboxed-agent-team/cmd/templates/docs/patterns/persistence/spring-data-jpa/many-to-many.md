# Many-to-Many — Neither Side Owns the Other

When two entities reference each other and neither controls the other's
lifecycle, map the relationship as `@ManyToMany` with a join table so both
entities exist independently and the join table records only the pairing.

## Example

The example uses an employee ↔ project relationship: an employee can be on
many projects; a project has many employees; deleting either should not
delete the other.

```java
// ProjectEntity (owning side — manages the join table)
@Getter(AccessLevel.NONE)
@ManyToMany
@JoinTable(
    name = "employee_project",
    joinColumns        = @JoinColumn(name = "project_key"),
    inverseJoinColumns = @JoinColumn(name = "employee_key"))
private List<EmployeeEntity> employees = new ArrayList<>();

// collection chaperone methods
public List<EmployeeEntity> getEmployees() {
    return Collections.unmodifiableList(employees);
}

public void clearEmployees() {
    employees.clear();
}

public void setEmployees(List<EmployeeEntity> employees) {
    clearEmployees();
    employees.forEach(this::addEmployee);
}

public void addEmployee(EmployeeEntity employee) {
    employees.add(employee);
}

public void removeEmployee(EmployeeEntity employee) {
    employees.remove(employee);
}
```

```java
// EmployeeEntity (non-owning side — does not write the join table)
@Getter(AccessLevel.NONE)
@ManyToMany(mappedBy = "employees")
private List<ProjectEntity> projects = new ArrayList<>();

// read-only collection chaperone method
public List<ProjectEntity> getProjects() {
    return Collections.unmodifiableList(projects);
}
```

`@ManyToMany` is lazy by default — no override needed.

## Mapping

One side must be designated the technical owning side — it carries `@JoinTable`
and is the side JPA reads when writing join table rows. The other side uses
`mappedBy`. Neither entity's lifecycle depends on the other.

Choose the side through which the relationship is managed in the domain: the
entity you naturally say "add X to Y" from. Here, employees are assigned to
projects — not the other way around — so `ProjectEntity` owns the join table.

## Cascade strategy

Use `{PERSIST, MERGE}` on the owning side, or no cascade at all. Never use
`CascadeType.ALL`: it includes `REMOVE`, so removing an employee from a
project's collection would delete the employee entity — almost certainly not
the intent.

Deleting a project removes its join table rows automatically; the employees
are unaffected. Deleting an employee removes their join table rows; the
projects are unaffected.

## Collection chaperone methods

JPA writes join table rows by reading the owning side's collection
(`project.employees`). Adding to the non-owning side (`employee.projects`)
writes nothing to the join table.

The owning side's helpers (`clearEmployees`, `setEmployees`, `addEmployee`,
`removeEmployee`) manipulate only the owning collection — syncing the
non-owning side's in-memory collection would trigger a hidden lazy load of
all the employee's projects. Callers that need `employee.projects` to reflect
changes within the same transaction must reload after the flush.

The non-owning side exposes only a getter — mutations always go through the
owning side's helpers.

## When you need extra columns on the join

If the relationship needs attributes (e.g., a `joinedDate` on project
membership), model the join table as its own entity with two `@ManyToOne`
relationships — effectively an association from each entity to the join entity:

```java
@Entity
@Table(name = "employee_project")
public class EmployeeProjectEntity extends BaseEntity<Long> {

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "employee_key", nullable = false)
    private EmployeeEntity employee;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "project_key", nullable = false)
    private ProjectEntity project;

    private LocalDate joinedDate;
}
```

Replace the `@ManyToMany` on both sides with `@OneToMany` → `EmployeeProjectEntity`,
following the association pattern.

## Related

- `persistence/spring-data-jpa/association.md` — one-to-many where the child has independent existence
- `persistence/spring-data-jpa/composition.md` — one-to-many where the child belongs entirely to the parent
- `persistence/spring-data-jpa/many-to-one-lazy.md` — why `@ManyToOne` is always `FetchType.LAZY`
