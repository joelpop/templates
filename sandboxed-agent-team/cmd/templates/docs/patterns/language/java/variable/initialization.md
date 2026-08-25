# Member Variable Initialization

When declaring member variables, declare without initialization and assign in
the constructor body so all configuration for each field is in one place — since
creation can itself include configuration, splitting initialization across the
declaration and the constructor body scatters related setup.

```java
// Avoid — creation at declaration, configuration in constructor; each field split across two places
public class ReportWriter {
    private PrintWriter writer = new PrintWriter(System.out);
    private NumberFormat formatter = NumberFormat.getCurrencyInstance();

    public ReportWriter() {
        writer.println("Report started");
        formatter.setMaximumFractionDigits(2);
    }
}
```

```java
// Preferred
public class ReportWriter {
    private final PrintWriter writer;
    private final NumberFormat formatter;

    public ReportWriter() {
        writer = new PrintWriter(System.out);
        writer.println("Report started");

        formatter = NumberFormat.getCurrencyInstance();
        formatter.setMaximumFractionDigits(2);
    }
}
```