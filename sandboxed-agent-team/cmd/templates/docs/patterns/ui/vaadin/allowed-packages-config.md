# Vaadin Allowed Packages Configuration

When configuring a Vaadin application, declare `vaadin.allowed-packages` in
`application.properties` to include the application's own UI package and all
Vaadin component package prefixes so all components resolve at runtime.

```properties
vaadin.allowed-packages=com.vaadin,org.vaadin,com.example.app.ui
```

The application's UI package must be listed alongside the Vaadin packages.
When adding a new Vaadin add-on, add its root package prefix to this property.
Without it, components in that package may fail to render.
