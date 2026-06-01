# Code Comment Discipline

When writing code, add a comment only when a non-obvious constraint, invariant,
or workaround would surprise or stump a future reader — all other explanations belong in
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
```

```java
// Preferred — captures a non-obvious invariant the code alone cannot convey

// `authProvider` may be null during the first request after a hot reload;
// this branch is what `SecurityFilter` falls back to in that window.
if (authProvider == null) { ... }
```

For doc comments on public types and methods, see `javadoc.md` — that is a
separate concern; this rule governs inline narrative comments inside method
bodies.
