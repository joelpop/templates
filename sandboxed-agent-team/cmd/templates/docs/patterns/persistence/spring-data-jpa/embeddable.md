# @Embeddable for Multi-Field Value Objects

When a value object has multiple fields belonging to the same table row — no separate table, no join, no independent identity — use `@Embeddable` rather than a separate entity.

```java
@Embeddable
@NoArgsConstructor @Getter @Setter
public class Address {
    private String street;
    private String city;
    private String state;
    private String postalCode;
}

@Entity
public class EmployeeEntity extends BaseEntity<Long> {
    @Embedded
    private Address homeAddress;
}
```

When the same `@Embeddable` is embedded twice, use `@AttributeOverrides` to
disambiguate column names.
