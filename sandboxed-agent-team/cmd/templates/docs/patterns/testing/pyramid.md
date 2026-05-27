# Testing Pyramid

When planning test coverage, structure tests across three tiers — unit,
browserless UI, and E2E — so fast in-process tests run on every commit while
slow browser tests are reserved for the pre-PR gate.

## Tiers

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

Every production class has a corresponding test class named with the production
class name plus a `Test` suffix:

```
EmployeeService       →  EmployeeServiceTest
EmployeeView          →  EmployeeViewTest
EmployeeMapper        →  EmployeeMapperTest
```

No test class covers multiple unrelated production classes.

## Coverage Targets

| Layer | Target |
|-------|--------|
| Service layer | ≥ 80% line coverage |
| Utility classes | ≥ 80% line coverage |
| UI views | All form interactions, validation errors, grid interactions covered by browserless tests |

Coverage is measured per module. The UI module is covered by browserless tests,
not line-coverage tools (which cannot easily instrument Vaadin component
interactions).
