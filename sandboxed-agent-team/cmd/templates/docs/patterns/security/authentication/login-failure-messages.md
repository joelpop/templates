# Login Failure Messages

When reporting a login failure, display the same generic message for every
failure condition so accounts cannot be enumerated and lockout state is not
disclosed.

- Account not found → display: "Incorrect [identifier] or password"
- Wrong password    → display: "Incorrect [identifier] or password"
- Account locked    → display: "Incorrect [identifier] or password"

Replace `[identifier]` with whatever the application uses as the account ID
(email address, username, employee ID, etc.).

All three conditions produce the same message. Switching to a distinct "Account
is locked" message at lockout reveals that the account exists and confirms to
an attacker that their brute-force triggered a lockout.
