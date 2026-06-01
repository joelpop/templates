# Grid Loading State

When displaying a collection in a Grid, use `CachingDataProvider` instead of
`grid.setItems(list)` so the Grid shows its own built-in loading indicator on
the first page request, and subsequent pages, filter changes, and sort changes
are instant because they operate on the in-memory cache. Data that is fast
today may not be fast tomorrow — using `CachingDataProvider` universally
avoids conversion work as data grows.

See the [CachingDataProvider recipe](recipes/caching-data-provider.md) for the
full implementation and usage examples.
