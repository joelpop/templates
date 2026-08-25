# AppLayout Scrolling Behavior

When building the main layout, call `setSizeFull()` on `AppLayout` so the content
area scrolls rather than `<body>` — this keeps the navbar and drawer fixed and lets
content components like `Grid` or `Scroller` fill the available height.

```java
// Avoid — <body> scrolls; navbar collapses on scroll; Grid cannot fill content area
@Layout
@PermitAll
public class MainLayout extends AppLayout {

    public MainLayout(/* dependencies */) { }
}
```

```java
// Preferred — content area scrolls; navbar stays fixed; Grid fills content area
@Layout
@PermitAll
public class MainLayout extends AppLayout {

    public MainLayout(/* dependencies */) {
        setSizeFull();
    }
}
```

Vaadin's default CSS already sets 100% height on `<html>` and `<body>`, so
`setSizeFull()` on the layout is sufficient.

**Related:** [AppLayout — Scrolling Behavior](https://vaadin.com/docs/latest/components/app-layout#scrolling-behavior)
