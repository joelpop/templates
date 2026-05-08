# Testing Patterns

Unit, browserless UI, E2E, and test data patterns for Vaadin 24+ with Spring Boot 3+
projects. `@DataJpaTest`, `@SpringBootTest`, `@Transactional` rollback, AssertJ, Mockito,
and Playwright usage shown below are stable across every supported Vaadin and Spring Boot
line. Vaadin's browserless UI test API (`SpringBrowserlessTest` / `$(...).id(...)` etc.)
is also consistent across supported versions.

## Testing Pyramid

```
         /\
        /E2E\          Playwright — browser-based, pre-PR gate only
       /------\
      /Browser-\       Vaadin browserless UI — in-process, per-commit
     / less UI  \
    /------------\
   /  Unit Tests  \    JUnit + Mockito — per-commit
  /--------------\
```

- **Unit tests** run on every commit via Maven surefire (`*Test.java` suffix)
- **Browserless UI tests** run on every commit via Maven surefire (`*Test.java` suffix)
- **E2E tests** run only at the pre-PR gate via Maven failsafe (`*IT.java` suffix or
  Playwright via Node.js test runner in `e2e/`)

## One Test Class Per Production Class

Every production class has a corresponding test class. The test class name is the
production class name with a `Test` suffix:

```
EmployeeService       →  EmployeeServiceTest
EmployeeView          →  EmployeeViewTest
EmployeeMapper        →  EmployeeMapperTest
```

No test class covers multiple unrelated production classes.

## Unit Tests — JUnit + Mockito

Write unit tests for every non-UI public method. Service methods, utility methods, and
business logic methods each have at least one corresponding test.

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

Use AssertJ for assertions — prefer `assertThat(...)` over JUnit's `assertEquals`.

## Browserless UI Tests — Vaadin TestBench

Every UI feature has a browserless UI test using Vaadin's browserless testing
(`SpringBrowserlessTest`). These tests exercise form submission, validation errors,
and grid interactions in-process without a real browser.

```java
@SpringBootTest
class EmployeeViewTest extends SpringBrowserlessTest {

    @Autowired EmployeeService employeeService;

    @Test
    void saveButton_savesEmployee_whenFormIsValid() {
        open(EmployeeView.class);

        var nameField = $(TextFieldElement.class).id("name-field");
        nameField.setValue("Jane Smith");

        var saveButton = $(ButtonElement.class).caption("Save");
        saveButton.click();

        assertThat(employeeService.listAll())
            .extracting(EmployeeListItem::getName)
            .contains("Jane Smith");
    }

    @Test
    void saveButton_showsValidationError_whenNameIsEmpty() {
        open(EmployeeView.class);

        $(ButtonElement.class).caption("Save").click();

        assertThat($(TextFieldElement.class).id("name-field").getErrorMessage())
            .isEqualTo("Name is required");
    }
}
```

Browserless UI tests live in the same package as the view they test and use the `*Test.java`
suffix (run by surefire, not failsafe).

## Page Object Pattern

Tests that look up components by deep tree traversal — `view.getChildren().flatMap(...)`
chains in browserless tests, or `page.locator('vaadin-text-field >> nth=2')` chains in
E2E tests — break whenever the layout is rearranged, even if no user-visible behavior
changed. The fix is the [Page Object Pattern][fowler-pageobject]: hide the traversal
behind a class whose public surface is the user-visible contract (heading text, button
labels, named inputs).

When the layout moves, only the page object changes; tests stay green because the
contract didn't.

[fowler-pageobject]: https://martinfowler.com/bliki/PageObject.html

### Page Objects in Browserless UI Tests

The page object walks the server-side component tree starting from the view root, and
exposes high-level lookups (`headingText()`, `buttonWithText("Save")`,
`paragraphMatching(predicate)`). Tests then assert against the user-visible contract,
not the layout's nesting depth.

```java
final class EmployeeFormPageObject {

    private final Component root;

    EmployeeFormPageObject(Component root) {
        this.root = root;
    }

    String headingText() {
        return descendants(H2.class).findFirst()
                .map(H2::getText)
                .orElseThrow(() -> new AssertionError("No heading found"));
    }

    Optional<Button> buttonWithText(String label) {
        return descendants(Button.class)
                .filter(b -> label.equals(b.getText()))
                .findFirst();
    }

    private <T extends Component> Stream<T> descendants(Class<T> type) {
        return walk(root).filter(type::isInstance).map(type::cast);
    }

    private static Stream<Component> walk(Component component) {
        return Stream.concat(
                Stream.of(component),
                component.getChildren().flatMap(EmployeeFormPageObject::walk));
    }
}
```

