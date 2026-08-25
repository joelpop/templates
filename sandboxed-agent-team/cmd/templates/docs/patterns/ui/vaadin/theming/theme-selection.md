# Theme and Utility Class Trade-offs

When a project's theme is being chosen, present the trade-offs between Lumo+LumoUtility,
Tailwind CSS, and hand-rolled approaches so the human can make an informed decision and
the project is configured to get the most from the chosen utility strategy.

| Approach                    | Theme        | Utility Classes              |
|-----------------------------|--------------|------------------------------|
| Lumo + LumoUtility          | Lumo         | Built-in, first-class        |
| Any theme + Tailwind CSS    | Any          | Experimental (feature flag)  |
| Any theme + hand-rolled CSS | Any          | None — write your own        |

`LumoUtility` and Tailwind CSS are mutually exclusive and cannot be reliably used together.

## Lumo + LumoUtility

Lumo is the only Vaadin theme with a built-in utility class library. `LumoUtility`
covers spacing, typography, color, display, flexbox, and responsive breakpoints —
enabling the complete styling progression (`theming/component-variants.md`) without
writing custom CSS.

Lock-in: the application is tied to Lumo design tokens and `LumoUtility` class names.

## Any Theme + Tailwind CSS

Tailwind CSS works with any Vaadin theme. Tailwind's utility classes have limited effect
on complex Vaadin components due to their nested shadow DOM structure — they work best on
native HTML elements and simple layout components (`HorizontalLayout`, `VerticalLayout`).

Lock-in: the application is tied to Tailwind class names; the feature is experimental.

See `theming/tailwind.md` for setup and usage patterns.

**Related:** `theming/tailwind.md` — Tailwind setup and usage;
`theming/lumo-utility.md` — LumoUtility usage;
`theming/component-variants.md` — styling progression;
`theming/stylesheet-and-theme-loading.md` — how to load themes.
