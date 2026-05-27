# JPA Application Configuration

When configuring Spring Data JPA, disable OSIV and declare explicit entity and
repository scan packages so data loading is transaction-scoped and all entities
are discovered regardless of package layout.

## Disable Open Session in View (OSIV)

```properties
spring.jpa.open-in-view=false
```

OSIV keeps the Hibernate session open through the entire HTTP request, including
view rendering. Disabling it:
- Forces all data loading into the service/transaction layer where it belongs
- Surfaces `LazyInitializationException` immediately in development (these are
  real bugs that OSIV was masking)
- Prevents N+1 queries from silently firing during serialization or template
  rendering

## Entity and Repository Scanning

When entity classes or repositories live in packages that are not sub-packages
of the application class, Spring Boot's default scan path does not cover them.
Declare a `@Configuration` class that explicitly names the scan packages:

```java
package com.example.app.jpaclient.config;

import org.springframework.boot.persistence.autoconfigure.EntityScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

/**
 * Declares entity and repository scan packages explicitly because the JPA
 * packages ({app}-jpamodel, {app}-jpaclient) are not sub-packages of the
 * application class and fall outside its default scan path.
 */
@Configuration
@EntityScan(basePackages = "com.example.app.jpamodel")
@EnableJpaRepositories(basePackages = "com.example.app.jpaclient")
public class JpaConfig {
}
```

This `@Configuration` class must itself be reachable from the application
class's scan path — place it in a package covered by `scanBasePackages`.
