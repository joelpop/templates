# SOLID Principles

- **Single Responsibility:** Each class has one reason to change. Views display; services
  enforce business rules; repositories store.
- **Open/Closed:** Extend behavior through interfaces and composition, not by modifying
  existing classes.
- **Liskov Substitution:** Subtypes must be substitutable for their base types.
- **Interface Segregation:** Prefer small, focused interfaces over large omnibus ones.
- **Dependency Inversion:** Depend on abstractions (`*Service` interfaces), not
  implementations (`Jpa*Service`).
