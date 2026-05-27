# Page Object Pattern

When writing UI tests, hide component traversal behind a Page Object class whose
public surface mirrors the user-visible contract so tests survive layout changes
without modification.

Tests that look up components by deep tree traversal — `view.getChildren().flatMap(...)`
chains in browserless tests, `page.locator('vaadin-text-field >> nth=2')` chains in
E2E tests — break when the layout is rearranged, even if no user-visible behavior
changed. The fix is the [Page Object Pattern][fowler-pageobject]: when the layout
moves, only the page object changes; tests stay green because the contract didn't.

[fowler-pageobject]: https://martinfowler.com/bliki/PageObject.html

## Browserless Page Objects — `ComponentTester<T>`

`ComponentTester<T>` (from `browserless-test-junit6`) is the page object base class
for browserless tests. Three-section structure:

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
  `ButtonTester`). The `get*Tester()` naming communicates the return type.
- **User input via tester API** — `getNameFieldTester().setValue(name)`, not
  `nameField.setValue(name)`. The tester checks usability before delegating.
- **Method chaining** — actions that add content return the newly added item's tester.
- **Separation of concerns** — visibility of a child is a concern of the parent container:
  `isCardVisible(card)` lives on the view tester, not the card tester.

**Slot children** (e.g., Vaadin `Card` header slot) are not reachable via `find()`. Navigate
through the component's Java API instead.

The test class works entirely through the tester's public API:

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

To trigger a keyboard shortcut without a browser:

```java
view.setName("Alice");
fireShortcut(Key.ENTER);
assertEquals(1, view.getCardCount());
```

## Playwright Page Objects

The same pattern in Playwright: the page object encapsulates DOM selectors;
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

When an action navigates to a new view, the page object should return the next
page's page object so tests can chain.

## TestBench Page Objects

For TestBench E2E tests (`BrowserTestBase`), the page object extends
`TestBenchElement` and carries an `@Element("vaadin-tag-name")` annotation that
identifies its root element. The same three-section structure and typed-accessor
conventions apply:

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
- **Focus** — check `hasAttribute("focused")` on the Vaadin field element.
- **Text selection** — use `executeScript()` against the inner `inputElement`.
- **Scroll visibility** — compare `getRect()` bounds with ±1 px tolerance.
- **Slot children** — use the slot-aware accessors (`getHeader()`, `getContents()`) rather
  than `$()`, which doesn't cross slot boundaries.

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

See `docs/patterns/testing/recipes/testbench-e2e-server.md` and
`docs/patterns/testing/recipes/testbench-e2e-parallel.md` for Maven `it` profile,
`ServerExtension`, and `TestBenchParallelLimiter` setup.
