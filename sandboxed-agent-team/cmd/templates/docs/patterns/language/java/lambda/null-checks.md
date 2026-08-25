# Null Check Policy

When writing a lambda or stream pipeline, do not add null guards for values the
framework or data source guarantees to be non-null — defensive checks for
impossible conditions obscure intent and spread distrust of the API contract.

- `TextField.getValue()` returns `""`, never `null`
- `MultiSelectComboBox.getValue()` returns an empty `Set`, never `null`
- `Signal` fields initialized with a non-null value remain non-null
- Stream sources over framework-managed collections do not produce null elements
