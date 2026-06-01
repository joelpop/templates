# Service Interface Contracts

When writing service interfaces, operate exclusively on UI model objects and group
queries before mutations with a comment separator so the contract is decoupled from
persistence and its shape is immediately scannable.

Service interfaces operate exclusively on UI model objects (POJOs from the `{app}-uimodel`
module). No service interface method signature references a JPA entity, interface
projection, repository, or MapStruct mapper. This is enforced at compile time by the
module structure (see `docs/patterns/structure/modules/layer-separation.md`).

```java
public interface EmployeeService {

    // --- Queries ---
    EmployeeDetail findByKey(long key);
    List<EmployeeListItem> listAll();
    boolean isEmailAvailable(String email, long excludeKey);

    // --- Mutations ---
    EmployeeDetail create(EmployeeDetail detail);
    EmployeeDetail update(EmployeeDetail detail);
    void deactivate(long key);
}
```

`create` and `update` both accept the UI model and return the updated UI model (with
server-assigned fields populated — key on create, refreshed audit fields on update).
`deactivate` takes just the key because there is no detail to send.
