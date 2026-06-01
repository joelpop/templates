# LumoUtility for Spacing, Color, and Layout

When applying padding, margin, color, flexbox, or sizing, use `LumoUtility` class constants rather than custom CSS so styling stays theme-aware and correct across light and dark modes.

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
