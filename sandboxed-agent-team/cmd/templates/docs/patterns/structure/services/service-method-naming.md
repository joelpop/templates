# Service Method Naming

When naming service methods, follow the standard patterns so callers can infer
query vs. mutation semantics from the name alone.

| Pattern           | Used for                        | Example                                  |
|-------------------|---------------------------------|------------------------------------------|
| `listAll*()`      | Query — returns full collection | `listAll()`, `listAllActive()`           |
| `findByKey()`     | Query — returns single result   | `findByKey(long key)`                    |
| `is*Available()`  | Query — uniqueness check        | `isDisplayIdAvailable(id, excludeKey)`   |
| `create*()`       | Mutation — insert               | `createEmployee(detail)`                 |
| `update*()`       | Mutation — update               | `updateEmployee(detail)`                 |
| `deactivate*()`   | Mutation — logical delete       | `deactivateEmployee(key)`                |
| `setActive()`     | Mutation — toggle active state  | `setActive(key, active)`                 |

`listAll*` returns a full `List<T>`. `findByKey` returns the UI model directly and
throws `EntityNotFoundException` when not found — no `Optional` wrapper at the service
boundary.
