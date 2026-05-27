# Entropy-Based Password Validation

When validating passwords at account creation or change, reject by entropy
rather than character-class rules so weak passwords are blocked without
frustrating users with arbitrary composition requirements.

## Thresholds

- Minimum entropy: 50 bits
- Minimum length: 8 characters
- Maximum length: 128 characters
- Common/breached passwords rejected via blocklist check
- Display a visual strength indicator (Weak / Fair / Good / Strong / Very Strong) during entry

## Service Layer Validation

```java
public void validatePasswordStrength(String password) {
    if (password.length() < 8 || password.length() > 128) {
        throw new ValidationException("Password must be 8–128 characters.");
    }
    if (entropyBits(password) < 50) {
        throw new ValidationException("Password is too weak. Try adding more variety.");
    }
    if (blocklist.contains(password.toLowerCase())) {
        throw new ValidationException("This password is too common. Please choose another.");
    }
}
```
