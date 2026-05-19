# Testing Patterns

Unit, browserless UI, E2E, and test-data patterns for Vaadin 24+ with
Spring Boot 3+. `@DataJpaTest`, `@SpringBootTest`, `@Transactional`
rollback, AssertJ, and Mockito usage below are stable across supported
Vaadin and Spring Boot lines. Vaadin's browserless UI API
(`BrowserlessTest` / `ComponentTester<T>`) and TestBench E2E API are
also consistent across supported versions.

## Testing Pyramid

```
          /\
         /  \
        /E2E \          TestBench / Playwright — browser-based, pre-PR gate only
       /------\
      /Browser-\       Vaadin browserless UI — in-process, per-commit
     / less UI  \
    /------------\
   /  Unit Tests  \    JUnit + Mockito — per-commit
  /----------------\
```

- **Unit tests** run on every commit via Maven surefire (`*Test.java` suffix)
- **Browserless UI tests** run on every commit via Maven surefire (`*Test.java` suffix)
- **E2E tests** run only at the pre-PR gate via Maven failsafe (`*IT.java` suffix with
  TestBench) or Node.js test runner in `e2e/` (Playwright)

## One Test Class Per Production Class

Every production class has a corresponding test class named with the
production class name plus a `Test` suffix:

```
EmployeeService       →  EmployeeServiceTest
EmployeeView          →  EmployeeViewTest
EmployeeMapper        →  EmployeeMapperTest
```

No test class covers multiple unrelated production classes.

## Unit Tests — JUnit + Mockito

Every non-UI public method has at least one unit test — service,
utility, and business-logic methods alike.

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

## Browserless UI Tests

Every UI feature has a browserless UI test using
`SpringBrowserlessTest`. These exercise form submission, validation
errors, and grid interactions in-process — no real browser.

Choose the base class based on whether the test needs Spring beans:

- **`BrowserlessTest`** — the default. No application context is started; `navigate(View.class)` returns a typed view instance. Use when the view doesn't need injected Spring beans.
- **`SpringBrowserlessTest`** — when the test needs `@Autowired` beans alongside the view. Annotate with `@SpringBootTest`.

Default to `BrowserlessTest`; switch to `SpringBrowserlessTest` only when the test must also assert against an injected service. The example below uses `SpringBrowserlessTest` to assert against `EmployeeService`:

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

Browserless UI tests live in the same package as the view they test
with the `*Test.java` suffix (surefire, not failsafe).

## Page Object Pattern

Tests that look up components by deep tree traversal —
`view.getChildren().flatMap(...)` chains in browserless tests,
`page.locator('vaadin-text-field >> nth=2')` chains in E2E tests —
break when the layout is rearranged, even if no user-visible behavior
changed. The fix is the [Page Object Pattern][fowler-pageobject]:
hide the traversal behind a class whose public surface is the
user-visible contract (heading text, button labels, named inputs).
When the layout moves, only the page object changes; tests stay green
because the contract didn't.

[fowler-pageobject]: https://martinfowler.com/bliki/PageObject.html

### Page Objects in Browserless UI Tests

`ComponentTester<T>` (from `browserless-test-junit6`) is the page object base class for
browserless tests. It wraps the server-side component tree and exposes `find(Type.class)`
queries so tests interact through stable, intention-revealing methods rather than raw tree
traversal.

**Three-section structure:**

```java
public class EmployeeViewTester extends ComponentTester<EmployeeView> {

    // PUBLIC API

    public EmployeeViewTester(EmployeeView component) {
        super(component);
    }

    public void setName(String name) {
        getNameFieldTester().setValue(name);
    }

    public void save() {
        getSaveButtonTester().click();
    }

    public String getValidationError() {
        return getNameFieldTester().getErrorMessage();
    }


    // INTERNAL component tester accessors

    private TextFieldTester<TextField, String> getNameFieldTester() {
        return new TextFieldTester<>(find(TextField.class).withCaption("Name").single());
    }

    private ButtonTester<Button> getSaveButtonTester() {
        return new ButtonTester<>(find(Button.class).withText("Save").single());
    }


    // INTERNAL helpers
    // ...
}
```

Key conventions:

- **Stable location** — locate by user-visible identifier: `find(TextField.class).withCaption("Name").single()`,
  not `find(TextField.class).single()`, which breaks when a second field is added.
- **Typed accessor naming** — private accessors return typed tester instances (`TextFieldTester`,
  `ButtonTester`, `SpanTester`, `DivTester`). The `get*Tester()` naming communicates the return
  type and prevents confusion.
- **User input via tester API** — `getNameFieldTester().setValue(name)`, not
  `nameField.setValue(name)`. The tester checks usability (visibility, enabled state) before
  delegating to the component.
- **Method chaining** — actions that add content return the newly added item's tester:
  `GreetingCardTester card = view.greet("Alice")` → `assertEquals("Hello, Alice!", card.getMessage())`.
