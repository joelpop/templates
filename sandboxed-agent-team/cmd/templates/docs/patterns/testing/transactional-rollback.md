# @Transactional Rollback in Integration Tests

When writing Spring Boot integration tests that insert data, annotate the test class with `@Transactional` so each test method starts in a known, clean state without shared seed data.

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
