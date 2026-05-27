# Login Rate Limiting

When a login endpoint receives repeated failures from the same IP, apply a
temporary block or CAPTCHA challenge so brute-force attacks are slowed without
exposing the threshold to the attacker.

## Threshold and Response

- Threshold: more than 10 failed attempts within 5 minutes from a single IP
- Response: temporary block or CAPTCHA challenge
- Log the rate-limit event server-side

Return HTTP 429 with a generic "Too many requests" message. Do not reveal
thresholds or countdowns in any header or body.

The implementation (in-memory, Redis, database) is a per-project architectural
decision.
