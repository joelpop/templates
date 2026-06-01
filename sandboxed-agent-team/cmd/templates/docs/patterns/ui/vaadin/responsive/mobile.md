# Mobile Layout Patterns (< 600px)

When rendering on mobile, use single-column layout, full-width dialogs, reduced grid columns, and a collapsed AppLayout drawer.

## Single-Column Layout

All primary views render as a single column. No horizontal overflow.

## Full-Width Dialogs

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

## Reduced Grid Columns

Grids reduce to essential columns below 600px. Full detail is accessible via row click.

```java
// Less important columns hidden on mobile
codeColumn.setVisible(!isMobile);
dateColumn.setVisible(!isMobile);
```

## Collapsed Navigation

The AppLayout drawer is collapsed by default on mobile, opened via `DrawerToggle`.
Consider a bottom tab bar for the 3–5 most accessed views on very small screens.
