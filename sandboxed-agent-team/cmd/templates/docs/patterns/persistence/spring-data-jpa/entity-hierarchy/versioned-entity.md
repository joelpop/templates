# VersionedEntity

`VersionedEntity<KEY>` extends `RootEntity` with a `@Version` column for optimistic locking — no manual version-checking code needed.

```java
@MappedSuperclass
public abstract class VersionedEntity<KEY> extends RootEntity<KEY> {
    @Version
    private Long version;
}
```

Hibernate increments the version column on every UPDATE; on a version mismatch at flush
time, Spring raises `ObjectOptimisticLockingFailureException`.
