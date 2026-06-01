# One Test Class Per Production Class

When creating a test class, name it after the single production class it tests with a `Test` suffix — no test class covers multiple unrelated production classes.

```
EmployeeService       →  EmployeeServiceTest
EmployeeView          →  EmployeeViewTest
EmployeeMapper        →  EmployeeMapperTest
```
