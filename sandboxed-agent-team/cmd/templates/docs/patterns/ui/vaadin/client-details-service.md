# ClientDetailsService

When any code needs browser details — whether in the service layer or the UI
layer — access them through a `ClientDetailsService` interface so callers get
a consistent API regardless of Vaadin version and the service module has no
Vaadin dependency.

## Why It Exists

`ExtendedClientDetails` has a different retrieval model in Vaadin 24 (async
callback, cached per-UI) and Vaadin 25 (synchronous). `ClientDetailsService`
encapsulates that difference: the interface lives in the service module, the
implementation in the UI module, and callers never deal with `UI.getCurrent()`
or `ComponentUtil` directly.

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

    /**
     * Whether the browser supports touch events.
     * Falls back to {@code false} when no Vaadin context is available
     * or the browser has not reported touch support.
     */
    boolean isTouchDevice();

    /**
     * The browser's reported inner viewport width in pixels.
     * Falls back to {@code 0} when no Vaadin context is available
     * or the browser has not reported a width.
     */
    int getWindowInnerWidth();

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

Add methods to this interface as new browser-context needs arise. Locale is
not needed here — it is always available directly on the UI thread via
`UI.getLocale()`.

**Related:** `client-details-impl.md` — Vaadin implementation of this interface;
`client-details-mapstruct.md` — wiring `ClientDetailsService` into MapStruct mappers.
