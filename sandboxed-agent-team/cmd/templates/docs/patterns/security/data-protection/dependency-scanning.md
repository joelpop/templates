# Dependency Security Scanning

When building the application in CI, run OWASP Dependency-Check or equivalent
so known-vulnerable libraries are caught before they ship to production.

- Fail the build on any dependency with CVSS score ≥ 9.0 (critical)
- Review and remediate or accept high-severity CVEs (7.0–8.9) within 30 days
