# Component Theme Variants

When styling a Vaadin component, work through this progression so built-in
framework styling is exhausted before reaching for custom CSS:

1. **Theme variants** — typed enum constants via `addThemeVariants()`
2. **Component theming** — named theme strings via `addThemeNames()` (`HasTheme`)
3. **Utility classes** — `LumoUtility` constants via `addClassNames()` (Lumo); Tailwind CSS classes via `addClassNames()` (Aura)
4. **Custom CSS class** — `addClassNames("my-class")` with a stylesheet rule
5. **Inline style** — `getStyle()` only for runtime-computed values with no other option

```java
// 1. Theme variants — prefer these first
button.addThemeVariants(ButtonVariant.PRIMARY, ButtonVariant.SMALL);
grid.addThemeVariants(GridVariant.ROW_STRIPES, GridVariant.LUMO_COMPACT);  // LUMO_ kept: Lumo-only variant
field.addThemeVariants(TextFieldVariant.SMALL);

// 2. Component theming — for named theme values without a typed variant
badge.addThemeNames("badge", "success");

// 3. Utility classes — layout, spacing, typography
div.addClassNames(LumoUtility.Display.FLEX, LumoUtility.Gap.MEDIUM);

// 4. Custom CSS class — when none of the above suffice
div.addClassNames("section-header");
```

Always use the plural form (`addThemeVariants`, `addThemeNames`, `addClassNames`)
even when passing a single value.

## Variant Constant Names by Version

Variant constant naming changed across Vaadin releases:

| Variant                    | v24 / v25.0                       | v25.1+                        |
|----------------------------|-----------------------------------|-------------------------------|
| Button primary             | `ButtonVariant.LUMO_PRIMARY`      | `ButtonVariant.PRIMARY`       |
| Button small               | `ButtonVariant.LUMO_SMALL`        | `ButtonVariant.SMALL`         |
| Grid row stripes           | `GridVariant.LUMO_ROW_STRIPES`    | `GridVariant.ROW_STRIPES`     |
| Grid compact (Lumo only)   | `GridVariant.LUMO_COMPACT`        | `GridVariant.LUMO_COMPACT`    |
| TextField small            | `TextFieldVariant.LUMO_SMALL`     | `TextFieldVariant.SMALL`      |

In v25.0, the Aura theme introduced a parallel `AURA_` prefix (e.g., `ButtonVariant.AURA_PRIMARY`).
In v25.1+, the theme-neutral unprefixed name is standard; Lumo-only variants retain `LUMO_`.

**Related:** `theming/lumo-utility.md`, `theming/custom-css.md`.
