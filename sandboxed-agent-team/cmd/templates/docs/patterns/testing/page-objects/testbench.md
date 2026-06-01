# TestBench Page Objects

When writing TestBench E2E tests, extend `TestBenchElement` with an `@Element` annotation for the page object so tests survive layout changes without modification.

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
