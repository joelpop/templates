# Project: <PROJECT_NAME>

## Stack
- Language: Java <version>
- Framework: Vaadin <version>, Spring Boot <version>
- Database: <DATABASE — e.g., PostgreSQL 16>
- Testing:
  - Unit & Browserless UI (Unit Tester): JUnit <JUNIT_VERSION — e.g.,
    5 or 6> for unit tests; Vaadin Browserless Testing
    (Vaadin 25.1+: `browserless-test-junit<JUNIT_VERSION>`, free /
    Apache 2.0, extends `SpringBrowserlessTest`; pre-25.1:
    `vaadin-testbench-junit<JUNIT_VERSION>`, commercial, extends
    `SpringUIUnitTest`) for in-process UI component and interaction
    tests (browser-less, container-less); Mockito for
    mocking. One test class per production class. Browserless UI
    tests live in the same package as the view they test (`*Test.java`
    suffix, run by surefire). Class name suffix distinguishes test
    type: `*Test.java` = surefire, `*IT.java` = failsafe.
  - End-to-End (E2E Tester): Node.js Playwright (`@playwright/test`)
    for browser-based end-to-end tests. E2E tests are written in
    TypeScript and live in `<e2e-test-dir>/` (e.g., `e2e/`). This is
    the Vaadin-recommended approach for E2E testing.
  - Testing pyramid: unit tests → browserless UI tests → end-to-end
    tests. E2E tests run only at the pre-PR gate, not per-commit.
- CI: <CI_PLATFORM — e.g., GitHub Actions, GitLab CI>

## Documentation Sources (MCP Servers)
The following MCP servers are configured in `.mcp.json` at the
project root and available to all teammates. Prefer these over training
data — training data may be outdated or describe deprecated patterns.

| Server | Provides | Primary Users |
|--------|----------|---------------|
| `java` | Java standard library and ecosystem Javadoc | Coder, Architect, Unit Tester |
| `vaadin` | Vaadin framework documentation and API | Coder, Architect, Unit Tester |
| `spring-docs` | Spring Boot and Spring Framework docs | Coder, Architect, Unit Tester |
| `playwright` | Playwright API docs and browser automation for visual debugging | E2E Tester, Coder, Architect |
| `fetch` | Fetch arbitrary web pages for documentation | All roles |

When in doubt about a framework API, query the relevant MCP server
before writing code. The "Primary Users" column is guidance — all
servers are available to all teammates.

**Visual debugging with `playwright`:** Any teammate can use the
`playwright` MCP server to interact with the running application
(requires the dev server to be running — see Key Commands). Navigate
to pages, take screenshots, click elements, and inspect visual state.
Use this to verify UI implementation, debug layout issues, or
investigate test failures. The Coder and E2E Tester are the primary
users; the Architect may use it when evaluating framework paradigm
compliance.

**Note:** Customize this table to match the `mcpServers` configured
in `.mcp.json`. Remove entries for servers not in use; add entries
for any project-specific servers.

## Canonical Claude Code References

This project uses Claude Code's (currently) experimental Agent Teams
feature (`CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` in
`.claude/settings.json`). The team mechanism — `TeamCreate`,
`SendMessage`, the `TaskCreate` / `TaskUpdate` / `TaskGet` /
`TaskList` shared task list, teammate persistence, the mailbox,
the `TeammateIdle` hook — is documented at:

- https://code.claude.com/docs/en/agent-teams.md
- https://code.claude.com/docs/en/sub-agents.md (related but distinct mechanism)
- https://code.claude.com/docs/en/hooks.md
- https://code.claude.com/docs/en/skills.md
- https://code.claude.com/docs/en/settings.md

A `SessionStart` hook (`.claude/hooks/session-start-fetch-docs.sh`,
registered in `.claude/settings.json` for `startup` / `resume` /
`clear` / `compact` matchers) fetches the Agent Teams doc into
context at every session start, including after resume and after
context compaction. Defer to its content over training-time recall
when answering any question about Agent Teams mechanics.

