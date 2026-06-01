# Code Enum Naming and Structure

When creating a JPA entity enum, suffix it with `Code` and keep it as plain
constants with no properties so the JPA enum stays free of
presentation concerns.

```java
public enum EquipmentTypeCode {
    VEHICLE, AIRCRAFT, MACHINERY, WATERCRAFT
}
```

Display properties belong in the corresponding UI type enum, not here.

Store as string, never ordinal — ordinal values break silently if declaration
order changes:

```java
@Enumerated(EnumType.STRING)
private EquipmentTypeCode equipmentTypeCode;
```