# @AttributeConverter for Custom Types

When a Java type maps to a single database column and is not a simple enum, write an `@AttributeConverter` with `autoApply = true` so the conversion is applied uniformly across all entities.

```java
@Converter(autoApply = true)
public class MoneyConverter implements AttributeConverter<Money, BigDecimal> {
    @Override
    public BigDecimal convertToDatabaseColumn(Money m) { return m == null ? null : m.amount(); }
    @Override
    public Money convertToEntityAttribute(BigDecimal v) { return v == null ? null : new Money(v); }
}
```

`autoApply = true` applies the converter to every field of that Java type across
all entities.
