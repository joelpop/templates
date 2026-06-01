# JPA Configuration

When configuring JPA, place `JpaConfig` in `{app}-jpaclient` with `@EntityScan`,
`@EnableJpaRepositories`, and `@EnableJpaAuditing`, and disable OSIV with
`spring.jpa.open-in-view=false`.

```java
@Configuration
@EntityScan(basePackages = "{base_package}.jpamodel")
@EnableJpaRepositories(basePackages = "{base_package}.jpaclient")
@EnableJpaAuditing
public class JpaConfig {

    @Bean
    public AuditorAware<UserEntity> auditorAware(EntityManager entityManager) {
        return () -> Optional.ofNullable(SecurityContextHolder.getContext().getAuthentication())
                .filter(Authentication::isAuthenticated)
                .map(Authentication::getPrincipal)
                .filter(AuditedPrincipal.class::isInstance)
                .map(AuditedPrincipal.class::cast)
                .map(AuditedPrincipal::getKey)
                .map(key -> entityManager.getReference(UserEntity.class, key));
    }
}
```

Each auth flow's principal implements `AuditedPrincipal` and carries the user's surrogate
key from its login-time validation step, so `AuditorAware` reads the key off the principal
without a per-write DB lookup. `EntityManager.getReference` returns a Hibernate proxy
holding just the key — JPA persists the FK from the proxy with no `SELECT` issued.

Set `spring.jpa.open-in-view=false` in `application.properties`. See
`docs/patterns/persistence/spring-data-jpa/jpa-config.md` for why OSIV must be disabled.
