# Unit Tests

When writing unit tests for service, utility, and business-logic classes, use
JUnit 5 with Mockito and AssertJ assertions named in `subject_verb_condition`
form so tests read as behavioral specifications and map directly to acceptance
criteria.

## Structure

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

Target code coverage for the service layer: **at least 80%**.

Use AssertJ for assertions — prefer `assertThat(...)` over JUnit's
`assertEquals`.

## Test Naming

Test method names describe the behavior, not the implementation. A reader
scanning the test class should be able to match a test name to an acceptance
criterion without reading the body:

```java
// Preferred — name describes the behavior
@Test void requestReset_sendsEmail_whenEmailIsRegistered() { ... }
@Test void consumeResetToken_failsAfter30Minutes() { ... }

// Avoid — name describes the test mechanics
@Test void testCase1() { ... }
@Test void emailServiceMockTest() { ... }
```

The `subject_verb_condition` pattern reads as a sentence and matches acceptance
criterion phrasing. When the requirement evolves, the test name evolves with it.
