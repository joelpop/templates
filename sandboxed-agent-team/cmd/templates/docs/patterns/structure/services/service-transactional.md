# @Transactional on Service Implementations

When annotating service implementations, put `@Transactional(readOnly = true)` on all
query methods and `@Transactional` on all mutation methods — annotations go on the
**implementation**, not the interface.

```java
@Service
public class JpaEmployeeService implements EmployeeService {

    @Override
    @Transactional(readOnly = true)
    public EmployeeDetail findByKey(long key) { /* ... */ }

    @Override
    @Transactional(readOnly = true)
    public List<EmployeeListItem> listAll() { /* ... */ }

    @Override
    @Transactional
    public EmployeeDetail create(EmployeeDetail detail) { /* ... */ }

    @Override
    @Transactional
    public EmployeeDetail update(EmployeeDetail detail) { /* ... */ }

    @Override
    @Transactional
    public void deactivate(long key) { /* ... */ }
}
```

`readOnly = true` tells Hibernate to skip dirty checking on flush — this reduces overhead
on every query-only service method.

Keep `@Transactional` at the service layer, not the repository layer. Business logic
boundaries are the service's responsibility.
