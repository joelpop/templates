# Recipe: `LoadingButton` — Button with Inline Loading State

When a button triggers an operation that could take longer than, say, 200 ms,
follow this recipe to produce a `LoadingButton` that disables itself on click,
shows a loading cue appropriate to its mode (text, icon, or both), then
restores its original state automatically.

## What this produces

- A `LoadingButton extends Button` that auto-detects whether it is text-only,
  icon-only, or icon+text and applies the appropriate loading indicator.
- Default loading text: original text + `"…"`. Override per instance with
  `setLoadingText(String)`.
- Default loading icon: `VaadinIcon.SPINNER`. Override per instance with
  `setLoadingIcon(Component)`.
- `addClickListener` overridden so every click listener automatically runs in
  loading state — no special API to discover or remember.
- `runWithLoading(SerializableRunnable)` — for programmatic invocation outside
  a click listener.

## Dependencies

- Vaadin 24+ (`Button`, `ClickEvent`, `ComponentEventListener`, `VaadinIcon`,
  `SerializableRunnable`).

## Step 1 — The `LoadingButton` class

Mode detection reads `getText()` (empty string when unset) and `getIcon()`
(null when unset) at the moment `runWithLoading` is called, so a button whose
text or icon changes after construction is handled correctly.

`VaadinIcon.SPINNER` is a static icon. To animate it, supply a custom component
via `setLoadingIcon()` — for example, an icon with a CSS `animation: spin`
class.

```java
package {base_package}.ui.component;

import com.vaadin.flow.component.ClickEvent;
import com.vaadin.flow.component.Component;
import com.vaadin.flow.component.ComponentEventListener;
import com.vaadin.flow.component.button.Button;
import com.vaadin.flow.component.icon.VaadinIcon;
import com.vaadin.flow.server.SerializableRunnable;
import com.vaadin.flow.shared.Registration;

public class LoadingButton extends Button {

    private String    loadingText;
    private Component loadingIcon;

    public LoadingButton() {
    }

    public LoadingButton(String text) {
        super(text);
    }

    public LoadingButton(Component icon) {
        super(icon);
    }

    public LoadingButton(String text, Component icon) {
        super(text, icon);
    }

    /**
     * @return the loading text override, or {@code null} to use the default
     *         (original text + "…")
     */
    public String getLoadingText() {
        return loadingText;
    }

    /**
     * @param loadingText text to display while loading; {@code null} restores
     *                    the default (original text + "…")
     */
    public void setLoadingText(String loadingText) {
        this.loadingText = loadingText;
    }

    /**
     * @return the loading icon override, or {@code null} to use the default
     *         ({@link VaadinIcon#SPINNER})
     */
    public Component getLoadingIcon() {
        return loadingIcon;
    }

    /**
     * @param loadingIcon icon to display while loading; {@code null} restores
     *                    the default ({@link VaadinIcon#SPINNER})
     */
    public void setLoadingIcon(Component loadingIcon) {
        this.loadingIcon = loadingIcon;
    }

    /**
     * Wraps every click listener in loading state so callers use the standard
     * {@link Button#addClickListener} API without needing to know about loading
     * behavior.
     *
     * @param listener the click listener to wrap
     * @return a registration for removing the listener
     */
    @Override
    public Registration addClickListener(ComponentEventListener<ClickEvent<Button>> listener) {
        return super.addClickListener(event -> runWithLoading(() -> listener.onComponentEvent(event)));
    }

    /**
     * Disables the button, applies the loading indicator, runs the action,
     * then restores the original text, icon, and enabled state.
     *
     * @param action the action to run
     */
    public void runWithLoading(SerializableRunnable action) {
        var originalText = getText();
        var originalIcon = getIcon();

        var hadText = !originalText.isEmpty();
        var hadIcon = originalIcon != null;

        var activeLoadingText = computeLoadingText(originalText, hadText, hadIcon);
        var activeLoadingIcon = computeLoadingIcon(originalIcon, hadText, hadIcon);

        setEnabled(false);
        if (activeLoadingText != null) {
            setText(activeLoadingText);
        }
        if (activeLoadingIcon != null) {
            setIcon(activeLoadingIcon);
        }

        try {
            action.run();
        } finally {
            if (activeLoadingText != null) {
                setText(originalText);
            }
            if (activeLoadingIcon != null) {
                setIcon(originalIcon);
            }
            setEnabled(true);
        }
    }

    private String computeLoadingText(String originalText, boolean hadText, boolean hadIcon) {
        if (loadingText != null) {
            return loadingText;
        }
        if (!hadText) {
            return null;
        }
        if (!hadIcon && loadingIcon != null) {
            return null;
        }
        return originalText + "…";
    }

    private Component computeLoadingIcon(Component originalIcon, boolean hadText, boolean hadIcon) {
        if (loadingIcon != null) {
            return loadingIcon;
        }
        if (!hadIcon) {
            return null;
        }
        if (!hadText && loadingText != null) {
            return null;
        }
        return VaadinIcon.SPINNER.create();
    }
}
```

## Step 2 — Use in views

```java
// Text-only: "Save" → "Saving…" while running
var saveButton = new LoadingButton("Save");
saveButton.addClickListener(this::onSaveClick);

// Icon-only: icon → spinner while running
var refreshButton = new LoadingButton(VaadinIcon.REFRESH.create());
refreshButton.setAriaLabel("Refresh");
refreshButton.addClickListener(this::onRefreshClick);

// Icon+text with custom loading text
var deleteButton = new LoadingButton("Delete", VaadinIcon.TRASH.create());
deleteButton.setLoadingText("Deleting…");
deleteButton.addClickListener(this::onDeleteClick);
```

## Decisions this recipe imposes

- **Mode detected at call time, not construction time.** Reading `getText()`
  and `getIcon()` inside `runWithLoading` means a button that changes its
  label after construction is handled correctly without any extra wiring.
- **`finally` for state restore.** The original text, icon, and enabled state
  are always restored — even if the action throws. The caller is responsible
  for surfacing errors (e.g., via `Notification`).
- **Synchronous action only.** `runWithLoading` runs the action on the UI
  thread. For background operations, call `runWithLoading` from a
  `UI.getCurrent().access(...)` block and restore state inside `access` after
  the background work completes.
- **`addClickListener` overridden, not a custom method.** A separate
  `addLoadingClickListener` would be overlooked — callers would reach for
  the standard API and get no loading behavior. Overriding ensures the
  loading state is automatic for every listener without any discovery burden.
- **`setDisableOnClick(true)` not used.** That built-in flag only disables the
  button; it cannot swap the label or icon. Managing enabled state in `finally`
  gives identical safety with full control over the loading indicator.

## What to verify

- A text-only button shows `originalText + "…"` while the action runs and
  reverts to the original text when done.
- An icon-only button shows the spinner while the action runs and reverts to
  the original icon when done.
- An icon+text button updates both simultaneously.
- A button with `setLoadingText` / `setLoadingIcon` overrides shows the custom
  values instead of the defaults.
- If the action throws, the button is still re-enabled and its original label
  restored.

## Related

- `docs/patterns/ui/vaadin/button-loading-state.md` — pattern-level guidance
  on when to use `LoadingButton`.
