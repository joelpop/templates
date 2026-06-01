# JPA Artifact Naming Conventions

When creating JPA entities, entity enums, repositories, or interface projections,
use these naming conventions so the persistence-layer role of each artifact is
unambiguous from its name alone.

| Type                     | Declaration                               | Convention                                                 | Example                       |
|:-------------------------|:------------------------------------------|:-----------------------------------------------------------|:------------------------------|
| JPA entity               | `@Entity` class                           | Suffix `Entity`                                            | `EquipmentEntity`             |
| JPA entity enum          | `enum`                                    | Suffix `Code`                                              | `EquipmentTypeCode`           |
| Spring Data repository   | `interface` extends `JpaRepository<E, K>` | Suffix `Repository`                                        | `EquipmentRepository`         |
| JPA interface projection | `interface`                               | Suffix `Projection`; prefix matches the UI model it serves | `EquipmentListItemProjection` |

**Rationale:**
- `Entity` signals a JPA-managed object with lifecycle implications (proxies,
  lazy loading, dirty checking). Callers that see `Entity` know it must not
  cross the service boundary.
- `Code` signals "this value is persisted as a code string." Code enums are
  plain constants — no properties — so the JPA enum stays free of
  presentation concerns.
- `Repository` identifies Spring Data JPA interfaces. Repositories are used
  only from service implementations, never from views or service interfaces.
- `Projection` signals a read-only, partial view of data with no persistence
  lifecycle. Projections are named for the UI context they serve so the pairing
  with the UI model they supply is obvious.
