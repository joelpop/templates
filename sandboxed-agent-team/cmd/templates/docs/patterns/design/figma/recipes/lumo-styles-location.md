---
name: lumo-styles-location
description: Where to write the styles.css file for Lumo CSS variable overrides, by Vaadin version.
---

# Styles.css Location by Vaadin Version

When generating CSS variable declarations for a Vaadin application, write them to
the path determined by the Vaadin version so the theme stylesheet is loaded correctly.

**Vaadin 24:**
- Path: `frontend/themes/[app]/styles.css`
- Theme folder is active by default
- Uses `@Theme` annotation on AppShellConfigurator

**Vaadin 25:**
- Path: `src/main/resources/META-INF/resources/styles.css`
- Uses `@StyleSheet("styles.css")` annotation on AppShellConfigurator
- The Vaadin 24 theme folder approach (`frontend/themes`) is deprecated but still supported using the `themeComponentStyles` feature flag
