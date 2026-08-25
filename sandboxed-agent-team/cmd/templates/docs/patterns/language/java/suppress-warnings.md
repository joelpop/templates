# Suppressing Warnings

Every `@SuppressWarnings` annotation must carry an inline `//` comment naming the
specific warning being suppressed and the reason. A reader should not have to leave
the file to understand why a warning was silenced.

```java
// Avoid — no context
@SuppressWarnings("java:S2160")
public class EmployeeEntity extends BaseEntity<Long> { /* ... */ }
```

```java
// Preferred
@SuppressWarnings("java:S2160") // key-based equality is inherited from RootEntity
public class EmployeeEntity extends BaseEntity<Long> { /* ... */ }

@SuppressWarnings("unchecked") // raw type comes from legacy third-party API
var bean = (Map<String, Object>) legacyApi.getProperties();
```

Without a comment, a suppression is indistinguishable from a cargo-culted one — future
readers can't tell whether it's still load-bearing or the underlying issue has been fixed.
