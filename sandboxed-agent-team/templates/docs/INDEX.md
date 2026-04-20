# Documentation Index

**Requirement status convention:** Every discrete requirement statement
in a doc carries a status checkbox: `[ ]` not started, `[-]` in
progress, `[x]` complete. See "Status Tracking" in CLAUDE.md for
transition rules. Example format inside a requirement doc:

    ## Authentication
    - [ ] Users can log in with SSO via SAML 2.0
      - Acceptance criteria: ...
    - [-] Passkey-based authentication is supported
      - Acceptance criteria: ...
    - [x] Session timeout after 30 minutes of inactivity
      - Acceptance criteria: ...

## Non-Functional Requirements
Quality attributes (ISO 25010). Every agent must re-read all of these
before starting any task. Files listed here that do not yet exist
should be skipped — their absence is expected early in the project
and does not indicate missing context.

| Tag | File | Description |
|-----|------|-------------|
| NON-FUNCTIONAL | `docs/non-functional/performance.md` | Response time, throughput, capacity |
| NON-FUNCTIONAL | `docs/non-functional/security/authentication.md` | Authentication mechanisms, providers, login flows |
| NON-FUNCTIONAL | `docs/non-functional/security/authorization.md` | Roles, permissions, access control |
| NON-FUNCTIONAL | `docs/non-functional/security/data-protection.md` | Encryption, PII handling, retention |
| NON-FUNCTIONAL | `docs/non-functional/security/hardening.md` | Headers, CORS, CSP, rate limiting |
| NON-FUNCTIONAL | `docs/non-functional/reliability.md` | Availability, fault tolerance, recoverability |
| NON-FUNCTIONAL | `docs/non-functional/usability.md` | Learnability, accessibility, user error protection |
| NON-FUNCTIONAL | `docs/non-functional/maintainability.md` | Modularity, testability, coding standards |
| NON-FUNCTIONAL | `docs/non-functional/portability.md` | Supported platforms, browsers, devices |
| NON-FUNCTIONAL | `docs/non-functional/compatibility.md` | Co-existence, interoperability |
<!-- Uncomment if applicable to this project:
| NON-FUNCTIONAL | `docs/non-functional/internationalization.md` | Language support, localization, text direction |
| NON-FUNCTIONAL | `docs/non-functional/observability.md` | Logging, monitoring, metrics, tracing |
-->

## Functional Requirements — Cross-Cutting
Behavioral requirements spanning multiple features. Every agent must
re-read all of these before starting any task.

| Tag | File | Description |
|-----|------|-------------|
| FUNCTIONAL-CROSS-CUTTING | `docs/functional/cross-cutting/error-handling.md` | Error handling and reporting standards |
| FUNCTIONAL-CROSS-CUTTING | `docs/functional/cross-cutting/data-validation.md` | Input validation rules and patterns |
| FUNCTIONAL-CROSS-CUTTING | `docs/functional/cross-cutting/api-standards.md` | API design conventions and contracts |
| FUNCTIONAL-CROSS-CUTTING | `docs/functional/cross-cutting/integration.md` | External APIs, third-party services, protocols |

## Functional Requirements — Data
Data model and persistence. Re-read when working on data-related tasks.

| Tag | File | Description |
|-----|------|-------------|
| FUNCTIONAL-DATA | `docs/functional/data/schema.md` | Entity model, relationships, constraints |
| FUNCTIONAL-DATA | `docs/functional/data/migration.md` | Migration strategy, seed data |

## Functional Requirements — Features
Re-read the primary doc and ALL supplementals for the feature you are
currently working on.

| Tag | File | Description |
|-----|------|-------------|
| FUNCTIONAL-FEATURE | `docs/functional/features/feature-a.md` | <Feature A — one-line summary> |
| FUNCTIONAL-FEATURE-SUPPLEMENTAL | `docs/functional/features/feature-a/views.md` | Views and dialogs for Feature A. Also read: `feature-a.md` |
| FUNCTIONAL-FEATURE-SUPPLEMENTAL | `docs/functional/features/feature-a/ux.md` | Interaction patterns for Feature A. Also read: `feature-a.md` |
| FUNCTIONAL-FEATURE | `docs/functional/features/feature-b.md` | <Feature B — one-line summary> |

## External Interface Requirements
System boundary and interface specifications.

| Tag | File | Description |
|-----|------|-------------|
| EXTERNAL-INTERFACE | `docs/external-interfaces/user-interfaces.md` | UI standards, interaction paradigms |
| EXTERNAL-INTERFACE | `docs/external-interfaces/software-interfaces.md` | OS, libraries, third-party software |
| EXTERNAL-INTERFACE | `docs/external-interfaces/communication-interfaces.md` | Network protocols, data exchange formats |

## Environmental Requirements
Operating environment and infrastructure.

| Tag | File | Description |
|-----|------|-------------|
| ENVIRONMENTAL | `docs/environmental/infrastructure.md` | Hosting, containers, CI/CD pipelines |
| ENVIRONMENTAL | `docs/environmental/platforms.md` | Supported browsers, OS, devices |
| ENVIRONMENTAL | `docs/environmental/deployment.md` | Deployment strategy, environments |
<!-- Uncomment if applicable to this project:
| ENVIRONMENTAL | `docs/environmental/configuration.md` | Configuration management, feature flags, environment-specific settings |
-->

## Technical Constraints
Design and implementation constraints.

| Tag | File | Description |
|-----|------|-------------|
| TECHNICAL | `docs/technical/stack.md` | Language, framework, DB versions |
| TECHNICAL | `docs/technical/build.md` | Build tools, dependency management |
| TECHNICAL | `docs/technical/constraints.md` | Regulatory, compliance, standards |

## Architectural
Known structural debt and design decisions. Every agent must re-read
before starting any task.

| Tag | File | Description |
|-----|------|-------------|
| ARCHITECTURAL | `docs/architecture-debt.md` | Known structural debt and recommended resolutions |
<!-- Uncomment if applicable to this project:
| ARCHITECTURAL | `docs/architecture-decisions.md` | Architecture Decision Records (ADRs) |
-->
