# NonComponent Event System

When a class does not extend `Component` — such as a delegating dialog — use
the `NonComponent` event infrastructure for event publishing so it supports
typed listener registration with the same semantics as Vaadin's
`ComponentEvent` system.

| Component-event equivalent              | Non-component counterpart                     |
| :-------------------------------------- | :-------------------------------------------- |
| `Component` (fires events)              | `NonComponent` (marker interface)             |
| `ComponentEvent<S>`                     | `NonComponentEvent<N extends NonComponent>`   |
| `ComponentEventBus` (held by Component) | `NonComponentEventSupport<N>` (held by field) |

Copy the three classes below into a shared event package in your UI module
(e.g. `{base_package}.ui.event`). They are small and have no dependencies
beyond `com.vaadin.flow.shared.Registration`.

```java
public interface NonComponent {
    <E extends NonComponentEvent<?>> Registration addListener(Class<E> eventType, Consumer<E> listener);
}
```

```java
public abstract class NonComponentEvent<N extends NonComponent> {
    private final N source;

    protected NonComponentEvent(N source) { this.source = source; }

    public N getSource() { return source; }
}
```

```java
public class NonComponentEventSupport<N extends NonComponent> {
    private final Map<Class<?>, List<Consumer<?>>> listeners = new HashMap<>();

    public <E extends NonComponentEvent<N>> Registration addListener(Class<E> eventType, Consumer<E> listener) {
        listeners.computeIfAbsent(eventType, k -> new ArrayList<>()).add(listener);
        return () -> {
            var list = listeners.get(eventType);
            if (list != null) list.remove(listener);
        };
    }

    @SuppressWarnings("unchecked")
    public <E extends NonComponentEvent<N>> void fireEvent(E event) {
        var list = listeners.get(event.getClass());
        if (list != null) {
            for (var listener : new ArrayList<>(list)) {
                ((Consumer<E>) listener).accept(event);
            }
        }
    }
}
```

Dispatch is keyed by the event's runtime class, matching `ComponentEventBus`
behavior: a listener on `SaveEvent.class` fires only for exact-class
`SaveEvent` instances.

## Caller-Side Pattern

Define typed event subclasses, expose convenience `add*Listener` methods, and
fire through the `NonComponentEventSupport` instance:

```java
public class EditItemDialog implements NonComponent {
    private final Dialog dialog;
    private final NonComponentEventSupport<EditItemDialog> eventSupport = new NonComponentEventSupport<>();

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

    @Override
    public <E extends NonComponentEvent<?>> Registration addListener(Class<E> eventType, Consumer<E> listener) {
        return eventSupport.addListener((Class) eventType, listener);
    }

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

The view attaches listeners and handles events:

```java
var editItemDialog = new EditItemDialog(...);
editItemDialog.addSaveListener(e -> handleSave(e.getItem()));
editItemDialog.addCancelListener(e -> editItemDialog.close());
```