# H2 in PostgreSQL Compatibility Mode

When configuring the test datasource, use H2 in PostgreSQL compatibility mode so production Flyway migration scripts also run in the test environment without a separate schema.

```properties
# application-test.properties
spring.datasource.url=jdbc:h2:mem:testdb;MODE=PostgreSQL;DB_CLOSE_DELAY=-1
spring.datasource.driver-class-name=org.h2.Driver
spring.jpa.database-platform=org.hibernate.dialect.H2Dialect
```
