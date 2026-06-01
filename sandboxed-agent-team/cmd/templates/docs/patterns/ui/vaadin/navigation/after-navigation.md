# AfterNavigationObserver for Per-Navigation Actions

When actions are needed after every navigation — closing the mobile drawer, updating a breadcrumb — implement `AfterNavigationObserver` on `MainLayout`.

```java
public class MainLayout extends AppLayout implements AfterNavigationObserver {

    @Override
    public void afterNavigation(AfterNavigationEvent event) {
        // Close drawer on mobile after navigation
        setDrawerOpened(false);
    }
}
```