- **Separation of concerns** — visibility of a child is a concern of the parent container:
  `isCardVisible(card)` lives on the view tester, not the card tester.

**Slot children** (e.g., Vaadin `Card` header slot) are not reachable via `find()`. Navigate
through the component's Java API instead: `getComponent().getContent().getHeader().getChildren()`.

The test class works entirely through the tester's public API — no component lookups in the body:

```java
class EmployeeViewTest extends BrowserlessTest {

    private EmployeeViewTester view;

    @BeforeEach
    void open() {
        view = new EmployeeViewTester(navigate(EmployeeView.class));
    }

    @Test
    void saveButton_savesEmployee_whenFormIsValid() {
        view.setName("Jane Smith");
        view.save();
        assertEquals(1, view.getCardCount());
    }

    @Test
    void saveButton_showsValidationError_whenNameIsEmpty() {
        view.save();
        assertEquals("Name is required", view.getValidationError());
    }
}
```

Page object classes live in the test source tree, in the same package as the view they cover.

To trigger a keyboard shortcut without a browser, call `fireShortcut()` from the test class:

```java
view.setName("Alice");
fireShortcut(Key.ENTER);
assertEquals(1, view.getCardCount());
```

### Page Objects in E2E Tests

Same pattern in Playwright: the page object encapsulates DOM selectors;
the test asserts user-visible behavior.

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

When an action navigates to a new view, the page object should return
the next page's page object so tests can chain. This also implicitly
asserts the navigation succeeded — if the next page isn't found, the
test fails at the chain point.

#### TestBench page objects

For TestBench E2E tests (`BrowserTestBase`), the page object extends `TestBenchElement` and
carries an `@Element("vaadin-tag-name")` annotation that identifies its root element. The
same three-section structure and typed-accessor conventions apply:

```java
@Element("vaadin-vertical-layout")
public class EmployeeViewElement extends TestBenchElement {

    // PUBLIC API

    public void setName(String name) {
        getNameFieldElement().click();
        getNameFieldElement().sendKeys(name);
    }

    public void save() {
        getSaveButtonElement().click();
    }

    public int getCardCount() {
        return $(EmployeeCardElement.class).all().size();
    }

    public boolean isNameFieldFocused() {
        return getNameFieldElement().hasAttribute("focused");
    }

    public boolean isNameFieldTextSelected() {
        return (Boolean) executeScript(
                "const el = arguments[0].inputElement;" +
                "return el.value.length > 0 && el.selectionStart === 0 && el.selectionEnd === el.value.length;",
                getNameFieldElement());
    }

    public boolean isCardVisible(EmployeeCardElement card) {
        var scrollerRect = $(ScrollerElement.class).single().getRect();
        var cardRect = card.getRect();
        return cardRect.getY() + cardRect.getHeight() <= scrollerRect.getY() + scrollerRect.getHeight() + 1
            && cardRect.getY() >= scrollerRect.getY() - 1;
    }


    // INTERNAL element accessors

    private TextFieldElement getNameFieldElement() {
        return $(TextFieldElement.class).withCaption("Name").single();
    }

    private ButtonElement getSaveButtonElement() {
        return $(ButtonElement.class).withText("Save").single();
    }
}
```

Differences from browserless page objects:

- **User input** — use `click() + sendKeys(name)`, not `TextFieldElement.setValue(name)`.
  `setValue()` leaves the prior text selected, masking whether autoselect behavior is present.
- **Root access** — retrieve the view element from the test with
  `$(EmployeeViewElement.class).onPage().get(0)`.
- **Focus** — check `hasAttribute("focused")` on the Vaadin field element (the web
  component, not the inner `<input>`).
- **Text selection** — use `executeScript()` against the inner `inputElement`;
  `selectionStart === 0 && selectionEnd === value.length` means all text is selected.
- **Scroll visibility** — compare `getRect()` bounds of the scroller and card with ±1 px
  tolerance to absorb sub-pixel rounding.
- **Slot children** (e.g., `CardElement` header/content slots) — use the slot-aware
  accessors (`getHeader()`, `getContents()`) rather than `$()`, which doesn't cross slot
  boundaries.

**Test class:**

```java
@ExtendWith(ServerExtension.class)
class EmployeeViewIT extends BrowserTestBase {

    private EmployeeViewElement view;

    @BeforeEach
    void open() {
        getDriver().get("http://localhost:" + System.getProperty("deployment.port") + "/");
        view = $(EmployeeViewElement.class).onPage().get(0);
    }

    @BrowserTest
    void saveButton_savesEmployee_whenFormIsValid() {
        view.setName("Jane Smith");
        view.save();
        assertEquals(1, view.getCardCount());
    }

    @BrowserTest
    void openView_nameFieldIsFocused() {
        waitUntil(_ -> view.isNameFieldFocused());
        assertTrue(view.isNameFieldFocused());
    }
}
```

