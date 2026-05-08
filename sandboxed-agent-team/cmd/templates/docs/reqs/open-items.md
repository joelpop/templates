# Open Items — Human Input Required

The single tracker for all questions requiring human input before
implementation can proceed. Each item has an ID, the question, which
documents it affects, and its current status.

When a question is resolved:

1. Update the **Status** to `Resolved` (with the date if helpful).
2. Record the answer in the **Notes** column.
3. Update the affected requirement documents accordingly.
4. The entry stays in this table — *do not delete resolved entries*.
   They are useful history when the same kind of question comes up
   again.

The Analyst owns this tracker. The Lead consults it when classifying
a human request to see whether the request resolves an open item.
Any teammate that surfaces a gap during implementation can request a
new entry through the Lead.

---

| ID | Question | Affects | Status | Notes |
|----|----------|---------|--------|-------|
| OI-001 | *(example — replace as real questions arise)* | — | — | — |

<!--
Example resolved entry:

| OI-002 | Should password reset emails include the user's full name in the greeting, or be neutral ("Hi there")? | login.md, account.md, email-templates.md | Resolved (2026-XX-XX) | Use the user's first name only ("Hi Alice"). Last name is omitted because the email is sent to the user's own address — full-name greeting feels formal/transactional. Falls back to "Hi" with no name if no first name on file. Updated `email-templates.md` accordingly. |
-->
