# @AttributeConverter for Custom Types

When a Java type maps to a single database column and is not a plain enum,
write an `@Converter(autoApply = true)` so the conversion is applied
uniformly across all entities without repeating `@Convert` at every field.

For example, a project might define a `Money` value object to hold a
monetary amount. The database stores it as a `DECIMAL` column, which JPA maps to
`BigDecimal`; the converter translates between `Money` and `BigDecimal`
on the way in and out.

```java
// Avoid — @Convert repeated at every entity field that uses Money
@Converter
public class MoneyConverter implements AttributeConverter<Money, BigDecimal> { /* ... */ }

@Entity
public class Invoice {

    /* ... */

    @Convert(converter = MoneyConverter.class)
    private Money price;

    @Convert(converter = MoneyConverter.class)
    private Money tax;
}
```

```java
// Preferred — autoApply = true; no @Convert needed at any field
@Converter(autoApply = true)
public class MoneyConverter implements AttributeConverter<Money, BigDecimal> {

    @Override
    public BigDecimal convertToDatabaseColumn(Money m) {
        return m == null ? null : m.amount();
    }

    @Override
    public Money convertToEntityAttribute(BigDecimal v) {
        return v == null ? null : new Money(v);
    }
}

@Entity
public class Invoice {

    /* ... */

    private Money price;
    private Money tax;
}
```

## Null handling

Always pass `null` through in both directions. JPA calls the converter even
for null column values; returning `null` lets the column's null constraint
remain the authority. Throwing on null input moves that responsibility into
the wrong layer.

## Placement

Place the converter in the persistence module (`{app}-data`) alongside other
JPA infrastructure. The type it converts (e.g., `Money`) can live in the
domain module with no JPA dependency; the converter is the only artifact that
bridges the two.

JPA must scan the converter's package. Verify the base packages in your JPA
configuration include it — see `persistence/spring-data-jpa/jpa-scan-config.md`.

## When not to use

- **Plain enums** — use `@Enumerated(EnumType.STRING)` or a named-column
  strategy; see `persistence/spring-data-jpa/enum-mapping.md`.
- **Multi-column types** — use `@Embeddable` / `@Embedded` instead.

## Related

- `persistence/spring-data-jpa/enum-mapping.md` — enum-specific column mapping
- `persistence/spring-data-jpa/jpa-scan-config.md` — ensuring converters are scanned
- `persistence/spring-data-jpa/entity-hierarchy/temporal-types.md` — `InstantConverter` and other temporal converters
