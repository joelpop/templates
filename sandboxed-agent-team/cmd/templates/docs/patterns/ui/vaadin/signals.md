# Vaadin Signals

When managing non-`Binder` component state, choose the mechanism based on
Vaadin version so reactive changes propagate without manual listener wiring.

> **Vaadin ≥25.1:** Vaadin Signals are the **preferred mechanism for non-`Binder`
> component state management** — UI state that is not backed by a JPA-style bean
> (visibility toggles, selection state, reactive counts, cross-session shared data,
> computed derivations). Use `ValueSignal`, `ListSignal`, `SharedNumberSignal`,
> `SharedMapSignal`, `Signal.computed(...)` as the default state-management
> primitives.
>
> **Vaadin 25.0:** Signals exist but are not yet positioned as the universal
> preference. Use them for reactive/shared state where they add clear value
> (cross-session updates, computed values); private state-management methods remain
> acceptable for simple cases.
>
> **Vaadin 24.x:** Signals are **not available**. Manage non-`Binder` state
> through private fields, manual listener wiring, and explicit rebind/refresh methods
> on the view or composite. Do not try to backport Signals.

**`Binder` is always preferred for bean-backed forms**, regardless of Vaadin
version — field-to-property binding with validation is the job `Binder` does
best. Signals complement `Binder`; they do not replace it.

For Signal field naming, see `signal-naming.md`.
