# Responsive Breakpoints

When writing responsive layout code, use these breakpoints so sizing decisions
are consistent across the application.

| Breakpoint  | Min Width | `LumoUtility` breakpoint class | CSS prefix |
|-------------|-----------|--------------------------------|------------|
| Extra Small | 0px       | —                              | —          |
| Small       | 640px     | `Breakpoint.Small`             | `sm:`      |
| Medium      | 768px     | `Breakpoint.Medium`            | `md:`      |
| Large       | 1024px    | `Breakpoint.Large`             | `lg:`      |
| Extra Large | 1280px    | `Breakpoint.XLarge`            | `xl:`      |

Breakpoints follow a mobile-first approach: styles apply at the named width and above.
Only six `LumoUtility` categories define a `Breakpoint` sub-class:

- `AlignItems` — aligning items along a flexbox's cross axis or a grid's block axis
- `Display` — setting the display property; determines block/inline and item layout
- `FlexDirection` — setting the flex direction of a flexbox layout
- `FontSize` — setting the font size of an element
- `Grid` — content flow on a grid layout
- `Position` — setting the position of an element

The breakpoint class is an intermediate path segment — usable constants are nested
inside it, e.g. `LumoUtility.Display.Breakpoint.Small.HIDDEN`. CSS prefixes apply
when writing raw utility classes in HTML or TSX.

Lumo does not expose the pixel values as Java constants. Define a `Breakpoints`
class in the UI module so server-side comparisons against `getWindowInnerWidth()`
avoid magic numbers:

```java
public enum Breakpoints {
    XS(0), SM(640), MD(768), LG(1024), XL(1280);

    public final int minWidthPx;
    public final String minWidth;

    Breakpoints(int minWidthPx) {
        this.minWidthPx = minWidthPx;
        this.minWidth = minWidthPx == 0 ? "0" : minWidthPx + "px";
    }
}
```

**Related:** `theming/lumo-utility.md`,
`responsive/server-side-detection.md`;
Vaadin docs — `https://vaadin.com/docs/latest/styling/lumo/utility-classes`
