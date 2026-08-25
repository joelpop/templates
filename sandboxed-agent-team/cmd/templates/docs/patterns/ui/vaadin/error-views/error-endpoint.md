# Suppress Spring Boot's Default Error Endpoint

When deploying a Vaadin application with custom error views, disable Spring
Boot's default `/error` endpoint so all error conditions are routed through
Vaadin's `HasErrorParameter` mechanism rather than the whitelabel page.

```properties
server.error.whitelabel.enabled=false
```