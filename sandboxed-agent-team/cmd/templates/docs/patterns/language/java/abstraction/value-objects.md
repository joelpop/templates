# Value Objects for State Bundles

When a small group of fields is consistently treated as a unit — passed together, validated together, displayed together — extract a value object rather than scattering primitives across the API.

## "Bytes plus content type" → a value object

```java
// Avoid — primitives scattered across the API; every consumer

// reassembles the relationship between bytes and type
public class Asset {
    private byte[] data;
    private String contentType;
}

void process(byte[] data, String contentType) { /* ... */ }
```

```java
// Preferred — the relationship is encoded in a type
public record ContentData(byte[] bytes, String contentType) { /* ... */ }

public class Asset {
    private ContentData content;
}

void process(ContentData content) { /* ... */ }
```

## "First name plus last name" → a name type

```java
// Avoid — components exposed as primitives; every consumer reimplements

// display order, locale choices, and casing independently
String firstName;
String lastName;
```

```java
// Preferred — components stay structured; rendering is the consumer's

// responsibility, not the model's
public record PersonName(String first, String last) {
    String firstLast() { return "%s, %s".formatted(first, last); }
    String lastFirst() { return "%s, %s".formatted(last, first); }
    String fullName()  { return firstLast(); }
}
```

When the same name-formatting helpers appear on multiple domain types, extract them as interface defaults:

```java
public interface HasNames {
    String getFirstName();
    String getLastName();
    default String firstLast() { return "%s, %s".formatted(getFirstName(), getLastName()); }
    default String lastFirst() { return "%s, %s".formatted(getLastName(), getFirstName()); }
    default String fullName()  { return firstLast(); }
}
```

Each type implements `HasNames` and declares its own fields; none re-implements the default methods.
