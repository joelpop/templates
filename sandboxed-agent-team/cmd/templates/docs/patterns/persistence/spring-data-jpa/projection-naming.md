# Projection Naming Convention

When naming a JPA interface projection, use the UI context it serves as the
name prefix so the pairing with the UI model it supplies is obvious at a glance.

Projections are named for the UI context they serve, not for the entity or
query method:

| UI context    | JPA interface projection      | UI model              |
|---------------|-------------------------------|-----------------------|
| Grid row      | `EquipmentListItemProjection` | `EquipmentListItem`   |
| Form / detail | `EquipmentDetailProjection`   | `EquipmentDetail`     |

The projection name mirrors the UI model name with `Projection` appended.
For new contexts, choose a UI-contextual name — not a generic `Summary`,
`Info`, or `Data` suffix.