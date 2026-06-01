# Interface Projections

When querying Spring Data JPA repositories, return interface projections rather
than full entities so only the fields needed for each UI context are fetched and
field renames are caught at compile time.

## Basic Projection

```java
// Projection for list views — only fields needed for the grid row
public interface EmployeeListItemProjection {
    Long getKey();
    String getFirstName();
    String getLastName();
    String getStatus();
}

// Projection for edit forms — all editable fields
public interface EmployeeDetailProjection {
    Long getKey();
    String getFirstName();
    String getLastName();
    String getEmail();
    String getStatus();
    Long getDepartmentKey();
}

// Entity implements its own projections — renaming a field produces a compile error
@Entity
@Table(name = "employees")
@NoArgsConstructor @Getter @Setter
public class EmployeeEntity extends BaseEntity<Long>
        implements EmployeeListItemProjection, EmployeeDetailProjection {
    // field names must satisfy all implemented projection interfaces
}
```

Name projections for their UI context, not generic suffixes:
- `EmployeePickerItemProjection` — for a `ComboBox` or `Select`
- `EmployeeListItemProjection` — for the grid
- `EmployeeDetailProjection` — for the edit form
