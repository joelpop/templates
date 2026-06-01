---
name: lumo-value-conversion
description: Rules for converting Figma token values (sizes, typography, colors, input fields) into CSS variable declarations so generated CSS matches Lumo conventions.
---

# Lumo Value Conversion Guidelines

When converting Figma token values to CSS variable declarations, follow these rules
so generated CSS matches Lumo's conventions and only non-default values are emitted.

## Size Conversion
- **Pixel to rem**: Divide by 16 (`36px` → `2.25rem`) — only if different from Lumo default
- **Unitless numbers**: Treat as pixels for border-radius, spacing
- **Font sizes**: Convert to rem; skip if matches Lumo default (16px = 1rem for font-size-m)
- **Line heights**: Unitless values (`1.6`, `1.25`)

## Typography Parsing
Parse Figma Font() definitions:
```javascript
"Normal body text": "Font(family: \"Inter\", style: Regular, size: 16, weight: 400, lineHeight: 1.61)"
```
Extract:
- family: "Inter" → font-family (only if different from system stack)
- size: 16 → Skip --lumo-font-size-m if it equals default 1rem
- weight: 400 → font-weight (handled by components, not CSS variables)
- lineHeight: 1.61 → --lumo-line-height-m: 1.6 (only if different from default)

## Color Values
Preserve color syntax as-is — a separate step handles format conversions. Map colors to
CSS variables only.

## Color Configurations
When assigning Primary, Error, Warning, Success, or Contrast colors, assign all related
values:
```css
/* Proper way to set colors is to take all related values into account */
  --lumo-primary-color: hsl(214, 100%, 48%);
  --lumo-primary-color-50pct: hsla(214, 100%, 49%, 0.76);
  --lumo-primary-color-10pct: hsla(214, 100%, 60%, 0.13);
  --lumo-primary-text-color: hsl(214, 100%, 43%);

/* Incorrect to only set some color values */
  --lumo-primary-color: hsl(214, 100%, 48%);
  --lumo-primary-color-10pct: hsla(214, 100%, 60%, 0.13);
```

## Suitable Values for Contrast Color Variables
- `--lumo-contrast-XXpct` color values can only have color values that have lightness channel value below 25%.
- Only exception is in the scope of `theme="dark"` where the lightness channel has to be above 75%.

## Input Field Styling
- If `--vaadin-input-field-border-color` is set, ensure `--vaadin-input-field-border-width` is also set at least 1px.
- Use CSS variables for styling whenever possible. If writing custom styles of input fields do not target specific fields. Target `::part(input-field)`, `::part(label)` and `::part(value)` to keep styling of all inputs consistent.

## Component-Specific Variables
Many component variables are directly available in Figma with exact CSS variable names:

```javascript
// Direct variable mapping from Figma
"vaadin-input-field-height" → --vaadin-input-field-height
"vaadin-input-field-background" → --vaadin-input-field-background  
"vaadin-input-field-border-color" → --vaadin-input-field-border-color
"vaadin-input-field-border-width" → --vaadin-input-field-border-width
"vaadin-user-color-0" → --vaadin-user-color-0
"vaadin-user-color-1" → --vaadin-user-color-1
// ... etc
```

Map Figma component tokens to variables:
```css
/* Input Field Component */
--vaadin-input-field-height: 2.25rem;
--vaadin-input-field-background: hsla(218, 31%, 35%, 0.1);
--vaadin-input-field-border-color: hsla(218, 31%, 20%, 0.52);
--vaadin-input-field-border-width: 1px;

/* User Colors */
--vaadin-user-color-0: hsla(320, 87%, 46%, 1);
--vaadin-user-color-1: hsla(258, 85%, 42%, 1);
--vaadin-user-color-2: hsla(193, 86%, 34%, 1);
--vaadin-user-color-3: hsla(35, 100%, 34%, 1);
```
