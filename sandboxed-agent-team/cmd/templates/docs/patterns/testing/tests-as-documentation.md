# Tests as Behavioral Documentation

When naming and writing tests, write them so a new contributor can infer the behavior of the subject from test names alone — without reading the production code.

A well-written test class doubles as a behavioral specification. A new
contributor reading `PasswordResetServiceTest` should be able to infer the reset
flow from test names alone, without reading the production code. If they can't —
if names are mechanical (`testEmptyInput`, `testEdgeCase2`) or generic
(`shouldWork`, `verifyBehavior`) — the tests fail as documentation.
