# CSRF for Non-Vaadin Endpoints

When adding non-Vaadin endpoints (actuator, custom REST) to a Vaadin application, implement CSRF protection independently — `VaadinWebSecurity` configures CSRF only for Vaadin's filter chain.
