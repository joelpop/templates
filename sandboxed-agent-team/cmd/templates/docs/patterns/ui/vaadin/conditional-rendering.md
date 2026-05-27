# Conditional Component Rendering

When deciding how to present a control the current user may not be able to
interact with, choose the mode based on the reason: "do not generate" for
authorization, `setVisible(false)` for contextually inapplicable, `setEnabled(false)`
with a tooltip for temporarily blocked — so each mode communicates the correct
signal and no authorization gate can be bypassed via the browser DOM.

## Three Modes

- **Do not generate** — **authorization** (the current user's role does not permit
  the action) or any "this user should never see this at all" condition. Never
  constructed, never in the DOM. Nothing to discover via dev tools, no attribute to
  re-enable by tampering.
- **Hide** (`setVisible(false)`) — **not applicable to the current situation** but
  could apply in another state the same user encounters (e.g., a Reactivate button
  hidden on an active record, shown on an inactive one). Component is constructed
  and lives in the layout so it can be revealed without a rebuild.
- **Disable** (`setEnabled(false)` + tooltip) — **applicable and authorized, but
  not possible right now** (e.g., cannot deactivate the last remaining admin, cannot
  submit while a save is in flight). Tooltip must explain *why*.

```java
// Do not generate — e.g., security gating by role. The button is never constructed for
// users who cannot perform the action; there is no attribute to tamper with in the DOM.
if (currentUser.hasRole(UserRole.ROLE_ADMIN)) {
    var deactivateButton = new Button("Deactivate");
    deactivateButton.addClickListener(_ -> confirmDeactivate());
    toolbar.add(deactivateButton);
}

// Hide — contextually not applicable. The button exists but is not visible for this
// record's current state; it re-appears when the state changes.
reactivateButton.setVisible(!record.isActive());

// Disable — applicable and authorized, temporarily blocked. Show a tooltip explaining why.
deactivateButton.setEnabled(!isLastAdmin);
deactivateButton.setTooltipText(isLastAdmin
    ? "Cannot deactivate the only admin account" : null);
```

**Do not** use `setVisible(false)` for authorization gating — that is the "do not
generate" case. `setEnabled(false)` is not a substitute for the other two: a
permanently-disabled control communicates "try again later" and invites futile
interaction.

## Server-Side State Governs Interactivity

Vaadin state lives on the server. Interactivity is governed by the server-side
component state (`setEnabled(false)`, `setVisible(false)`), not by client-side CSS
or HTML attributes. Any purely-CSS concealment (`visibility: hidden`, `display: none`
via `getStyle().set(...)`, a client-side `disabled` attribute toggled in the browser
inspector) can be reverted by the user in dev tools; if the server still considers
the component enabled and visible, the server will happily process clicks on the
"re-revealed" element. Use server-side state, not CSS, to prevent interaction.

## Layout Preservation — When a Placeholder Is Needed

"Do not generate" and `setVisible(false)` both remove the component from layout
flow; surrounding elements collapse to fill the space. "Disable" preserves layout
at its normal size and position. When absence would cause jarring layout shift or
misalignment (e.g., a toolbar with positionally-dependent buttons), the absent
component needs a **placeholder** that occupies the same space without being
interactive.

Options, in preference order:

1. **Rethink the mode.** If a missing control would disrupt layout, "disable with
   tooltip" preserves the slot by design, blocks interaction on the server, and
   communicates intent more clearly than an empty placeholder.
2. **An inert placeholder** the same size as the control — a `Div`, empty `Span`,
   or subdued "—" label. A separate component, not the real control styled invisible
   — no enabled element hiding in the DOM.
3. **A neutral affordance** ("No action" label, subdued status pill) when the
   absence itself is informative.

Do not construct the real interactive component and hide it with CSS alone
(`getStyle().set("visibility", "hidden")` on an enabled button). That leaves a
fully-interactive server-side component reachable from the browser inspector. If the
real component must occupy the slot, `setEnabled(false)` on the server brings you
back to disable mode — the slot is already preserved.

The three-mode rubric answers *whether* the component exists or functions. Placeholder
decisions answer *what occupies the space* — they compose with the mode choice and
must never weaken server-enforced interactivity state.
