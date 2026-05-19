# Business Glossary

The project's canonical vocabulary for business and user-facing
concepts. Every term defined here has:

- A **canonical form** that requirement docs and other internal
  documents use exclusively.
- An **anchor** (the heading) so requirements can link to it
  inline via Markdown links.
- A **definition** that distinguishes it from neighboring terms.
- An optional **slang variants** list — informal or alternative
  names a human might use in conversation, mapped to the canonical
  form. The team understands these in conversation; documents use
  the canonical form.

The Analyst curates and commits this file. Terms are proposed during
requirement drafting and approved alongside the requirements that
introduce them.

The glossary is small by design. If a term is widely understood and
unambiguous (e.g., "user," "create," "save"), it does not belong
here — only ambiguous or project-specific business terms are entries.

---

*(Seed file — replace examples with real entries as the project's
vocabulary develops.)*

## Edit affordance

A UI surface for editing a record. Examples include modal dialogs,
slide-over panels, full-page edit views, or inline-editable rows.
Requirement statements use this term when the choice between those
implementations is open.

**Slang variants:** "edit form", "edit pane", "edit view".

## Action trigger

A UI element that invokes an action. Examples include buttons, menu
items, swipe actions, or keyboard shortcuts. Requirement statements
use this term when the specific affordance is open.

**Slang variants:** "button" (when the human means an action
trigger generally, not a specific button widget).

## Impersonate

The act of an administrator assuming another user's identity for
the duration of a session, for diagnostic or support purposes.
Requirement statements use this canonical form.

**Slang variants:** "spoof", "masquerade", "log in as".
