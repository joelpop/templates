# Button Loading State

When a button triggers an operation that could take longer than, say, 200 ms,
use `LoadingButton` instead of `Button` so the user receives immediate feedback
that the action is in progress and cannot double-submit.

`LoadingButton` auto-detects text-only, icon-only, and icon+text modes and
restores the original state — including re-enabling the button — when the
operation completes or throws.

See the [LoadingButton recipe](recipes/loading-button.md) for the full
implementation and usage examples.