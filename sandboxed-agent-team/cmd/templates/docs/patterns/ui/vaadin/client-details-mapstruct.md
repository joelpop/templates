# ClientDetailsService MapStruct Integration

When a MapStruct mapper needs `Instant` → `LocalDateTime` conversion, declare `ClientDetailsService` in the `uses` list so the timezone-aware conversion is applied automatically without explicit `@Mapping` annotations.

```java
@Mapper(
        componentModel = MappingConstants.ComponentModel.SPRING,
        injectionStrategy = InjectionStrategy.CONSTRUCTOR,
        uses = {ClientDetailsService.class})
public interface EquipmentMapper {
    EquipmentDetail toDetail(EquipmentDetailProjection projection);
}
```

MapStruct injects `ClientDetailsService` as a constructor dependency and routes
`Instant` → `LocalDateTime` field conversions through its `toBrowserTime` method
automatically — no explicit `@Mapping` is needed for those fields.

**Related:** `client-details-service.md` — the `ClientDetailsService` interface definition;
`client-details-impl.md` — the Vaadin implementation.
