# Application Icon Catalog

When adding icons to a Vaadin application, define an intent-named `AppIcon` enum catalog
rather than using `VaadinIcon` or third-party icon constants directly at call sites so
icons can be swapped centrally and call sites are decoupled from the backing icon library.

Two constants can share the same visual today and diverge independently later — naming by
intent (`EDIT`, `IMPERSONATE`) rather than by appearance (`PERSON_ICON`) makes that
change a one-line update per constant with no effect on call sites.

```java
// Avoid — icon names describe the visual object, not the application intent
editButton.setIcon(VaadinIcon.PENCIL.create());
sendButton.setIcon(VaadinIcon.ENVELOPE.create());
```

```java
// Preferred — intent at call site; visual choice is one line in AppIcon
editButton.setIcon(AppIcon.EDIT.create());
sendButton.setIcon(AppIcon.SEND.create());
```

**Related:** `ui/vaadin/recipes/app-icon.md` — full implementation recipe for `AppIcon`,
adapter enums, and custom iconsets.
