# Fix Discipline

How to fix a problem without making the codebase worse. The temptation
under pressure is to make the failing test pass, the error message go
away, the warning disappear — and call it done. That's how a codebase
accumulates the kind of debt that's expensive to unwind later.

## Diagnose before fixing

When a build error, test failure, or unexpected runtime behavior shows
up:

1. **Stop.** Read the full error output. Identify the *root cause*, not
   just the symptom.
2. **Classify the failure** before touching code:
   - **TRIVIAL** — typo, missing import, wrong method name. The fix is
     mechanical. Proceed.
   - **LOCALIZED** — logic error within the current method or class.
     Approach is sound; implementation has a bug. Proceed, but if the
     fix needs changes outside the originating method/class,
     reclassify.
   - **STRUCTURAL** — the error suggests the current approach won't
     work, or the fix requires modifying interfaces, adding parameters,
     changing data flow, or working around a framework constraint. Do
     not patch. Escalate to the Architect; the approach itself is what
     needs to change.

The classification step is the discipline. Rushing past it leads to
patching symptoms, which compounds.

## The fix-attempt limit

If `{{FIX_ATTEMPT_LIMIT}}` consecutive fix attempts target the same
root cause and the problem persists, **stop**. Escalate, regardless
of classification. Reaching the limit on one root cause means the
diagnosis is wrong; another attempt is unlikely to succeed and is
very likely to deepen the mess.

This rule counts *root causes*, not error messages. Patches
addressing the same underlying issue count toward the limit even if
the symptoms (stack traces, error strings) differ. When in doubt,
treat consecutive attempts as targeting the same root cause and
escalate sooner rather than later.

Examples:

- **Same root cause across attempts (counts toward the limit):**
  Attempt 1 — the "admin can log out" test fails; Coder adjusts
  `LogoutHandler` so admin sessions clear correctly. The admin test
  passes, but the previously-passing "regular user can log out" test
  now fails. Attempt 2 — Coder reworks the handler to also handle the
  regular-user case; the regular-user test passes, but the admin test
  breaks again. Each fix is downstream of one root cause: the handler
  conflates two session shapes. Patching one role at a time will keep
  ping-ponging — escalate so the Architect can revise the handler's
  shape (e.g., dispatch by session type, or split into two handlers).
- **Distinct root causes (each counts separately):** attempt 1 fixes
  a parser bug; attempt 2 fixes an unrelated retry-loop bug that the
  parser bug was masking. Independent defects.

## Workaround signatures

These patterns are smells. If you're reaching for one, the
classification is **structural** — escalate, don't apply.

- `@SuppressWarnings`, `noinspection`, `// eslint-disable`, or
  equivalent suppression annotations or comments. They silence the
  tool that was telling you something useful.
- `try`/`catch` blocks that swallow exceptions to make a test pass or
  to keep code "running". The error is the message, not the noise.
- Type casts or `instanceof` checks added to bypass type-system
  errors. The type system is reporting a real conflict; casting
  defers it to runtime where it shows up worse.
- Null checks that mask incorrect data flow. If a value should never
  be null at this point but is, the bug is upstream — the null check
  hides the bug rather than fixing it.
- Copying code rather than fixing the shared abstraction the original
  came from. The first time is "fast"; every later change has to
  remember both copies.

When a workaround signature shows up in a diff, the appropriate review
response is "this is a structural fix, not a localized one — what is
the underlying problem and how do we address it?"

## Right thing vs. working thing

A fix that makes the code *work* is not the same as a fix that makes
the code *right*. The goal is the latter. The temptation is the
former.

Right-thing fixes:

- Address the root cause. The symptom disappears as a consequence.
- Leave the codebase as understandable as it was before, or more so.
- Don't add cleverness. The fix should look obvious in retrospect.
- Don't introduce special cases that break the surrounding model.

Working-thing fixes (avoid):

- Address the symptom. The root cause persists.
- Add scaffolding — flags, branches, fallbacks, retry loops — that
  makes the immediate failure go away but increases the surface area
  of the code.
- "If it works, ship it." Works under what conditions? With what
  assumptions? At what cost to the next change?

When the right thing is too big for the current task, the answer is to
*scope the right thing* — surface it as a separate item, sized
appropriately — not to ship the working thing and call it done.

## Revert before rework

When the Architect responds to a mid-task escalation with an approach
revision:

1. Identify all uncommitted changes that were part of the abandoned
   approach.
2. **Revert them** before starting the revised approach. `git
   checkout` or `git stash`. Do not try to salvage partial work by
   adapting it, unless the Architect explicitly identifies specific
   pieces to keep.
3. The revised approach starts from the last clean commit, not from
   the failed state.

Salvaging partial work from an abandoned approach is one of the
fastest ways to end up with the worst of both designs. The new
approach absorbs constraints from the old one without their context,
and quirks become invariants without anyone noticing.

## Fix-mode comment discipline

A specific failure mode worth its own callout: when fixing code, do
**not** paste the conversational explanation of the fix into a
comment block in the source. The fix's *why* belongs in the commit
message and PR description, not the code. See
[`comments.md`](comments.md) for the full rule.

## Tests are part of the fix

A bug fix without a test that fails before the fix and passes after is
incomplete. The test is what prevents regression and what documents
that the bug existed. Skipping the test because "it's obvious" or
"the manual smoke test passed" forfeits both benefits.

Where the fix exposes a missing acceptance criterion (the spec didn't
require this behavior, but the bug shows it should have), the AC goes
into the requirement first; the test then verifies the AC; the fix
satisfies the test. See the doc-first-fix process for the broader
flow.