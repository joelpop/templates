# Responsive Layout Patterns

Breakpoints, device detection, and adaptive behavior for Vaadin 24+ across phone, tablet,
and desktop. `FormLayout` responsive steps, `AppLayout` drawer behavior, and LumoUtility
classes are stable across all supported Vaadin lines. The `retrieveExtendedClientDetails`
API for server-side width detection differs between Vaadin 24 and 25 — see
`docs/patterns/README.md` → "Version Compatibility".

## Breakpoints

| Breakpoint | Width | Layout behavior |
|------------|-------|-----------------|
| Mobile | < 600px | Single column, full-width dialogs, collapsed nav |
| Tablet | 600px–1023px | Two-column where appropriate, drawer overlays content |
| Desktop | 1024px+ | Full layout, persistent nav sidebar |

Layouts must not produce horizontal scrolling or clipped content at any breakpoint.

## Vaadin Responsive Tools — Not CSS Frameworks

Use Vaadin's built-in responsive features — not Bootstrap, Tailwind, or any other CSS
framework grid system. Achieve breakpoint changes with `LumoUtility` classes and Vaadin
component responsive APIs.

```java
// Preferred — Vaadin responsive step on FormLayout
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),       // 1 column below 600px
    new FormLayout.ResponsiveStep("600px", 2),   // 2 columns at 600px+
    new FormLayout.ResponsiveStep("1024px", 3)   // 3 columns at 1024px+
);
```

## Mobile (< 600px)

### Single-Column Layout

All primary views render as a single column. No horizontal overflow.

### Full-Width Dialogs

Edit dialogs expand to full viewport width on mobile:

```java
dialog.setWidth("100vw");
dialog.setMaxWidth("100vw");
// or via responsive CSS
dialog.addClassName("mobile-full-width");
```

```css
@media (max-width: 599px) {
    .mobile-full-width {
        width: 100vw !important;
        max-width: 100vw !important;
        margin: 0 !important;
    }
}
```

### Reduced Grid Columns

Grids reduce to essential columns below 600px. Full detail is accessible via row click.

```java
var nameColumn = grid.addColumn(Item::getName).setHeader("Name").setFlexGrow(1);
var statusColumn = grid.addColumn(Item::getStatus).setHeader("Status").setAutoWidth(true);

// Less important columns hidden on mobile
var codeColumn = grid.addColumn(Item::getCode).setHeader("Code");
var dateColumn = grid.addColumn(Item::getCreatedAt).setHeader("Created");

// Hide on mobile — Vaadin Grid column visibility
addAttachListener(_ -> {
    var isMobile = getWidth() < 600;
    codeColumn.setVisible(!isMobile);
    dateColumn.setVisible(!isMobile);
});
```

### Collapsed Navigation

The AppLayout drawer is collapsed by default on mobile, opened via `DrawerToggle`.
Consider a bottom tab bar for the 3–5 most accessed views on very small screens.

## Tablet (600px–1023px)

### Two-Column Forms

`FormLayout` with `ResponsiveStep` produces two-column forms at tablet width:

```java
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),
    new FormLayout.ResponsiveStep("600px", 2)
);
```

### Drawer Overlay

The AppLayout drawer overlays content at tablet width — default Vaadin behavior, no
additional configuration.

### Touch Targets

Vaadin's default sizing is touch-friendly; avoid custom compact styles that reduce hit
areas below 44px.

## Desktop (1024px+)

### Persistent Navigation

The navigation sidebar is persistent at desktop width — AppLayout handles this at its
built-in breakpoint.

### Multi-Column Layouts

Use `HorizontalLayout`, `SplitLayout`, or CSS grid (via `LumoUtility`) for side-by-side
content panels where appropriate.

## Server-Side Breakpoint Detection

For server-side decisions based on viewport width, use Vaadin's `Page` API:

```java
UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> {
    var isMobile = details.getWindowInnerWidth() < 600;
    // adjust layout server-side
});
```

Use sparingly — prefer CSS-based responsive rules via `LumoUtility` or `@media` queries
for visual-only changes.

## Accessibility at All Breakpoints

- Tab order remains logical at all breakpoints
- No interactive element becomes unreachable by keyboard when layout changes
- Touch targets are large enough at mobile (aim for 44px minimum)
- Text remains readable without zooming (minimum 16px effective font size for body text)
