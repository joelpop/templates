# Custom CSS

When no theme variant, `HasTheme` name, or `LumoUtility` class covers the need, add a named CSS class with `addClassNames()` and a matching stylesheet rule so custom CSS stays targeted and framework capabilities are not bypassed.

Acceptable uses:
- Styles requiring pseudo-selectors (`:hover`, `:focus`, `::before`) that cannot be expressed as a static class or inline style
- Component tweaks that variants and `LumoUtility` do not cover
- Application-specific layout rules for non-standard compositions
- Scoped overrides for third-party component styles

A large custom CSS file signals the styling progression (`theming/component-variants.md`) is being skipped.

**Related:** `theming/component-variants.md` — full styling progression; `theming/brand-customization.md` — Lumo property overrides; `theming/stylesheet-and-theme-loading.md` — where stylesheets live and how they are loaded.