The test stays focused on behavior:

```java
@Test
void saveButton_isVisible_afterFormLoads() {
    var view = new EmployeeView();
    var page = new EmployeeFormPageObject(view);

    assertThat(page.headingText()).isEqualTo("New Employee");
    assertThat(page.buttonWithText("Save")).isPresent();
}
```

When `SpringBrowserlessTest` is used, the framework's `$()` and `$view()` queries
already provide much of this abstraction (filter by component type, label, attribute,
predicate). The page object is most valuable when those queries aren't enough — for
example, when several components share a tag and only their position in the tree
distinguishes them, or when a single semantic concept ("the details paragraph")
needs a single named accessor across many tests.

Page objects live in the test source tree, package-private, in the same package as
the views they cover. A working in-tree example is `ErrorViewPageObject` in
`fleet-acuity-ui/src/test/java/.../ui/view/error/`, used by all four error-view tests.

### Page Objects in E2E Tests

The same pattern applies in Playwright tests: the page object encapsulates DOM
selectors, and the test asserts against user-visible behavior.

```typescript
// e2e/page-objects/EmployeeFormPage.ts
import { Page, Locator } from '@playwright/test';

export class EmployeeFormPage {
    constructor(private readonly page: Page) {}

    readonly heading = (): Locator => this.page.getByRole('heading');
    readonly nameField = (): Locator => this.page.getByLabel('Name');
    readonly emailField = (): Locator => this.page.getByLabel('Email');
    readonly saveButton = (): Locator => this.page.getByRole('button', { name: 'Save' });

    async submit(name: string, email: string): Promise<void> {
        await this.nameField().fill(name);
        await this.emailField().fill(email);
        await this.saveButton().click();
    }
}
```

```typescript
// e2e/employee.spec.ts
test('creates an employee successfully', async ({ page }) => {
    await page.goto('/employees/new');
    const form = new EmployeeFormPage(page);

    await form.submit('Alice Smith', 'alice@example.com');

    await expect(page.getByText('Alice Smith')).toBeVisible();
});
```

When an action navigates to a new view, the page object should return the next page's
page object so tests can chain calls. This also implicitly asserts that the navigation
succeeded — if the next page isn't found, the test fails at the chain point.

