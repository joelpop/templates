# What Browserless Tests Cannot Cover

When a test needs to assert focus, text selection, scroll position, CSS rendering, or shadow DOM behavior, use TestBench (browser-based) — browserless tests cannot exercise anything the browser engine mediates.

| Capability | Browserless | TestBench |
|---|---|---|
| Component state, server events | ✓ | ✓ |
| Focus | — | ✓ (`hasAttribute("focused")`) |
| Text selection | — | ✓ (`executeScript()` on `selectionStart`/`selectionEnd`) |
| Scroll position | — | ✓ (`getRect()` comparison) |
| CSS rendering, hover states | — | ✓ |
| Web component internals (slots, shadow DOM) | — | ✓ |

For browser-only assertions, use `waitUntil()` to let the browser settle before
asserting:

```java
waitUntil(_ -> view.isNameFieldFocused());
assertTrue(view.isNameFieldFocused());
```

See `docs/patterns/testing/recipes/testbench-browserless.md` for setup.
