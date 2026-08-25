# Variable Naming

When naming a variable, field, or parameter, use the full semantic name — do not
shorten or collapse a type into a generic alias — so the role is unambiguous at
the declaration site without reading surrounding context:

```java
// Avoid — "preference" could be the dialog, binder, service, model, or component
var preference = new UserPreferenceSettingDialog();
```

```java
// Preferred
var userPreferenceSettingDialog  = new UserPreferenceSettingDialog();
var userPreferenceSettingBinder  = new Binder<>(UserPreferenceSetting.class);
var userPreferenceSettingService = context.getBean(UserPreferenceSettingService.class);
```

When multiple artifacts exist for the same domain concept (dialog, binder, component,
service, utility), a shortened name forces the reader to guess the role. The full name
makes the role unambiguous at the declaration site.

**UI component fields** follow the same rule — suffix with the component's role, which
is usually the component type but may be a more descriptive semantic name when the type
alone is too generic:

```java
private final Span                   totalValueSpan;
private final Button                 saveButton;
private final TextField              displayIdField;
private final Grid<ItemListItem>     itemGrid;
private final Div                    buttonBar;    // role is clearer than buttonDiv
```

**Signal fields** follow the same rule — suffix with the signal type so the
reactive type is unambiguous alongside non-signal fields of the same domain:

```java
// Avoid — "items" and "editing" could be lists, booleans, or signals
private final ListSignal<ItemListItem>  items;
private final ValueSignal<Boolean>      editing;
```

```java
// Preferred
private final ListSignal<ItemListItem>  itemsListSignal;
private final ValueSignal<Boolean>      editingSignal;
```

**Exception — lambda parameters:** Single-use lambda parameters may use a short
abbreviation because their scope is a single expression and the type is visible on the
right-hand side or inferred:

```java
preferences.forEach(pref -> grid.addItem(pref.getLabel()));
```
