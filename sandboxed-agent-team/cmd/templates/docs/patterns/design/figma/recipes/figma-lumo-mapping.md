---
name: figma-lumo-mapping
description: How to map Figma design tokens to Vaadin LumoUtility class constants and CSS custom properties so generated code stays framework-aligned and avoids inline styles.
---

# Lumo Theme Mapping Guidelines

When mapping Figma design tokens to Vaadin/Java code, use LumoUtility class constants
before CSS custom properties so the implementation stays framework-aligned and type-safe.

## Examples of Figma → Lumo Color Mappings
```java
// Colors
"Semantic colors/Primary" → "var(--lumo-primary-color)" or LumoUtility.Background.PRIMARY
"Semantic colors/Primary, Text" → "var(--lumo-primary-text-color)" or LumoUtility.TextColor.PRIMARY
"Header Text" → "var(--lumo-header-text-color)" or LumoUtility.TextColor.HEADER
"Body Text" → "var(--lumo-body-text-color)" or LumoUtility.TextColor.BODY

// Typography
"Typography/Font-family" → "var(--lumo-font-family)"
"Typography/Font-size-m" → "var(--lumo-font-size-m)"
```

## Implementation Examples — LumoUtility (preferred)
```java
// Using Lumo Utility Classes (preferred when available)
title.addClassNames(LumoUtility.TextColor.HEADER);
card.addClassNames(LumoUtility.Background.CONTRAST_5);

// Proper way to configure component styles
span.addClassNames(LumoUtility.TextColor.SECONDARY, LumoUtility.FontSize.SMALL);

// ❌ INCORRECT way to configure component styles
span.getStyle().set("color", "rgba(27,43,65,0.69)").set("font-size", "15px").set("line-height", "1.34");
```

## Alternative Using CSS Styles
Add a custom classname to the element:
```java
span.addClassName("secondary-text");
```

Target the classname in CSS. Use `styles.css` unless there is a component- or view-specific
stylesheet that is more appropriate. Whenever possible use existing CSS custom properties
instead of defining new values:
```css
.secondary-text {
    color: var(--lumo-secondary-text-color);
    font-size: var(--lumo-font-size-s)
}
```
