# No Task, Fix, or Caller References in Code Comments

When adding a comment that explains why code exists or what was changed, put
that context in the commit message and PR description — not in the source code.

The current task, the bug being fixed, the developer asking for the change, the
PR — none of that belongs in code comments. Those facts live in commit messages,
PR descriptions, and issue trackers, where they're discoverable via `git blame`,
`git log`, and the platform UI without polluting the source.

```java
// Avoid

// Added per Joel's request in PR #482 to handle the case where
// the dashboard widget loads before auth completes.
if (authProvider == null) { /* ... */ }

// Avoid

// FIX: removed the null check that was causing NPEs in the
// reporting flow when the customer field is empty
public void process(Order order) { /* ... */ }
```
