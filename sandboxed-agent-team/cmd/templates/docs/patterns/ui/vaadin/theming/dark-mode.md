# Dark Mode

When adding dark mode support, configure `@ColorScheme` on `AppShellConfigurator` and use
semantic color tokens in custom CSS so the color scheme takes effect globally and custom
styles adapt automatically.

## Activation

Static application-wide color scheme:

```java
@ColorScheme(ColorScheme.Value.LIGHT_DARK)
public class Application implements AppShellConfigurator { /* ... */ }
```

| Value        | Behavior                                                       |
|--------------|----------------------------------------------------------------|
| `LIGHT`      | Always light (Vaadin default if `@ColorScheme` is omitted)     |
| `DARK`       | Always dark                                                    |
| `LIGHT_DARK` | Follows OS/browser preference, defaults to light (recommended) |
| `DARK_LIGHT` | Follows OS/browser preference, defaults to dark                |

Dynamic toggle at runtime:

```java
UI.getCurrentOrThrow().getPage().setColorScheme(ColorScheme.Value.DARK);
```

## Token Overrides in Dark Mode

Overriding a Lumo token (see `theming/brand-customization.md`) replaces its built-in
light/dark adaptation. Use `light-dark(lightValue, darkValue)` so the override adapts to the active color scheme
regardless of which scheme is configured:

```css
html {
    --lumo-primary-color: light-dark(#1a73e8, #7baaf7);
}
```

Do not use the `[theme~="dark"]` selector for token overrides — it only works with the
forced `DARK` scheme and silently has no effect with `LIGHT_DARK` or `DARK_LIGHT`.

Non-color token overrides (font family, border radius) do not need paired dark values.

**Related:** `theming/brand-customization.md` — Lumo property overrides.