# ItemButton<T>

When a button is associated with a specific item — in a grid column, a virtual list row, a card, or any item-presenting component — use `ItemButton<T>` to carry the item directly on the button so event handlers call `event.getItem()` rather than querying an external selection or casting.

```java
public class ItemButton<T> extends Button {

    private final T item;

    public ItemButton(String text, T item) {
        super(text);
        this.item = item;
    }

    public T getItem() {
        return item;
    }

    @SuppressWarnings("unchecked")
    public Registration addClickListener(ComponentEventListener<ItemClickEvent<T>> listener) {
        return addListener(ItemClickEvent.class, (ComponentEventListener) listener);
    }

    public static class ItemClickEvent<T> extends ComponentEvent<ItemButton<T>> {

        public ItemClickEvent(ItemButton<T> source, boolean fromClient) {
            super(source, fromClient);
        }

        public T getItem() {
            return getSource().getItem();
        }
    }
}
```

The item is set at construction and fixed for the button's lifetime:

```java
var deactivateButton = new ItemButton<>("Deactivate", item);
deactivateButton.addClickListener(this::onDeactivateButtonClick);
```

```java
private void onDeactivateButtonClick(ItemButton.ItemClickEvent<ItemListItem> event) {
    var item = event.getItem();
    // ...
}
```

`ItemButton<T>` deliberately extends `Button` (IS-A) rather than wrapping it — see `composite-base.md`.