---
name: figma-to-lumo-theme
description: Map Figma design tokens to Lumo CSS variables by extracting tokens, categorizing them, and generating CSS declarations in the styles.css file with only non-default values.
---

# Figma to Lumo CSS Variables Mapping

When translating a Figma design system into the application's Vaadin theme,
extract Figma variables via `get_variable_defs`, map each token to a Lumo CSS
custom property, and emit only the non-default declarations into `styles.css`
so the implementation tracks the design source without overriding Lumo defaults
that already match.

## Required Workflow for CSS Variable Mapping
Create TODOs based on these steps.
- Step 1: Extract Figma Design Tokens
- Step 2: Categorize Design Tokens
- Step 3: Extract Component Styles From Figma with `get_design_context` tool
- Step 4. Review Component Styling Documentation
- Step 5: Generate CSS Variable Declarations (Only Non-Default Values)

This approach ensures:
- ✅ Only custom values are set, preserving Lumo defaults
- ✅ Component-specific variables are included when needed  
- ✅ Dark theme only included if Figma design shows dark variants
- ✅ Systematic conversion from Figma tokens to CSS variables
- ✅ No validation/testing steps - pure value transfer process

### Step 1: Extract Figma Design Tokens
**Always start with `get_variable_defs`** — returns tokens like:
```javascript
{
  "lumo-primary-color": "#006af5",
  "lumo-header-text-color": "#192434",
  "Normal body text": "Font(family: \"Inter\", style: Regular, size: 16, weight: 400, lineHeight: 1.61)",
  "lumo-border-radius-m": "8",
  "vaadin-input-field-height": "36",
  "vaadin-input-field-border-color": "#1c304a85"
}
```

### Step 2: Categorize Design Tokens
Organize Figma tokens by Lumo categories:

#### **Color Tokens** → Map to Lumo Color Properties
- `lumo-primary-color` → `--lumo-primary-color`
- `lumo-header-text-color` → `--lumo-header-text-color`
- `lumo-body-text-color` → `--lumo-body-text-color`
- `lumo-base` → `--lumo-base-color`
- `lumo-contrast-XXpct` → `--lumo-contrast-XXpct`
- `lumo-error-color`, `lumo-warning-color`, `lumo-success-color` → Error/warning/success colors

#### **Typography Tokens** → Map to Lumo Typography Properties
- `lumo-font-family` → `--lumo-font-family`
- `lumo-font-size-xxxl`, `lumo-font-size-xxl`, `lumo-font-size-xl`, etc. → `--lumo-font-size-{xxxl|xxl|xl|l|m|s|xs|xxs}`
- Font() definitions from `Heading 1`, `Heading 2`, `Normal body text` → Extract size and line-height values

#### **Sizing Tokens** → Map to Lumo Size & Space Properties
- Component sizes → `--lumo-size-{xs|s|m|l|xl}`
- Spacing values → `--lumo-space-{xs|s|m|l|xl}`
- Icon sizes → `--lumo-icon-size-{s|m|l}`
- `Component sizes/Fields/Field - Label gap` → Custom spacing values

#### **Shape Tokens** → Map to Lumo Shape Properties
- `lumo-border-radius-s`, `lumo-border-radius-m`, `lumo-border-radius-l` → `--lumo-border-radius-{s|m|l}`

#### **Elevation Tokens** → Map to Lumo Elevation Properties
- Shadow definitions → `--lumo-box-shadow-{xs|s|m|l|xl}`

#### **Component-Specific Tokens** → Direct Variable Mapping
Many component variables are directly available in Figma:
- `vaadin-input-field-height` → `--vaadin-input-field-height`
- `vaadin-input-field-background` → `--vaadin-input-field-background`
- `vaadin-input-field-border-color` → `--vaadin-input-field-border-color`
- `vaadin-input-field-border-width` → `--vaadin-input-field-border-width`
- Use Vaadin MCP for additional component variables if needed

### Step 3: Extract Component Styles From Figma with `get_design_context`
- Most detailed component information
- Check `className` for Tailwind classnames
- Identify theme/styling variable names and values

### Step 4: Review Component Styling Documentation
Use Vaadin MCP to identify component-specific CSS variables:
```javascript
// Search for additional component variables if needed
search_vaadin_docs("button styling CSS variables")
get_full_document("components/button/styling-flow.md")
```

Use Vaadin MCP to identify Lumo theme defaults
```javascript
// Search for Lumo style properties
search_vaadin_docs("lumo style properties")
get_full_document("styling/lumo/lumo-style-properties.md")
```


### Step 5: Generate CSS Variable Declarations (Only Non-Default Values)
**Only set variables that differ from Lumo defaults.** Create CSS declarations in
`styles.css`:

```css
/* Imported from Figma */

html {
  /* Primary colors */
  --lumo-primary-color: hsla(218, 100%, 48%, 1);
  --lumo-primary-text-color: hsla(218, 100%, 44%, 1);

  /* Text colors */
  --lumo-header-text-color: hsla(218, 31%, 20%, 1);
  --lumo-body-text-color: hsla(218, 31%, 20%, 0.94);
  --lumo-secondary-text-color: hsla(218, 31%, 26%, 0.69);

  /* Typography */
  --lumo-font-family: "Inter", "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  --lumo-font-size-xl: 1.375rem;

  /* Component styles */
  --vaadin-input-field-height: 2.25rem;
  --vaadin-input-field-background: hsla(218, 31%, 35%, 0.1);
  --vaadin-input-field-border-color: hsla(218, 31%, 20%, 0.52);
}

/* Dark theme */
[theme~="dark"] {
  --lumo-primary-color: hsla(177, 35%, 55%, 1); /* Darker variant of primary */
}
```
