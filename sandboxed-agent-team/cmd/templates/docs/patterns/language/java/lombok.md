# Lombok Guidelines

When writing an enum with properties, a data POJO, or a class that needs a
logger, apply Lombok's annotation subset for the class type so boilerplate is
eliminated without proxy failures, lazy-load exceptions, or identity errors.

## for JPA Entities

See `docs/patterns/persistence/spring-data-jpa/entity-lombok.md` for the full
treatment — safe vs. unsafe annotations and managed collection suppression.

Short form: use `@NoArgsConstructor`, `@Getter`, and `@Setter`. Never use
`@Data`, `@EqualsAndHashCode`, or `@ToString` on an entity.

## for Enumerations

Enums that carry additional per-constant properties (labels, sort order, symbols,
categories, etc.) are a good fit for `@Getter` + `@RequiredArgsConstructor` — Lombok
generates the constructor Java requires for enum constants and the accessors for the
fields, removing boilerplate without touching any identity semantics.

```java
@Getter
@RequiredArgsConstructor
public enum PriorityCode {
    LOW("Low", 1),
    MEDIUM("Medium", 2),
    HIGH("High", 3),
    URGENT("Urgent", 4);

    private final String label;
    private final int sortOrder;
}
```

Enum constants are singletons, so `equals`/`hashCode`/`toString` identity is already
correct by default — **do not** apply `@Data`, `@EqualsAndHashCode`, or `@ToString` to an
enum. Override `toString()` by hand only if you specifically need a display form
different from the constant name (most code should call the dedicated getter — e.g.,
`PriorityCode.HIGH.getLabel()` — rather than rely on `toString()`).

Fields on an enum should be `private final` so each constant is immutable; `@Setter`
has no place here.

## for POJOs

For mutable POJOs (DTOs, request/response objects, service data containers), use
`@Data`. It generates getters, setters, `equals`/`hashCode`, `toString`, and a
required-args constructor in one annotation — correct for plain Java objects where
field-based identity is expected:

```java
@Data
public class EmployeeDetail {
    private Long key;
    private String firstName;
    private String lastName;
    private EmploymentStatus status;
}
```

For immutable POJOs, use `@Value`. It declares all fields `private final`, generates
only getters (no setters), and produces a stable `equals`/`hashCode`:

```java
@Value
public class EmployeeFilter {
    String lastName;
    EmploymentStatus status;
}
```

When a POJO needs a builder for complex construction, add `@Builder`. On a plain POJO
(not a JPA entity), `@Builder` works standalone — it generates the required all-args
constructor automatically:

```java
@Data
@Builder
public class EmployeeDetail {
    private Long key;
    private String firstName;
    private String lastName;
}
```

`@Data` and `@EqualsAndHashCode` are safe on POJOs because there are no Hibernate
proxies, lazy-loading side effects, or bidirectional relationships to navigate. The
same annotations are unsafe on JPA entities for exactly those reasons — see
`docs/patterns/persistence/spring-data-jpa/entity-lombok.md`.

## for Logging

Use Lombok's `@Slf4j` annotation on any class that needs a logger. It generates a
ready-to-use `log` field with no boilerplate and no copy-paste hazard of declaring the
logger against the wrong class:

```java
@Service
@Slf4j
public class JpaEmployeeService implements EmployeeService {

    @Override
    @Transactional
    public EmployeeDetail create(EmployeeDetail detail) {
        log.info("Creating employee with id {}", detail.getKey());
        // ...
    }
}
```

Use SLF4J parameterized logging (`log.debug("loaded {} records", count)`), not string
concatenation. Parameters are stringified only when the level is enabled, which matters
for DEBUG/TRACE paths.

See `docs/patterns/security/data-protection/pii-logging.md` for the rule against
logging user-identifying information at INFO level and below.