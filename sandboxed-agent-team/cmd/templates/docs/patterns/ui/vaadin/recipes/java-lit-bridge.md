# Recipe: Java–Lit Bridge — Wrapping Browser APIs as Flow Components

When implementing a Vaadin component that needs to invoke a browser-only API
(WebAuthn, clipboard, geolocation, MediaDevices), follow this recipe to produce
a Java `Component` / Lit `LitElement` pair that hides the browser side behind
typed Java events so Flow code stays in pure server-side Java and the Lit code
stays in pure TypeScript.

## What this produces

- A Java `Component` subclass annotated with `@Tag` and `@JsModule`
  that passes configuration via HTML attributes and receives results
  via `@DomEvent`-mapped typed events.
- A Lit `LitElement` that owns all browser-side logic, dispatches
  `CustomEvent`s on completion, and renders in light DOM so the host
  page's Lumo theme styles apply.

## When to use this pattern

Use when a browser API or third-party JS library has no Vaadin
component and cannot be called server-side:

- Browser APIs: `navigator.credentials` (WebAuthn), `navigator.geolocation`,
  `navigator.clipboard`, `window.print`, `MediaDevices`, etc.
- Third-party JS libraries with no Java wrapper: charting, mapping,
  rich text editors.
- Async browser-initiated flows where the result must return to server
  state.

**Don't use for** browser *details* (timezone, locale, screen size)
that only need to cross the module-boundary once — the
`ClientDetailsService` pattern in
`docs/patterns/conventions/vaadin/datetime.md` covers that case with less
overhead.

## Dependencies

- Vaadin 24+ (`Component`, `ComponentEvent`, `DomEvent`, `EventData`,
  `JsModule`, `Tag`).
- A Lit 3 development environment (the project's frontend build
  compiles `.ts` files under
  `src/main/resources/META-INF/resources/frontend/`).

## Step 1 — Define the Java wrapper

The Java class extends `Component` directly — not `Composite` — because
it wraps a custom HTML element, not a Vaadin layout. The constructor is
private; static factory methods enforce the valid configurations:

```java
package {base_package}.ui.component;

import com.vaadin.flow.component.Component;
import com.vaadin.flow.component.ComponentEvent;
import com.vaadin.flow.component.ComponentEventListener;
import com.vaadin.flow.component.DomEvent;
import com.vaadin.flow.component.EventData;
import com.vaadin.flow.component.Tag;
import com.vaadin.flow.component.dependency.JsModule;
import com.vaadin.flow.shared.Registration;
import java.util.Optional;

@Tag("{app}-geolocation-button")                     // must match @customElement in the TS file
@JsModule("./{app}-geolocation-button.ts")           // relative to src/main/resources/META-INF/resources/frontend/
public class GeolocationButton extends Component {

    private GeolocationButton() { }

    /** Create a button that requests the user's current position. */
    public static GeolocationButton create(String label) {
        var button = new GeolocationButton();
        button.getElement().setAttribute("label", label);
        return button;
    }

    public Registration addSuccessListener(
            ComponentEventListener<GeolocationSuccessEvent> listener) {
        return addListener(GeolocationSuccessEvent.class, listener);
    }

    public Registration addErrorListener(
            ComponentEventListener<GeolocationErrorEvent> listener) {
        return addListener(GeolocationErrorEvent.class, listener);
    }

    @DomEvent("geolocation-success")                 // must match the CustomEvent name dispatched in TS
    public static class GeolocationSuccessEvent extends ComponentEvent<GeolocationButton> {

        private final double latitude;
        private final double longitude;

        public GeolocationSuccessEvent(
                GeolocationButton source,
                boolean fromClient,
                @EventData("event.detail.latitude") double latitude,   // JS expression evaluated client-side
                @EventData("event.detail.longitude") double longitude) {
            super(source, fromClient);
            this.latitude = latitude;
            this.longitude = longitude;
        }

        public double getLatitude()  { return latitude; }
        public double getLongitude() { return longitude; }
    }

    @DomEvent("geolocation-error")
    public static class GeolocationErrorEvent extends ComponentEvent<GeolocationButton> {

        private final String message;

        public GeolocationErrorEvent(
                GeolocationButton source,
                boolean fromClient,
                @EventData("event.detail.message") String message) {
            super(source, fromClient);
            this.message = message;
        }

        public String getMessage() {
            return message != null ? message : "Location request failed";
        }
    }
}
```

### Coordination points

| Java | TypeScript | Rule |
|---|---|---|
| `@Tag("x-foo")` | `@customElement('x-foo')` | Must match exactly — Flow uses the tag name to connect the two |
| `@JsModule("./x-foo.ts")` | file path | Relative to `src/main/resources/META-INF/resources/frontend/` |
| `getElement().setAttribute("label", value)` | `@property({type: String}) label` | Attribute names are kebab-case in HTML; camelCase in TS `@property` with `attribute: 'kebab-name'` when they differ |
| `@DomEvent("x-success")` | `new CustomEvent('x-success', ...)` | Must match exactly |
| `@EventData("event.detail.foo")` | `detail: { foo: value }` | The expression is evaluated in the browser against the DOM `Event` object |

## Step 2 — Write the Lit component

Place the TypeScript file at
`src/main/resources/META-INF/resources/frontend/{app}-geolocation-button.ts`:

