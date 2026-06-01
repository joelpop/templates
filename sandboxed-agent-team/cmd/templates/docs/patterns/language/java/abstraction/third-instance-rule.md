# The Third-Instance Rule

When the same code shape appears for the second time, extract the abstraction — not at the first instance (premature), and not at the third (too late; the duplication has hardened).

> One conditional is fine. Two is a pattern. Three is a framework that
> doesn't exist yet. **Catch it at two.**

When the same shape appears in two places, you have a candidate abstraction.
The window for cheap extraction is *right after the second instance lands.*

The rule applies to:

- **Branching logic** — repeated `if`/`switch` shapes that key off the same field
- **Validation** — repeated null checks, range checks, format checks
- **Mapping** — repeated translations between the same two representations
- **Error handling** — repeated try/catch shapes producing the same outcome
- **State combinations** — repeated bundles of fields treated as a unit in multiple places

## Don't Extract Prematurely

Two callers with *similar-looking* code aren't necessarily the same abstraction.
Before extracting, check that the duplicates *really* share the same:

- **Inputs and outputs** (not just shape — meaning)
- **Failure modes**
- **Invariants** (what must be true before / after / during)

If the duplicates represent different domain ideas — they just look similar today —
extraction ties them together prematurely. Leave the duplication and revisit when
the third instance clarifies.

A useful signal: can you name the abstraction by *purpose*?
"ContentData", "PersonName", "Money", "DateRange" name purposes.
"ProcessTwoFields" or "DoTheThing" describe mechanics — the abstraction isn't ready.

## When to Leave Duplication Alone

Duplication is not always a flaw. Cases where it might stay:

- **The duplicated code is genuinely simple** — three or four lines with no state, no
  branching. Extracting adds an indirection without reducing complexity.
- **The duplicates have different change rates** — one is stable domain logic, the
  other is exploratory and changing weekly. Coupling them via a shared abstraction
  would slow the exploration.
- **Trying to extract revealed that the "duplicates" weren't really the same** —
  the differences are meaningful.

In each case, document *why* in the relevant `docs/solutions/` entry or a comment
at the duplication site — the next reviewer shouldn't have to relitigate.
