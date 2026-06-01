# Browserless UI Tests

When writing fast, in-process UI tests, use `BrowserlessTest` or
`SpringBrowserlessTest` so form submission, validation errors, and grid
interactions are exercised without starting a real browser.

## Choosing the Base Class

| Base class | When to use |
|---|---|
| `BrowserlessTest` | View does not need injected Spring beans |
| `SpringBrowserlessTest` | Test must also assert against an `@Autowired` bean |

Default to `BrowserlessTest`; switch to `SpringBrowserlessTest` only when the
test must also assert against an injected service:

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

Browserless UI tests live in the same package as the view they test with the
`*Test.java` suffix (surefire, not failsafe).

