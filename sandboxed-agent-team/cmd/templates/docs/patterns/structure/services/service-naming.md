# Service Class Naming

When naming service classes, use `{Domain}Service` for interfaces and
`{Technology}{Domain}Service` for implementations so the backing technology is
immediately apparent from the class name.

| Role                   | Convention                           | Example                                                                            |
|------------------------|--------------------------------------|------------------------------------------------------------------------------------|
| Service interface      | Suffix `Service`                     | `EmployeeService`, `DepartmentService`, `CurrentUserTenantService`                 |
| Service implementation | Technology prefix + suffix `Service` | `JpaEmployeeService`, `RestDepartmentService`, `VaadinSessionCurrentUserTenantService` |

The prefix names the **actual backing technology** — not a generic label:

| Implementation strategy                  | Prefix          | Example                                    |
|------------------------------------------|-----------------|--------------------------------------------|
| Spring Data JPA                          | `Jpa`           | `JpaEmployeeService`                       |
| Vaadin session / `AuthenticationContext` | `VaadinSession` | `VaadinSessionCurrentUserTenantService`    |
| REST / HTTP client                       | `Rest`          | `RestDepartmentService`                    |
| In-memory / test double                  | `Mock`          | `MockEmployeeService`                      |

`Jpa` is correct only for implementations backed by Spring Data JPA. Never apply it to a
service that has no JPA dependency.
