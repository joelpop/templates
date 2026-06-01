# ClientDetailsService UI-Module Implementation

When implementing `ClientDetailsService`, place `VaadinClientDetailsService` in the UI module — it only needs to implement `getBrowserTimezone()` since the conversion defaults are inherited.

## Vaadin ≥25

`Page.getExtendedClientDetails()` is synchronous — the details are available
immediately on the UI thread without a callback:

```java
@SpringComponent
public class VaadinClientDetailsService implements ClientDetailsService {

    @Override
    public ZoneId getBrowserTimezone() {
        var ui = UI.getCurrent();
        if (ui == null) {
            return ZoneId.systemDefault();
        }
        var tzId = ui.getPage().getExtendedClientDetails().getTimeZoneId();
        return tzId == null ? ZoneId.systemDefault() : ZoneId.of(tzId);
    }
}
```

## Vaadin 24.x

The API is asynchronous: `Page.retrieveExtendedClientDetails(callback)` fires
once the browser responds. Pre-cache the result in `VaadinSession` attributes
early in the session so subsequent calls — including those from Push threads —
can read it without a UI context.

**Trigger retrieval in the main layout:**

```java
@Override
protected void onAttach(AttachEvent event) {
    super.onAttach(event);
    event.getUI().getPage().retrieveExtendedClientDetails(details ->
            VaadinSession.getCurrent().setAttribute(ExtendedClientDetails.class, details));
}
```

**24.x-compatible implementation:**

```java
@SpringComponent
public class VaadinClientDetailsService implements ClientDetailsService {

    @Override
    public ZoneId getBrowserTimezone() {
        var session = VaadinSession.getCurrent();
        if (session != null) {
            var details = session.getAttribute(ExtendedClientDetails.class);
            if (details != null && details.getTimeZoneId() != null) {
                return ZoneId.of(details.getTimeZoneId());
            }
        }
        return ZoneId.systemDefault();
    }
}
```

Use `VaadinSession` attributes, not `@SessionScope` beans — `@SessionScope` is
backed by `HttpSession` and is inaccessible from Push threads. `VaadinSession`
works on both HTTP request and Push threads.

When `UI.getCurrent()` returns `null` — Push threads, background tasks, unit
tests — the implementation falls back to `ZoneId.systemDefault()`. Tests that
exercise timezone-dependent logic should mock `ClientDetailsService` directly.
