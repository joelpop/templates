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

See `docs/INDEX.md` for the master pointer to the four-tree
structure. Briefly:

- `docs/reqs/` — project-specific requirements (Analyst owns).
  IEEE 830 / ISO 29148 (SRS structure) and ISO 25010 (quality
  model).
- `docs/patterns/` — project-agnostic conventions, architecture
  patterns, and recipes for the project's stack. Curated by the
  Architect, committed by the Analyst. Designed to extract and
  reuse across projects.
- `docs/architecture/` — project-specific architecture and
  design — *how this codebase realizes the requirements*.
  Curated by the Architect, committed by the Analyst.
- `docs/guides/` — install / deploy / user / admin / operator
  guides. Owned by the Tech Writer; release-cadence updates.
- `docs/glossary.md` — the project's canonical vocabulary
  (implementation-agnostic terms with optional slang variants).
  Curated by the Architect, committed by the Analyst.

`docs/INDEX.md` is a sample file seeded at setup time — setup
does not overwrite it on re-runs, so edit it freely as the
project's docs evolve.

### Tags and teammate reading rules

- **PATTERN** — project-agnostic conventions, architecture
  patterns, and recipes under `docs/patterns/`. Every teammate
  reads relevant entries (per the Primary references in their
  agent definition).
- **NON-FUNCTIONAL** — quality attribute requirements under
  `docs/reqs/non-functional/` (ISO 25010). Every teammate reads
  all of these before starting any task. Files listed that do not
  yet exist should be skipped — their absence is expected early
  in the project and does not indicate missing context.
- **FUNCTIONAL-CROSS-CUTTING** — behavioral requirements spanning
  multiple features, under `docs/reqs/functional/cross-cutting/`.
  Every teammate reads all of these before any task.
- **FUNCTIONAL-DATA** — data model and persistence, under
  `docs/reqs/functional/data/`. Read when working on data-related
  tasks.
- **FUNCTIONAL-FEATURE** — primary doc for a specific feature,
  under `docs/reqs/functional/features/`. Read the primary doc
  AND all supplementals for the feature currently being worked
  on.
- **FUNCTIONAL-FEATURE-SUPPLEMENTAL** — additional detail for a
  feature (views, UX, feature-scoped NFRs, etc.); does not stand
  alone. Each entry must include an "Also read" pointer to its
  primary FEATURE doc, and vice versa.
- **EXTERNAL-INTERFACE** — system boundary and interface
  requirements, under `docs/reqs/external-interfaces/`. Read when
  touching those interfaces.
- **ENVIRONMENTAL** — infrastructure and deployment requirements,
  under `docs/reqs/environmental/`. Read when touching deployment
  or infrastructure.
- **TECHNICAL** — stack, tooling, and design constraints
  expressed as requirements, under `docs/reqs/technical/`. Read
  as relevant to the current task.
- **ARCHITECTURAL** — structural debt and design decisions; the
  project-specific architecture entries under
  `docs/architecture/`. The Architect curates; every teammate
  reads relevant entries during their task work (per the Primary
  references in their agent definition).
- **GUIDE** — user-facing guide content under `docs/guides/`.
  Owned by the Tech Writer; teammates rarely consult unless
  release notes or deployment behavior is in scope.
- **GLOSSARY** — entries in `docs/glossary.md`. Curated by the
  Architect, committed by the Analyst. Every teammate reads
  `docs/glossary.md` before starting any task (it is small) so
  requirement links resolve.

