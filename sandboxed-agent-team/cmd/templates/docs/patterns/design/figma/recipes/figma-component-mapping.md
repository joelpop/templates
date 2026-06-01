---
name: figma-component-mapping
description: Reference table mapping common Figma component and layout names to their Vaadin Flow equivalents, including variants and generic HTML elements.
---

# Common Figma → Vaadin Component Mappings

When selecting a Vaadin component for a Figma element, use this reference table as a
starting point — always verify against Vaadin documentation for the specific version
before implementing.

## Vaadin Components
- `Button` → `Button.class`
- `Button (tertiary, icon-only)` → `Button` + `ButtonVariant.LUMO_TERTIARY` + `ButtonVariant.LUMO_ICON`
- `Text Field` → `TextField.class`
- `Grid` → `Grid.class`
- `Message List` → `MessageList.class`
- `Avatar` → `Avatar.class`
- `Card` → `Card.class` (since v24.8)

## Vaadin Layouts
- Vertical auto layout → `VerticalLayout`
- Vertical auto layout with wrapping → `VerticalLayout` + `setWrap(true)`
- Horizontal auto layout → `HorizontalLayout`
- Layout → `FlexLayout` + `addClassNames(LumoUtility.FlexDirection.ROW, LumoUtility.AlignItems.BASELINE)`
- Master-Detail Layout → `MasterDetailLayout.class` (with feature flag)
- Form → `FormLayout`

## Generic HTML Elements
- Text layer → `com.vaadin.flow.component.html.Span`
- Heading 3 → `com.vaadin.flow.component.html.H3`
