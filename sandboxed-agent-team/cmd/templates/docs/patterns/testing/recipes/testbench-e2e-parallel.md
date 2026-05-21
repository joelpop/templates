# Recipe: Parallel TestBench E2E Tests

This recipe configures parallel execution for TestBench integration tests.
The non-obvious part: TestBench's `ParallelConfigurationStrategy` reads the
concurrent limit during JUnit launcher setup — before `@ExtendWith` extension
classes are loaded. `TestBenchParallelLimiter`, a `LauncherSessionListener` SPI,
fires at the right moment to set the limit in both Maven and IDE runs.

**Requires:** `testbench-e2e-server` recipe applied first (provides
`vaadin-testbench-junit6`, `ServerExtension`, and `it-test.properties`).

---

## Step 1 — Maven dependency

`junit-platform-launcher` is needed to compile `TestBenchParallelLimiter`:

```xml
<dependency>
    <groupId>org.junit.platform</groupId>
    <artifactId>junit-platform-launcher</artifactId>
    <scope>test</scope>
</dependency>
```

---

## Step 2 — Parallel limit property

Define the limit and pass it to failsafe before the forked JVM starts:

```xml
<properties>
    <integration-test.concurrent-limit>8</integration-test.concurrent-limit>
</properties>
```

In the `it` profile's failsafe `<systemPropertyVariables>`:

```xml
<com.vaadin.testbench.Parameters.testsInParallel>${integration-test.concurrent-limit}</com.vaadin.testbench.Parameters.testsInParallel>
```

In `src/test/resources/it-test.properties` (so `ServerExtension` can read it for IDE runs):

```properties
integration-test.concurrent-limit=${integration-test.concurrent-limit}
```

---

## Step 3 — `junit-platform.properties`

Create `src/test/resources/junit-platform.properties`:

```properties
junit.jupiter.execution.parallel.enabled = true
junit.jupiter.execution.parallel.mode.default = concurrent
junit.jupiter.execution.parallel.mode.classes.default = concurrent
junit.jupiter.execution.parallel.config.strategy = custom
junit.jupiter.execution.parallel.config.custom.class = com.vaadin.testbench.parallel.ParallelConfigurationStrategy
```

---

## Step 4 — TestBenchParallelLimiter

Create `src/test/java/.../it/TestBenchParallelLimiter.java`:

```java
public class TestBenchParallelLimiter implements LauncherSessionListener {

    @Override
    public void launcherSessionOpened(LauncherSession session) {
        ServerExtension.applyConcurrentLimit();
    }
}
```

`ServerExtension.applyConcurrentLimit()` is defined in the `testbench-e2e-server`
recipe — add the static initializer, `applyConcurrentLimit()`, and
`readIntegrationTestConcurrentLimit()` to `ServerExtension` as shown there.

Register via SPI — create
`src/test/resources/META-INF/services/org.junit.platform.launcher.LauncherSessionListener`
containing:

```
com.example.application.it.TestBenchParallelLimiter
```

Add a corresponding `static` block to `ServerExtension` as a safety net for IDEs
that may not fire the `LauncherSessionListener` SPI:

```java
static {
    applyConcurrentLimit();
}
```

This fires when JUnit loads the extension class. It is typically too late for
`ParallelConfigurationStrategy` to pick up the value, but it costs nothing and
covers edge cases.

---

## Verify

```
mvn verify -Pit
```

Check the failsafe reports — multiple IT classes should show overlapping
start/finish timestamps, confirming tests ran concurrently.