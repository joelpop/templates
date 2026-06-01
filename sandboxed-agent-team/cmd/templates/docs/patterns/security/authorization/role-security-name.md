# getSecurityName() on UserRole

When a `UserRole` enum constant is mapped to a Spring Security granted-authority string, expose it via `getSecurityName()` — the constant form for annotations and the getter form for programmatic integrations that need a method-of-an-instance accessor.

Spring Security represents authorities as strings and some integrations need a
method-of-an-instance accessor — e.g., when mapping a `UserEntity` field to a
granted-authority list. The getter is the seam; the constants are the
compile-time form for annotations.

```java
// Annotation usage — compile-time constant
@RolesAllowed(UserRole.ROLE_ADMIN)

// Programmatic usage — instance method
authorities.add(new SimpleGrantedAuthority(role.getSecurityName()));
```
