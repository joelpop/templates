# Stylesheet and Theme Loading

When loading themes and stylesheets, use the version-appropriate mechanism so
CSS is applied globally before any view renders.

## Vaadin ≥25 — Stylesheet Loading

Stylesheets go in `src/main/resources/META-INF/resources/` and are loaded via
`@StyleSheet` on the `AppShellConfigurator` class:

```java
@StyleSheet(Lumo.STYLESHEET)
@StyleSheet(Lumo.UTILITY_STYLESHEET)
@StyleSheet("styles.css")
@ColorScheme(ColorScheme.Value.LIGHT_DARK)
public class Application implements AppShellConfigurator { /* ... */ }
```

`@ColorScheme(LIGHT_DARK)` follows the OS/browser preference and is the recommended
default — see `theming/dark-mode.md` for color scheme options.

Do not use these approaches in Vaadin ≥25:

| Approach                               | Why not                                            |
|----------------------------------------|----------------------------------------------------|
| `@Theme("my-theme")`                   | Deprecated in Vaadin 25; use `@StyleSheet` instead |
| `theme.json` to load utility/theme CSS | No longer supported in Vaadin 25                   |
| `src/main/frontend/themes/<name>/`     | Vaadin 24 theme folder; not the Vaadin 25 pattern  |
| `@CssImport` for application styles    | For add-ons only; app styles use `@StyleSheet`     |
| `@JsModule` for stylesheets            | For JS/TS modules only, not CSS                    |

## Vaadin 24.x — Theme-Based Loading

A single `@Theme` on `AppShellConfigurator` loads the theme folder. The theme
folder is the root for all application styles — `styles.css` is the entry point
and pulls in additional stylesheets via CSS `@import`. Do not use `@StyleSheet`
or `@CssImport` for application styles alongside `@Theme`; they are for custom
standalone components only.

```java
@Theme("<theme-name>")
public class Application implements AppShellConfigurator { /* ... */ }
```

The theme folder at `src/main/frontend/themes/<theme-name>/` contains
`styles.css` (required), optional sub-stylesheets imported from it, and
`theme.json` for asset declarations and parent theme configuration.

**Related:** `theming/brand-assets.md` — where to place fonts, images, and logo;
Vaadin docs — `https://vaadin.com/docs/latest/styling/stylesheets`,
`https://vaadin.com/docs/latest/flow/advanced/loading-resources`
