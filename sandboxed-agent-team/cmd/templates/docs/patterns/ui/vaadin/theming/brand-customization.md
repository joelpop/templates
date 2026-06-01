# Brand Customization via CSS Custom Properties

When customizing the application's brand colors or typography, override Lumo CSS custom properties in a dedicated stylesheet — do not modify Lumo component internals directly.

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
