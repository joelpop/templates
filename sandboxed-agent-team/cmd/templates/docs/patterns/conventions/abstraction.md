# Abstraction Recognition

When to extract a value object, type, helper, or shared mechanism — and
when not to. The cost of catching this late is high; the cost of
catching it early is small.

## The third-instance rule

> One conditional is fine. Two is a pattern. Three is a framework that
> doesn't exist yet. **Catch it at two.**

When the same shape appears in two places, you have a candidate
abstraction. Don't extract on the first instance — that's premature.
Don't wait for the third instance — by then the duplication has hardened
and the cost of unifying is high. The window for cheap extraction is
*right after the second instance lands.*

The rule applies to:

- **Branching logic** — repeated `if`/`switch` shapes that key off the
  same field
- **Validation** — repeated null checks, range checks, format checks
- **Mapping** — repeated translations between the same two
  representations
- **Error handling** — repeated try/catch shapes producing the same
  outcome
- **State combinations** — repeated bundles of fields treated as a unit
  in multiple places

## State shapes argue for value objects

When a small group of fields is consistently treated as a unit — passed
together, validated together, displayed together — that's a value
object, not a scatter of primitives. The pull to use primitives is real
(less code now), but the cost compounds (every consumer reimplements the
combination).

Two recurring examples worth recognizing on sight:

### "Bytes plus content type" → a value object

```java
// Avoid — primitives scattered across the API; every consumer
// reassembles the relationship between bytes and type
public class Asset {
    private byte[] data;
    private String contentType;
}

void process(byte[] data, String contentType) { ... }

// Preferred — the relationship is encoded in a type
public record ContentData(byte[] bytes, String contentType) { ... }

public class Asset {
    private ContentData content;
}

void process(ContentData content) { ... }
```

The `ContentData` type localizes the "these two fields belong together"
invariant. Every consumer benefits: signatures shrink, validation is
single-sourced, and downstream code can ask `content.contentType()`
without worrying whether some caller passed mismatched arguments.

### "First name plus last name" → a name type, not five string variants

```java
// Avoid — the model exposes the components, then five string
// variants get added across views as needs arise
String firstName;
String lastName;
// ... and elsewhere ...
String fullName;       // "Alice Smith"
String firstLast;      // "Alice Smith"
String lastFirst;      // "Smith, Alice"
String displayName;    // varies by locale

// Preferred — components stay structured; rendering is the consumer's
// job, not the model's
public record PersonName(String first, String last) {
    String full()      { return "%s %s".formatted(first, last); }
    String lastFirst() { return "%s, %s".formatted(last, first); }
}
```

The model carries structure (first, last). Display order, locale-aware
composition, and abbreviation choices live in renderers — formatting
code, not the data model. Every "we need yet another variant" pull is a
sign the rendering responsibility has leaked into the model.

## Don't extract prematurely

Two callers with *similar-looking* code aren't necessarily the same
abstraction. Before extracting, check that the duplicates *really*
share the same:

- **Inputs and outputs** (not just shape — meaning)
- **Failure modes**
- **Invariants** (what must be true before / after / during)

If the duplicates are conceptually different — they happen to look the
same today but represent different domain ideas — extracting will tie
them together prematurely. The next change to one will be hampered by
having to keep the other working. Better to leave the duplication and
revisit when the third instance reveals which way it really wants to
go.

A useful signal: can you give the abstraction a name that describes its
*purpose*, not its mechanics? "ContentData", "PersonName", "Money",
"DateRange" name purposes. "ProcessTwoFields" or "DoTheThing" don't —
they describe mechanics, which is a sign the abstraction isn't ready.

## Where extracted abstractions live

| Kind of duplication | Where the abstraction usually goes |
|---------------------|-------------------------------------|
| State bundle (fields used together) | A value type / record / Java class |
| Logic shape (algorithmic duplication) | A method on the relevant class, or a utility |
| Cross-cutting concern (validation, logging, mapping, etc.) | A shared mechanism (interface + implementation, or aspect) |
| Branching on type | Polymorphism — split the type, dispatch via virtual call |
| Repeated try/catch + error wrapping | A wrapper method or higher-order helper |

When the answer is "shared mechanism," that mechanism deserves an
[architecture entry](../../architecture/) describing its contract — not
just a piece of code that exists. Future implementers should be able to
find the mechanism without grepping.

## Architect's role

Pattern recognition during code review is one of the Architect's
highest-leverage activities. Specifically:

- When reading a Coder commit, watch for the second-instance moment.
  Flag it before the third instance ships.
- Name the candidate abstraction in the review feedback. Naming makes
  the proposal concrete.
- Record the pattern in `docs/architecture/` if the project commits to
  it. Subsequent implementers should read the pattern, not re-derive
  it.
- If the pattern is project-agnostic (like the `ContentData` shape
  above), the entry also goes in `docs/patterns/architecture/` so it
  carries to other projects.

## When to leave duplication alone

Duplication is not always a flaw. Cases where it stays:

- **The duplicated code is genuinely simple** — three or four lines
  with no state, no branching. Extracting adds an indirection without
  reducing complexity.
- **The duplicates have different change rates** — one is stable
  domain logic, the other is exploratory and changing weekly. Coupling
  them via a shared abstraction would slow the exploration.
- **Trying to extract revealed that the "duplicates" weren't really
  the same** — the differences are meaningful.

In each of those cases, document *why* the duplication stays — a short
note in the relevant `docs/architecture/` entry, or a comment at the
duplication site if the reason is non-obvious. The next reviewer
shouldn't have to relitigate the decision.