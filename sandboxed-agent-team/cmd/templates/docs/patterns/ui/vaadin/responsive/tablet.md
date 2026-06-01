# Tablet Layout Patterns (600px–1023px)

When rendering on tablet, use two-column forms, allow the AppLayout drawer to overlay content (default behavior), and ensure touch targets remain at least 44px.

## Two-Column Forms

`FormLayout` with `ResponsiveStep` produces two-column forms at tablet width:

```java
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),
    new FormLayout.ResponsiveStep("600px", 2)
);
```

## Drawer Overlay

The AppLayout drawer overlays content at tablet width — default Vaadin behavior, no
additional configuration.

## Touch Targets

Vaadin's default sizing is touch-friendly; avoid custom compact styles that reduce hit
areas below 44px.
