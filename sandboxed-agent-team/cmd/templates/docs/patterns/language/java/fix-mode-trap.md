# The Fix-Mode Trap

When you've just been told what was wrong and why, and you're now changing the
code to fix it, resist the temptation to paste that explanation into a comment
block. The explanation belonged in the conversation that produced it; it belongs
next in the commit message; it does not belong in the code.

The mental check: did the existing code have a comment block describing why it
was the way it was? If no, your fix shouldn't add one. The fix either makes the
code self-explanatory (preferred) or, in the rare case the fix introduces a
non-obvious invariant, gets a short comment about that **invariant** — not about
the path that led to discovering it.

```java
// Avoid — fix-mode comment explaining the history

// We used to combine firstname and lastname into a single
// "fullName" string, but that made it impossible to sort by
// last name and caused issues when localizing the display order.
public record Person(String firstname, String lastname) { }
```

```java
// Preferred — the code is self-explanatory; rationale lives in the commit message
public record Person(String firstname, String lastname) { }
```
