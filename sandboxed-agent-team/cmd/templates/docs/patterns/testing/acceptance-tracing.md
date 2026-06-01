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

