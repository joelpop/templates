# Dependency Injection

Use **constructor injection** for all mandatory dependencies. Declare fields `private
final` and let Spring call the single public constructor — no `@Autowired` annotation
needed (Spring 4.3+ auto-detects the single constructor).

```java
// Avoid — field injection
@SpringComponent
public class OrderService {
    @Autowired private OrderRepository orderRepository;
    @Autowired private CustomerRepository customerRepository;
}
```

```java
// Preferred
@SpringComponent
public class OrderService {
    private final OrderRepository orderRepository;
    private final CustomerRepository customerRepository;

    public OrderService(OrderRepository orderRepository,
                        CustomerRepository customerRepository) {
        this.orderRepository = orderRepository;
        this.customerRepository = customerRepository;
    }
}
```

Constructor injection makes dependencies explicit, allows `final` fields, works in
plain unit tests without reflection, and fails fast at instantiation. Setter / field
injection is acceptable only for truly optional or reconfigurable dependencies.