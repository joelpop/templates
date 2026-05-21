# Recipe: Parallel Browserless UI Tests with BrowserlessTest + ComponentTester

This recipe adds the `browserless-test-junit6` dependency and configures Maven
surefire for parallel execution. Tests run entirely in the JVM — no browser,
no servlet container, no Spring startup cost — and belong in the surefire
(unit) tier: `*Test.java` suffix, run on every commit.

For `BrowserlessTest` vs `SpringBrowserlessTest` base class selection and
`ComponentTester<T>` page object conventions, see
`docs/patterns/testing/patterns.md` → "Browserless UI Tests" and
"Page Objects in Browserless UI Tests".

**Requires:** Vaadin 25.1+. Uses `browserless-test-junit6` (free, Apache 2.0).

---

## Step 1 — Maven dependency

Add to `<dependencies>`:

```xml
<dependency>
    <groupId>com.vaadin</groupId>
    <artifactId>browserless-test-junit6</artifactId>
    <scope>test</scope>
    <optional>true</optional>
</dependency>
```

If the project does not already have an SLF4J binding (e.g., via
`spring-boot-starter-test`), add one so Vaadin log output is visible
during test runs:

```xml
<dependency>
    <groupId>org.slf4j</groupId>
    <artifactId>slf4j-simple</artifactId>
    <scope>test</scope>
</dependency>
```

---

## Step 2 — Parallel execution configuration

Configure surefire to run browserless tests in parallel. Add a
`<unit-test.concurrent-limit>` property and a surefire configuration block:

```xml
<properties>
    <unit-test.concurrent-limit>32</unit-test.concurrent-limit>
</properties>
```

```xml
<plugin>
    <groupId>org.apache.maven.plugins</groupId>
    <artifactId>maven-surefire-plugin</artifactId>
    <configuration>
        <properties>
            <configurationParameters>
                junit.jupiter.execution.parallel.enabled=true
                junit.jupiter.execution.parallel.mode.default=concurrent
                junit.jupiter.execution.parallel.mode.classes.default=concurrent
                junit.jupiter.execution.parallel.config.strategy=fixed
                junit.jupiter.execution.parallel.config.fixed.parallelism=${unit-test.concurrent-limit}
            </configurationParameters>
        </properties>
    </configuration>
</plugin>
```

Vaadin's first-time class initialization (~1 s) is paid once per JVM and shared
across all parallel tests. Running 32 tests wall-clock costs about the same as
running 1. Vaadin's browserless infrastructure handles `ThreadLocal` session
isolation automatically — do not share UI or session state across test methods.

---

## Verify

```
mvn test
```

Tests should run in parallel (32 concurrent by default). Total wall-clock time
for the full suite should be close to the time for a single test.
