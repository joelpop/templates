# Touch-Optimized Navigation on Mobile

On mobile (< 600px), use a bottom tab bar for the most frequently accessed views and an accordion for secondary navigation within a view — nested sidebar items are difficult to tap accurately on small screens.

Vaadin's `Tabs` component with `HORIZONTAL` orientation and bottom positioning serves
as the bottom tab bar. For secondary navigation within a mobile view, use an
accordion rather than nested sidebar items.
