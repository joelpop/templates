# UI-Layer Date and Time Utilities

When formatting or processing date or time values in the UI layer,
place the logic in `DateTimeUtil` so UI-side temporal logic is not
duplicated and callers have a clear API for it.

## `DateTimeUtil`

`DateTimeUtil` provides UI-side date/time operations using the user's locale
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

    public static String formatElapsed(LocalDateTime since) {
        if (since == null) {
            return "";
        }
        var now = LocalDateTime.now();
        var minutes = ChronoUnit.MINUTES.between(since, now);
        if (minutes < 1) {
            return "just now";
        }
        if (minutes < 60) {
            return minutes + " min ago";
        }
        var hours = ChronoUnit.HOURS.between(since, now);
        if (hours < 24) {
            return hours + " hr ago";
        }
        var days = ChronoUnit.DAYS.between(since, now);
        if (days == 1) {
            return "yesterday";
        }
        var months = ChronoUnit.MONTHS.between(since, now);
        if (months < 1) {
            return days + " days ago";
        }
        if (months < 12) {
            return months + " months ago";
        }
        var years = ChronoUnit.YEARS.between(since, now);
        return years + " years ago";
    }

    private static Locale currentLocale() {
        return Optional.ofNullable(UI.getCurrent())
                .map(UI::getLocale)
                .orElse(Locale.ROOT);
    }
}
```

### When to use each method

| Method                         | Use for                                             |
|--------------------------------|-----------------------------------------------------|
| `format(LocalDateTime)`        | Grid columns and labels rendering a date-time value |
| `format(LocalDate)`            | Grid columns and labels rendering a date-only value |
| `format(LocalTime)`            | Grid columns and labels rendering a time-only value |
| `formatElapsed(LocalDateTime)` | "Last updated" and other elapsed-time labels        |

`DateTimePicker`, `DatePicker`, and `TimePicker` bind directly to their
respective local types — `DateTimeUtil` is only needed when rendering a value
for display, not for `Binder` bindings.

`Instant` values from the database arrive at the UI already converted to the
user's local type by `ClientDetailsService` in the mapper layer. See
`docs/patterns/ui/vaadin/client-details-mapstruct.md`.

### Why per-call resolution

Locale is resolved from `UI.getCurrent()` on each call rather than captured at
construction time. A locale preference change (`UI.setLocale(...)`) flows through
every existing `format()` call automatically without re-wiring.

### Fallback behavior

When `UI.getCurrent()` returns `null` (background threads, unit tests) the locale
falls back to `Locale.ROOT`.
