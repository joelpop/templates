# BCrypt Password Storage

When storing user passwords, hash them with BCrypt at a work factor of 10 or
higher so credentials cannot be recovered if the database is compromised.

## Configuration

```java
@Bean
public PasswordEncoder passwordEncoder() {
    return new BCryptPasswordEncoder(10); // work factor >= 10
}
```

## Usage

```java
// Store
String hash = passwordEncoder.encode(rawPassword);

// Verify
boolean valid = passwordEncoder.matches(rawPassword, storedHash);
```

The `users` table contains no plaintext password column. The hash column stores
BCrypt format (`$2a$...`). Work factor must be ≥ 10.
