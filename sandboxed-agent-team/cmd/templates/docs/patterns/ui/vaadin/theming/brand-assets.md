# Logo, Fonts, and Brand Assets

When adding brand assets (logo, custom fonts, images), place them in the
version-appropriate location so they are resolvable from stylesheets and the
Vaadin `Image` component.

## Vaadin ≥25

Assets go in `src/main/resources/META-INF/resources/`:

```
src/main/resources/META-INF/resources/
├── fonts/
│   └── brand-font.woff2   ← referenced from styles.css via @font-face
└── images/
    └── logo.svg            ← referenced from CSS or Vaadin Image component
```

Reference fonts from the master stylesheet:

```css
@font-face {
    font-family: 'BrandFont';
    src: url('fonts/brand-font.woff2') format('woff2');
}
```

## Vaadin 24.x

Assets go in the theme folder and must be declared in `theme.json`:

```
src/main/frontend/themes/<theme-name>/
├── fonts/
│   └── brand-font.woff2
└── images/
    └── logo.svg
```

```json
{
  "assets": {
    "fonts/brand-font.woff2": "fonts/brand-font.woff2",
    "images/logo.svg": "images/logo.svg"
  }
}
```

**Related:** `theming/stylesheet-and-theme-loading.md` — how to load themes and
stylesheets per Vaadin version.