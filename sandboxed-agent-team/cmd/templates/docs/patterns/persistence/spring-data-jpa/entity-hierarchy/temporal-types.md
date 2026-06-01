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

**MapStruct conversion:** An abstract `InstantMapper` class handles the conversion,
injected with a service that provides the current user's timezone. Include
`InstantMapper.class` in each mapper's `uses` clause:

```java
@Mapper(componentModel = SPRING, uses = {InstantMapper.class})
public abstract class OrderMapper {
    abstract OrderDetail toDetail(OrderDetailProjection projection);
}
```

See `docs/patterns/ui/vaadin/client-details-service.md` for the full
`ClientDetailsService` pattern.
