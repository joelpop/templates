# Dialog Delegation and NonComponent Events

When building a Vaadin dialog class, wrap `Dialog` via delegation (never extend
it) and implement the `NonComponent` event infrastructure for event publishing so
the class exposes a focused public API and supports typed listener registration
without inheriting `Dialog`'s full interface.

## Custom Dialogs Use Delegation, Not Inheritance

No application class may extend `Dialog`. Wrapping `Dialog` via delegation exposes
a focused, intentional API rather than Dialog's 50+ public methods.

```java
// Preferred
public class EditItemDialog {
    private final Dialog dialog;

    public EditItemDialog(...) {
        dialog = new Dialog();
        // configure dialog contents
    }

    public void open() { dialog.open(); }
    public void close() { dialog.close(); }
}

// Avoid
public class EditItemDialog extends Dialog { ... }
```

## NonComponent Event System

Classes that do not extend `Component` — such as delegating dialogs — must use
the `NonComponent` event infrastructure for event publishing. This provides the
same listener registration and event firing semantics as Vaadin's `ComponentEvent`
system.

| Component-event equivalent              | Non-component counterpart                         |
|-----------------------------------------|---------------------------------------------------|
| `Component` (fires events)              | `NonComponent` (marker interface)                 |
| `ComponentEvent<S>`                     | `NonComponentEvent<N extends NonComponent>`       |
| `ComponentEventBus` (held by Component) | `NonComponentEventSupport<N>` (held by field)     |

Copy the three classes below into a shared event package in your UI module (for
example, `{base_package}.ui.event`). They are small and have no dependencies
beyond `com.vaadin.flow.shared.Registration`.

```java
/**
 * Interface for classes that can fire events but don't extend Component.
 * Analogous to how Component provides event capabilities.
 */
public interface NonComponent {

    /**
     * Adds a listener for events of the given type.
     *
     * @param eventType the type of event to listen for
     * @param listener  the listener to add
     * @param <E>       the event type
     * @return a registration that can be used to remove the listener
     */
    <E extends NonComponentEvent<?>> Registration addListener(Class<E> eventType, Consumer<E> listener);
}
```

```java
/**
 * Base class for events fired by non-component sources.
 * Analogous to ComponentEvent but for classes that don't extend Component.
 *
 * @param <N> the type of the event source, must implement NonComponent
 */
public abstract class NonComponentEvent<N extends NonComponent> {
    private final N source;

    /**
     * Creates a new event with the given source.
     *
     * @param source the source of the event
     */
    protected NonComponentEvent(N source) {
        this.source = source;
    }

    /**
     * Returns the source of the event.
     *
     * @return the source of the event
     */
    public N getSource() {
        return source;
    }
}
```

```java
/**
 * Helper class that provides event listener management for NonComponent implementations.
 * Use via composition to add event support to classes that don't extend Component.
 *
 * <p>Designed for UI-thread use, matching ComponentEventBus conventions.
 * Callers dispatching or registering from a Push thread should wrap the call
 * in {@code UI.access(...)}.</p>
 *
 * @param <N> the type of the event source
 */
public class NonComponentEventSupport<N extends NonComponent> {
    private final Map<Class<?>, List<Consumer<?>>> listeners = new HashMap<>();

    /**
     * Adds a listener for events of the given type.
     *
     * @param eventType the type of event to listen for
     * @param listener  the listener to add
     * @param <E>       the event type
     * @return a registration that can be used to remove the listener
     */
    public <E extends NonComponentEvent<N>> Registration addListener(Class<E> eventType, Consumer<E> listener) {
        listeners.computeIfAbsent(eventType, k -> new ArrayList<>()).add(listener);
        return () -> {
            var list = listeners.get(eventType);
            if (list != null) {
                list.remove(listener);
            }
        };
    }

    /**
     * Fires an event to all registered listeners of the event's type.
     *
     * @param event the event to fire
     * @param <E>   the event type
     */
    @SuppressWarnings("unchecked")
    public <E extends NonComponentEvent<N>> void fireEvent(E event) {
        var list = listeners.get(event.getClass());
        if (list != null) {
            // Snapshot before iterating in case a listener adds or removes listeners during dispatch.
            for (var listener : new ArrayList<>(list)) {
                ((Consumer<E>) listener).accept(event);
            }
        }
    }
}
```

Dispatch is keyed by the event's runtime class, so a listener registered on
`SaveEvent.class` only fires for exact-class `SaveEvent` instances — matching
`ComponentEventBus` behavior.

## Caller-Side Pattern

Define typed event subclasses, expose convenience `add*Listener` methods, and
fire through the `NonComponentEventSupport` instance:

```java
public class EditItemDialog implements NonComponent {
    private final Dialog dialog;
    private final NonComponentEventSupport<EditItemDialog> eventSupport = new NonComponentEventSupport<>();

    // Event class — extend NonComponentEvent<SourceType>
    public static class SaveEvent extends NonComponentEvent<EditItemDialog> {
        private final ItemDetail item;

        public SaveEvent(EditItemDialog source, ItemDetail item) {
            super(source);
            this.item = item;
        }

        public ItemDetail getItem() { return item; }
    }

    public static class CancelEvent extends NonComponentEvent<EditItemDialog> {
        public CancelEvent(EditItemDialog source) { super(source); }
    }

    // Implement NonComponent
    @Override
    public <E extends NonComponentEvent<?>> Registration addListener(Class<E> eventType,
                                                                      Consumer<E> listener) {
        return eventSupport.addListener((Class) eventType, listener);
    }

    // Convenience listener registration methods
    public Registration addSaveListener(Consumer<SaveEvent> listener) {
        return eventSupport.addListener(SaveEvent.class, listener);
    }

    public Registration addCancelListener(Consumer<CancelEvent> listener) {
        return eventSupport.addListener(CancelEvent.class, listener);
    }

    protected void fireEvent(NonComponentEvent<EditItemDialog> event) {
        eventSupport.fireEvent(event);
    }
}
```

The caller attaches listeners and handles events:

```java
var dialog = new EditItemDialog(...);
dialog.addSaveListener(e -> handleSave(e.getItem()));
dialog.addCancelListener(e -> dialog.close());
```
