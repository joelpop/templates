# Acceptance Criterion Tracing

When writing tests, trace every acceptance criterion to at least one passing
test so coverage gaps are visible and "done" means the specified behavior is
verified, not just that lines were executed.

Line coverage measures *whether code ran*; it doesn't measure *whether the code
did what was specified*. Both matter, and they answer different questions.

## One AC, One or More Tests

A requirement looks like this: one `implementation` child plus one checkbox per
AC, with the parent as a roll-up:

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

The Coder marks `implementation` when the requirement's implementation is
committed and ready for testing. The Tester marks each AC when an automated test
that verifies that AC is passing. The Analyst maintains the parent: `[x]` only
when `implementation` and every AC are `[x]`; `[-]` when any child is `[-]` or
`[x]` but not all are `[x]`; `[ ]` when all children are `[ ]`.

Each AC needs a test (or several, when the AC has parameters):

| AC | Test class / method |
|----|---------------------|
| AC1 | `LoginViewTest.forgotPasswordLink_isVisible` |
| AC2 | `LoginViewTest.forgotPasswordLink_opensEmailPrompt` |
| AC3 | `PasswordResetServiceTest.requestReset_sendsEmail_whenEmailIsRegistered` |
| AC4 | `PasswordResetServiceTest.consumeResetToken_failsAfter30Minutes` |
| AC5 | `LoginViewTest.expiredResetLink_displaysFriendlyError` |

When an AC parameterizes (e.g., "valid email formats are accepted"), parameterize
the test accordingly.

## Test Names Should Identify What They Verify

A reader scanning the test class should be able to match a test name to an AC
without reading the body. The `subject_verb_condition` pattern reads as a
sentence and matches AC phrasing. When the requirement evolves, the test name
evolves with it.

## Cross-Reference, But Don't Depend on It

Some teams put the AC ID in the test name (e.g., `AC3_requestReset_sendsEmail`)
or in a Javadoc comment. That makes traceability searchable but couples test
names to external identifiers that rename over time.

Pragmatic stance: name tests by *behavior* (stable); maintain the AC ↔ test
mapping in the requirement document if the project benefits from explicit
traceability. The mapping table lives alongside the requirement statement; the
Unit Tester updates it during the per-commit cycle.

## Coverage Gaps Surface ACs Without Tests

At the pre-PR gate, the Unit Tester asks "which ACs have a passing test?" — not
"what's the line-coverage percentage?". A line of code without a test for the AC
it implements is a hole; coverage of that line under an unrelated test does not
fill it.

If an AC has no test, the Unit Tester reports the gap to the Coder and
Architect. Closing it may mean (a) writing the missing test, (b) recognizing the
AC was never implemented and implementing it, or (c) recognizing the AC was
misstated and revising it. All three are legitimate; shipping without one of them
is not.

## Tests Are Documentation of Behavior

A well-written test class doubles as a behavioral specification. A new
contributor reading `PasswordResetServiceTest` should be able to infer the reset
flow from test names alone, without reading the production code. If they can't —
if names are mechanical (`testEmptyInput`, `testEdgeCase2`) or generic
(`shouldWork`, `verifyBehavior`) — the tests fail as documentation.
