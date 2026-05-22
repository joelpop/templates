# Code Comment Discipline

When writing code, add a comment only when a non-obvious constraint, invariant,
or workaround would surprise a future reader — all other explanations belong in
commit messages and PR descriptions, not in the source unless otherwise noted by
another pattern.

## Default: no comment

Well-named identifiers and small functions explain *what* the code does.
Add a comment only when the **why** is non-obvious from the code alone:
a hidden constraint, a subtle invariant, a workaround for a specific bug,
behavior that would surprise a reader. If removing the comment wouldn't
confuse a future reader, don't write it.

```java
// Avoid — restates what the code already says
// Increment the counter
counter++;

// Avoid — describes WHAT, not WHY
// Loop through users and send email
for (User u : users) sendEmail(u);

// OK — captures a non-obvious invariant
// `authProvider` may be null during the first request after a hot reload;
// this branch is what `SecurityFilter` falls back to in that window.
if (authProvider == null) { ... }
```

## Don't reference the task, fix, caller, or commit

The current task, the bug being fixed, the developer asking for the
change, the PR — none of that belongs in code comments. Those facts live
in commit messages, PR descriptions, and issue trackers, where they're
discoverable via `git blame`, `git log`, and the platform UI without
polluting the source.

```java
// Avoid
// Added per Joel's request in PR #482 to handle the case where
// the dashboard widget loads before auth completes.
if (authProvider == null) { ... }

// Avoid
// FIX: removed the null check that was causing NPEs in the
// reporting flow when the customer field is empty
public void process(Order order) { ... }
```

## The fix-mode trap

A specific failure pattern: when you've just been *told* what was wrong
and why, and you're now changing the code to fix it, the temptation is to
paste that explanation into a comment block above the change. **Resist
this.** The explanation belonged in the conversation that produced it; it
belongs next in the commit message and PR description; it does **not**
belong in the code.

The mental check: *did the existing code have a comment block describing
why it was the way it was?* If no, your fix shouldn't add one. The fix
either makes the code self-explanatory (preferred) or, in the rare case
the *fix* introduces a non-obvious invariant, gets a short comment about
that **invariant** — not about the path that led to discovering it.

```java
// Avoid — fix-mode comment explaining the history
// We used to combine firstname and lastname into a single
// "fullName" string, but that made it impossible to sort by
// last name and caused issues when localizing the display order.
// Joel pointed out that we should instead keep them as separate
// fields and let the rendering layer compose them. Now we expose
// firstname and lastname directly.
public record Person(String firstname, String lastname) { }

// Preferred — the code is self-explanatory; rationale lives in the
// commit message
public record Person(String firstname, String lastname) { }
```

## Where the fix's *why* should go

- **Commit message body** — the primary home. One paragraph on what was
  wrong, why the fix works, what alternatives were considered if any.
  This is what `git blame` surfaces years later.
- **PR description** — for cross-cutting context, screenshots, links to
  issues. PR descriptions are reviewer-facing.
- **A `docs/solutions/` entry** — if the fix changes a pattern
  (e.g., the value-object refactor in the Person example above), the new
  pattern goes in the architecture docs. The code itself stays clean.
- **A new `docs/patterns/` entry** — if the fix surfaced a project-
  agnostic rule worth carrying across projects, capture it there.

## TODOs, FIXMEs, XXX

Avoid in committed code. They almost never get cleaned up and become
noise. If something is genuinely unfinished, leave it broken and obvious
(don't compile, throw, or fail a test) so it can't be missed; or open an
issue and reference the issue in the commit message. A `// TODO: revisit
this` comment is information-free.

## Doc comments (Javadoc / JSDoc / docstrings)

A separate concern. Doc comments on public types, methods, and functions
*are* expected (per the Coder's responsibilities) — they describe the
contract for callers, which is information the call site doesn't have.
That's not what this rule is about. The rule above is about inline
narrative comments inside method bodies.