# Response Header Technology Disclosure

When deploying the application, ensure HTTP response headers do not reveal the server technology stack — suppress or omit headers that expose the framework or server version.

- `Server` header: absent or generic
- `X-Powered-By` header: absent

These headers allow attackers to target known vulnerabilities for a specific server or framework version.
