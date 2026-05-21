# Date and Time in the UI

Converting and formatting `Instant` values (stored as UTC) for display using the browser's zone and locale.

## Browser Client Details — Bridging the SoC Wall

Browser-level details (timezone, screen dimensions, locale) are available on the UI
thread but not in the service or persistence layers, which have no access to Vaadin
APIs (enforced by module boundaries). When service-layer code needs browser context —
most commonly timezone for `InstantMapper` (see `persistence.md`) — a bridging
service is required.

### The Pattern

1. **Service interface** (in the service module, no Vaadin dependency):

```java
public interface ClientDetailsService {
    ZoneId getBrowserTimezone();
    // add other browser details as needed
}
```

2. **UI-module implementation** (has Vaadin access):

The implementation obtains browser details via Vaadin's client-detail APIs and
caches them for the duration of the session. Cache in `VaadinSession` attributes —
not via `@SessionScope`, which is not compatible with Push threads.

On the UI thread, browser details may be directly accessible (no bridging needed).
The service exists primarily so that **non-UI code** (service implementations,
MapStruct mappers) can access browser context without a Vaadin dependency.

### Version-Specific Notes

> **Vaadin 24.x:** `UI.getCurrent().getPage().retrieveExtendedClientDetails(details -> ...)`
> — asynchronous callback pattern. The callback-based API makes caching in `VaadinSession`
> essential, since the details are not synchronously available.
>
> **Vaadin ≥25:** Consult the current documentation. The mechanism may differ — for
> instance, if browser details are synchronously available on the UI thread, the
> `VaadinSession` cache is primarily useful for the persistence side of the module
> boundary (where `UI.getCurrent()` is not available).

### Why Not @SessionScope?

`@SessionScope` beans are tied to `HttpSession`. Vaadin's Push threads don't run in
an HTTP request context, so `@SessionScope` beans are inaccessible from them.
`VaadinSession` attributes work on both HTTP request and Push threads.

### Fallback Behavior

Return a sensible default when details aren't yet available — typically
`ZoneId.systemDefault()`.

## `DateTimeUtil` — UI-Layer Instant Formatting

`DateTimeUtil` is the canonical conversion seam for the UI module. It converts
`Instant` values (stored as UTC) to display-ready forms using the browser's zone
and locale, resolved **per call** from the active `UI` context:

```java
public final class DateTimeUtil {

    private DateTimeUtil() { }

    public static LocalDateTime toLocalDateTime(Instant when) {
        return (when == null) ? null : when.atZone(currentZone()).toLocalDateTime();
    }

    public static Instant toInstant(LocalDateTime when) {
        return (when == null) ? null : when.atZone(currentZone()).toInstant();
    }

    public static String formatShort(Instant when) {
        return format(when, FormatStyle.SHORT);
    }

    public static String format(Instant when) {
        return format(when, FormatStyle.MEDIUM);
    }

    public static String formatLong(Instant when) {
        return format(when, FormatStyle.LONG);
    }

    private static String format(Instant when, FormatStyle style) {
        if (when == null) {
            return "";
        }
        return DateTimeFormatter
                .ofLocalizedDateTime(style)
                .withLocale(currentLocale())
                .format(when.atZone(currentZone()));
    }

    private static Locale currentLocale() {
        return Optional.ofNullable(UI.getCurrent())
                .map(UI::getLocale)
                .orElse(Locale.ROOT);
    }

    private static ZoneId currentZone() {
        var ui = UI.getCurrent();
        if (ui == null) {
            return ZoneOffset.UTC;
        }
        var details = ui.getPage().getExtendedClientDetails();
        if (details == null || details.getTimeZoneId() == null) {
            return ZoneOffset.UTC;
        }
        return ZoneId.of(details.getTimeZoneId());
    }
}
```

### When to use each method

| Method                                                  | Use for                                               |
|---------------------------------------------------------|-------------------------------------------------------|
| `toLocalDateTime(Instant)` + `toInstant(LocalDateTime)` | `Binder` conversions for `DateTimePicker` fields      |
| `formatShort(Instant)`                                  | Compact cells where space is constrained              |
| `format(Instant)`                                       | Default Grid column renderers, read-only text labels  |
| `formatLong(Instant)`                                   | Detail views where timezone clarity is important      |

### Why per-call resolution

Zone and locale are resolved fresh on each call rather than captured at construction
time. This means a user-preference feature (locale or timezone settable mid-session)
flows through every existing call site without any code changes — `UI.setLocale(...)`
updates the locale for all subsequent `format()` calls automatically.

### Relationship to `ClientDetailsService`

`DateTimeUtil` is for **UI code** that already has `UI.getCurrent()` available.
`ClientDetailsService` is for **non-UI code** (service implementations, MapStruct
`InstantMapper`) that sits behind a module boundary and cannot import Vaadin. Both
read from the same `ExtendedClientDetails` source; the bridging service exists only
to cross the module boundary.

### Fallback behaviour

When `UI.getCurrent()` returns `null` (background threads, unit tests) the zone
falls back to `ZoneOffset.UTC` and the locale to `Locale.ROOT`. Tests that exercise
time formatting can supply a real UI mock or accept the UTC/ROOT fallback explicitly.
