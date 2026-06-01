# Key and Identifier Property Naming

When naming entity properties that represent keys or business identifiers, use
`Key` for surrogate PKs/FKs and `Id` for business identifiers so the semantic
distinction is clear from the property name alone.

| Property suffix | Meaning                                                          | Examples                    |
|-----------------|------------------------------------------------------------------|-----------------------------|
| `Key`           | Surrogate PK/FK — system-generated, opaque, no business meaning  | `equipmentKey`, `tenantKey` |
| `Id`            | Business identifier — human-meaningful or externally assigned    | `displayId`, `employeeId`   |

Never use `Id` for a surrogate primary key — that suffix is reserved for
business identifiers.

The recommended database column naming mirrors this convention (`equipment_key`,
`display_id`). When the database is shared or externally managed, column names
may differ from this recommendation. Use `@Column(name = "...")` to map a
well-named entity property to whatever the database column is actually named:

```java
@Column(name = "equip_id")   // external DB standard — but the entity property stays meaningful
private Long equipmentKey;
```