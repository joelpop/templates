# Recipe: `@ViewIcon` — Type-Safe View Icon Declaration

A runtime annotation that declares a view's icon as a type-safe enum
value — one source of truth, no string literals, no knowledge of icon
collection names required.

## What this produces

- A `@ViewIcon` annotation targeting view classes, holding an
  `IconFactory`-implementing enum value (`AppIcon` when the project uses
  the full catalog; any other `IconFactory` enum for simpler projects).
- A reflection pattern for reading it in navigation components, view
  headers, or any other surface that needs the view's icon.

## Dependencies

- Vaadin 24+ (`IconFactory`, `Icon`, `@Menu`, `MenuConfiguration`).
- The [app-icon recipe](app-icon.md) if using `AppIcon` as the value
  type (recommended for projects with multiple icon libraries or a
  curated catalog). Projects that only use `VaadinIcon` do not need
  `AppIcon`.

## The problem

Vaadin's `@Menu` annotation accepts the view's icon as a string:

```java
@Menu(order = 1, icon = "vaadin:group")
```

This is fragile in three ways:
1. The string encodes two pieces of knowledge — the collection name
   (`vaadin`) and the icon name (`group`) — that the developer must look
   up and remember.
2. It silently breaks if the icon collection or name changes; the
   compiler never sees it.
3. It is a second declaration of what `AppIcon` already knows. If
   `AppIcon.USER_MANAGEMENT` delegates to `VaadinIcon.USERS`, the
   developer must independently discover and maintain that `"vaadin:users"`
   string in every view that uses this icon.

`@ViewIcon` makes the icon declaration once, in the `AppIcon` enum, and
reads it anywhere via reflection.

## Step 1 — Define `@ViewIcon`

The annotation lives in the UI module (`{app}-ui`), in whichever package
holds your icon types.

Java annotation values must be a concrete type — primitives, `String`,
`Class`, another annotation, or an **enum**. `IconFactory` is an
interface and cannot be the value type directly. Instead, use whichever
`IconFactory`-implementing enum the project already defines:

**With `AppIcon`** (projects using the full icon catalog):

```java
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface ViewIcon {
    AppIcon value();
}
```

**With `VaadinIcon`** (projects that only use Vaadin's built-in icons and
do not want to define an `AppIcon` catalog):

```java
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
public @interface ViewIcon {
    VaadinIcon value();
}
```

Any other enum implementing `IconFactory` works the same way. The value
type is the project's choice; the consuming code treats it as
`IconFactory` via the interface regardless.

`RUNTIME` retention is required — the annotation is read reflectively at
runtime by navigation and layout components.

## Step 2 — Annotate views

Place `@ViewIcon` on the view class alongside `@Route` and `@Menu`:

```java
@Route("users")
@Menu(order = 2, title = "Users")
@RolesAllowed(UserRole.ROLE_ADMIN)
@ViewIcon(AppIcon.USER_MANAGEMENT)
public class UserManagementView extends BaseView { ... }
```

`@Menu` is still present for Vaadin's auto-discovered sidebar
navigation; `@ViewIcon` provides the icon declaration for any surface
that reads it programmatically. If your project uses a fully custom
navigation component (not `MenuConfiguration`), `@Menu(icon = ...)` can
be omitted entirely — `@ViewIcon` carries the icon.

## Step 3 — Read `@ViewIcon` at runtime

Any component that needs the view's icon reads the annotation via
reflection. The null check covers views without an icon:

```java
private Icon iconForView(Class<?> viewClass) {
    ViewIcon annotation = viewClass.getAnnotation(ViewIcon.class);
    return annotation != null ? annotation.value().create() : null;
}
```

`annotation.value()` returns the enum constant. Because every
`IconFactory`-implementing enum has `create()`, calling `.create()` is
valid regardless of which enum type the project chose as the value type.

Typical call sites:

```java
// Navigation item in a custom sidebar
Icon icon = iconForView(viewClass);
if (icon != null) {
    sideNavItem.setPrefixComponent(icon);
}

// In-view subheader (e.g., in BaseView or a page header component)
ViewIcon annotation = getClass().getAnnotation(ViewIcon.class);
if (annotation != null) {
    headerIcon.add(annotation.value().create());
}
```

## Decisions this recipe imposes

- **The value type is the project's choice — any `IconFactory` enum
  works.** Use `AppIcon` when the project has a full icon catalog (the
  usual choice); use `VaadinIcon` for simpler projects that only need
  built-in icons. The consumer always calls `.create()` through the
  `IconFactory` contract regardless.
- **`@ViewIcon` is the icon source of truth.** Changing the delegate on
  an `AppIcon` constant updates every surface that reads `@ViewIcon`
  automatically — drawer nav, subheader, breadcrumb, wherever else the
  annotation is consumed.
- **`@Menu(icon = ...)` is a secondary concern.** Vaadin's
  `MenuConfiguration` reads `@Menu` to build the auto-discovered
  sidebar. If you keep both, keep them consistent. If you build a custom
  nav component that reads `@ViewIcon` directly, you can drop `@Menu(icon
  = ...)` and eliminate the string literal.
- **Views without a meaningful icon omit `@ViewIcon`.** The null check in
  Step 3 is the right approach; a missing annotation is not an error.
- **The annotation is in the UI module.** Its value type creates Vaadin
  components; both belong where Vaadin dependencies live.

## What to verify

- A view annotated with `@ViewIcon(AppIcon.DASHBOARD)` renders the
  correct icon in the navigation and in any in-view header surface.
- Changing `AppIcon.DASHBOARD`'s delegate to a different `IconFactory`
  updates the icon on every surface without modifying any view class or
  nav component.
- A view without `@ViewIcon` does not cause a `NullPointerException` in
  the reading component.
- A typo in `@ViewIcon(AppIcon.TYPO)` is a compile error — not a silent
  runtime miss.

## Related

- [app-icon.md](app-icon.md) — `AppIcon` and `IconFactory` adapters;
  the enum that `@ViewIcon` carries as its value.
- `docs/patterns/ui/navigation.md` — `@Menu` and sidebar navigation;
  the string-literal icon pattern that `@ViewIcon` replaces.
