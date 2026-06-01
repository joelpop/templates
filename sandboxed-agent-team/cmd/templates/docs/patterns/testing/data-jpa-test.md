# @DataJpaTest for Repository Tests

When testing Spring Data JPA repositories, use `@DataJpaTest` so only the JPA layer is loaded and each test runs in an automatically rolled-back transaction without the full application context.

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
