# Account Lockout Notification

When an account is locked due to repeated failed login attempts, notify the
account owner out-of-band so they know to contact an administrator — without
revealing lockout state at the login form.

Send a notification to the account owner's registered contact (email, etc.)
at the moment of lockout:

```
"Your account has been locked due to repeated failed login attempts.
Contact an administrator to regain access."
```

The login form continues to display the generic failure message ("Incorrect
[identifier] or password") so the attacker receives no confirmation that their
attempts triggered a lockout. Server logs record the lockout reason, actor, and
timestamp.