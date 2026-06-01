# Managed Collection Bidirectional Sync

When a `@OneToMany` or `@ManyToMany` collection requires bidirectional synchronization, suppress Lombok's getter/setter and provide manual implementations that maintain both sides of the relationship.

```java
@Getter(AccessLevel.NONE)
@Setter(AccessLevel.NONE)
@OneToMany(mappedBy = "employee", cascade = CascadeType.ALL, orphanRemoval = true)
private List<PhoneEntity> phones = new ArrayList<>();

// Manual getter — returns unmodifiable view so callers cannot bypass sync helpers
public List<PhoneEntity> getPhones() {
    return Collections.unmodifiableList(phones);
}

// Manual setter — routes through addPhone to maintain both sides of the relationship
public void setPhones(List<PhoneEntity> phones) {
    this.phones.clear();
    phones.forEach(this::addPhone);
}

// Varargs add helper — maintains the back-reference
public void addPhone(PhoneEntity... phones) {
    Stream.of(phones).forEach(p -> {
        this.phones.add(p);
        p.setEmployee(this);
    });
}

// Varargs remove helper
public void removePhone(PhoneEntity... phones) {
    Stream.of(phones).forEach(this.phones::remove);
}
```

This pattern ensures:
- Callers cannot mutate the collection and bypass back-reference synchronization
- Orphan removal works correctly (items removed from the collection are deleted)
- The bidirectional relationship is always consistent