Vaadin's TestBench documentation covers the same pattern for Java E2E tests using
`@Element("tag-name")` annotated classes that extend `TestBenchElement` — see
[Tests with Page Objects][vaadin-page-objects]. The shape it describes is the same
abstraction this section recommends; the locator mechanism differs (TestBench's
`$(LoginViewElement.class)` vs. Playwright's `page.getByLabel(...)`), but the goal
is identical: tests stay coupled to the user-visible contract, not the DOM.

[vaadin-page-objects]: https://vaadin.com/docs/latest/flow/testing/end-to-end/page-objects

## Repository Tests — @DataJpaTest

Use `@DataJpaTest` for lightweight repository tests. It loads only the JPA layer and wraps
each test in a transaction that rolls back automatically:

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

## Test Data Management

### @Transactional Rollback

Integration tests manage their own test data using `@Transactional` rollback rather than
shared seed data:

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

No test depends on data created by another test. Each test starts in a known state.

### H2 in PostgreSQL Compatibility Mode

Test configuration uses H2 in-memory database in PostgreSQL compatibility mode so that
migration scripts written for PostgreSQL are compatible with the test environment:

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

Pair with `datasource-proxy` to assert exact query counts in performance-sensitive paths:

```java
@Test
void listAll_issuesExactlyOneQuery() {
    // setup test data ...
    var queryCounter = new QueryCountHolder();

    var results = service.listAll();

    assertThat(queryCounter.getCount()).isEqualTo(1);
}
```

## E2E Tests — Playwright (TypeScript)

End-to-end tests are written in TypeScript using `@playwright/test` and live in the `e2e/`
directory. They run against the full application stack (started via Docker or Maven failsafe).

```typescript
// e2e/employee.spec.ts
import { test, expect } from '@playwright/test';

test('creates an employee successfully', async ({ page }) => {
    await page.goto('/employees');
    await page.getByRole('button', { name: 'Add Employee' }).click();

    await page.getByLabel('Name').fill('Alice Smith');
    await page.getByLabel('Email').fill('alice@example.com');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Alice Smith')).toBeVisible();
});
```

E2E tests run only at the pre-PR gate — not on every commit. They exercise user-visible
flows end-to-end and catch integration failures that browserless tests cannot detect.

## Coverage Targets

| Layer | Target |
|-------|--------|
| Service layer | ≥ 80% line coverage |
| Utility classes | ≥ 80% line coverage |
| UI views | All form interactions, validation errors, grid interactions covered by browserless tests |

Coverage is measured per module. The UI module is covered by browserless tests, not line
coverage tools (which cannot easily instrument Vaadin component interactions).

## Tests Trace to Acceptance Criteria

Line coverage measures *whether code ran*. It does not measure
*whether the code did what was specified*. Both matter, but they
answer different questions.

Every requirement's acceptance criteria (ACs) must have at least one
automated test that exercises and verifies that AC. If an AC has no
test, the requirement is not implemented — regardless of what line
coverage reports say.

### One AC, one or more tests

A requirement looks like this — one `implementation` child plus one
checkbox per AC, with the parent as a roll-up:

```markdown
- [ ] User can reset password via email link
      ... additional detail and description ...
  - [ ] implementation
  - [ ] AC1: A "Forgot password" link is visible on the login view.
  - [ ] AC2: Clicking the link prompts for an email address.
  - [ ] AC3: A valid registered email triggers a reset email within 60 seconds.
  - [ ] AC4: The reset link expires 30 minutes after issue.
  - [ ] AC5: An expired reset link displays a clear, non-technical error.
```

The Coder marks `implementation` when the requirement's
implementation is committed and ready for testing. The Tester marks
each AC when an automated test that verifies that AC is passing.
The Analyst maintains the parent: `[x]` only when `implementation`
and every AC are `[x]`; `[-]` when any child is `[-]` or `[x]` but
not all are `[x]`; `[ ]` when all children are `[ ]`.

Each AC needs a test (or several, when the AC has parameters):

| AC | Test class / method |
|----|---------------------|
| AC1 | `LoginViewTest.forgotPasswordLink_isVisible` |
| AC2 | `LoginViewTest.forgotPasswordLink_opensEmailPrompt` |
| AC3 | `PasswordResetServiceTest.requestReset_sendsEmail_whenEmailIsRegistered` |
| AC4 | `PasswordResetServiceTest.consumeResetToken_failsAfter30Minutes` |
| AC5 | `LoginViewTest.expiredResetLink_displaysFriendlyError` |

When an AC parameterizes (e.g., "valid email formats are accepted"),
the test is parameterized accordingly.

### Test names should identify what they verify

Test method names describe the behavior under test, not the
implementation. A reader scanning the test class should be able to
match a test name to an AC without reading the test body:

```java
// Preferred — name describes the behavior
@Test void requestReset_sendsEmail_whenEmailIsRegistered() { ... }
@Test void consumeResetToken_failsAfter30Minutes() { ... }

// Avoid — name describes the test mechanics
@Test void testCase1() { ... }
@Test void emailServiceMockTest() { ... }
```

This naming pattern (`subject_verb_condition`) reads as a sentence
and matches naturally to AC phrasing. When the requirement evolves,
the test name evolves with it.

### Cross-reference, but don't depend on it

Some teams put the AC ID in the test name (e.g.,
`AC3_requestReset_sendsEmail`) or in a Javadoc comment on the test
method. This makes traceability searchable but couples test names to
external identifiers that rename over time.

The pragmatic stance: name tests by *behavior* (which is stable);
maintain the AC ↔ test mapping in the requirement document if the
project benefits from explicit traceability. The mapping table can
live alongside the requirement statement and is updated by the Unit
Tester at the per-commit cycle.

### Coverage gaps surface ACs without tests

When the Unit Tester checks AC coverage at the pre-PR gate, the
question is "which ACs have a passing test?" — not "what's the line
coverage percentage?". A line of code without a test for the AC it
implements is a hole; line coverage of that line under an unrelated
test does not fill the hole.

If an AC has no test, the Unit Tester reports the gap to the Coder
and Architect. Closing the gap may mean (a) writing the missing test,
(b) recognizing the AC was never implemented and adding the
implementation, or (c) recognizing the AC was misstated and revising
it (back through the requirement gate). All three are legitimate
outcomes; shipping without one of them is not.

### Tests are documentation of behavior

A well-written test class doubles as a behavioral specification. A
new contributor reading `PasswordResetServiceTest` should be able to
infer the password-reset flow from the test names alone, without
reading the production code. If they can't — if the test names are
mechanical (`testEmptyInput`, `testEdgeCase2`) or generic
(`shouldWork`, `verifyBehavior`) — the tests are failing as
documentation, not just failing as discipline.
