# JPA Entity and Repository Scan Configuration

When entity classes or repositories live in packages that are not sub-packages
of the application class, declare a `@Configuration` class that explicitly
names the scan packages so Spring Boot discovers all entities and repositories
regardless of package layout.

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