Feature-scoped non-functional requirements (e.g., "dashboard loads in
2 seconds") live under the feature as FUNCTIONAL-FEATURE-SUPPLEMENTAL,
not under `docs/reqs/non-functional/`.

### Requirement status convention

Every discrete requirement statement in a doc carries a status
checkbox: `[ ]` not started, `[-]` in progress, `[x]` complete. Under
each requirement sit two kinds of child checkboxes:

- **`implementation`** (one per requirement) — the **Coder** marks
  `[x]` when the requirement's implementation is committed and ready
  for testing. This is the Coder's record of work.
- **`AC1`, `AC2`, ...** (one per acceptance criterion) — the **Tester**
  (Unit Tester or E2E Tester, as applicable) marks `[x]` when an
  automated test that verifies that AC is passing. One test per AC at
  minimum; parameterized ACs may have multiple.

The parent's checkbox is a roll-up that the **Analyst** maintains:
`[x]` only when `implementation` is `[x]` AND every AC is `[x]`;
`[-]` when any child is `[-]` or `[x]` but not all children are `[x]`;
`[ ]` when all children are `[ ]`. See "Status Tracking" below for
transition rules. Example format inside a requirement doc:

    ## Authentication
    - [ ] Users can log in with SSO via SAML 2.0
          ... additional detail and description ...
      - [ ] implementation
      - [ ] AC1: SAML 2.0 metadata exchange supported.
      - [ ] AC2: Authenticated user lands on the post-login redirect target.
      - [ ] AC3: Failed authentication displays a non-technical error.
    - [-] Passkey-based authentication is supported
          ... additional detail and description ...
      - [x] implementation
      - [x] AC1: Passkey registration available from account settings.
      - [-] AC2: Passkey login available on the login view.
      - [ ] AC3: Lost-passkey recovery flow.
    - [x] Session timeout after 30 minutes of inactivity
          ... additional detail and description ...
      - [x] implementation
      - [x] AC1: Inactive session ends after 30 minutes.
      - [x] AC2: User is redirected to login on timeout.
      - [x] AC3: Active sessions are not affected.

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
- `docs/reqs/`                → Analyst
- `docs/patterns/`            → Architect curates / Analyst commits
- `docs/architecture/`        → Architect curates / Analyst commits
- `docs/guides/`              → Tech Writer
- `docs/glossary.md`          → Architect curates / Analyst commits
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
  Requirements Clarification Escalation in Team Coordination
  Procedures below)
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
Each requirement carries a parent status checkbox and three kinds of
child checkboxes — one `implementation` child plus one per
acceptance criterion. All checkboxes use the same notation:
- `[ ]` — not started
- `[-]` — in progress
- `[x]` — complete

**Ownership of each checkbox:**
- `implementation` — **Coder** marks. `[-]` when implementation
  begins; `[x]` when the requirement's implementation is committed
  and ready for testing.
- `AC1`, `AC2`, ... — **Tester** (Unit Tester or E2E Tester per
  scenario type) marks. `[-]` when test authoring or execution
  begins; `[x]` when the automated test that verifies the AC passes.
- Parent requirement — **Analyst** rolls up from the children. The
  Analyst does not author leaf statuses; only the rollup.

**Roll-up rule for the parent (computed by the Analyst):**
- All children `[x]` → parent `[x]`
- Any child `[-]` or `[x]` (but not all `[x]`) → parent `[-]`
- All children `[ ]` → parent `[ ]`

**Status transitions:**
- `[ ]` → `[-]`: the owning role marks the leaf (Coder for
  `implementation`, Tester for an AC) as work begins. The Analyst
  updates the parent roll-up to `[-]` once any child is `[-]` or
  `[x]`.
- `[-]` → `[x]`: the owning role marks the leaf when its work is
  done (Coder when implementation is committed; Tester when the
  test passes). The Analyst updates the parent roll-up to `[x]`
  only when all children are `[x]`. The squash merge carries the
  final state to `<DEV_BRANCH_NAME>`. Dev only ever sees
  `[ ]` → `[x]` transitions.
