# Detecting N+1 Queries in Tests

When testing query-heavy service methods, enable Hibernate statistics and pair with `datasource-proxy` to assert exact query counts and catch N+1 regressions before they reach production.

```properties
# application-test.properties
spring.jpa.properties.hibernate.generate_statistics=true
logging.level.org.hibernate.stat=DEBUG
```

```java
@Test
void listAll_issuesExactlyOneQuery() {
    // setup test data ...
    var queryCounter = new QueryCountHolder();

    var results = service.listAll();

    assertThat(queryCounter.getCount()).isEqualTo(1);
}
```
