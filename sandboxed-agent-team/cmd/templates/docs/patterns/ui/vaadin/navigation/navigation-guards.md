# Navigation Guards

When a view requires entry conditions beyond `@RolesAllowed` — such as
validating a route parameter or confirming the current user is permitted to
access a specific record — implement `BeforeEnterObserver` so the check runs
before the view renders and unauthorized or invalid navigation is redirected
cleanly.

```java
@Route("items/:key/edit")
@RolesAllowed(UserRole.ROLE_USER)
public class ItemEditView extends Composite<VerticalLayout> implements BeforeEnterObserver {

    @Override
    public void beforeEnter(BeforeEnterEvent event) {
        var key = event.getRouteParameters().get("key")
            .map(Long::parseLong)
            .orElse(null);
        // Route parameter missing or record does not exist
        if (key == null || !itemService.exists(key)) {
            event.rerouteToError(NotFoundException.class);
            return;
        }
        // Record exists but is not accessible to the current user
        if (!itemService.isOwnedBy(key, currentUserService.currentUserId())) {
            event.rerouteToError(AccessDeniedException.class); // logged as access denial; view returns 404
        }
    }
}
```

**Related:** `security/authorization/view-access.md` — static role annotation via `@RolesAllowed`;
`security/authorization/self-editing.md` — service-layer ownership enforcement;
`ui/vaadin/error-views/access-denied.md` — why access violations return 404 not 403.
