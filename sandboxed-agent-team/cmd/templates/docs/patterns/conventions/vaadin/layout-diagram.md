# Vaadin Layout Diagram Convention

Text-based layout diagrams in Javadoc `<pre>` blocks: placement, box labeling,
drawing rules, component content, width calculation, custom components, and
repeated items.

## Placement

Add a `<pre>...</pre>` block to the class Javadoc, after the one-line description
and a blank `*` line:

```java
/**
 * One-line description.
 *
 * <pre>
 * ...diagram...
 * </pre>
 */
```

## Box Labeling

Each box is labeled `+-varName(ClassName)--+`:

- Use the field name or local variable name as `varName`.
- For a `Composite<T>`'s root element (returned by `getContent()`), always use
  `content` as the name regardless of what the local variable is called.
- Omit `(ClassName)` and the parentheses when the variable name's suffix matches
  the class name (case-insensitive): `scroller` → `Scroller` ✓,
  `greetButton` → `Button` ✓, `content` wrapping a `Card` ✗ → `content(Card)`.
- Unnamed components use only the class name: `Div`.

## Drawing Characters

- `-` horizontal lines
- `|` vertical lines
- `+` line intersections
- Single space gap between every pair of nested `|` characters

## Component Content

| Component | Content shown |
|---|---|
| `Button` (text) | Label text, centered |
| `Button` (icon only) | `[X]`, centered |
| `TextField` | Placeholder or label text, left-aligned |
| `Checkbox` | `[ ] Label text`, left-aligned |
| `ComboBox` | value with `\| V` dropdown cell on right |
| `Span` / `Div` | Sample/representative text, left-aligned |

Individual component examples:

```
+-saveButton-+         +-userId(TextField)-------+
|    Save    |         | name@company.com        |
+------------+         +-------------------------+

+-showTimestampsOption(Checkbox)-+    +-defaultName(ComboBox)----------+---+
| [ ] Show timestamps            |    | Anonymous                      | V |
+--------------------------------+    +--------------------------------+---+
```

## Width Calculation

Build widths inside-out:

- **Leaf box:** `W = max(len(label) + 4, len(content) + 3)` (content gets 1 leading space)
- **Vertical layout wrapping one child:** `W_parent = W_child + 4`
- **Horizontal layout with two children:**
  `W_parent = W_child1 + 1 + W_child2 + 4` (the `+1` is the gap between children;
  the `+4` is the parent's own walls + spaces)
- Make all children of a vertical layout the same width as the widest sibling.
- A child with `setWidthFull()` or `setFlexGrow()` fills the available width; size it
  to consume all remaining space after fixed-size siblings.

## Custom Components

Each custom component class gets its own diagram in its own class comment. At the
use site, show a single placeholder box with no internals:

```
+-varName(ClassName)--+
| (see ClassName)     |
+---------------------+
```

## Dynamic / Repeated Items

Show two instances of the repeated item followed by three dot rows to convey
zero-or-more multiplicity:

```
| | | +-item(Type)--+ | | |
| | | | (see Type)  | | | |
| | | +-------------+ | | |
| | | +-item(Type)--+ | | |
| | | | (see Type)  | | | |
| | | +-------------+ | | |
| | |        .        | | |
| | |        .        | | |
| | |        .        | | |
```

## Verification

After drawing, confirm every vertical wall is straight: pick each `|` column and
trace it top-to-bottom — it must be `+` at every box corner and `|` on all
interior rows throughout.