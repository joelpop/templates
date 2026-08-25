# Test Method Naming

When naming a test method, use the `subject_verb_condition` pattern so a reader scanning the test class can match each test to an acceptance criterion without reading the body — and the name evolves naturally when the requirement changes.

```java
// Avoid — name describes the test mechanics
@Test void testCase1() { /* ... */ }
@Test void emailServiceMockTest() { /* ... */ }
```

```java
// Preferred — name describes the behavior
@Test void requestReset_sendsEmail_whenEmailIsRegistered() { /* ... */ }
@Test void consumeResetToken_failsAfter30Minutes() { /* ... */ }
```

The `subject_verb_condition` pattern reads as a sentence and matches acceptance criterion phrasing. When the requirement evolves, the test name evolves with it.
