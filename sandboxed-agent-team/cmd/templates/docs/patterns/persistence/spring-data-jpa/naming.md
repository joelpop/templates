# Spring Data JPA Naming Conventions

When creating JPA entities, entity enums, repositories, or interface
projections, use these naming conventions so the persistence-layer role
of each artifact is unambiguous from its name alone.

## Artifact Naming

| Type                     | Convention                                                 | Example                                         |
|--------------------------|------------------------------------------------------------|-------------------------------------------------|
| JPA entity               | Suffix `Entity`                                            | `EquipmentEntity`, `TenantEntity`               |
| JPA entity enum          | Suffix `Code`                                              | `EquipmentTypeCode`, `UserRoleCode`             |
| Spring Data repository   | Suffix `Repository`                                        | `EquipmentRepository`                           |
| JPA interface projection | Suffix `Projection`; prefix matches the UI model it serves | `EquipmentListItemProjection`                   |

**Rationale:**
- `Entity` signals a JPA-managed object with lifecycle implications (proxies,
  lazy loading, dirty checking). Callers that see `Entity` know it must not
  cross the service boundary.
- `Code` signals "this value is persisted as a code string." Code enums are
  plain constants — no properties, no Lombok — so the JPA enum stays free of
  presentation concerns.
- `Repository` identifies Spring Data JPA interfaces. Repositories are used
  only from service implementations, never from views or service interfaces.
- `Projection` signals a read-only, partial view of data with no persistence
  lifecycle. Projections are named for the UI context they serve so the pairing
  with the UI model they supply is obvious.

## Projection Naming

Projections are named for the UI context they serve, not for the entity or
query method:

| UI context    | JPA interface projection      | UI model              |
|---------------|-------------------------------|-----------------------|
| Grid row      | `EquipmentListItemProjection` | `EquipmentListItem`   |
| Form / detail | `EquipmentDetailProjection`   | `EquipmentDetail`     |

The projection name mirrors the UI model name with `Projection` appended.
For new contexts, choose a UI-contextual name — not a generic `Summary`,
`Info`, or `Data` suffix.

## Property and Column Naming for Keys and Identifiers

Entity properties that represent keys or business identifiers should use these suffixes — we
always control these names, regardless of what the database column is named:

| Property suffix | Meaning                                                          | Examples                    |
|-----------------|------------------------------------------------------------------|-----------------------------|
| `Key`           | Surrogate PK/FK — system-generated, opaque, no business meaning  | `equipmentKey`, `tenantKey` |
| `Id`            | Business identifier — human-meaningful or externally assigned    | `displayId`, `employeeId`   |

Never use `Id` for a surrogate primary key — that suffix is reserved for business identifiers.

The recommended database column naming mirrors this convention (`equipment_key`, `display_id`).
When the database is shared or externally managed, column names may differ from this recommendation.
Use `@Column(name = "...")` to map a well-named entity property to whatever the database column
is actually named:

```java
@Column(name = "equip_pk")   // external DB standard — but the entity property stays meaningful
private Long equipmentKey;
```

## Enum Column Storage

Code enums are plain constants — no properties, no Lombok necessary:

```java
public enum EquipmentTypeCode {
    VEHICLE, AIRCRAFT, MACHINERY, WATERCRAFT
}
```

Display properties belong in the corresponding UI type enum, not here.

Store as string, never ordinal — ordinal values break silently if declaration
order changes:

```java
@Enumerated(EnumType.STRING)
private EquipmentTypeCode equipmentTypeCode;
```