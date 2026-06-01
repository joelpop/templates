# Lombok on POJOs

When a POJO needs getters, setters, equals/hashCode, and toString, use `@Data`; for immutable POJOs use `@Value`; add `@Builder` when complex construction is needed.

For mutable POJOs (DTOs, request/response objects, service data containers):

```java
@Data
public class EmployeeDetail {
    private Long key;
    private String firstName;
    private String lastName;
    private EmploymentStatus status;
}
```

For immutable POJOs, prefer a Java record — it is the idiomatic choice on Java 16+:

```java
public record EmployeeFilter(String lastName, EmploymentStatus status) {}
```

Use `@Value` when a record does not fit: inheritance required, a framework needs a
no-arg constructor, or the type must interoperate with a Lombok-based hierarchy.

```java
@Value
public class EmployeeFilter {
    String lastName;
    EmploymentStatus status;
}
```

When a POJO needs a builder for complex construction, add `@Builder`. On a plain POJO
(not a JPA entity), `@Builder` generates the required all-args constructor automatically:

```java
@Data
@Builder
public class EmployeeDetail {
    private Long key;
    private String firstName;
    private String lastName;
}
```

`@Data` and `@EqualsAndHashCode` are typically safe on POJOs — there are no Hibernate proxies,
lazy-loading side effects, or bidirectional relationships to navigate. If a mutable
POJO needs to serve as a `HashMap` key, prefer making it immutable (record or
`@Value`) rather than hand-writing `equals`/`hashCode`. In inheritance hierarchies,
set `@EqualsAndHashCode(callSuper = true)` or parent fields are silently excluded.
The same annotations are unsafe on JPA entities — see
`docs/patterns/persistence/spring-data-jpa/entity-lombok.md`.
