# Tailwind CSS

When using Tailwind CSS instead of `LumoUtility`, enable the feature flag and apply class
names to native HTML elements and simple layouts so Tailwind's utility classes are
compiled into the application stylesheet.

## Setup

```properties
# src/main/resources/vaadin-featureflags.properties
com.vaadin.experimental.tailwindCss=true
```

Tailwind classes are applied the same way as any CSS class name:

```java
var warningBox = new Div("Warning!");
warningBox.addClassNames("bg-orange-400 p-4");
```

## Where It Works

Tailwind utility classes are effective on native HTML elements (`Div`, `Span`, etc.) and
simple layout components:

```java
var layout = new HorizontalLayout();
layout.addClassNames("gap-4 p-6 items-center");
```

Complex Vaadin components (`Grid`, `ComboBox`, etc.) have nested shadow DOM structures that
limit how much utility classes can reach — theme variants and `addThemeNames()` are more
effective there.

## Build Constraint: Class Names Must Appear Literally in Source

The Tailwind compiler scans source files to collect used class names and generates a static
stylesheet containing only those classes. Class names constructed at runtime are not
detected and will be missing from the compiled output:

```java
// Avoid — "bg-" + color is not detected by the Tailwind compiler
var color = isDanger ? "red-500" : "green-500";
div.addClassNames("bg-" + color);

// Preferred — full class names appear literally and are compiled
div.addClassNames(isDanger ? "bg-red-500" : "bg-green-500");
```

## Java Constants for Class Names

Defining Tailwind class names as Java string constants — modeled after `LumoUtility` —
satisfies the compiler's literal-in-source requirement while adding compile-time safety
and IDE autocompletion. The scanner detects the string literal at the constant definition
site regardless of where the constant is used:

```java
public final class Tw {

    public static final class Background {
        public static final String WHITE  = "bg-white";
        public static final String SUBTLE = "bg-gray-100";
        public static final String DANGER = "bg-red-50";
    }

    public static final class Padding {
        public static final String NONE   = "p-0";
        public static final String SMALL  = "p-2";
        public static final String MEDIUM = "p-4";
        public static final String LARGE  = "p-6";
    }

    public static final class Text {
        public static final String SMALL  = "text-sm";
        public static final String BASE   = "text-base";
        public static final String LARGE  = "text-lg";
        public static final String MUTED  = "text-gray-500";
        public static final String DANGER = "text-red-600";
    }

    public static final class Rounded {
        public static final String SMALL  = "rounded-sm";
        public static final String MEDIUM = "rounded-md";
        public static final String LARGE  = "rounded-lg";
        public static final String FULL   = "rounded-full";
    }

    public static final class Shadow {
        public static final String SMALL  = "shadow-sm";
        public static final String MEDIUM = "shadow-md";
        public static final String LARGE  = "shadow-lg";
    }

    private Tw() {}
}
```

Usage mirrors `LumoUtility`:

```java
var card = new Div();
card.addClassNames(Tw.Background.WHITE, Tw.Padding.MEDIUM, Tw.Rounded.LARGE, Tw.Shadow.MEDIUM);

var errorLabel = new Span("Required");
errorLabel.addClassNames(Tw.Text.SMALL, Tw.Text.DANGER);
```

**Related:** `theming/theme-selection.md` — theme and utility class trade-offs;
`theming/lumo-utility.md` — LumoUtility usage (mutually exclusive with Tailwind).
