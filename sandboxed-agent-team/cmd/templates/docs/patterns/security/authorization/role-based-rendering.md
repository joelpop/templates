# Role-Based Rendering — Do Not Generate

When a component or navigation item is not authorized for the current user's
role, never construct it, never add it to the layout, and never leave it in the
DOM — so there is no artifact the user can discover or re-enable via the browser
inspector.

`setVisible(false)` and CSS concealment are **not** acceptable: the underlying
component is still present and (if the server considers it enabled) interactive.

For the full three-mode rubric (do not generate / hide / disable),
layout-preservation guidance, and the Vaadin-server-state authority rule, see:

- `docs/patterns/ui/vaadin/components.md` → "Conditional Component Rendering —
  Do Not Generate vs. Hide vs. Disable"
- `docs/patterns/ui/vaadin/navigation.md` → "Conditional Navigation Rendering"
