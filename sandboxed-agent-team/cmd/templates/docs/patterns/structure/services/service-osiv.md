# OSIV Disabled — Service Layer Must Load All Data

When OSIV is disabled (`spring.jpa.open-in-view=false`), load all data the view needs
inside the service method's transaction boundary so lazy associations never reach the
view layer and produce `LazyInitializationException`.

The Hibernate session closes at the end of the service method. Any lazy association not
loaded within the transaction will throw `LazyInitializationException` if accessed later.

Service methods must load all data the view needs before returning. Use projections (not
full entities) for list views to avoid lazy-loading issues and over-fetching.

Never pass JPA entities to the view layer — return UI models only.