- `[x]` → `[ ]` or `[-]` → `[ ]`: the Analyst resets a leaf when
  adding a new AC, substantively changing an AC's intent, or
  invalidating prior implementation. (The Analyst owns reset
  authority across all leaves and the parent because a reset is a
  scope-management call, not the leaf-owner's call.) Resetting any
  child may cascade to the parent (a previously-`[x]` requirement
  becomes `[-]` if a child resets to `[ ]`). The Analyst must
  notify the Lead on any reset so the Lead can assess impact on
  active or completed tasks.
- Renaming or moving a requirement or AC does not reset status, but
  the Analyst must update all cross-references (`INDEX.md`, active
  task files in `.claude/.tasks/`).

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
  - The Analyst has no sub-branch — it works on `requirement/<slug>`
    branches and commits status marks directly on the task branch.
  - The Architect has no branch — it reads code on other teammates'
    branches but does not commit.
  - The Tech Writer has no task sub-branch — it works on
    `guide/<slug>` branches on the release cadence, not the task
    cadence.
- Teammate sub-branch operations: each teammate creates their
  sub-branch once at the start of the task and reuses it for all
  commit cycles within that task. Merge (not rebase) in both
  directions — merge FROM the task branch to stay current, merge
  INTO the task branch using the Task Branch Merge Protocol (see
  Team Coordination Procedures below). No teammate commits to
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

### Cost report destinations

The Integrator computes a per-model token/cost delta report at task
conclusion (via `ccusage`, see T.6 in `team-start.md`). It is
**always reported verbally to the human** at task wrap-up
regardless of the destination settings below; the settings control
whether the report is also recorded durably.

- **Include cost report in commit message:** `<COST_IN_COMMIT>`
  (`yes` or `no`). When `yes`, the Integrator appends the report
  as a trailing block of the final squash-merge commit message so
  it persists in git history.
- **Append cost report to project log:** `<COST_IN_LOG>` (`yes` or
  `no`). When `yes`, the Integrator appends the report (with task
  ID, date, and per-model breakdown) to `.claude/.cost-log.md` —
  a project-local log file (gitignored) for cumulative
  developer-local cost tracking. Useful when the team wants
  cost visibility but does not want it in git history.

Both settings can be `yes` (record in both places), one `yes` and
one `no` (single destination), or both `no` (verbal only). Other
destinations may be added over time.

To change either setting later, ask the Lead ("change `Include
cost report in commit message` to `yes`/`no`"). The Lead delegates
to the Integrator, which creates a working branch off
`<DEV_BRANCH_NAME>`, updates the line in `CLAUDE.md`, commits, and
finalizes per the project's merge method above. **Caveat:**
`CLAUDE.md` is read by every role at task start, so a
mid-engagement edit can collide with any task/requirement branch
in flight. Prefer to make this change during a quiet period — no
active tasks, no open requirement branches — so the new value is
in effect uniformly for all subsequent work. If a change must
happen while work is in flight, the Lead should pause new task
creation until the edit is merged and all teammates have pulled
the latest `<DEV_BRANCH_NAME>`.

## Team Coordination Procedures

These are the cross-cutting procedures every teammate may need to
execute or participate in. The Lead's standing instructions
(`team-start.md`) reference these from team-side workflows; teammates
load them via this file (which their CLAUDE.md auto-loads at
session start).

### Mid-Task Architect Escalation

When the Coder encounters a problem during implementation that
requires Architect involvement before committing (see the Coder's
DIAGNOSIS-FIRST FIX PROTOCOL for triggers), use this procedure.

**Triggers** (Coder MUST escalate, not MAY):

- Failure classified as Structural
- 2-attempt fix limit reached for the same root cause (see FIX
  ATTEMPT LIMIT in the Coder's DIAGNOSIS-FIRST FIX PROTOCOL)
- Task requires modifying files or interfaces not identified in
  the task file's scope or the Architect's kickoff guidance (if
  any)
- Need to add a dependency or change a method signature in a
  shared interface

**Escalation message format** (Coder → Architect via `SendMessage`):

```
BLOCKER: [one-sentence description of what failed]
ROOT CAUSE: [one-sentence diagnosis of why it failed]
APPROACH IMPACT: [does this suggest the current approach needs to
  change, or is it a gap in the plan?]
ATTEMPTED: [list fix attempts already made, if any]
FILES TOUCHED SO FAR: [list of files modified since last commit]
```

**Architect response** (one of three outcomes):

1. **TARGETED GUIDANCE** — the approach is sound; here is the
   correct fix with rationale. Coder proceeds.
2. **APPROACH REVISION** — the approach needs to change. Architect
   provides a revised implementation plan for the remaining work.
   Coder reverts uncommitted changes that conflict with the
   revised plan (see REVERT-BEFORE-REWORK in Coder rules), then
   proceeds with the new plan.
3. **SCOPE FLAG TO LEAD** — the problem reveals a gap in
   requirements or a cross-cutting concern that affects other
   tasks. Architect notifies the Lead for task re-scoping.

**Priority:** Mid-task escalations take priority over post-commit
reviews. The Architect should respond before the Coder's next
commit, not after it.

**Coder behavior while waiting:** Do NOT continue building on top
of the blocked code path. Work on an independent part of the task
if one exists, or wait.

### Requirements Clarification Escalation

When any teammate identifies a requirement that is unclear,
ambiguous, conflicting, or insufficiently specified (see
"Requirements Ambiguity — Do Not Guess" above), use this
procedure.

**Step 1 — Teammate raises the ambiguity to the Architect** via
`SendMessage`:

```
AMBIGUITY: [which requirement, with file path and line/section]
CONFLICT/GAP: [what is unclear or contradictory]
OPTIONS: [2-3 concrete interpretations, each with a one-sentence
  consequence for the implementation]
BLOCKED WORK: [what cannot proceed until this is resolved]
```

**Step 2 — Architect attempts internal resolution.**
The Architect searches all project documentation (`docs/`,
`CLAUDE.md`, code comments, commit messages) for evidence that
resolves the ambiguity. If the docs collectively make the answer
clear — even if no single doc states it explicitly — the Architect
records the resolution and rationale in the task file. Work
proceeds.

**Step 3 — Lead escalates to human (if Architect cannot resolve).**
If the Architect cannot resolve the ambiguity from existing docs,
the Architect routes the question to the Lead, who presents it to
the human using the teammate's original format plus the Architect's
research summary. The Lead records the human's answer in the task
file. If the answer reveals a docs gap, the Lead assigns the
Analyst to draft an update to the relevant requirement doc; the
Analyst submits the draft to the Architect for pre-review and
incorporates feedback before submitting to the Lead. Because
requirement docs are human-owned, the Architect-reviewed draft
must be presented to the human for approval before it is committed
(see Analyst rules in `.claude/agents/analyst.md`).

**While waiting for resolution:**

- The teammate may continue working on unambiguous parts of the
  task.
- The teammate MUST NOT implement the ambiguous part using a
  guess, placeholder, or TODO comment — incomplete implementations
  create false progress and mislead other teammates.
- If the ambiguous part blocks ALL remaining work, the teammate
  signals this in the escalation so the Lead knows to prioritize
  the question.

### Subtask Discovery

During implementation, the team may discover that satisfying the
in-scope requirements also requires work not originally in the
plan steps — but this work does NOT require a separate full task
lifecycle. Examples: an additional validation rule, a missing data
migration step, a UI state that wasn't anticipated.

**Procedure:**

1. The teammate reports the discovery to the Lead via
   `SendMessage`.
2. If the work maps to an existing documented requirement: the
   Lead asks the Integrator to add the requirement cross-reference
   to the task file's "Requirements in Scope" section and adds
   new plan steps.
3. If the work requires a new requirement: follow the ad-hoc
   discovery flow (Analyst drafts → Architect pre-review → human
   approves). Once approved, the Lead asks the Integrator to add
   the cross-reference and plan steps.
4. The Analyst marks the newly in-scope ACs as `[-]` and rolls up
   the parent on the task branch.
5. Work proceeds within the same task branch — no suspension
   needed.

### Task Branch Merge Protocol

When any teammate merges their sub-branch into the task branch,
they must follow this protocol to prevent concurrent merges from
creating conflicts:

1. **Announce:** Message all teammates on the task: "I'm merging
   to the task branch."
2. **Hold:** All other teammates hold off on their own merges
   until the announcement in step 5.
3. **Sync:** Merge from the task branch into your sub-branch
   first to pick up any recent changes. Resolve conflicts if
   necessary.
4. **Merge:** Merge your sub-branch into the task branch.
5. **Release:** Message all teammates: "I'm done merging to the
   task branch."

This protocol applies to all teammates that merge into the task
branch (Coder, Unit Tester, E2E Tester), not just during parallel
Coder work. Teammates waiting to merge proceed in the order they
announced.

**Crash recovery:** If a teammate does not post the release
message (step 5) within 5 minutes of the announce (step 1), the
Lead investigates:

1. Check `git log` on the task branch to determine whether the
   merge commit was created.
2. If the merge completed: the Lead posts the release message on
   behalf of the unresponsive teammate, then attempts recovery
   per Teammate Recovery (resume the teammate via `SendMessage`
   first; if resume fails, spawn a replacement from the same
   agent definition).
3. If the merge did not complete (or is partial): the Lead
   reverts any partial merge state on the task branch, then
   recovers the teammate per Teammate Recovery.
4. The Lead notifies all holding teammates before they proceed.

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
   architecture documents
3. `docs/glossary.md` — agnostic vocabulary referenced inline by
   requirement docs (Markdown links). Read this before any
   requirement doc so the linked terms make sense.
4. Every file tagged NON-FUNCTIONAL, FUNCTIONAL-CROSS-CUTTING, or
   ARCHITECTURAL in `docs/INDEX.md`, plus any TECHNICAL,
   ENVIRONMENTAL, or EXTERNAL-INTERFACE docs relevant to your
   current task. Also any PATTERN entry your role's Primary
   references list points at (see your agent definition).
5. `docs/architecture/architecture-debt.md` — known structural debt
6. The FEATURE doc in `docs/INDEX.md` matching your current task, plus
   all FEATURE-SUPPLEMENTAL docs linked from it. Follow inline
   Markdown links from the requirements into `docs/glossary.md` and
   `docs/architecture/` as you encounter them — those are part of the
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
