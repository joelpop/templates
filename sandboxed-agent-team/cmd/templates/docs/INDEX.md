# Documentation Index

*Sample file.* Setup writes this once if `docs/INDEX.md` does not
exist. Re-setup and remove do not touch it. Edit freely to reflect
the project's actual documents.

The kit's documentation convention splits `docs/` into two trees:

- `docs/agnostic/` — project-agnostic patterns, preferences,
  and standing guidance that apply regardless of which project
  the team is working on.
- `docs/reqs/` — project-specific requirements following
  IEEE 830 / ISO 29148 (SRS structure) and ISO 25010 (quality
  model).

## Project-Agnostic (`docs/agnostic/`)

Patterns, preferences, and cross-project requirements. Not
"requirements" in the ISO sense — more like standing guidance the
team applies everywhere. Each entry below is a seed example; add,
remove, or restructure freely.

| Tag | File | Description |
|-----|------|-------------|
| AGNOSTIC | `docs/agnostic/preferences.md` | Developer / team preferences that persist across projects |
| AGNOSTIC | `docs/agnostic/patterns.md` | Design patterns and idioms the team favors |
| AGNOSTIC | `docs/agnostic/standards.md` | Coding standards and conventions |

## Project Requirements (`docs/reqs/`)

### Non-Functional Requirements

| Tag | File | Description |
|-----|------|-------------|
| NON-FUNCTIONAL | `docs/reqs/non-functional/performance.md` | Response time, throughput, capacity |
| NON-FUNCTIONAL | `docs/reqs/non-functional/security/authentication.md` | Authentication mechanisms, providers, login flows |
| NON-FUNCTIONAL | `docs/reqs/non-functional/security/authorization.md` | Roles, permissions, access control |
| NON-FUNCTIONAL | `docs/reqs/non-functional/security/data-protection.md` | Encryption, PII handling, retention |
| NON-FUNCTIONAL | `docs/reqs/non-functional/security/hardening.md` | Headers, CORS, CSP, rate limiting |
| NON-FUNCTIONAL | `docs/reqs/non-functional/reliability.md` | Availability, fault tolerance, recoverability |
| NON-FUNCTIONAL | `docs/reqs/non-functional/usability.md` | Learnability, accessibility, user error protection |
| NON-FUNCTIONAL | `docs/reqs/non-functional/maintainability.md` | Modularity, testability, coding standards |
| NON-FUNCTIONAL | `docs/reqs/non-functional/portability.md` | Supported platforms, browsers, devices |
| NON-FUNCTIONAL | `docs/reqs/non-functional/compatibility.md` | Co-existence, interoperability |
<!-- Uncomment if applicable to this project:
| NON-FUNCTIONAL | `docs/reqs/non-functional/internationalization.md` | Language support, localization, text direction |
| NON-FUNCTIONAL | `docs/reqs/non-functional/observability.md` | Logging, monitoring, metrics, tracing |
-->

### Functional Requirements — Cross-Cutting

| Tag | File | Description |
|-----|------|-------------|
| FUNCTIONAL-CROSS-CUTTING | `docs/reqs/functional/cross-cutting/error-handling.md` | Error handling and reporting standards |
| FUNCTIONAL-CROSS-CUTTING | `docs/reqs/functional/cross-cutting/data-validation.md` | Input validation rules and patterns |
| FUNCTIONAL-CROSS-CUTTING | `docs/reqs/functional/cross-cutting/api-standards.md` | API design conventions and contracts |
| FUNCTIONAL-CROSS-CUTTING | `docs/reqs/functional/cross-cutting/integration.md` | External APIs, third-party services, protocols |

### Functional Requirements — Data

| Tag | File | Description |
|-----|------|-------------|
| FUNCTIONAL-DATA | `docs/reqs/functional/data/schema.md` | Entity model, relationships, constraints |
| FUNCTIONAL-DATA | `docs/reqs/functional/data/migration.md` | Migration strategy, seed data |

### Functional Requirements — Features

| Tag | File | Description |
|-----|------|-------------|
| FUNCTIONAL-FEATURE | `docs/reqs/functional/features/feature-a.md` | <Feature A — one-line summary> |
| FUNCTIONAL-FEATURE-SUPPLEMENTAL | `docs/reqs/functional/features/feature-a/views.md` | Views and dialogs for Feature A. Also read: `feature-a.md` |
| FUNCTIONAL-FEATURE-SUPPLEMENTAL | `docs/reqs/functional/features/feature-a/ux.md` | Interaction patterns for Feature A. Also read: `feature-a.md` |
| FUNCTIONAL-FEATURE | `docs/reqs/functional/features/feature-b.md` | <Feature B — one-line summary> |

### External Interface Requirements

| Tag | File | Description |
|-----|------|-------------|
| EXTERNAL-INTERFACE | `docs/reqs/external-interfaces/user-interfaces.md` | UI standards, interaction paradigms |
| EXTERNAL-INTERFACE | `docs/reqs/external-interfaces/software-interfaces.md` | OS, libraries, third-party software |
| EXTERNAL-INTERFACE | `docs/reqs/external-interfaces/communication-interfaces.md` | Network protocols, data exchange formats |

### Environmental Requirements

| Tag | File | Description |
|-----|------|-------------|
| ENVIRONMENTAL | `docs/reqs/environmental/infrastructure.md` | Hosting, containers, CI/CD pipelines |
| ENVIRONMENTAL | `docs/reqs/environmental/platforms.md` | Supported browsers, OS, devices |
| ENVIRONMENTAL | `docs/reqs/environmental/deployment.md` | Deployment strategy, environments |
<!-- Uncomment if applicable to this project:
| ENVIRONMENTAL | `docs/reqs/environmental/configuration.md` | Configuration management, feature flags, environment-specific settings |
-->

### Technical Constraints

| Tag | File | Description |
|-----|------|-------------|
| TECHNICAL | `docs/reqs/technical/stack.md` | Language, framework, DB versions |
| TECHNICAL | `docs/reqs/technical/build.md` | Build tools, dependency management |
| TECHNICAL | `docs/reqs/technical/constraints.md` | Regulatory, compliance, standards |

### Architectural

| Tag | File | Description |
|-----|------|-------------|
| ARCHITECTURAL | `docs/reqs/architecture-debt.md` | Known structural debt and recommended resolutions |
<!-- Uncomment if applicable to this project:
| ARCHITECTURAL | `docs/reqs/architecture-decisions.md` | Architecture Decision Records (ADRs) |
-->
