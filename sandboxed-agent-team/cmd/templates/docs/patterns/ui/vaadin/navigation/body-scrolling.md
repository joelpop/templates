# Body Scrolling in AppLayout

When views do not scroll correctly on mobile, add `setSizeFull()` or explicit height constraints — `AppLayout` manages scrolling via the content area, not `<body>`.

In standard `AppLayout`, scrolling is managed by the content area, not `<body>`. Views
must use `setSizeFull()` or explicit height constraints — otherwise the content area may
not scroll correctly on mobile.
