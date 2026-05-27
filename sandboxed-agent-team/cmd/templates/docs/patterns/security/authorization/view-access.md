# View Access Annotations

When creating a `@Route`-annotated view, add exactly one of the four Jakarta
Security annotations so `AnnotatedViewAccessChecker` has an unambiguous access
rule for every route and no view is accidentally left open or closed.

## The Four Annotations

| Annotation | Meaning |
|------------|---------|
| `@AnonymousAllowed` | Public — no authentication required |
| `@PermitAll`        | Any authenticated user              |
| `@RolesAllowed(...)`| Only users with one of the listed roles |
| `@DenyAll`          | No one (used to explicitly block a view) |

A view without an access annotation is a security defect.

## Do Not Use `@PreAuthorize`

Do not use `@PreAuthorize` (Spring Security method security) on Vaadin views.
Access control is enforced at the view level via Jakarta Security annotations
and Vaadin's `AnnotatedViewAccessChecker`. `@PreAuthorize` adds confusion and
may not behave as expected with Vaadin's navigation model.
