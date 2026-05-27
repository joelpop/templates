# Repository Tests

When testing Spring Data JPA repositories, use `@DataJpaTest` with
`@Transactional` rollback and H2 in PostgreSQL-compatibility mode so each test
runs in isolation against a production-equivalent schema without spinning up the
full application context.

## @DataJpaTest

Use `@DataJpaTest` for lightweight repository tests. It loads only the JPA layer
(no web layer, no services) and wraps each test in a transaction that rolls back
automatically:

```java
@DataJpaTest
class EmployeeRepositoryTest {

    @Autowired EmployeeRepository repository;
    @Autowired TestEntityManager em;

    @Test
    void findActiveByDepartment_excludesInactiveEmployees() {
        var dept = em.persist(new DepartmentEntity("Engineering"));
        em.persist(new EmployeeEntity("Alice", true, dept));
        em.persist(new EmployeeEntity("Bob", false, dept));  // inactive
        em.flush();

        var results = repository.findActiveByDepartment(dept.getKey());

        assertThat(results).extracting(EmployeeListItemProjection::getName)
            .containsExactly("Alice");
    }
}
```

## @Transactional Rollback

Integration tests manage their own test data via `@Transactional` rollback
rather than shared seed data:

```java
@SpringBootTest
@Transactional  // rolls back after each test method
class EmployeeServiceIntegrationTest {

    @Autowired EmployeeService service;
    @Autowired EmployeeRepository repository;

    @Test
    void create_persistsEmployee() {
        var detail = new EmployeeDetail("Alice Smith", "alice@example.com");
        service.create(detail);

        assertThat(repository.count()).isEqualTo(1);
    }
}
```

No test depends on data created by another test. Each test starts in a known
state.

## H2 in PostgreSQL Compatibility Mode

Tests use H2 in-memory in PostgreSQL compatibility mode so production PostgreSQL
migration scripts also run in tests:

```properties
# application-test.properties
spring.datasource.url=jdbc:h2:mem:testdb;MODE=PostgreSQL;DB_CLOSE_DELAY=-1
spring.datasource.driver-class-name=org.h2.Driver
spring.jpa.database-platform=org.hibernate.dialect.H2Dialect
```

## Detecting N+1 in Tests

Enable Hibernate statistics in tests to catch N+1 query regressions:

```properties
# application-test.properties
spring.jpa.properties.hibernate.generate_statistics=true
logging.level.org.hibernate.stat=DEBUG
```

Pair with `datasource-proxy` to assert exact query counts in
performance-sensitive paths:

```java
@Test
void listAll_issuesExactlyOneQuery() {
    // setup test data ...
    var queryCounter = new QueryCountHolder();

    var results = service.listAll();

    assertThat(queryCounter.getCount()).isEqualTo(1);
}
```
