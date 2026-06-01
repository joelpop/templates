# JavaDoc

Write JavaDoc on all public types and methods — they describe the contract for
callers, which is information the call site cannot see. This is distinct from
inline narrative comments inside method bodies (see `comments.md`).

For constructors and methods, include as appropriate:
- `@param` for every parameter
- `@return` describing the return value (omit for `void` and constructors)
- `@throws` for every checked and documented unchecked exception

```java
/**
 * Finds the entity with the given key.
 *
 * @param key the surrogate primary key
 * @return the matching entity data
 * @throws EntityNotFoundException if no entity with that key exists in the current context
 */
ItemDetail findByKey(long key);
```
