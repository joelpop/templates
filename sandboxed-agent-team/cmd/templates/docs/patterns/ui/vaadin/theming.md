# Theming

When styling a Vaadin component or view, apply Lumo CSS custom properties via
`LumoUtility` class names and component theme variants rather than hand-written CSS
so visual choices stay aligned with the design tokens and brand overrides remain
isolated to the central theme module.

## Lumo Is the Only Theme

The application must use the Vaadin Lumo theme. The Aura theme is not used.

Configure in `Application.java`:

```java
@Theme(Lumo.class)
// or via @StyleSheet(Lumo.STYLESHEET) + @StyleSheet(Lumo.UTILITY_STYLESHEET)
public class Application implements AppShellConfigurator { ... }
```

Lumo CSS custom properties are applied globally and define the visual foundation for all
components.

## LumoUtility for Spacing, Color, and Layout

Use `LumoUtility` class constants for padding, margin, color, flexbox, and sizing instead
of custom CSS. This ensures consistency, theme-awareness, and correct behavior across
light and dark modes.

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

// Avoid — custom CSS for things Lumo provides
layout.getStyle().set("padding", "var(--lumo-space-m)");
```

When elements with padding cause horizontal overflow (e.g., 100% width + padding exceeds
container), apply `LumoUtility.BoxSizing.BORDER` to include padding within the declared width:

```java
element.addClassNames(LumoUtility.BoxSizing.BORDER);
```

## Component Theme Variants

Use Vaadin component theme variants before resorting to custom styling:

```java
// Button variants
button.addThemeVariants(ButtonVariant.LUMO_PRIMARY);
button.addThemeVariants(ButtonVariant.LUMO_ERROR);
button.addThemeVariants(ButtonVariant.LUMO_TERTIARY);

// Grid variants
grid.addThemeVariants(GridVariant.LUMO_ROW_STRIPES);
grid.addThemeVariants(GridVariant.LUMO_NO_BORDER);

// TextField variants
field.addThemeVariants(TextFieldVariant.LUMO_SMALL);

// Badge variants (for status indicators)
badge.getElement().getThemeList().add("badge success");
badge.getElement().getThemeList().add("badge error");
```

Check available variants in the Vaadin documentation before writing custom CSS.

## Brand Customization via CSS Custom Properties

Override Lumo CSS custom properties in a dedicated stylesheet — do not modify Lumo
component internals directly.

```css
/* styles.css — brand overrides */
html {
    --lumo-primary-color: #1a73e8;
    --lumo-primary-text-color: #1a73e8;
    --lumo-primary-color-10pct: #1a73e810;
    --lumo-primary-color-50pct: #1a73e880;
    --lumo-font-family: 'Inter', sans-serif;
    --lumo-border-radius-m: 6px;
}
```

Common overridable properties:

| Property | Controls |
|----------|---------|
| `--lumo-primary-color` | Primary action color (buttons, links, focus rings) |
| `--lumo-primary-text-color` | Text color for primary actions |
| `--lumo-font-family` | Application font family |
| `--lumo-border-radius-m` | Default component corner radius |
| `--lumo-base-color` | Page background |
| `--lumo-body-text-color` | Default text color |
| `--lumo-header-text-color` | Heading text color |
| `--lumo-secondary-text-color` | Subdued text (labels, captions) |
| `--lumo-success-color` | Success state color |
| `--lumo-error-color` | Error state color |
| `--lumo-warning-color` | Warning state color |

## Custom CSS — When It Is Acceptable

Custom CSS files are acceptable only for:
- Brand overrides not expressible via `LumoUtility` (custom font loading, logo placement)
- Component tweaks that variants do not cover
- Application-specific layout rules for non-standard compositions

Keep custom CSS minimal. A large custom CSS file signals framework capabilities are
being bypassed.

## Logo and Brand Assets

Place brand assets (logo, fonts) in `src/main/frontend/` and reference via `@StyleSheet`
or `@JsModule` on `Application.java`, or from the custom CSS stylesheet.

## Dark Mode

Lumo supports dark mode via the `dark` theme variant on `<html>`. Use Lumo's semantic
color tokens (`--lumo-primary-color`, `--lumo-base-color`, etc.) — they adapt
automatically. Do not hardcode hex values in component styles.
