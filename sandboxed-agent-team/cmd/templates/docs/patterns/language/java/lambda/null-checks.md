# Null Check Policy

Do not null-check where the framework guarantees non-null results:

- `TextField.getValue()` returns `""`, never `null`
- `MultiSelectComboBox.getValue()` returns an empty `Set`, never `null`
- `Signal` fields initialized with a non-null value remain non-null

Do not filter Streams for nulls when the data source guarantees non-null elements.
