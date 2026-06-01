# ClientDetailsService MapStruct Integration

When a MapStruct mapper needs `Instant` → `LocalDateTime` conversion, declare `ClientDetailsService` in the `uses` list so the timezone-aware conversion is applied automatically without explicit `@Mapping` annotations.

```java
@Mapper(
        componentModel = MappingConstants.ComponentModel.SPRING,
        injectionStrategy = InjectionStrategy.CONSTRUCTOR,
        uses = {AuditMapper.class, ClientDetailsService.class})
public interface EquipmentMapper {
    EquipmentDetail toDetail(EquipmentDetailProjection projection);
}
```

MapStruct injects `ClientDetailsService` as a constructor dependency and routes
`Instant` → `LocalDateTime` field conversions through `toBrowserTime` automatically — no
explicit `@Mapping` is needed for those fields.
