---
name: figma-quality-standards
description: Code quality rules for Figma-to-Vaadin implementations — use component APIs not getElement(), addThemeVariants() for sizes, LumoUtility for spacing, and know when to ask for clarification.
---

# Figma-to-Vaadin Quality Standards

When implementing a Figma design in Vaadin, use Vaadin component APIs (not
`getElement()` or inline styles), use `addThemeVariants()` for sizes, and use
LumoUtility for spacing so generated code follows framework conventions.

## Accuracy Over Speed
- Read all metadata through Figma MCP before implementing
- Verify component choice against documentation

## Semantic Correctness
- Use proper Vaadin components, not generic HTML
- Follow Vaadin component APIs and patterns
- Preserve component semantics and accessibility

## Code Style
- Avoid tiny wrapper methods that only delegate without adding logic — inline or
  generalize with parameters.

## Follow Vaadin Patterns
```java
// Proper way to configure components is to use component API's when available
textField.setReadOnly(true);

// ❌ INCORRECT ways to configure components is to use getComponent()
textField.getElement().setAttribute("readonly", "true");
button.getElement().getStyle().set("background", "transparent");

// Proper way to set component theme variants
button.addThemeVariants(ButtonVariant.LUMO_TERTIARY);

// Proper way to set styles using Lumo Utility classes
layout.addClassNames(LumoUtility.Padding.Horizontal.LARGE, LumoUtility.Padding.Vertical.MEDIUM);

// Proper ways to set sizing
layout.setSizeFull();
layout.setWidth("600px");
layout.setHeight("50%");

// ❌ INCORRECT way to set sizing
layout.getStyle().set("width", "600px");

// Proper way to set space around layout, always use padding
layout.addClassName(LumoUtility.Padding.Bottom.MEDIUM);

// ❌ INCORRECT way to set space around layout, never use margin
layout.getStyle().set("margin-bottom", "36px");

// Proper way to set component size is to first use available size variants
avatar.addThemeVariants(AvatarVariant.LUMO_LARGE);

// ❌ INCORRECT way to set component size
avatar.getStyle().set("--vaadin-avatar-size", "48px");

// Proper accessibility
iconButton.setAriaLabel("Close");

// Proper way to set input field label if component implements HasLabel
input.setLabel("Label");

// ❌ INCORRECT way to set input field label
Span label = new Span("Label");
VerticalLayout.add(label, input);

// Proper way to set border
layout.addClassNames(LumoUtility.Border.TOP, LumoUtility.BorderColor.CONTRAST_10);

// ❌ INCORRECT way to set border
layout.getStyle().set("border-top", "1px solid var(--lumo-contrast-10pct)");
```

## When to Ask for Clarification

Ask when multiple Vaadin components could fit the design, the Figma name doesn't
clearly map to a Vaadin component, or you're uncertain about variants, styling, or
interaction patterns.

**Form**: "Should this be a [ComponentA] or [ComponentB]? The Figma shows [description]"
