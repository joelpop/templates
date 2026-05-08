# Persistence Patterns

JPA best practices for Spring Data JPA projects using Hibernate. Adapt `{base_package}`
to your root Java package.

## Entity Base Class Hierarchy

All JPA entities extend a layered `@MappedSuperclass` hierarchy. Each layer adds concerns
incrementally. Most principal entities extend `BaseEntity`; reference tables that need no
audit trail may extend `RootEntity` or `VersionedEntity` directly.

```
RootEntity<KEY>              @MappedSuperclass
                             - key: Long (@Id, @GeneratedValue IDENTITY)
                             - equals/hashCode: key-based with HibernateProxy handling
  └─ VersionedEntity<KEY>    @MappedSuperclass
                             - version: Long (@Version — optimistic locking)
       └─ AuditedEntity<KEY> @MappedSuperclass
                             - @EntityListeners(AuditingEntityListener.class)
                             - created_at: Instant (@CreatedDate, non-nullable)
                             - updated_at: Instant (@LastModifiedDate, non-nullable)
                             - created_by: UserEntity (@CreatedBy, @ManyToOne LAZY, FK)
                             - updated_by: UserEntity (@LastModifiedBy, @ManyToOne LAZY, FK)
            └─ BaseEntity<KEY> @MappedSuperclass
                             - extends AuditedEntity; most principal entities extend this
```

Note: The field is named `key`, not `id`. See `docs/patterns/conventions/naming.md` for
the `_key` vs `_id` distinction.

### RootEntity Implementation

```java
@MappedSuperclass
@Getter
public abstract class RootEntity<KEY> {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private KEY key;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof RootEntity<?> other)) return false;
        return key != null && key.equals(other.key)
            && effectiveClass(this) == effectiveClass(other);
    }

    @Override
    public int hashCode() {
        return getClass().hashCode(); // stable across managed/detached/proxy states
    }

    private static Class<?> effectiveClass(Object o) {
        return o instanceof HibernateProxy hp
            ? hp.getHibernateLazyInitializer().getPersistentClass()
            : o.getClass();
    }
}
```

**Why this pattern:**
- `instanceof RootEntity<?>` passes for Hibernate proxies without initializing them
- `effectiveClass()` unwraps proxies to get the real entity class — prevents entities of
  different types with the same key from being considered equal
- Only persisted entities (non-null key) are considered equal
- Constant `hashCode` ensures stability across managed, detached, and proxy states
- No entity ever writes `equals`/`hashCode` — inherited via `BaseEntity`

### VersionedEntity

```java
@MappedSuperclass
public abstract class VersionedEntity<KEY> extends RootEntity<KEY> {
    @Version
    private Long version;
}
```

`@Version` provides optimistic locking. Hibernate increments the version column on every
UPDATE. If the stored version does not match at flush time, Spring raises
`ObjectOptimisticLockingFailureException`. No manual version-checking code is needed.

### AuditedEntity

```java
@MappedSuperclass
@EntityListeners(AuditingEntityListener.class)
public abstract class AuditedEntity<KEY> extends VersionedEntity<KEY> {
    @CreatedDate
    private Instant createdAt;

    @LastModifiedDate
    private Instant updatedAt;

    @CreatedBy
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "created_by_key", updatable = false)
    private UserEntity createdBy;

    @LastModifiedBy
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "updated_by_key")
    private UserEntity updatedBy;
}
```

`@EnableJpaAuditing` must be present on a `@Configuration` class, and an
`AuditorAware<UserEntity>` bean must be registered (see
`docs/patterns/architecture/modules.md`).

`created_by` and `updated_by` are `@ManyToOne` references to the users table. The
`AuditorAware` bean reads the surrogate key off the authenticated principal (which carries
it from login-time validation) and resolves it via `EntityManager.getReference` — a
Hibernate proxy holding just the key, persisted as the FK with no `SELECT` issued. The
audit columns (`created_by_key`, `updated_by_key`) are nullable so system-originated
writes (no authenticated principal) leave them NULL.

### Temporal Types — Instant for Storage, LocalDateTime for Display

All persisted date/time fields **must** use `java.time.Instant` (UTC). No exceptions.
`LocalDateTime`, `ZonedDateTime`, `LocalDate`, or `LocalTime` must never be used on
JPA entity fields. This ensures all stored timestamps are timezone-safe regardless of
server location.

**Entity layer (JPA):** `Instant` for all timestamp fields.

```java
@Entity
public class OrderEntity extends BaseEntity<Long> {
    // Inherited: createdAt (Instant), updatedAt (Instant)
    private Instant activationDate;
    private Instant deactivationDate;
}
```

