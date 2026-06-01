# Lumo Is the Only Theme

The application must use the Vaadin Lumo theme — the Aura theme is not used.

```java
@Theme(Lumo.class)
// or via @StyleSheet(Lumo.STYLESHEET) + @StyleSheet(Lumo.UTILITY_STYLESHEET)
public class Application implements AppShellConfigurator { ... }
```

Lumo CSS custom properties are applied globally and define the visual foundation for all
components.
