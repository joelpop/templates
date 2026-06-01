# Entity Base Class Hierarchy Overview

When creating JPA entities, extend the layered `@MappedSuperclass` hierarchy so identity, optimistic locking, auditing, and `equals`/`hashCode` are inherited consistently.

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
