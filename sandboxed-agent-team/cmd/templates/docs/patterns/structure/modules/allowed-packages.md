# vaadin.allowed-packages

When adding a Vaadin add-on, add its root package to `vaadin.allowed-packages` in
`application.properties` so add-on components render correctly — without this entry,
add-on components may silently fail to render.

```properties
vaadin.allowed-packages=com.vaadin,org.vaadin
```

Add the root package of each Vaadin add-on to this list when it is introduced.
