# UI Model Naming Conventions

When defining types in the `{app}-uimodel` module, name them for the UI
context they serve — not for the entity or data source — so the name stays
stable even if the backing technology changes.

| Type          | Convention                                      | Example                              |
| :------------ | :---------------------------------------------- | :----------------------------------- |
| Data POJO     | No suffix — named for its UI context            | `EmployeeListItem`, `EmployeeDetail` |
| Enum          | No suffix — named for the concept it represents | `EmploymentStatus`, `PhoneType`      |
| Picker record | Suffix `PickerItem`                             | `EmployeePickerItem`                 |

Avoid generic suffixes like `Summary`, `Info`, or `Data` — `EmployeeListItem`
is self-explaining; `EmployeeSummary` is not.

## Pairing with the Implementation Layer

UI model POJOs sit at the service boundary. What they pair with on the
implementation side depends on the backing technology:

| Implementation   | Counterpart                                                                                                      |
| :--------------- |:-----------------------------------------------------------------------------------------------------------------|
| Spring Data JPA  | JPA interface projection (`EmployeeListItemProjection`) — see `persistence/spring-data-jpa/projection-naming.md` |
| REST client      | Response DTO or deserialized payload mapped by the service impl                                                  |
| Internal/derived | Assembled in the service impl from other sources — no counterpart                                                |