# AfterNavigationObserver for Per-Navigation Actions

When actions are needed after every navigation — closing the mobile drawer, updating a breadcrumb — implement `AfterNavigationObserver` on `MainLayout`.

```java
public class MainLayout extends AppLayout implements AfterNavigationObserver {

    private final ClientDetailsService clientDetailsService;

    public MainLayout(ClientDetailsService clientDetailsService) {
        this.clientDetailsService = clientDetailsService;
    }

    @Override
    public void afterNavigation(AfterNavigationEvent event) {
        if (clientDetailsService.isTouchDevice()) {
            setDrawerOpened(false);
        }
    }
}
```

**Related:** `client-details-service.md` — `ClientDetailsService` interface including `isTouchDevice()`;
`client-details-impl.md` — version-specific implementation.
