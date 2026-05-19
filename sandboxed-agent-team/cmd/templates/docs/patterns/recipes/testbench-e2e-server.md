# Recipe: TestBench E2E Server — Maven `it` Profile and IDE-Compatible Startup

This recipe wires Jetty start/stop around Maven failsafe and a `ServerExtension`
JUnit extension so integration tests run identically from Maven
(`mvn verify -Pit`) and from the IDE without prestarting anything.
`ServerExtension` detects whether Maven already started a server on the
configured port and is a no-op in that case.

**Requires:** Vaadin 25.1+. For parallel execution, apply the
`testbench-e2e-parallel` recipe alongside this one.

---

## Step 1 — Maven dependencies

```xml
<!-- TestBench (commercial) -->
<dependency>
    <groupId>com.vaadin</groupId>
    <artifactId>vaadin-testbench-junit6</artifactId>
    <scope>test</scope>
    <optional>true</optional>
</dependency>

<!-- Embedded Jetty for ServerExtension (IDE runs) -->
<dependency>
    <groupId>org.eclipse.jetty.ee10</groupId>
    <artifactId>jetty-ee10-webapp</artifactId>
    <version>${jetty.version}</version>
    <scope>test</scope>
    <optional>true</optional>
</dependency>

<!-- Annotation scanning — required for Vaadin's servlet SCI registration -->
<dependency>
    <groupId>org.eclipse.jetty.ee10</groupId>
    <artifactId>jetty-ee10-annotations</artifactId>
    <version>${jetty.version}</version>
    <scope>test</scope>
    <optional>true</optional>
</dependency>
```

---

## Step 2 — Maven `it` profile

The IT port (9090) must differ from the dev server's default (8080) so
`mvn verify -Pit` can run alongside a running dev server.

```xml
<properties>
    <it-deployment.port>9090</it-deployment.port>
</properties>
```

```xml
<profile>
    <id>it</id>
    <build>
        <plugins>
            <plugin>
                <groupId>org.eclipse.jetty.ee10</groupId>
                <artifactId>jetty-ee10-maven-plugin</artifactId>
                <configuration>
                    <scan>0</scan>
                    <httpConnector>
                        <port>${it-deployment.port}</port>
                    </httpConnector>
                    <stopPort>8081</stopPort>
                    <stopWait>5</stopWait>
                    <stopKey>${project.artifactId}</stopKey>
                </configuration>
                <executions>
                    <execution>
                        <id>start-jetty</id>
                        <phase>pre-integration-test</phase>
                        <goals><goal>start</goal></goals>
                    </execution>
                    <execution>
                        <id>stop-jetty</id>
                        <phase>post-integration-test</phase>
                        <goals><goal>stop</goal></goals>
                    </execution>
                </executions>
            </plugin>

            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-failsafe-plugin</artifactId>
                <executions>
                    <execution>
                        <goals>
                            <goal>integration-test</goal>
                            <goal>verify</goal>
                        </goals>
                    </execution>
                </executions>
                <configuration>
                    <trimStackTrace>false</trimStackTrace>
                    <enableAssertions>true</enableAssertions>
                    <systemPropertyVariables>
                        <deployment.port>${it-deployment.port}</deployment.port>
                    </systemPropertyVariables>
                </configuration>
            </plugin>
        </plugins>
    </build>
</profile>
```

If the parallel recipe is also applied, add
`<com.vaadin.testbench.Parameters.testsInParallel>` to
`<systemPropertyVariables>` as described there.

---

## Step 3 — Maven-filtered test resource

`ServerExtension` reads the port (and the parallel limit, if applied) from a
Maven-filtered properties file so IDE runs use the same values as Maven runs.

Enable filtering in `<build>`:

```xml
<testResources>
    <testResource>
        <directory>src/test/resources</directory>
        <filtering>true</filtering>
    </testResource>
</testResources>
```

Create `src/test/resources/it-test.properties`:

```properties
deployment.port=${it-deployment.port}
```

---

## Step 4 — ServerExtension

`ServerExtension` is a reference-counted `BeforeAllCallback`/`AfterAllCallback`.
The first IT class to run brings the server up; the last to finish takes it down.
Concurrent classes wait on a `CountDownLatch` until the server is ready.

```java
public class ServerExtension implements BeforeAllCallback, AfterAllCallback {

    private static final Object LOCK = new Object();
    private static Server server;
    private static int useCount;
    private static CountDownLatch ready;

    @Override
    public void beforeAll(ExtensionContext context) throws Exception {
        boolean isFirst;
        CountDownLatch signal;
        synchronized (LOCK) {
            isFirst = ++useCount == 1;
            if (isFirst) {
                ensureDeploymentPortSet();
                ready = new CountDownLatch(1);
            }
            signal = ready;
        }
        if (isFirst) {
            try {
                if (!isPortInUse(deploymentPort())) {
                    startServer(deploymentPort());
                }
            } finally {
                signal.countDown();
            }
        } else {
            signal.await();
        }
    }

    @Override
    public void afterAll(ExtensionContext context) throws Exception {
        synchronized (LOCK) {
            if (--useCount == 0 && server != null) {
                server.stop();
                server = null;
            }
        }
    }
}
```

**`waitUntilReady()`** — poll for HTTP 200 AND the presence of `type="module"`
in the response body. A bare Jetty startup page returns 200 without that marker;
only Vaadin's real HTML shell includes it:

```java
private static void waitUntilReady(int port) throws Exception {
    var url = URI.create("http://localhost:" + port + "/").toURL();
    var deadline = System.currentTimeMillis() + 30_000;
    while (System.currentTimeMillis() < deadline) {
        try {
            var conn = (HttpURLConnection) url.openConnection();
            conn.setConnectTimeout(1_000);
            conn.setReadTimeout(10_000);
            if (conn.getResponseCode() == 200) {
                var body = new String(conn.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
                if (body.contains("type=\"module\"")) return;
            }
        } catch (Exception ignored) {}
        Thread.sleep(250);
    }
    throw new IllegalStateException(
            "Vaadin did not become ready on port " + port + " within 30 seconds");
}
```

**Jetty classpath** — `AnnotationConfiguration` must find Vaadin's
`LookupServletContainerInitializer` to register `VaadinServlet`. Build the
classpath from both `java.class.path` and the `URLClassLoader` hierarchy — in
Maven the JARs live in the Plexus `URLClassLoader`, not `java.class.path`:

```java
webApp.setExtraClasspath(collectClasspath());
webApp.setParentLoaderPriority(true);
```

---

## Verify

```
mvn verify -Pit
```

Failsafe starts Jetty on port 9090, runs `*IT.java` tests, and stops Jetty.

For IDE runs: run any `*IT` class directly — `ServerExtension` starts the
embedded server automatically and shuts it down when tests finish.