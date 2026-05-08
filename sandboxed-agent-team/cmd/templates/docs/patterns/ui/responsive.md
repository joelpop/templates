# Responsive Layout Patterns

Breakpoints, device detection, and adaptive behavior for Vaadin 24+ applications across
phone, tablet, and desktop form factors. `FormLayout` responsive steps, `AppLayout` drawer
behavior, and LumoUtility classes are stable across all supported Vaadin lines. The
`retrieveExtendedClientDetails` API used for server-side width detection differs between
Vaadin 24 and 25 — see `docs/agnostic/conventions/vaadin.md` → "Version-Specific Notes" in
the Browser Client Details section.

## Breakpoints

| Breakpoint | Width | Layout behavior |
|------------|-------|-----------------|
| Mobile | < 600px | Single column, full-width dialogs, collapsed nav |
| Tablet | 600px–1023px | Two-column where appropriate, drawer overlays content |
| Desktop | 1024px+ | Full layout, persistent nav sidebar |

These breakpoints apply to all primary views. Layouts must not produce horizontal scrolling
or clipped content at any breakpoint.

## Vaadin Responsive Tools — Not CSS Frameworks

Use Vaadin's built-in responsive features. Do not use Bootstrap, Tailwind, or any other CSS
framework grid system. Layout changes at breakpoints are achieved with `LumoUtility` classes
and Vaadin component responsive APIs.

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

Grids reduce to essential columns on screens below 600px. Full detail is accessible via
row click opening a detail dialog.

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

On mobile, the AppLayout drawer is collapsed by default and opened via `DrawerToggle`.
Consider a bottom tab bar for the 3–5 most frequently accessed views as an alternative or
supplement to the drawer on very small screens.

## Tablet (600px–1023px)

### Two-Column Forms

`FormLayout` with `ResponsiveStep` naturally produces two-column forms at tablet width:

```java
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),
    new FormLayout.ResponsiveStep("600px", 2)
);
```

### Drawer Overlay

The AppLayout drawer overlays the content at tablet width rather than pushing it aside.
This is the default Vaadin AppLayout behavior — no additional configuration needed.

### Touch Targets

Ensure interactive elements (buttons, grid rows, form fields) have adequate touch target
size. Vaadin's default component sizing is designed for touch, but avoid custom compact
styles that reduce hit areas below 44px.

## Desktop (1024px+)

### Persistent Navigation

At desktop width the navigation sidebar is visible by default and does not overlay content.
AppLayout handles this automatically at its built-in breakpoint.

### Multi-Column Layouts

Use `HorizontalLayout`, `SplitLayout`, or CSS grid (via `LumoUtility`) for side-by-side
content panels where appropriate.

## Server-Side Breakpoint Detection

For Java-level decisions based on viewport width, use Vaadin's `Page` API:

```java
UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> {
    var isMobile = details.getWindowInnerWidth() < 600;
    // adjust layout server-side
});
```

Use this sparingly — prefer CSS-based responsive rules via `LumoUtility` or `@media`
queries for visual-only changes.

## Accessibility at All Breakpoints

- Tab order remains logical at all breakpoints
- No interactive element becomes unreachable by keyboard when layout changes
- Touch targets are large enough at mobile (aim for 44px minimum)
- Text remains readable without zooming (minimum 16px effective font size for body text)
