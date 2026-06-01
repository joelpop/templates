# Coverage Targets

When measuring test coverage, target ≥ 80% line coverage for service and utility layers; cover UI views through browserless tests rather than line-coverage tools.

| Layer | Target |
|-------|--------|
| Service layer | ≥ 80% line coverage |
| Utility classes | ≥ 80% line coverage |
| UI views | All form interactions, validation errors, grid interactions covered by browserless tests |

Coverage is measured per module. The UI module is covered by browserless tests,
not line-coverage tools (which cannot easily instrument Vaadin component
interactions).
