# Browserless Page Objects

When writing browserless UI tests, extend `ComponentTester<T>` for the page object so component traversal is hidden behind a stable public API that mirrors the user-visible contract.

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
