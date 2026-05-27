# Self-Editing Restrictions

When an action on a user's own record would leave the system inconsistent,
enforce the restriction in the service layer so it holds against programmatic
callers regardless of what the UI shows.

## Service Layer Guard

```java
@Transactional
public void deactivate(long key) {
    if (key == currentUser.getKey()) {
        throw new ValidationException("You cannot deactivate your own account.");
    }
    // ...
}
```

Mirror these in the UI (disabled button with tooltip) but enforce them in the
service layer — UI-only guards are bypassable.