### Verify before asserting Claude Code mechanics

Claude's training data describes earlier versions of Claude Code
and is known to conflate distinct features (most notably Agent
Teams and subagents). Before asserting how *anything* in Claude
Code works — hooks, skills, slash commands, the `Agent` tool,
Agent Teams, settings fields, file conventions — verify against
`code.claude.com/docs/...` first. Use the `fetch` MCP server,
`WebFetch`, or the `claude-code-guide` subagent if available. If
you cannot verify, say "I don't know — let me look it up" rather
than answering from recall.

This rule exists because past sessions on this kit spent many
hours and significant tokens on designs built atop
confidently-asserted mechanics that later turned out to be wrong.
Verifying takes seconds. Drift compounds.

## Documentation Index

The kit's documentation convention splits `docs/` into the
following trees:

- `docs/agnostic/` — project-agnostic patterns, preferences, and
  standing guidance (not tied to a specific project's requirements).
- `docs/reqs/` — project-specific requirements following IEEE 830 /
  ISO 29148 (SRS structure) and ISO 25010 (quality model).
- `docs/tech/` — project-specific technical reference (the
  "pattern playbook"). Describes HOW the team implements recurring
  patterns, at a level above code but below requirements. Curated
  by the Architect, committed by the Analyst. Entries can describe
  planned or already-implemented patterns and grow as new patterns
  are decided.
- `docs/glossary.md` — the project's implementation-agnostic
  vocabulary used by requirements via inline Markdown links (e.g.,
  "edit affordance" instead of "dialog"). Curated by the Architect,
  committed by the Analyst. Note: the term "agnostic" here means
  *implementation-agnostic* — separate from the `AGNOSTIC` tag,
  which marks project-agnostic content under `docs/agnostic/`.

The master index is `docs/INDEX.md`. It is a sample file seeded at
setup time — setup does not overwrite it on re-runs, so edit it
freely as the project's docs evolve.

### Tags and teammate reading rules

- **AGNOSTIC** — project-agnostic patterns, preferences, and
  standing guidance under `docs/agnostic/`. Every teammate re-reads
  all of these before starting any task.
- **NON-FUNCTIONAL** — quality attribute requirements under
  `docs/reqs/non-functional/` (ISO 25010). Every teammate re-reads all
  of these before starting any task. Files listed that do not yet
  exist should be skipped — their absence is expected early in the
  project and does not indicate missing context.
- **FUNCTIONAL-CROSS-CUTTING** — behavioral requirements spanning
  multiple features, under `docs/reqs/functional/cross-cutting/`.
  Every teammate re-reads all of these before any task.
- **FUNCTIONAL-DATA** — data model and persistence, under
  `docs/reqs/functional/data/`. Re-read when working on
  data-related tasks.
- **FUNCTIONAL-FEATURE** — primary doc for a specific feature,
  under `docs/reqs/functional/features/`. Re-read the primary doc
  AND all supplementals for the feature currently being worked on.
- **FUNCTIONAL-FEATURE-SUPPLEMENTAL** — additional detail for a
  feature (views, UX, feature-scoped NFRs, etc.); does not stand
  alone. Each entry must include an "Also read" pointer to its
  primary FEATURE doc, and vice versa.
- **EXTERNAL-INTERFACE** — system boundary and interface
  requirements, under `docs/reqs/external-interfaces/`. Re-read when
  touching those interfaces.
- **ENVIRONMENTAL** — infrastructure and deployment requirements,
  under `docs/reqs/environmental/`. Re-read when touching deployment
  or infrastructure.
- **TECHNICAL** — stack, tooling, and design constraints, under
  `docs/reqs/technical/`. Re-read as relevant to the current task.
- **ARCHITECTURAL** — structural debt and design decisions, under
  `docs/reqs/`. Every teammate re-reads all of these before any task.
- **TECH-REF** — implementation patterns the team uses (the
  "pattern playbook"), under `docs/tech/`. Curated by the
  Architect, committed by the Analyst. Re-read entries linked from
  the in-scope requirements of your current task; the Architect
  also reads these at task kickoff to anchor structural proposals
  in existing patterns.
- **GLOSSARY** — entries in `docs/glossary.md`, the project's
  implementation-agnostic vocabulary used by requirements via
  inline Markdown links. Curated by the Architect, committed by the
  Analyst. Every teammate reads `docs/glossary.md` before starting
  any task (it is small) so requirement links resolve.

Feature-scoped non-functional requirements (e.g., "dashboard loads in
2 seconds") live under the feature as FUNCTIONAL-FEATURE-SUPPLEMENTAL,
not under `docs/reqs/non-functional/`.

### Requirement status convention

Every discrete requirement statement in a doc carries a status
checkbox: `[ ]` not started, `[-]` in progress, `[x]` complete. See
"Status Tracking" below for transition rules. Example format inside
a requirement doc:

    ## Authentication
    - [ ] Users can log in with SSO via SAML 2.0
      - Acceptance criteria: ...
    - [-] Passkey-based authentication is supported
      - Acceptance criteria: ...
    - [x] Session timeout after 30 minutes of inactivity
      - Acceptance criteria: ...

## Repository Structure
```
<paste or refine what Claude discovered>
```

## Directory Ownership Rules (for Agent Teams)
These rules prevent teammates from overwriting each other's work:
- Files marked COORDINATE: message the Lead before editing.
- Each teammate owns their assigned directories only.
- Shared config files (pom.xml, etc.): Lead approves all edits.
- These rules are auto-derived from the project structure. When the
  structure changes significantly (e.g., single to multi-module), the
  Lead updates this section to reflect the new layout.

Ownership map (auto-derived — adjust after structural changes):
- `src/main/java/`            → Coder
- `src/main/resources/`       → Coder
- `src/main/frontend/`        → Coder
- `src/test/java/`            → Unit Tester
- `<e2e-test-dir>/`           → E2E Tester
- `docs/`                     → Analyst
- `pom.xml`                   → COORDINATE (Lead approves)
- `README.md`                 → COORDINATE (Lead approves)
- CI/CD config (e.g., `.github/workflows/`) → COORDINATE (Lead approves)
- `Dockerfile` / `docker-compose.yml` → COORDINATE (Lead approves)
- DB migrations (e.g., `src/main/resources/db/migration/`) → Coder (Architect reviews)

**Multi-module projects:** Replace the map above with per-module
entries (e.g., `module-a/src/main/java/` → Coder). Each module's
`pom.xml` is COORDINATE. The root `pom.xml` is COORDINATE. The Lead
may assign different Coders to different modules for parallel subtasks.

## Key Commands
```bash
# Build
mvn clean package

# Run all tests
mvn test

# Run targeted tests (specific class, module, or pattern)
mvn test -Dtest=AuthServiceTest

# Lint (static analysis)
<LINT_COMMAND — e.g., mvn sonar:sonar>

# Format (auto-fix style)
<FORMAT_COMMAND — e.g., mvn spotless:apply>

# Start dev server
mvn spring-boot:run
```

## Implementation Philosophy
Prefer elegant, idiomatic solutions over verbose ones, AS LONG AS the code
remains readable to a mid-level developer without special explanation.

Specifically:
- Use enum properties (fields, methods, lambdas) instead of switch statements
  or if/else chains on enum values. The behavior belongs on the enum, not
  scattered across consumers.
- Use polymorphism and strategy patterns over type-checking conditionals.
- Use composition over inheritance when extending behavior.
- Use functional idioms (map, filter, Optional chaining, Stream pipelines)
  when they make intent clearer than imperative loops.
- If a "clever" solution requires a comment to explain it, it's too clever.
  Refactor until the code explains itself.

## Framework Identity: Vaadin Is Not Traditional Web Development
Vaadin is a server-side UI framework. The UI is built in Java, runs on
the server, and Vaadin handles all client-server communication
automatically. This is fundamentally different from traditional web
development, and teammates MUST use Vaadin idioms — not patterns from
their general web training data.

**Core paradigm:**
- UI is built with Java component classes, not HTML templates
- UI state lives on the server, not in the browser
- Navigation uses Vaadin's `@Route` system, not REST endpoints
- Data binding uses `Binder`, not manual form handling
- Styling uses Lumo theme and `LumoUtility` classes, not CSS
  frameworks (Bootstrap, Tailwind, etc.)
- Server push replaces client-side polling and WebSocket management
- `DataProvider` handles lazy loading, not pagination APIs

**Anti-patterns to reject — these indicate traditional web thinking:**
- REST controllers (`@RestController`, `@GetMapping`) for UI data —
  Vaadin views call service interfaces directly from Java
- JavaScript/TypeScript for business logic — logic belongs in Java
  on the server; JS is only for low-level browser interop
- Client-side state management (Redux, stores, signals in JS) — state
  is server-side Java fields and Vaadin Signals
- HTML/template files for layout — use Java component composition
- `fetch()` / AJAX / JSON APIs between "frontend" and "backend" —
  there is no separate frontend; it is all one server-side application
- CSS frameworks or custom CSS for things Lumo provides — check
  `LumoUtility` and component theme variants first
- Manual DOM manipulation — use Vaadin's component API
- Servlet filters for auth — use Vaadin's view-level access control
  (`@RolesAllowed`, `@PermitAll`, `@AnonymousAllowed`)

**Before starting any Vaadin-related task**, every teammate must consult
the `vaadin` MCP server to get current information about modern Vaadin
development. For Spring-related work, consult `spring-docs`. For Java
API questions, consult `java`. Do not rely on training data —
framework APIs evolve between versions and training data may describe
deprecated patterns. See "Documentation Sources (MCP Servers)" above
for the full list of available servers.

**Note for Claude Code:** Customize this section for the project's
framework. The principle — use the framework's idioms, not generic
web patterns — applies to any framework. Replace the Vaadin-specific
content above with the relevant paradigm, anti-patterns, and
documentation sources for the project's actual framework.

## CRITICAL: Requirements Are Not Negotiable
**Teammates must NEVER change project requirements to match their implementation.**

If a requirement specifies a version, library, framework, or approach:
- Use that exact version/library/framework, even if you have limited
  training data on it. Search documentation, read source code, experiment.
- If you genuinely cannot make it work after a thorough attempt, MESSAGE
  THE LEAD and explain what you tried and what failed. Do NOT silently
  downgrade, substitute, or rewrite the requirement.
- If you encounter a conflict between your training data ("conventional
  wisdom") and what the project's own documentation or code comments say,
  THE PROJECT'S DOCUMENTATION WINS. Always. Your training data may be
  outdated, inapplicable, or wrong for this context.
- Before applying patterns from general knowledge, CHECK whether the
  project's docs, README, or code comments explicitly warn against that
  pattern. Grep for "do not", "don't", "avoid", "WARNING", "NOTE" in
  relevant source files and documentation.

Violations of this rule — silently changing requirements or ignoring
in-project documentation in favor of general training — are treated as
the highest-severity issue and must be escalated to the Lead immediately
by any teammate that notices.

## Requirements Ambiguity — Do Not Guess
Requirements will sometimes be unclear, ambiguous, conflicting, or
insufficiently specified. When this happens, teammates must escalate —
not guess.

**Recognize these ambiguity signals:**
- A requirement says WHAT but not HOW, and multiple valid approaches
  exist with different trade-offs that would affect the public API,
  data model, or user-visible behavior
- Two documents (or two sections of the same document) describe
  contradictory behavior for the same scenario
- A requirement uses vague terms ("appropriate", "as needed", "handle
  gracefully") without defining what that means concretely
- An edge case or boundary condition is not addressed — the requirement
  covers the happy path but is silent on error cases, empty states,
  concurrent access, etc.
- A requirement references a concept, entity, or workflow that is not
  defined elsewhere in the docs
- You find yourself choosing between two reasonable interpretations and
  cannot determine which one the human intended

**What you MUST NOT do:**
- Fill in gaps using your training data or "common sense" — your
  assumptions may contradict the human's intent
- Pick the simplest interpretation because it's easier to implement
- Treat documentation silence as permission OR prohibition — silence
  on a requirement-level question (WHAT the system must do, WHAT
  constraints it must satisfy) is ambiguity, and is neither "go
  ahead" nor "don't". This rule covers only the WHAT; silence at
  the HOW level (implementation details within the bounds an
  existing requirement defines) is professional judgment, not
  ambiguity — see the Requirement Gate Workflow's refinement rule
  in `team-start.md`.
- Implement both interpretations and "let the human choose later" —
  this creates dead code and doubles the test surface

**What you MUST do:**
- STOP implementation of the ambiguous part (you may continue working
  on unambiguous parts of the same task)
- Escalate to the Architect with the specific ambiguity (see
  Requirements Clarification Escalation in Coordination Rules in
  team-start.md)
- The Architect will attempt to resolve from existing docs; if not
  possible, the Lead will escalate to the human
- Do not proceed with the ambiguous part until a resolution is recorded
  in the task file

## Conventions
- Commit messages: conventional commits
  (https://www.conventionalcommits.org/) —
  `<type>(<scope>): <description>`. Types: `feat`, `fix`, `docs`,
  `test`, `refactor`, `chore`, `perf`. Scope is the feature or
  component affected (e.g., `feat(auth): add SSO integration`).
  The description should explain *what changed and why*, not
  itemize files or lines touched.
- All PRs require passing tests before merge.
- Do NOT commit directly to `<DEV_BRANCH_NAME>`.

## Status Tracking

### Requirement Status
Every discrete requirement statement in `docs/` carries a status checkbox:
- `[ ]` — not started (no task has begun implementing this requirement)
- `[-]` — in progress (a task is actively implementing this requirement)
- `[x]` — complete (implemented and verified through the full task lifecycle)

Acceptance criteria beneath a requirement inherit the requirement's status
and do not carry their own checkboxes.

**Status transitions:**
- `[ ]` → `[-]`: Analyst marks on the task branch at task kickoff
  (first commit on the branch, before sub-branches are created).
- `[-]` → `[x]`: Analyst marks on the task branch at the pre-PR gate
  (after confirming requirement coverage). The squash merge carries
  these to `<DEV_BRANCH_NAME>`. Dev only ever sees `[ ]` → `[x]`.
- `[x]` → `[ ]` or `[-]` → `[ ]`: Analyst resets when adding or
  substantively changing a requirement. Analyst must notify Lead on any
  reset so Lead can assess impact on active or completed tasks.
- Renaming or moving a requirement does not reset its status, but the
  Analyst must update all cross-references (INDEX.md, active task files
  in `.claude/.tasks/`).

### Task Plan Status
Each task file in `.claude/.tasks/<task-id>.md` tracks progress at the
plan-step level. Steps are role-assigned and use the same checkbox
notation. Each teammate marks their own steps as `[-]` when starting
and `[x]` when done.

### Project Status
`.claude/.progress.md` is a minimal dispatcher — it exists solely so the
Lead can recover current state after context compaction. It answers two
questions: "which task am I working on, and what else is parked?" and
"which requirement branches are in flight?"

`progress.md` is gitignored local metadata. It is not affected by branch
operations and persists across branch switches and task
suspension/resumption. It carries only IDs and one-line labels for
recognition — all detail lives in the task files and requirement docs.

**Single writer:** Only the Integrator writes `.claude/.progress.md`.
No other role edits it directly. When state changes (a new task
becomes active, a task is suspended or resumed, a requirement branch
is created or merged, etc.), the Lead directs the Integrator to
update `progress.md`; the Integrator also updates it proactively as
part of workflows it owns (e.g., the Integration Merge Workflow).
This single-writer rule prevents concurrent writes to the file.

Structure:
```markdown
# Progress

## Active Task
- <task-id>: <one-line description>
  <!-- Optional indented annotation if the task is resumed but held: -->
  - blocked on `<DEV_BRANCH_NAME>` health since <ISO 8601 UTC>

## Suspended Tasks
- <task-id>: Blocked by <prerequisite task-id or description>

## Requirement Branches
- requirement/<slug>: <status> — <one-line description>
```

Requirement branch statuses:
- `drafting` — Analyst is actively working on this branch
- `awaiting-approval` — draft submitted to the Lead for human review
- `approved` — human approved; ready to merge to `<DEV_BRANCH_NAME>`
- `merged` — merged to `<DEV_BRANCH_NAME>`; branch can be deleted

## Branching
- Development branch: `<DEV_BRANCH_NAME>` (e.g., `develop`)
- Requirement branches: `requirement/<slug>` — branched off `<DEV_BRANCH_NAME>`
  by the Integrator for the Analyst to draft requirement docs. One branch per
  topic or related group (e.g., `requirement/authentication`,
  `requirement/dashboard-v2`), not per individual requirement — the
  Analyst freely splits, merges, and cross-references requirements
  within a group. Multiple requirement branches can exist simultaneously
  at different stages. Squash-merged back to `<DEV_BRANCH_NAME>` after human
  approval. Tracked in `.claude/.progress.md`.
- Task branches: `task/<task-id>` — branched off `<DEV_BRANCH_NAME>` by the
  Integrator for each implementation task.
- Teammate sub-branches: `task/<task-id>/<role>` — each teammate
  branches off the task branch to do their work:
  - `task/<task-id>/coder` (or `coder-a`, `coder-b` when the Lead
    splits a task across parallel Coders — see Parallel Subtask
    Coders in Coordination Rules in team-start.md)
  - `task/<task-id>/unit-tester`
  - `task/<task-id>/e2e-tester`
  - `task/<task-id>/janitor`
  - The Analyst has no sub-branch — it works on `requirement/<slug>`
    branches and commits status marks directly on the task branch.
  - The Architect has no branch — it reads code on other teammates'
    branches but does not commit.
- Teammate sub-branch operations: each teammate creates their
  sub-branch once at the start of the task and reuses it for all
  commit cycles within that task. Merge (not rebase) in both
  directions — merge FROM the task branch to stay current, merge
  INTO the task branch using the Task Branch Merge Protocol (see
  Coordination Rules in team-start.md). No teammate commits to
  another teammate's branch.
  Sub-branches are local only — they are never pushed to the remote. Only
  `<DEV_BRANCH_NAME>` interacts with the remote (via the Integration Merge
  Workflow).
- Merge strategy: squash merge for all branch-to-`<DEV_BRANCH_NAME>` merges.
  This keeps `<DEV_BRANCH_NAME>` history clean but loses per-commit
  granularity — ensure the squash commit message captures key decisions
  and affected components (see Integration Merge Workflow T.5 in
  team-start.md).
- Merge method: `<MERGE_METHOD>`
- Include cost report in commit message: `<COST_IN_COMMIT>` (`yes` or
  `no`). When `yes`, the Integrator appends a per-model token/cost
  delta report (via `ccusage`) to the final squash-merge commit
  message so it persists in git history. When `no`, the report is
  still shown to the human at task conclusion but is not recorded
  in git. See T.6 in `team-start.md` for the exact flow. To change
  this setting later, ask the Lead ("change `Include cost report
  in commit message` to `yes`/`no`"). The Lead delegates to the
  Integrator, which creates a working branch off `<DEV_BRANCH_NAME>`,
  updates this line in `CLAUDE.md`, commits, and finalizes per the
  project's merge method above. **Caveat:** `CLAUDE.md` is read by
  every role at task start, so a mid-engagement edit can collide
  with any task/requirement branch in flight. Prefer to make this
  change during a quiet period — no active tasks, no open
  requirement branches — so the new value is in effect uniformly
  for all subsequent work. If a change must happen while work is
  in flight, the Lead should pause new task creation until the
  edit is merged and all teammates have pulled the latest
  `<DEV_BRANCH_NAME>`.

## What NOT to do
- Do not add new dependencies without messaging the Lead.
- Do not modify CI/CD pipeline files without explicit approval.
- Do not store secrets in code. Use environment variables.

## Architecture Debt
See `docs/reqs/architecture-debt.md` for known structural debt and
recommended resolutions.

## Non-Functional Requirements
See `docs/reqs/non-functional/` for performance, security, reliability,
usability, and other quality attribute requirements.

## Context Compaction Warning
<!-- SYNC NOTE: The file list below is duplicated in the Pre-Task
     Context Check in team-start.md. If you update one, update both. -->
This file is read at session start but may be LOST during long sessions
when context compaction occurs. You cannot reliably detect whether
compaction has occurred. Therefore: before starting ANY task, you MUST
verify you still have the context needed to work safely. Do this by
explicitly re-reading the following files in order:

1. `CLAUDE.md` (this file) — stack, ownership rules, critical constraints
2. `docs/INDEX.md` — master list of all requirement, glossary, and
   tech-ref documents
3. `docs/glossary.md` — agnostic vocabulary referenced inline by
   requirement docs (Markdown links). Read this before any
   requirement doc so the linked terms make sense.
4. Every file tagged NON-FUNCTIONAL, FUNCTIONAL-CROSS-CUTTING, or
   ARCHITECTURAL in `docs/INDEX.md`, plus any TECHNICAL,
   ENVIRONMENTAL, EXTERNAL-INTERFACE, or TECH-REF docs relevant to
   your current task
5. `docs/reqs/architecture-debt.md` — known structural debt
6. The FEATURE doc in `docs/INDEX.md` matching your current task, plus
   all FEATURE-SUPPLEMENTAL docs linked from it. Follow inline
   Markdown links from the requirements into `docs/glossary.md` and
   `docs/tech/` as you encounter them — those are part of the
   requirement's intent.
7. `.claude/.tasks/<your-task>.md` — your specific assignment
8. `.claude/.progress.md` — which task is active, which are suspended.
   Verify you are working on the correct active task.

**Worktree note:** Items 1–6 are version-controlled and exist in every
worktree. Items 7–8 are gitignored and exist only in the main project
root. Teammates in worktrees must use the absolute project root path
(provided by the Lead at spawn time) to read these files — do not use
relative paths.

If any of these files are missing or their content does not match your
understanding of the project, STOP and message the Lead before
proceeding. Do not work from memory. Do not assume your context is
intact.

Critical rules that MUST survive compaction (re-read if in doubt):
1. Requirements are not negotiable. Do not change versions or substitute
   libraries. (See "Requirements Are Not Negotiable" above.)
2. Ambiguous requirements: do not guess. Escalate to the Architect.
   (See "Requirements Ambiguity — Do Not Guess" above.)
3. Use framework idioms, not traditional web patterns. Consult MCP
   servers for framework documentation, not training data. (See
   "Framework Identity" and "Documentation Sources" above.)
4. Project documentation overrides your training data. Always.
5. Check Directory Ownership before editing any file.
6. Lint and format only the files you have touched, using the commands
   in Key Commands above, before every commit.
7. Mark your own task plan steps as `[-]` when starting and `[x]` when
   done. (See "Status Tracking" above.)
8. Requirement docs (`docs/`) are human-owned. Never commit changes to
   `docs/` without human approval relayed through the Lead.

Keep `.claude/.progress.md` current: the Integrator updates it when a
task becomes active, is suspended, or completes. No other role writes
to this file — see "Single writer" under Project Status above.
