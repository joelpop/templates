# Date and Time Display

When displaying date or time values in a Vaadin view, format them through
`DateTimeUtil` with the locale resolved per-call from `UI.getCurrent()` so
rendered values reflect the user's locale without requiring explicit locale
management at every call site.

## `DateTimeUtil` — Locale-aware Display Formatting

`DateTimeUtil` formats local date/time types for display using the user's locale
resolved per-call from the active `UI` context:

```java
public final class DateTimeUtil {

    private DateTimeUtil() { }

    public static String format(LocalDateTime when) {
        if (when == null) return "";
        return DateTimeFormatter
                .ofLocalizedDateTime(FormatStyle.MEDIUM)
                .withLocale(currentLocale())
                .format(when);
    }

    public static String format(LocalDate when) {
        if (when == null) return "";
        return DateTimeFormatter
                .ofLocalizedDate(FormatStyle.MEDIUM)
                .withLocale(currentLocale())
                .format(when);
    }

    public static String format(LocalTime when) {
        if (when == null) return "";
        return DateTimeFormatter
                .ofLocalizedTime(FormatStyle.MEDIUM)
                .withLocale(currentLocale())
                .format(when);
    }

    private static Locale currentLocale() {
        return Optional.ofNullable(UI.getCurrent())
                .map(UI::getLocale)
                .orElse(Locale.ROOT);
    }
}
```

### When to use each method

| Method                  | Use for                                                 |
|-------------------------|---------------------------------------------------------|
| `format(LocalDateTime)` | Grid columns and labels rendering a date-time value     |
| `format(LocalDate)`     | Grid columns and labels rendering a date-only value     |
| `format(LocalTime)`     | Grid columns and labels rendering a time-only value     |

`DateTimePicker`, `DatePicker`, and `TimePicker` bind directly to their
respective local types — `DateTimeUtil` is only needed when rendering a value
for display, not for `Binder` bindings.

`Instant` values from the database arrive at the UI already converted to the
user's local type by `InstantMapper` in the service layer. See
`docs/patterns/ui/vaadin/client-details-service.md`.

### Why per-call resolution

Locale is resolved from `UI.getCurrent()` on each call rather than captured at
construction time. A locale preference change (`UI.setLocale(...)`) flows through
every existing `format()` call automatically without re-wiring.

### Fallback behaviour

When `UI.getCurrent()` returns `null` (background threads, unit tests) the locale
falls back to `Locale.ROOT`.