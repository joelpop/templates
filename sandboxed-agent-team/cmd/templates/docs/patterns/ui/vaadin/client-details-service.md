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
