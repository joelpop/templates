# ClientDetailsService

When service-layer code needs browser details (such as the user's timezone) that
are only available on the Vaadin UI thread, expose them through a
`ClientDetailsService` interface so the service module has no Vaadin dependency
and non-UI callers can access browser context through a plain Spring bean.

## Why It Exists

Vaadin's `Page.getExtendedClientDetails()` is accessible only on the UI thread.
Service and persistence layers have no Vaadin dependency (enforced by module
boundaries) and cannot call it directly. `ClientDetailsService` bridges the gap:
the interface lives in the service module, the implementation lives in the UI
module.

## The Interface

Declare the interface in the service module (no Vaadin import). The `default`
conversion methods derive directly from `getBrowserTimezone()` — any
implementation inherits them for free:

```java
package com.example.app.service.client;

import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;

/**
 * Browser-level client details exposed to callers outside the Vaadin
 * thread context.
 */
public interface ClientDetailsService {

    /**
     * The browser's reported timezone for the current user's session.
     * Falls back to {@link ZoneId#systemDefault()} when no Vaadin
     * context is available or the browser has not reported a zone.
     */
    ZoneId getBrowserTimezone();

    /** UTC → user's local time for display. {@code null} in, {@code null} out. */
    default LocalDateTime toBrowserTime(Instant utc) {
        return utc == null ? null : LocalDateTime.ofInstant(utc, getBrowserTimezone());
    }

    /** User's local time → UTC for storage. {@code null} in, {@code null} out. */
    default Instant toServerTime(LocalDateTime local) {
        return local == null ? null : local.atZone(getBrowserTimezone()).toInstant();
    }
}
```

Add methods to this interface only as new cross-boundary needs arise. Locale is
not needed here — it is always available directly on the UI thread via
`UI.getLocale()`.

## The UI-Module Implementation

`VaadinClientDetailsService` only needs to implement `getBrowserTimezone()` —
the conversion defaults are inherited:

```java
package com.example.app.ui.service;

import com.vaadin.flow.component.UI;
import com.vaadin.flow.spring.annotation.SpringComponent;
import java.time.ZoneId;
import com.example.app.service.client.ClientDetailsService;

/**
 * Vaadin-aware {@link ClientDetailsService}. Reads the browser zone
 * from {@code Page.getExtendedClientDetails()}; falls back to
 * {@link ZoneId#systemDefault()} when no Vaadin UI context is
 * available or the browser has not reported a zone.
 */
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

### Vaadin ≥25

`Page.getExtendedClientDetails()` is synchronous — the details are available
immediately on the UI thread without a callback. No caching is needed.

### Vaadin 24.x

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

**24.x-compatible `VaadinClientDetailsService`:**

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

Any call to `getBrowserTimezone()` that arrives before the async callback fires
falls back to `ZoneId.systemDefault()`.

### Fallback behaviour

When `UI.getCurrent()` returns `null` — Push threads, background tasks, unit
tests — the implementation falls back to `ZoneId.systemDefault()`. Tests that
exercise timezone-dependent logic should mock `ClientDetailsService` directly
rather than rely on the fallback.

## MapStruct Integration

MapStruct mappers that need `Instant` → `LocalDateTime` conversion declare
`ClientDetailsService` in their `uses` list. MapStruct injects it as a
constructor dependency and routes `Instant` → `LocalDateTime` field conversions
through `toBrowserTime` automatically — no explicit `@Mapping` is needed for
those fields:

```java
@Mapper(
        componentModel = MappingConstants.ComponentModel.SPRING,
        injectionStrategy = InjectionStrategy.CONSTRUCTOR,
        uses = {AuditMapper.class, ClientDetailsService.class})
public interface EquipmentMapper {
    EquipmentDetail toDetail(EquipmentDetailProjection projection);
}
```
