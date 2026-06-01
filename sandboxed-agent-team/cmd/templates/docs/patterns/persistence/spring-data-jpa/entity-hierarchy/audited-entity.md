# AuditedEntity

`AuditedEntity<KEY>` extends `VersionedEntity` with JPA auditing fields — `createdAt`, `updatedAt`, `createdBy`, `updatedBy` — populated automatically by Spring Data's auditing listener.

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

The `AuditorAware` bean reads the surrogate key off the authenticated principal
and resolves it via `EntityManager.getReference` — a Hibernate proxy holding
just the key, persisted as the FK with no `SELECT`. The audit columns
(`created_by_key`, `updated_by_key`) are nullable so system-originated writes
(no authenticated principal) leave them NULL.