**UI model layer:** `LocalDateTime` (or `LocalDate` / `LocalTime` where appropriate)
for all timestamps shown to users. Conversion from `Instant` happens in the
service/mapper layer using the user's configured timezone.

```java
public class OrderDetail {
    private LocalDateTime activationDate;
    private LocalDateTime deactivationDate;
}
```

**MapStruct conversion:** An abstract `InstantMapper` class handles the conversion.
It is injected with a service that provides the current user's timezone (from their
profile or session configuration) and exposes two methods:

- `toBrowserTime(Instant)` — converts UTC to the user's local time for display
- `toServerTime(LocalDateTime)` — converts user-entered local time back to UTC for storage

Include `InstantMapper.class` in each mapper's `uses` clause so temporal conversions
are automatic:

```java
@Mapper(componentModel = SPRING, uses = {InstantMapper.class})
public abstract class OrderMapper {
    abstract OrderDetail toDetail(OrderDetailProjection projection);
}
```

**Timezone source:** The `InstantMapper` is injected with a `ClientDetailsService`
that provides the user's timezone (typically resolved from the user's profile or
session configuration). See `conventions/vaadin.md` → "Browser Client Details —
Bridging the SoC Wall" for the full `ClientDetailsService` pattern, including
version-specific notes and `VaadinSession` caching guidance.

This pattern ensures:
- Database values are always UTC — portable, comparable, sortable
- Display values match the user's expectations — no timezone confusion
- Conversion logic is centralized in one mapper — not scattered across views

## Disable Open Session in View (OSIV)

```properties
spring.jpa.open-in-view=false
```

OSIV keeps the Hibernate session open through the entire HTTP request, including view
rendering. Disabling it:
- Forces all data loading into the service/transaction layer where it belongs
- Surfaces `LazyInitializationException` immediately in development (good — these are real
  bugs that OSIV was masking)
- Prevents N+1 queries from silently firing during serialization or template rendering

## @ManyToOne Always Lazy

JPA's default for `@ManyToOne` is `FetchType.EAGER` — a spec decision that causes surprise
queries everywhere. Always override it:

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

`ORDINAL` stores the position of the constant in the declaration. Adding or reordering
constants silently maps existing data to the wrong value.

## Insert vs. Update Patterns

**Insert (new entity):** Call `repository.save(newEntity)` on a transient entity (null key).
The null key triggers `persist()`.

**Update (managed entity):** Load the entity by key within a transaction, mutate its fields,
and let the transaction flush. Hibernate's dirty checking detects the changes and issues a
targeted UPDATE. Do not call `save()` on a managed entity — it is redundant.

**Never call `save()` on a detached entity** (one that was loaded in a previous transaction).
This triggers `merge()`, which issues a SELECT + full-column UPDATE that can silently overwrite
concurrent changes. Instead, load a fresh managed copy by key and apply changes selectively.

```java
// Insert
@Transactional
public EmployeeDetail create(EmployeeDetail detail) {
    var entity = new EmployeeEntity();
    mapper.toEntity(detail, entity);
    return mapper.toDetail(repository.save(entity));
}

// Update — dirty checking, no save() call
@Transactional
public void update(long key, EmployeeDetail detail) {
    var entity = repository.findById(key).orElseThrow(() -> new EntityNotFoundException(key));
    mapper.toEntity(detail, entity);
    // transaction flush performs the UPDATE automatically
}
```

## @DynamicUpdate for Wide Tables

For entities with many columns, large columns (TEXT, BLOB), or high write volume, annotate
with `@DynamicUpdate` to limit UPDATE statements to only changed columns:

```java
@Entity
@Table(name = "employees")
@DynamicUpdate
public class EmployeeEntity extends BaseEntity<Long> { ... }
```

Without `@DynamicUpdate`, Hibernate generates a static full-column UPDATE at startup that
overwrites all columns on every flush regardless of what changed. `@DynamicUpdate` adds a
small per-flush overhead to compute the dirty field diff — use it only when profiling confirms
the full-column writes are a meaningful cost.

## Interface Projections

Fetch only the fields needed for a given UI context. Entity implementing its own projections
provides compile-time contract verification: renaming a field on the entity causes an
immediate compiler error.

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

