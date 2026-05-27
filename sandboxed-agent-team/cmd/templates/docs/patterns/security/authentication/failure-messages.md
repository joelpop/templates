# Authentication Failure Messages

When reporting authentication failures or account lockouts, use the same
generic message for all failure reasons and communicate lockout without
revealing thresholds so accounts cannot be enumerated via login probing and
lockout policy cannot be reverse-engineered.

## Login Failure

Do not reveal whether an email address exists in the system:

- "Email not found" → display: "Incorrect email or password"
- "Wrong password" → display: "Incorrect email or password"

Both conditions produce the same generic message. Identical messages prevent
account enumeration via login probing.

## Account Lockout

Communicate lockout status without revealing security details:

```
"Account is locked. Contact an administrator."
```

Server logs record the lockout reason and actor. The user never sees threshold
values, attempt counts, or countdown timers.
