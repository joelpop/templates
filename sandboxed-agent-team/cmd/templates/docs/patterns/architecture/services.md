# Service Layer Patterns

Conventions for service interfaces, service implementations, transaction management,
data loading, and MapStruct mapping in Vaadin 24+ with Spring Boot 3+ and Spring Data
JPA projects. The patterns in this document are identical across every supported Vaadin
and Spring Boot line.

## Service Interface Contracts

Service interfaces operate exclusively on UI model objects (POJOs from the `{app}-uimodel`
module). No service interface method signature references a JPA entity, interface
projection, repository, or MapStruct mapper. This is enforced at compile time by the
module structure (see `docs/agnostic/architecture/modules.md`).

### Query / Mutation Separation

Group query operations separately from mutation operations within a single service
interface — queries together, mutations together, with a comment separator.

```java
public interface EmployeeService {

    // --- Queries ---
    EmployeeDetail findByKey(long key);
    List<EmployeeListItem> listAll();
    boolean isEmailAvailable(String email, long excludeKey);

    // --- Mutations ---
    EmployeeDetail create(EmployeeDetail detail);
    EmployeeDetail update(EmployeeDetail detail);
    void deactivate(long key);
}
```

`create` and `update` both accept the UI model and return the updated UI model (with
server-assigned fields populated — key on create, refreshed audit fields on update).
`deactivate` takes just the key because there is no detail to send.

### @Transactional Annotations on Implementations

Service implementations apply `@Transactional(readOnly = true)` to all query methods and
`@Transactional` to all mutation methods. Annotations go on the **implementation**, not
the interface.

```java
@Service
public class JpaEmployeeService implements EmployeeService {

    @Override
    @Transactional(readOnly = true)
    public EmployeeDetail findByKey(long key) { ... }

    @Override
    @Transactional(readOnly = true)
    public List<EmployeeListItem> listAll() { ... }

    @Override
    @Transactional
    public EmployeeDetail create(EmployeeDetail detail) { ... }

    @Override
    @Transactional
    public EmployeeDetail update(EmployeeDetail detail) { ... }

    @Override
    @Transactional
    public void deactivate(long key) { ... }
}
```

`readOnly = true` tells Hibernate to skip dirty checking on flush. This reduces overhead
on every query-only service method.

Keep `@Transactional` at the service layer, not the repository layer. Business logic
boundaries are the service's responsibility.

## Update Pattern — Dirty Checking

For updates, load the managed entity using the key carried by the UI model, apply
changes via MapStruct, and let the transaction flush. The method returns the refreshed
UI model so the caller has the latest state without a re-query:

```java
@Override
@Transactional
public EmployeeDetail update(EmployeeDetail detail) {
    var entity = repository.findById(detail.getKey())
        .orElseThrow(() -> new EntityNotFoundException(detail.getKey()));
    mapper.toEntity(detail, entity);
    // transaction flush performs the UPDATE automatically — no save() call needed
    return mapper.toDetail(entity);
}
```

The entity itself implements the detail projection interface (see
`docs/agnostic/architecture/persistence.md` → Interface Projections), so
`mapper.toDetail(entity)` is valid without a separate projection query.

Do not call `save()` on a managed (already-loaded) entity — it is redundant.
Do not call `save()` on a detached entity — it triggers a full-column overwrite.
See `docs/agnostic/architecture/persistence.md` for full explanation.

## Insert Pattern — save() for New Entities Only

```java
@Override
@Transactional
public EmployeeDetail create(EmployeeDetail detail) {
    var entity = new EmployeeEntity();
    mapper.toEntity(detail, entity);
    var saved = repository.save(entity);   // null key triggers persist()
    return mapper.toDetail(saved);
}
```

## Grid Data Loading Pattern

Grids load the full dataset for the current context into memory and perform sorting and
filtering in-memory using Vaadin's `ListDataProvider`. This eliminates the complexity of
server-side pagination for datasets that fit comfortably in memory.

```java
// Service method returns the complete list
@Transactional(readOnly = true)
public List<EmployeeListItem> listAll() {
    return mapper.toListItems(repository.findAll(...));
}
```

```java
// View uses ListDataProvider
var items = employeeService.listAll();
var dataProvider = new ListDataProvider<>(items);
grid.setDataProvider(dataProvider);
```

Quick Filter and column sorting are applied by the in-memory `ListDataProvider`, not by
database queries. No `Pageable`, `Page<T>`, `CallbackDataProvider` size/fetch callbacks,
or offset/limit parameters are used for grid display.

If data volume eventually requires server-side pagination, migrate to `CallbackDataProvider`
at that point — do not add the complexity preemptively.

## Caching at the Service Layer

Use Spring `@Cacheable` at the service layer for frequently read, rarely changing data.
This is the preferred caching approach — simpler and more predictable than Hibernate's
L2 cache, and works with projections and UI models (not just entities).

```java
@Cacheable("departments")
@Transactional(readOnly = true)
public List<DepartmentListItem> listAllDepartments() { ... }

@CacheEvict(value = "departments", allEntries = true)
@Transactional
public void createDepartment(DepartmentDetail detail) { ... }
```

Enable caching with `@EnableCaching` on a configuration class. For most applications,
`@Cacheable` at the service layer covers the majority of caching needs.

## Error Contracts

Service methods throw typed exceptions for predictable failure modes. Views catch these
exceptions and display appropriate error messages — not raw exception text.

Define exceptions in the `{app}-service` module so both service implementations and views
can reference them without violating layer boundaries:

```java
// Entity not found (or belongs to a different context)
public class EntityNotFoundException extends RuntimeException {
    public EntityNotFoundException(long key) {
        super("Entity not found: " + key);
    }
}

// Business rule violation — may include multiple error messages
public class ValidationException extends RuntimeException {
    private final List<String> errors;

    public ValidationException(String message) {
        super(message);
        this.errors = List.of(message);
    }

    public ValidationException(List<String> errors) {
        super(String.join("; ", errors));
        this.errors = List.copyOf(errors);
    }

    public List<String> getErrors() { return errors; }
}
```

Only these two typed exceptions are defined. Uniqueness violations, required-field
omissions, cross-entity checks, and every other predictable business-rule failure surface
as a `ValidationException` with a field-specific message — not as a dedicated subclass
per field. A dedicated `DuplicateKeyException` / `DuplicateDisplayIdException` / etc.
proliferates exception types for no benefit; callers distinguish cases by the message,
not by the class.

Service implementations catch `DataIntegrityViolationException` from the database and
translate them to `ValidationException` before propagating to the UI. Raw database
error messages must never reach the user.

## MapStruct Mapper Pattern

Mappers live in `{app}-jpaservice`. Each mapper converts between the three representations:
JPA interface projection → UI model, and UI model → JPA entity.

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

The `@MappingTarget` pattern for updates is critical: it overwrites only the fields present
in the UI model and leaves other entity fields (audit fields, version, etc.) untouched.

## OSIV Disabled — Service Layer Must Load All Data

With `spring.jpa.open-in-view=false`, the Hibernate session closes at the end of the
service method. Any lazy association not loaded within the transaction will throw
`LazyInitializationException` if accessed later.

Service methods must load all data the view needs before returning. Use projections (not
full entities) for list views to avoid lazy-loading issues and over-fetching.

Never pass JPA entities to the view layer — return UI models only.
