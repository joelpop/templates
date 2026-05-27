# Entity Base Class Hierarchy

When creating Spring Data JPA entities, extend the layered `@MappedSuperclass`
hierarchy so identity, optimistic locking, auditing, and `equals`/`hashCode` are
inherited consistently and never re-implemented per entity.

## Hierarchy Overview

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

Most principal entities extend `BaseEntity`. Reference tables that need no audit
trail may extend `RootEntity` or `VersionedEntity` directly.

Note: the field is named `key`, not `id`. See
`docs/patterns/persistence/spring-data-jpa/naming.md` for the `_key` vs `_id`
distinction.

## RootEntity

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

- `instanceof RootEntity<?>` passes for Hibernate proxies without initializing them
- `effectiveClass()` unwraps proxies — prevents entities of different types with the
  same key from being considered equal
- Only persisted entities (non-null key) are considered equal
- Constant `hashCode` ensures stability across managed, detached, and proxy states
- No entity ever writes `equals`/`hashCode` — inherited via `BaseEntity`

## VersionedEntity

```java
@MappedSuperclass
public abstract class VersionedEntity<KEY> extends RootEntity<KEY> {
    @Version
    private Long version;
}
```

`@Version` provides optimistic locking. Hibernate increments the version column
on every UPDATE; on a version mismatch at flush time, Spring raises
`ObjectOptimisticLockingFailureException`. No manual version-checking code
needed.

## AuditedEntity

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
`docs/patterns/structure/modules.md`).

`created_by` and `updated_by` are `@ManyToOne` references to the users table.
The `AuditorAware` bean reads the surrogate key off the authenticated principal
and resolves it via `EntityManager.getReference` — a Hibernate proxy holding
just the key, persisted as the FK with no `SELECT`. The audit columns
(`created_by_key`, `updated_by_key`) are nullable so system-originated writes
(no authenticated principal) leave them NULL.

## Temporal Types — Instant for Storage, LocalDateTime for Display

All persisted date/time fields **must** use `java.time.Instant` (UTC). No
exceptions. `LocalDateTime`, `ZonedDateTime`, `LocalDate`, or `LocalTime` must
never be used on JPA entity fields.

**Entity layer:** `Instant` for all timestamp fields.

```java
@Entity
public class OrderEntity extends BaseEntity<Long> {
    // Inherited: createdAt (Instant), updatedAt (Instant)
    private Instant activationDate;
    private Instant deactivationDate;
}
```

**UI model layer:** `LocalDateTime` (or `LocalDate` / `LocalTime` where
appropriate) for all timestamps shown to users. Conversion from `Instant`
happens in the service/mapper layer using the user's configured timezone.

**MapStruct conversion:** An abstract `InstantMapper` class handles the
conversion, injected with a service that provides the current user's timezone.
It exposes:

- `toBrowserTime(Instant)` — converts UTC to the user's local time for display
- `toServerTime(LocalDateTime)` — converts user-entered local time back to UTC
  for storage

Include `InstantMapper.class` in each mapper's `uses` clause:

```java
@Mapper(componentModel = SPRING, uses = {InstantMapper.class})
public abstract class OrderMapper {
    abstract OrderDetail toDetail(OrderDetailProjection projection);
}
```

See `docs/patterns/ui/vaadin/datetime.md` → "Browser Client Details —
Bridging the SoC Wall" for the full `ClientDetailsService` pattern.