```typescript
import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('{app}-geolocation-button')     // must match @Tag on the Java class
export class AppGeolocationButton extends LitElement {

    @property({ type: String }) label = 'Share location';

    // Render in light DOM so the host page's Lumo theme styles the
    // inner <vaadin-button>. Shadow DOM would isolate the component
    // from Lumo variables, breaking theme consistency.
    createRenderRoot() {
        return this;
    }

    render() {
        return html`
            <vaadin-button @click="${this._onClick}">
                ${this.label}
            </vaadin-button>
        `;
    }

    private _onClick() {
        navigator.geolocation.getCurrentPosition(
            (position) => {
                this.dispatchEvent(new CustomEvent('{app}-geolocation-success', {
                    detail: {
                        latitude:  position.coords.latitude,
                        longitude: position.coords.longitude,
                    },
                    bubbles: true,    // let the event propagate up the DOM
                    composed: true,   // cross shadow DOM boundaries (defensive; required if
                                      // any ancestor uses shadow DOM)
                }));
            },
            (error) => {
                this.dispatchEvent(new CustomEvent('{app}-geolocation-error', {
                    detail: { message: error.message },
                    bubbles: true,
                    composed: true,
                }));
            }
        );
    }
}

// TypeScript global type registration — enables IDE autocompletion
// for the custom element name. Not required by Flow.
declare global {
    interface HTMLElementTagNameMap {
        '{app}-geolocation-button': AppGeolocationButton;
    }
}
```

### Light DOM vs shadow DOM

`createRenderRoot() { return this; }` opts into **light DOM** rendering.
The rendered content (the `<vaadin-button>`) becomes a direct child of
the custom element, visible to the host page's CSS and Lumo variables.

Omitting `createRenderRoot` uses Lit's default **shadow DOM** rendering:
the content is isolated in a shadow root and Lumo theme variables do not
apply. Use shadow DOM only for components that are fully self-contained
and have no Vaadin element children.

### Event dispatch requirements

Always set `bubbles: true` and `composed: true` when dispatching events:

- `bubbles: true` — the event propagates up the DOM tree so Flow's
  listener (registered on the element) can receive it regardless of
  where in the template it originates.
- `composed: true` — the event crosses shadow DOM boundaries. Defensive
  even with light DOM rendering; required if any ancestor component
  uses shadow DOM.

## Step 3 — Use in a view

```java
var locationButton = GeolocationButton.create("Share my location");

locationButton.addSuccessListener(event -> {
    locationService.save(event.getLatitude(), event.getLongitude());
    Notification.show("Location saved.").addThemeVariants(NotificationVariant.LUMO_SUCCESS);
});

locationButton.addErrorListener(event -> {
    Notification.show("Could not get location: " + event.getMessage())
            .addThemeVariants(NotificationVariant.LUMO_ERROR);
});

add(locationButton);
```

The view stays in pure Java. All browser-API code — the async ceremony,
error handling, event formatting — lives in the TypeScript file.

## Passing complex data

HTML attributes are strings. For richer configuration, use
`getElement().setPropertyList(...)` (Flow serialises a `List` or `Map`
to a JSON array / object) or encode as a JSON string attribute and parse
in the Lit component.

For data flowing the other direction (from Lit to Java), `@EventData`
expressions are evaluated client-side against the `Event` object — any
serialisable value in `event.detail` is extractable. Stick to
primitives and plain JSON objects; `@EventData` cannot extract class
instances.

## Decisions this recipe imposes

- **One Lit component per browser API.** Don't accumulate multiple
  unrelated APIs in one component — one tag, one responsibility.
- **The Java class owns the configuration contract.** Static factory
  methods communicate the valid modes (as `PasskeyButton` does with
  `forAuthentication` / `forRegistration`); callers don't set
  attributes directly.
- **The Lit component owns all browser-side logic.** No inline JS in
  Java (no `getElement().executeJs(...)`), no raw `Page.executeJs(...)`
  calls for browser-API ceremonies. These bypass Flow's type system
  and are unverifiable by the compiler.
- **Light DOM by default when rendering Vaadin elements.** Shadow DOM
  is opt-in for fully isolated components.
- **Events are named with an app prefix.** `{app}-success` rather than
  `success` avoids collisions with browser-native events and events
  from other components on the same element.

## What to verify

- Adding the Java component to a view renders the Lit component in
  the browser (inspect the DOM; the custom element tag should be
  present with the configured attributes).
- A successful browser-API call fires the Java success listener with
  the correct event data values.
- A failed browser-API call fires the Java error listener; the error
  message is non-null.
- Removing the `createRenderRoot` override (reverting to shadow DOM)
  visually breaks Lumo styling on any inner Vaadin component — confirms
  light DOM is needed.
- Renaming the `@Tag` or `@customElement` value without updating the
  other causes the component to silently render as an empty unknown
  element with no JS behavior — confirms the names must match.

## Related

- `docs/patterns/conventions/vaadin/datetime.md` — "Browser Client Details —
  Bridging the SoC Wall": the related pattern for passing browser
  *details* (timezone, locale) across the module boundary without a
  full Lit component.
- [passkey.md](passkey.md) — the WebAuthn `PasskeyButton` is a
  concrete instance of this pattern with additional Spring Security
  7 CSRF coordination and backend coupling specifics.
- [app-icon.md](app-icon.md) — `@JsModule` is also used by
  `UntitledUiIcon` to register a custom Vaadin iconset; a simpler
  use of the same annotation without the `@Tag` / event machinery.