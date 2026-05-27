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

## Projection Inheritance — Detail Extends List-Item

A detail projection extends the list-item projection rather than re-declaring
its fields. The detail query returns everything the grid already shows, plus the
detail-only fields:

```java
public interface EmployeeListItemProjection {
    Long getKey();
    String getFirstName();
    String getLastName();
    String getStatus();
}

// Inherits all list-item getters; adds only what the edit form needs beyond the grid
public interface EmployeeDetailProjection extends EmployeeListItemProjection {
    String getEmail();
    Long getDepartmentKey();
}
```

The entity implements both — satisfying `EmployeeDetailProjection` implies
satisfying `EmployeeListItemProjection` by inheritance:

```java
public class EmployeeEntity extends AuditedEntity<Long>
        implements EmployeeDetailProjection { ... }
```

## Shared Slices — `AuditProjection`

Audit columns (`createdAt`, `updatedAt`, and the audit-by user references) recur
on every audited entity's detail view. Extract a shared slice:

```java
public interface AuditProjection {
    Instant getCreatedAt();
    Instant getUpdatedAt();
    UserAuditProjection getCreatedBy();
    UserAuditProjection getUpdatedBy();
}
```

Any detail projection for an audited entity mixes it in:

```java
public interface EmployeeDetailProjection
        extends EmployeeListItemProjection, AuditProjection {
    String getEmail();
    Long getDepartmentKey();
}
```

## Nested Projections

Returning another projection interface from a getter causes Spring Data to
generate a JOIN and populate the nested object in the same query — no separate
SELECT per row.

`UserAuditProjection` is the canonical example: it exposes only the fields
needed to display an audit-by label ("Created by Alice Anderson"), avoiding a
full `UserEntity` fetch:

```java
public interface UserAuditProjection {
    String getFirstName();
    String getLastName();
}
```

Spring Data resolves `getCreatedBy()` → `UserAuditProjection` via a JOIN on the
audit FK column. The audit-by references can be `null` for system-originated
writes; null-check before rendering.

The same technique applies to any `@ManyToOne` relationship where the caller
only needs a narrow slice of the related entity.

## Picker Projections

A picker projection fetches only the key and the field(s) needed to render the
display label. It is independent of the list-item and detail projections — it
does not extend either:

```java
public interface EmployeePickerItemProjection {
    Long getKey();
    String getFirstName();
    String getLastName();
}
```

The corresponding UI model record implements `HasCaption` to carry the display
label. See `docs/patterns/ui/vaadin/uimodel.md` for the UI-side pattern.
