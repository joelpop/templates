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

