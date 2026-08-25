# Brand Customization via CSS Custom Properties

When customizing the application's brand colors or typography, override Lumo CSS custom properties in a dedicated stylesheet so changes propagate across all Vaadin components and theme CSS stays separate from component CSS.

```css
/* styles.css — master entry point */
@import "theme.css";
```

```css
/* theme.css — Lumo custom property overrides */
html {
    --lumo-primary-color:       light-dark(#1a73e8, #7baaf7);
    --lumo-primary-text-color:  light-dark(#1a73e8, #7baaf7);
    --lumo-primary-color-10pct: light-dark(#1a73e810, #7baaf710);
    --lumo-primary-color-50pct: light-dark(#1a73e880, #7baaf780);
    --lumo-font-family:         'Inter', sans-serif;
    --lumo-border-radius-m:     6px;
}
```

Color overrides replace Lumo's built-in light/dark adaptation — pair each with
`light-dark(lightValue, darkValue)` when the application uses an adaptive color scheme.
Non-color overrides (font family, border radius) do not need paired values.

Common overridable properties:

| Property                      | Controls                                           |
|-------------------------------|----------------------------------------------------|
| `--lumo-primary-color`        | Primary action color (buttons, links, focus rings) |
| `--lumo-primary-text-color`   | Text color for primary actions                     |
| `--lumo-font-family`          | Application font family                            |
| `--lumo-border-radius-m`      | Default component corner radius                    |
| `--lumo-base-color`           | Page background                                    |
| `--lumo-body-text-color`      | Default text color                                 |
| `--lumo-header-text-color`    | Heading text color                                 |
| `--lumo-secondary-text-color` | Subdued text (labels, captions)                    |
| `--lumo-success-color`        | Success state color                                |
| `--lumo-error-color`          | Error state color                                  |
| `--lumo-warning-color`        | Warning state color                                |
