# Unit Test Structure

When writing unit tests for service and business-logic classes, use `@ExtendWith(MockitoExtension.class)` with `@Mock` / `@InjectMocks` and AssertJ assertions so dependencies are isolated and test failures produce clear, readable messages.

```java
@ExtendWith(MockitoExtension.class)
class EmployeeServiceTest {

    @Mock EmployeeRepository repository;
    @Mock EmployeeMapper mapper;

    @InjectMocks JpaEmployeeService service;

    @Test
    void findByKey_returnsDetail_whenEntityExists() {
        var projection = mock(EmployeeDetailProjection.class);
        var expected = new EmployeeDetail();
        when(repository.findProjectedById(42L)).thenReturn(Optional.of(projection));
        when(mapper.toDetail(projection)).thenReturn(expected);

        var result = service.findByKey(42L);

        assertThat(result).isSameAs(expected);
    }

    @Test
    void findByKey_throwsEntityNotFoundException_whenNotFound() {
        when(repository.findProjectedById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.findByKey(99L))
            .isInstanceOf(EntityNotFoundException.class);
    }
}
```

Use AssertJ (`assertThat(...)`) over JUnit's `assertEquals`. Target ≥ 80% line coverage for the service layer.
