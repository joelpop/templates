# Requirement Writing

How requirement statements are phrased in `docs/reqs/`. The form
matters: a well-formed requirement is testable, unambiguous, and
binding on the implementation. A poorly formed requirement reads
like a description of what users will probably do, which leaves the
implementer guessing.

## System-facing imperative form

A requirement states what **the system** must do. Use the system as
the subject; use a modal verb of obligation (`must`, `shall`,
`will`) as the predicate.

```markdown
Preferred:
- The system must provide a means for users to log in with their
  username and password.

Avoid (this is a capability statement, not a requirement):
- User can log in with username and password.
```

The "User can ..." form *describes* a capability that may exist; it
doesn't *require* the system to provide it. As a result it is hard
to fail a test against. The "system must" form is binding.

## Modal verbs

| Verb | Meaning | Use for |
|------|---------|---------|
| **must** / **shall** | The system is *required* to do this. Failure is a defect. | Normal requirements. |
| **must not** / **shall not** | The system is required *not* to do this. | Prohibitions, security boundaries, safety constraints. |
| **will** | A statement of fact about behavior the system performs. Equivalent to *shall* in many SRS conventions; consistency within a project matters more than which one is chosen. | Same as *must*/*shall*; pick one and use it consistently. |
| ~~should~~ / ~~may~~ / ~~could~~ | **Avoid.** These are recommendations, not requirements. If the behavior is required, use *must*/*shall*. If it is genuinely optional, it is probably an out-of-scope item, not a requirement. |

When a project commits to one verb (`must` or `shall`), use it
everywhere — mixing them creates ambiguity about whether some
requirements are weaker than others.

## One concept per requirement

A requirement covers exactly one binding obligation. If the
statement uses "and" to join two obligations, it is two
requirements.

```markdown
Avoid:
- The system must validate the input and persist it to the database.

Preferred:
- The system must validate user input against the rules in
  `data-validation.md`.
- The system must persist validated input to the user store.
```

Atomicity matters because acceptance criteria attach to
requirements; a compound requirement gets compound ACs that are
harder to trace and easier to partially-implement.

## Testable

Every requirement must be verifiable through observation —
typically a concrete test, sometimes a manual check. The acceptance
criteria under each requirement are what makes the verification
explicit. If you cannot imagine a test that would fail in the
absence of this requirement being satisfied, the statement is too
vague.

```markdown
Avoid (untestable):
- The system must be user-friendly.

Preferred (testable; ACs cover the verification):
- The system must allow first-time users to complete the primary
  workflow without consulting external documentation.
  - implementation
  - AC1: A new user, given only the login URL, can complete the
    "create item" flow within 3 minutes in usability testing.
  - AC2: All form fields display inline help text that names the
    expected input.
  - AC3: All error messages name the field, the violated rule, and
    a remediation step.
```

If a requirement resists ACs, it usually needs to be split, made
more concrete, or escalated as a clarification (see "Requirements
Ambiguity — Do Not Guess" in CLAUDE.md).

## Unambiguous, agnostic vocabulary

Use the project's glossary terms. Avoid jargon, slang, and
implementation-specific component names — the requirement says
*what*, not *how*.

```markdown
Avoid (implementation leaks in):
- The system must show a modal dialog when the user clicks "Save."

Preferred (agnostic; links to glossary):
- The system must present a confirmation [affordance](../glossary.md#edit-affordance)
  before persisting changes triggered by the
  [Save action](../glossary.md#action-trigger).
```

If a hard constraint (regulatory, accessibility) requires a
specific component, link the concrete term to a justification entry
in `docs/architecture/` or a tech-ref entry — see the Markdown link
convention in CLAUDE.md.

Words that almost always introduce ambiguity, and what to do about
them:

| Word | Why it is ambiguous | Replacement |
|------|---------------------|-------------|
| "fast" / "quick" | No threshold | Specify the latency budget (e.g., "within 200 ms") |
| "secure" | Not measurable | Specify the property (e.g., "encrypted at rest", "rate-limited at N requests/min") |
| "user-friendly" / "intuitive" | Not measurable | Specify the observable behavior (see the "first-time user" example above) |
| "appropriate" / "reasonable" | Defers the decision | State the rule directly |
| "etc." / "and so on" | Open-ended scope | Enumerate the cases or scope down |
| "if possible" / "preferably" | Optional, not required | Either remove (if optional) or commit (if required) |

## Distinguish requirements from user stories

User stories ("As a user, I want X so that Y") are **drivers** of
requirements, not requirements themselves. They capture *intent*
and *value* in a form that's useful for prioritization and
discussion. The Analyst translates them into formal requirements
when committing to `docs/reqs/`.

A user story might live in a discussion, a request, or a design
document. The requirement that comes from it lives in `docs/reqs/`
in system-facing imperative form, with acceptance criteria. Don't
file user stories as requirements; translate them.

## Distinguish requirements from acceptance criteria

The requirement is the *what*. The ACs are the *observable
verifications* that the *what* is in place. ACs sit under the
requirement they verify; they are not separate requirements
themselves. See "Requirement Status" in CLAUDE.md for the doc
structure (parent requirement + `implementation` + per-AC
checkboxes).

When an AC starts to read like a requirement on its own ("the
system must also..."), it is probably a sibling requirement that
needs its own statement, not an AC under another requirement.

## Active voice; concrete subjects

Active voice with the system as subject. Avoid passive constructions
that obscure who does what.

```markdown
Avoid:
- It must be possible for the data to be exported.

Preferred:
- The system must export the dataset in CSV format on user request.
```

When the subject is genuinely a subsystem or a specific component,
name it: "The export service must ...", "The login view must ...".
Specificity helps the implementer locate the correct module.
