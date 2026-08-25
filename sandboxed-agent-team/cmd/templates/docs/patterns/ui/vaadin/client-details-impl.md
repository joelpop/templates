# ClientDetailsService UI-Module Implementation

When implementing `ClientDetailsService`, place `VaadinClientDetailsService` in the UI module — it has access to Vaadin APIs (`UI`, `Page`, `ComponentUtil`) that only the UI module may depend on.

## Vaadin ≥25

`Page.getExtendedClientDetails()` is synchronous — the details are available
immediately on the UI thread without a callback:

```java
@SpringComponent
public class VaadinClientDetailsService implements ClientDetailsService {

    @Override
    public ZoneId getBrowserTimezone() {
        var details = getDetails();
        if (details == null) {
            return ZoneId.systemDefault();
        }
        var tzId = details.getTimeZoneId();
        return tzId == null ? ZoneId.systemDefault() : ZoneId.of(tzId);
    }

    @Override
    public boolean isTouchDevice() {
        var details = getDetails();
        return details != null && details.isTouchDevice();
    }

    @Override
    public int getWindowInnerWidth() {
        var details = getDetails();
        return details != null ? details.getWindowInnerWidth() : 0;
    }

    private ExtendedClientDetails getDetails() {
        var ui = UI.getCurrent();
        return ui == null ? null : ui.getPage().getExtendedClientDetails();
    }
}
```

## Vaadin 24.x

The API is asynchronous: `Page.retrieveExtendedClientDetails(callback)` fires
once the browser responds. Cache the result per-UI using `ComponentUtil` so
each browser tab has its own details and subsequent calls can read them without
a round-trip.

Trigger retrieval via a `VaadinServiceInitListener` so it fires for every new
UI regardless of which layout or view the user lands on first:

```java
@SpringComponent
public class ExtendedClientDetailsInitializer implements VaadinServiceInitListener {

    @Override
    public void serviceInit(ServiceInitEvent event) {
        event.getSource().addUIInitListener(uiEvent -> {
            var ui = uiEvent.getUI();
            ui.getPage().retrieveExtendedClientDetails(details ->
                    ComponentUtil.setData(ui, ExtendedClientDetails.class, details));
        });
    }
}
```

**24.x-compatible implementation:**

```java
@SpringComponent
public class VaadinClientDetailsService implements ClientDetailsService {

    @Override
    public ZoneId getBrowserTimezone() {
        var details = getDetails();
        if (details != null && details.getTimeZoneId() != null) {
            return ZoneId.of(details.getTimeZoneId());
        }
        return ZoneId.systemDefault();
    }

    @Override
    public boolean isTouchDevice() {
        var details = getDetails();
        return details != null && details.isTouchDevice();
    }

    @Override
    public int getWindowInnerWidth() {
        var details = getDetails();
        return details != null ? details.getWindowInnerWidth() : 0;
    }

    private ExtendedClientDetails getDetails() {
        var ui = UI.getCurrent();
        return ui == null ? null : ComponentUtil.getData(ui, ExtendedClientDetails.class);
    }
}
```

When `UI.getCurrent()` returns `null` — Push threads, background tasks, unit
tests — the implementation falls back to `ZoneId.systemDefault()`. Tests that
exercise timezone-dependent logic should mock `ClientDetailsService` directly.

**Related:** `client-details-service.md` — the interface this implements;
`client-details-mapstruct.md` — wiring `ClientDetailsService` into MapStruct mappers.
