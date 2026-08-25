# LumoUtility for Spacing, Color, and Layout

When using the Lumo theme and neither theme variants nor `addThemeNames()` covers the need, use `LumoUtility` class constants (step 3 of the styling progression) rather than custom CSS so styling stays theme-aware without leaving the framework.

```java
// Avoid — custom CSS for things Lumo provides
layout.getStyle().set("padding", "var(--lumo-space-m)");
```

```java
// Preferred — LumoUtility constants
layout.addClassNames(
    LumoUtility.Padding.MEDIUM,
    LumoUtility.Gap.SMALL,
    LumoUtility.Display.FLEX,
    LumoUtility.FlexDirection.COLUMN,
    LumoUtility.AlignItems.CENTER
);

span.addClassNames(
    LumoUtility.FontSize.LARGE,
    LumoUtility.FontWeight.BOLD,
    LumoUtility.TextColor.PRIMARY
);
```

When elements with padding cause horizontal overflow, apply `LumoUtility.BoxSizing.BORDER`
to include padding within the declared width:

```java
element.addClassNames(LumoUtility.BoxSizing.BORDER);
```

**Related:** `theming/component-variants.md` — full styling progression;
`theming/theme-selection.md` — Lumo vs Aura trade-offs;
`theming/tailwind.md` — utility classes for Aura.
