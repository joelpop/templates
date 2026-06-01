# Where Extracted Abstractions Live

When deciding where to place an extracted abstraction, match the size of the abstraction to where it is documented and registered.

| Duplication shape                                          | Right size               | Where the abstraction goes                                                         |
|------------------------------------------------------------|--------------------------|------------------------------------------------------------------------------------|
| State bundle (fields used together)                        | Value object             | Record / value type / Java class                                                   |
| Logic shape (algorithmic duplication)                      | Method                   | On the relevant class, or a utility                                                |
| Repeated try/catch + error wrapping                        | Method                   | Wrapper method or higher-order helper                                              |
| Branching on type                                          | Polymorphism             | Split the type, dispatch via virtual call                                          |
| Repeated structural template (chrome, lifecycle, fixture)  | Shared base class        | Abstract or concrete base class                                                    |
| N classes realizing one named capability                   | Component-family package | A package; capability boundary = package boundary                                  |
| Cross-cutting concern (validation, logging, mapping, etc.) | Shared mechanism         | Interface + implementation (or aspect); document the contract in `docs/solutions/` |

When the answer is "shared mechanism," that mechanism deserves a `docs/patterns/`
entry describing its contract — not just a piece of code that exists.
