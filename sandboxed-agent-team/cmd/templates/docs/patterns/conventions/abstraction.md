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

## Beyond value objects: sizing the abstraction

Not every duplication wants a record. The third-instance rule still
applies, but the *size* of the abstraction has to match the shape of
the duplication. Two larger shapes recur often enough to be worth
recognizing on sight: a shared structural template (a base class) and
a cohesive multi-class capability (a package).

### Repeated structural template → shared base class

When several classes share *structure* — layout, chrome, lifecycle,
fixture wiring — rather than state, the abstraction is a base class,
not a value object. Recognition signal: every implementer rebuilds
the same scaffolding before doing the part that's actually different.

```java
// Avoid — every view re-implements the same chrome
public class CustomersView extends Composite<VerticalLayout> {
    public CustomersView() {
        var title = new H2("Customers");
        title.addClassNames(/* standard view-title styling */);
        var header = new HorizontalLayout(title, /* actions slot */);
        // ... shared scaffolding repeated in every view ...
        getContent().add(header, body);
    }
}
public class OrdersView extends Composite<VerticalLayout> { /* same shape, copied */ }

// Preferred — chrome lives in a base; subclasses say only what differs
public class BaseView extends Composite<VerticalLayout> {
    private final HorizontalLayout headerActions = new HorizontalLayout();
    private final VerticalLayout body = new VerticalLayout();
    protected BaseView(String title) { /* builds header + body once */ }
    protected void setHeader(Component c)  { headerActions.removeAll(); headerActions.add(c); }
    protected void setContent(Component c) { body.removeAll(); body.add(c); }
}

public class CustomersView extends BaseView {
    public CustomersView() { super("Customers"); setContent(new CustomersGrid()); }
}
```

Wrong-size signs:

- The "base" defines a single hook and nothing else → callers are
  still doing the work; the base earns nothing.
- Subclasses override most methods → the shared shape isn't real;
  the duplication was apparent, not structural.
- A new feature requires changing the base to accommodate one
  caller → the base is leaking caller-specific knowledge; either
  parameterise or split.

### Cohesive multi-class capability → component-family package

When N classes together realise one capability the requirements name
as a unit, the abstraction is a *package*, not a class. Each class
is a partial answer; the capability is the package. Recognition
signal: the package boundary maps to a requirement boundary, and the
public face is one or two types — the rest are package-private.

```
// Avoid — three list-view features each rebuild the same
// toolbar + quick-filter + filter-popover + grid skeleton
ui/customers/CustomersGrid.java       + toolbar, filter, popover code
ui/orders/OrdersGrid.java             + toolbar, filter, popover code
ui/products/ProductsGrid.java         + toolbar, filter, popover code

// Preferred — the capability is one package, parameterised over the
// row type; features use the package, not the parts inside
ui/component/itembrowser/
    ItemBrowser.java<T>           // public face: parameterised list view
    Toolbar.java                  // package-private; an ItemBrowser part
    FilterPopover.java            // package-private; an ItemBrowser part
    filter/
        CustomFilter.java         // public; features pass these in

ui/customers/CustomersView.java extends BaseView {
    setContent(new ItemBrowser<Customer>(/* configured for customers */));
}
```

Wrong-size signs:

- A package with one class → it's a class, not a package.
- Internals reused independently outside the package → either the
  capability isn't really cohesive (split the package) or you've
  built a feature-specific subsystem where a generic mechanism
  belonged (broaden, move to `docs/patterns/architecture/`).
- Extending the package requires deep cross-knowledge of its
  internals → it's a directory, not a package; the public surface
  is too narrow or the parts are too coupled.

## Don't extract prematurely

Two callers with *similar-looking* code aren't necessarily the same
abstraction. Before extracting, check that the duplicates *really*
share the same:

- **Inputs and outputs** (not just shape — meaning)
- **Failure modes**
- **Invariants** (what must be true before / after / during)

If the duplicates represent different domain ideas — they just look
similar today — extraction ties them together prematurely. Leave the
duplication and revisit when the third instance clarifies.

A useful signal: can you name the abstraction by *purpose*?
"ContentData", "PersonName", "Money", "DateRange" name purposes.
"ProcessTwoFields" or "DoTheThing" describe mechanics — the abstraction
isn't ready.

## Where extracted abstractions live

| Duplication shape | Right size | Where the abstraction goes |
|-------------------|-----------|-----------------------------|
| State bundle (fields used together) | Value object | Record / value type / Java class |
| Logic shape (algorithmic duplication) | Method | On the relevant class, or a utility |
| Repeated try/catch + error wrapping | Method | Wrapper method or higher-order helper |
| Branching on type | Polymorphism | Split the type, dispatch via virtual call |
| Repeated structural template (chrome, lifecycle, fixture) | Shared base class | Abstract or concrete base class — see "Sizing the abstraction" |
| N classes realising one named capability | Component-family package | A package; capability boundary = package boundary — see "Sizing the abstraction" |
| Cross-cutting concern (validation, logging, mapping, etc.) | Shared mechanism | Interface + implementation (or aspect); document the contract in `docs/solutions/` |

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
- Record the pattern in `docs/solutions/` if the project commits to
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

In each case, document *why* in the relevant `docs/solutions/` entry
or a comment at the duplication site — the next reviewer shouldn't have
to relitigate.