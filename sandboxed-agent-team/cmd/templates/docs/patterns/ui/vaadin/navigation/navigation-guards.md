# Navigation Guards

When access control beyond `@RolesAllowed` is needed, implement `BeforeEnterObserver` to perform programmatic checks at navigation time.

```java
@Route("items/:key/edit")
@RolesAllowed(UserRole.ROLE_ADMIN)
public class ItemEditView extends Composite<VerticalLayout> implements BeforeEnterObserver {

    @Override
    public void beforeEnter(BeforeEnterEvent event) {
        var key = event.getRouteParameters().get("key")
            .map(Long::parseLong)
            .orElse(null);
        if (key == null || !itemService.exists(key)) {
            event.rerouteToError(NotFoundException.class);
        }
    }
}
```