Use `@BrowserTest` (not `@Test`) to opt into TestBench's parallel executor. Use
`waitUntil()` for any assertion on state that may not be immediate after a server
round-trip.

See `testbench-e2e-server` and `testbench-e2e-parallel` recipes for Maven `it`
profile, `ServerExtension`, and `TestBenchParallelLimiter` setup.

[vaadin-page-objects]: https://vaadin.com/docs/latest/flow/testing/end-to-end/page-objects

## What Only a Real Browser Can Test

Browserless tests cover server-side component state and event handling without a browser.
A real browser (TestBench) is required for anything the browser engine mediates:

| Capability | Browserless | TestBench |
|---|---|---|
| Component state, server events | ✓ | ✓ |
| Focus | — | ✓ (`hasAttribute("focused")`) |
| Text selection | — | ✓ (`executeScript()` on `selectionStart`/`selectionEnd`) |
| Scroll position | — | ✓ (`getRect()` comparison) |
| CSS rendering, hover states | — | ✓ |
| Web component internals (slots, shadow DOM) | — | ✓ |

For browser-only assertions, use `waitUntil()` to let the browser settle before asserting:

```java
waitUntil(_ -> view.isNameFieldFocused());
assertTrue(view.isNameFieldFocused());
```

## Repository Tests — @DataJpaTest

Use `@DataJpaTest` for lightweight repository tests. It loads only the
JPA layer and wraps each test in an auto-rollback transaction:

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

Integration tests manage their own test data via `@Transactional`
rollback rather than shared seed data:

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

Tests use H2 in-memory in PostgreSQL compatibility mode so production
PostgreSQL migration scripts also run in tests:

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

E2E tests use `@playwright/test` (TypeScript) and live in `e2e/`. They
run against the full application stack (started via Docker or Maven
failsafe). For Java-based E2E tests, see the TestBench page objects
section above and the `testbench-e2e-server` / `testbench-e2e-parallel`
recipes.

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

E2E tests run only at the pre-PR gate, not per-commit. They exercise
user-visible flows end-to-end and catch integration failures that
browserless tests can't detect.

## Coverage Targets

| Layer | Target |
|-------|--------|
| Service layer | ≥ 80% line coverage |
| Utility classes | ≥ 80% line coverage |
| UI views | All form interactions, validation errors, grid interactions covered by browserless tests |

Coverage is measured per module. The UI module is covered by
browserless tests, not line-coverage tools (which can't easily
instrument Vaadin component interactions).

## Tests Trace to Acceptance Criteria

Line coverage measures *whether code ran*; it doesn't measure *whether
the code did what was specified*. Both matter, and they answer
different questions.

Every requirement's acceptance criteria (ACs) must have at least one
automated test exercising and verifying it. If an AC has no test, the
requirement isn't implemented — regardless of what line coverage
reports say.

### One AC, one or more tests

A requirement looks like this: one `implementation` child plus one
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
parameterize the test accordingly.

### Test names should identify what they verify

Test method names describe the behavior, not the implementation. A
reader scanning the test class should be able to match a test name
to an AC without reading the body:

```java
// Preferred — name describes the behavior
@Test void requestReset_sendsEmail_whenEmailIsRegistered() { ... }
@Test void consumeResetToken_failsAfter30Minutes() { ... }

// Avoid — name describes the test mechanics
@Test void testCase1() { ... }
@Test void emailServiceMockTest() { ... }
```

The `subject_verb_condition` pattern reads as a sentence and matches
AC phrasing. When the requirement evolves, the test name evolves
with it.

### Cross-reference, but don't depend on it

Some teams put the AC ID in the test name (e.g.,
`AC3_requestReset_sendsEmail`) or in a Javadoc comment. That makes
traceability searchable but couples test names to external
identifiers that rename over time.

Pragmatic stance: name tests by *behavior* (stable); maintain the
AC ↔ test mapping in the requirement document if the project benefits
from explicit traceability. The mapping table lives alongside the
requirement statement; the Unit Tester updates it during the
per-commit cycle.

### Coverage gaps surface ACs without tests

At the pre-PR gate, the Unit Tester asks "which ACs have a passing
test?" — not "what's the line-coverage percentage?". A line of code
without a test for the AC it implements is a hole; coverage of that
line under an unrelated test does not fill it.

If an AC has no test, the Unit Tester reports the gap to the Coder
and Architect. Closing it may mean (a) writing the missing test, (b)
recognizing the AC was never implemented and implementing it, or (c)
recognizing the AC was misstated and revising it (back through the
requirement gate). All three are legitimate; shipping without one of
them is not.

### Tests are documentation of behavior

A well-written test class doubles as a behavioral specification. A
new contributor reading `PasswordResetServiceTest` should be able to
infer the reset flow from test names alone, without reading the
production code. If they can't — if names are mechanical
(`testEmptyInput`, `testEdgeCase2`) or generic (`shouldWork`,
`verifyBehavior`) — the tests fail as documentation, not just as
discipline.
