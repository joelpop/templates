# Vaadin Responsive Tools Only

When achieving responsive breakpoints, use Vaadin's built-in responsive APIs and `LumoUtility` classes — not Bootstrap, Tailwind, or any CSS framework grid system.

```java
// Preferred — Vaadin responsive step on FormLayout
form.setResponsiveSteps(
    new FormLayout.ResponsiveStep("0", 1),       // 1 column below 600px
    new FormLayout.ResponsiveStep("600px", 2),   // 2 columns at 600px+
    new FormLayout.ResponsiveStep("1024px", 3)   // 3 columns at 1024px+
);
```
