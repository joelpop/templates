# RootEntity

`RootEntity<KEY>` is the base of the entity hierarchy, providing key-based identity and Hibernate-proxy-safe `equals`/`hashCode`.

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
