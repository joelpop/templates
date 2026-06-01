# AC Cross-Reference in Test Names

When tracing tests to acceptance criteria, name tests by behavior rather than by AC identifier so test names remain stable as requirement documents evolve.

Some teams put the AC ID in the test name (e.g., `AC3_requestReset_sendsEmail`)
or in a Javadoc comment. That makes traceability searchable but couples test
names to external identifiers that rename over time.

Pragmatic stance: name tests by *behavior* (stable); maintain the AC ↔ test
mapping in the requirement document if the project benefits from explicit
traceability. The mapping table lives alongside the requirement statement; the
Unit Tester updates it during the per-commit cycle.
