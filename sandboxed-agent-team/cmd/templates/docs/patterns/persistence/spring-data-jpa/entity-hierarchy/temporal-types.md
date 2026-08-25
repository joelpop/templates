# Temporal Types — Instant for Storage, LocalDateTime for Display

When persisting date/time fields, always use `java.time.Instant` (UTC) on entity fields; convert to `LocalDateTime` (or `LocalDate`/`LocalTime`) only at the service/mapper layer for display.

**Entity layer:** `Instant` for all timestamp fields. `LocalDateTime`, `ZonedDateTime`,
`LocalDate`, or `LocalTime` must never be used on JPA entity fields.

```java
@Entity
public class OrderEntity extends BaseEntity<Long> {
    // Inherited: createdAt (Instant), updatedAt (Instant)
    private Instant activationDate;
    private Instant deactivationDate;
}
```

**UI model layer:** `LocalDateTime` (or `LocalDate` / `LocalTime` where appropriate)
for all timestamps shown to users. Conversion from `Instant` happens in the
service/mapper layer using the user's configured timezone.

**MapStruct conversion:** Include `ClientDetailsService.class` in each mapper's
`uses` clause. MapStruct injects it as a constructor dependency and routes
`Instant` → `LocalDateTime` conversions through `toBrowserTime` and
`LocalDateTime` → `Instant` conversions through `toServerTime` automatically —
no explicit `@Mapping` is needed for those fields:

```java
@Mapper(componentModel = SPRING, uses = {ClientDetailsService.class})
public interface OrderMapper {
    OrderDetail toDetail(OrderDetailProjection projection);  // Instant → LocalDateTime
    OrderEntity toEntity(OrderDetail detail);               // LocalDateTime → Instant
}
```

**Related:** `docs/patterns/ui/vaadin/client-details-service.md` — the
`ClientDetailsService` interface; `docs/patterns/ui/vaadin/client-details-mapstruct.md`
— MapStruct wiring detail.