// Entity implements its own projections
@Entity
@Table(name = "employees")
@NoArgsConstructor @Getter @Setter
public class EmployeeEntity extends BaseEntity<Long>
        implements EmployeeListItemProjection, EmployeeDetailProjection {
    // field names must satisfy all implemented projection interfaces
    // renaming a field produces an immediate compile error
}
```

Name projections for their UI context, not generic suffixes:
- `EmployeeListItemProjection` — for the grid
- `EmployeeDetailProjection` — for the edit form

Prefer interface projections over loading full entities for list views and search results.

## Relationship Cascade Strategies

Choose cascade strategy based on the nature of the relationship:

### Association — child has independent existence

Child can exist without the parent; parent deletion should not cascade-delete children.

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

`orphanRemoval = false` (the default). Removing from the collection does not delete the child.
Attempting to delete a parent that still has children fails at the DB level (FK constraint) —
this is the correct behavior.

### Composition — child has no meaning without parent

Child lifecycle belongs entirely to the parent. Deleting the parent deletes the children.
Removing from the collection deletes from the database.

```java
// EmployeeEntity — one-to-many, composition
@Getter(AccessLevel.NONE)
@Setter(AccessLevel.NONE)
@OneToMany(mappedBy = "employee",
           cascade = CascadeType.ALL,
           orphanRemoval = true)
private List<PhoneEntity> phones = new ArrayList<>();

// Manual collection management — see docs/patterns/conventions/lombok.md
public List<PhoneEntity> getPhones() { return Collections.unmodifiableList(phones); }
public void addPhone(PhoneEntity... phones) { ... }
public void removePhone(PhoneEntity... phones) { ... }
```

Rule of thumb: reach for `{PERSIST, MERGE}` first. Escalate to `ALL` + `orphanRemoval = true`
only when the child truly cannot exist without the parent.

## @AttributeConverter for Custom Types

For any Java type that maps to a single DB column and isn't a simple enum:

```java
@Converter(autoApply = true)
public class MoneyConverter implements AttributeConverter<Money, BigDecimal> {
    @Override
    public BigDecimal convertToDatabaseColumn(Money m) { return m == null ? null : m.amount(); }
    @Override
    public Money convertToEntityAttribute(BigDecimal v) { return v == null ? null : new Money(v); }
}
```

`autoApply = true` applies the converter to every field of that Java type across all entities.

## @Embeddable for Multi-Field Value Objects

When a value object has multiple fields that belong to the same table row (no separate table,
no join, no independent identity):

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

When the same `@Embeddable` is embedded twice, use `@AttributeOverrides` to disambiguate
column names.

## Batch and Bulk Operations

For high-volume inserts, updates, or deletes, bypass entity loading:

```java
// Batch inserts — configure batch_size in application.properties
// spring.jpa.properties.hibernate.jdbc.batch_size=50
// spring.jpa.properties.hibernate.order_inserts=true
employeeRepository.saveAll(employees);   // grouped into batches

// Bulk update — single SQL UPDATE, no entity loading
@Modifying
@Query("UPDATE EmployeeEntity e SET e.status = :status WHERE e.department.key = :deptKey")
int updateStatusByDepartment(@Param("status") String status, @Param("deptKey") Long deptKey);

// Bulk delete by IDs — single DELETE WHERE key IN (...)
employeeRepository.deleteAllByIdInBatch(keys);
// AVOID: deleteAll() loads every entity first, then issues N individual DELETEs
```

## Testing Patterns

### @DataJpaTest

Use `@DataJpaTest` for lightweight repository tests. It loads only the JPA layer (no web
layer, no services) and wraps each test in a transaction that rolls back automatically:

```java
@DataJpaTest
class EmployeeRepositoryTest {
    @Autowired EmployeeRepository repository;

    @Test
    void findByDepartment_returnsMatchingEmployees() {
        // H2 in-memory by default; @Transactional rollback after each test
    }
}
```

### Detecting N+1 with Hibernate Statistics

Enable Hibernate statistics in tests to catch N+1 regressions before production:

```properties
# application-test.properties
spring.jpa.properties.hibernate.generate_statistics=true
logging.level.org.hibernate.stat=DEBUG
```

Pair with `datasource-proxy` to assert exact query counts in tests on critical paths.

## N+1 Prevention

The N+1 problem: loading N entities then firing N additional queries to load a lazy
collection. Solutions:

- `@EntityGraph` — declarative, specified on the repository method:

```java
@EntityGraph(attributePaths = {"phones"})
List<EmployeeEntity> findAllWithPhones();
```

- `JOIN FETCH` in `@Query` — explicit, when the query is already in JPQL:

```java
@Query("SELECT e FROM EmployeeEntity e JOIN FETCH e.department WHERE e.status = :status")
List<EmployeeEntity> findActiveWithDepartment(@Param("status") String status);
```

Prefer `@EntityGraph` as the first choice — it keeps fetch semantics separate from the query
string. Use `JOIN FETCH` when a `@Query` is already required.
