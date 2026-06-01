# MapStruct Mapper Pattern

When mapping between JPA projections, UI models, and entities, declare a MapStruct
mapper in `{app}-jpaservice` with the `@MappingTarget` pattern for updates so unrelated
entity fields are preserved.

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface EmployeeMapper {

    // Projection → UI model (for reads)
    EmployeeListItem toListItem(EmployeeListItemProjection projection);
    List<EmployeeListItem> toListItems(List<EmployeeListItemProjection> projections);
    EmployeeDetail toDetail(EmployeeDetailProjection projection);

    // UI model → entity (for creates and updates)
    // @MappingTarget updates the existing entity in place — leaves unrelated fields untouched
    EmployeeEntity toEntity(EmployeeDetail detail, @MappingTarget EmployeeEntity entity);
}
```

The `@MappingTarget` pattern for updates is critical: it overwrites only the fields
present in the UI model and leaves other entity fields (audit fields, version, etc.)
untouched.
