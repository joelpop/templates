# Disable Open Session in View (OSIV)

When configuring Spring Data JPA, disable OSIV so data loading is
transaction-scoped and lazy-loading bugs surface immediately rather than
silently firing at an unexpected call site outside the transaction.

```properties
spring.jpa.open-in-view=false
```

OSIV keeps the Hibernate session open for the entire HTTP request lifetime.
Disabling it:
- Forces all data loading into the service/transaction layer where it belongs
- Surfaces `LazyInitializationException` immediately in development (these are
  real bugs that OSIV was masking)
- Prevents N+1 queries from silently firing anywhere outside the expected
  transaction boundary
