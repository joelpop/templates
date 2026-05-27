# XSS Prevention

When rendering user-supplied data in the UI, rely on Vaadin's component layer
to escape output so no additional escaping is required for text shown via
standard Vaadin components.

Vaadin's component layer escapes output by default, preventing XSS from
user-supplied data rendered in the UI. No additional escaping is required for
text shown via standard Vaadin components.
