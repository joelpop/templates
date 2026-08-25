# Version-Specific Code Isolation

When code must branch on a library version, isolate the version-specific logic in
a single private method so a future upgrade touches one place.

```java
// Avoid — version-specific API scattered across every method
@Override
public ZoneId getBrowserTimezone() {
    var ui = UI.getCurrent();
    if (ui == null) {
        return ZoneId.systemDefault();
    }
    var tzId = ui.getPage().getExtendedClientDetails().getTimeZoneId();
    return tzId == null ? ZoneId.systemDefault() : ZoneId.of(tzId);
}

@Override
public boolean isTouchDevice() {
    var ui = UI.getCurrent();
    return ui != null && ui.getPage().getExtendedClientDetails().isTouchDevice();
}
```

```java
// Preferred — version-specific API isolated in one private method
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

private ExtendedClientDetails getDetails() {
    var ui = UI.getCurrent();
    return ui == null ? null : ui.getPage().getExtendedClientDetails();
}
```

Upgrading the library changes only `getDetails()` — the business logic in each
public method is untouched.