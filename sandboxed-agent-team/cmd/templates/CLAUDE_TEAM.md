# Project: {{PROJECT_NAME}}

## Stack
- Language: Java <version>
- Framework: Vaadin <version>, Spring Boot <version>
- Database: <DATABASE — e.g., PostgreSQL 16>
- Testing:
  - Unit & Browserless UI (Unit Tester): JUnit <JUNIT_VERSION — e.g.,
    5 or 6>; Vaadin Browserless Testing (Vaadin 25.1+:
    `browserless-test-junit{{JUNIT_VERSION}}`, free / Apache 2.0,
    extends `SpringBrowserlessTest`; pre-25.1:
    `vaadin-testbench-junit{{JUNIT_VERSION}}`, commercial, extends
    `SpringUIUnitTest`) for in-process UI component tests
    (browser-less, container-less); Mockito for mocking. One test
    class per production class. Browserless UI tests live in the
    same package as the view (`*Test.java` suffix). Class suffix:
    `*Test.java` = surefire, `*IT.java` = failsafe.
  - End-to-End (E2E Tester): Node.js Playwright
    (`@playwright/test`), TypeScript, in `<e2e-test-dir>/` (e.g.,
    `e2e/`). Vaadin-recommended for E2E.
  - Testing pyramid: unit → browserless UI → E2E. E2E runs only
    at the pre-PR gate, not per-commit.
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
before writing code (Vaadin first; see "Framework Identity" below
for why). The "Primary Users" column is guidance — all servers are
available to all teammates.

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

Training data describes earlier Claude Code versions and conflates
distinct features (notably Agent Teams and subagents). Before
asserting how *anything* in Claude Code works — hooks, skills,
slash commands, the `Agent` tool, Agent Teams, settings fields,
file conventions — verify against `code.claude.com/docs/...`
first. Use the `fetch` MCP server, `WebFetch`, or the
`claude-code-guide` subagent if available. If you cannot verify,
say "I don't know — let me look it up" rather than answering from
recall.

Past sessions burned many hours on designs built atop
confidently-asserted mechanics that turned out wrong. Verifying
takes seconds; drift compounds.

`docs/claude-code-facts.md` records facts that have been verified
against the installed binary alongside the false beliefs they
correct (e.g., `TeamCreate` is real even though the public docs
describe team creation only by behavior; `isolation: worktree` is
valid for spawned teammates; teammates *can* spawn synchronous
subagents). Read it before "fixing" anything that looks wrong in
agent frontmatter or team-related tool calls.

## Documentation Index

See `docs/README.md` for the master pointer to the four-tree
structure:

- `docs/reqs/` — project-specific requirements (Analyst owns).
  IEEE 830 / ISO 29148 (SRS structure) and ISO 25010 (quality
  model).
- `docs/patterns/` — project-agnostic conventions, architecture
  patterns, and recipes for the project's stack. Architect owns;
  designed to extract across projects.
- `docs/solutions/` — how this codebase realizes its requirements:
  non-obvious implementation choices, project-specific pattern
  applications, and known debt. Coder owns.
- `docs/guides/` — install / deploy / user / admin / operator
  guides. Tech Writer owns; release-cadence updates.
- `docs/glossary/business.md` — business and user-facing vocabulary
  (implementation-agnostic terms with optional slang variants).
  Analyst owns.
- `docs/glossary/technical.md` — technical and implementation
  vocabulary. Architect owns.

`docs/README.md` is seeded at setup and is human-owned — edit it
freely as the docs structure evolves; agents do not write to it.

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
- **SOLUTION** — non-obvious implementation choices and known debt
  under `docs/solutions/`. Coder owns; every teammate reads
  relevant entries during task work (per the Primary references in
  their agent definition).
- **GUIDE** — user-facing guide content under `docs/guides/`.
  Owned by the Tech Writer; teammates rarely consult unless
  release notes or deployment behavior is in scope.
- **GLOSSARY** — entries in `docs/glossary/business.md` (business
  terms, Analyst owns) and `docs/glossary/technical.md` (technical
  terms, Architect owns). Every teammate reads both before starting
  any task so requirement links and technical references resolve.

Feature-scoped non-functional requirements (e.g., "dashboard loads in
2 seconds") live under the feature as FUNCTIONAL-FEATURE-SUPPLEMENTAL,
not under `docs/reqs/non-functional/`.

### Requirement status convention

Every requirement statement and every acceptance criterion carries a
two-column status marker — `[D]` for the **draft lifecycle** and `[C]`
for the **code lifecycle**:

    `[D][C]` Requirement or AC text

**Status key:**

| Char | D — draft | C — code |
|------|-----------|----------|
| `[ ]` | Identified — drafting not started | Not started |
| `[-]` | Drafting in progress | Coding in progress |
| `[x]` | Drafted — awaiting approval | Coded — ACs pending verification |
| `[*]` | Approved | Coded and verified |
| `[!]` | Stale — spec needs revision | Stale — code needs update |

**Valid combinations:**

```
[D][C]  meaning
[ ][ ]  identified
[-][ ]  drafting
[x][ ]  drafted
[*][ ]  approved
[*][-]  coding
[*][x]  coded
[*][*]  verified

[!][x]  draft stale, code current
[!][*]  draft stale, code verified
[-][!]  redrafting, code stale
[x][!]  redrafted (pending approval), code stale
[*][!]  re-approved, code stale
```

**Requirement D precondition:** A requirement's D can only reach `[*]`
after all of its ACs' D are already `[*]`. Adding or substantively
changing any AC resets the requirement's D from `[*]` to `[x]`, forcing
a fresh review that the requirement text is still in sync with its AC
set before re-approval.

**Test pass/fail** is tracked by CI, not in requirement status. `[*]`
in C marks a milestone — implementation was verified at the time of
marking — not a live health indicator. A subsequent CI failure does not
reset C status.

Example format inside a requirement doc:

    ## Authentication
    - `[*][*]` Users can log in with SSO via SAML 2.0
               ... requirement description ...
      - `[*][*]` AC1: SAML 2.0 metadata exchange supported
      - `[*][*]` AC2: Authenticated user lands on the post-login redirect target
      - `[*][*]` AC3: Failed authentication displays a non-technical error
    - `[*][-]` Passkey-based authentication is supported
               ... requirement description ...
      - `[*][x]` AC1: Passkey registration available from account settings
      - `[*][-]` AC2: Passkey login available on the login view
      - `[*][ ]` AC3: Lost-passkey recovery flow
    - `[x][ ]` Users can reset their password via a magic link
               ... requirement description ...
      - `[x][ ]` AC1: Magic link sent to registered email
      - `[-][ ]` AC2: Link expires after 15 minutes

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
- `docs/patterns/`            → Architect
- `docs/solutions/`           → Coder
- `docs/guides/`              → Tech Writer
- `docs/glossary/business.md` → Analyst
- `docs/glossary/technical.md`→ Architect
- `pom.xml`                   → COORDINATE (Lead approves)
- `README.md`                 → COORDINATE (Lead approves)
- CI/CD config (e.g., `.github/workflows/`) → COORDINATE (Lead approves)
- `Dockerfile` / `docker-compose.yml` → COORDINATE (Lead approves)
- DB migrations (e.g., `src/main/resources/db/migration/`) → Coder (Architect reviews)

Human-owned — no agent edits; agent-assisted editing happens outside the sandbox only:
- `CLAUDE.md`
- `CLAUDE_TEAM.md`
- `ONBOARDING.md`
- `TEAM_GUIDE.md`
- `docs/README.md`
- `docs/glossary/INDEX.md`
- `.claude/agents/`
- `.claude/commands/`
- `.claude/hooks/`
- `.claude/settings.json`

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
Prefer elegant, idiomatic solutions over verbose ones, AS LONG AS
the code remains readable to a mid-level developer without
special explanation.

- Use enum properties (fields, methods, lambdas) instead of switch
  statements or if/else chains on enum values — behavior belongs
  on the enum, not scattered across consumers.
- Use polymorphism and strategy patterns over type-checking
  conditionals.
- Use composition over inheritance when extending behavior.
- Use functional idioms (map, filter, Optional chaining, Stream
  pipelines) when they make intent clearer than imperative loops.
- If a "clever" solution needs a comment to explain it, it's too
  clever. Refactor until the code explains itself.

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

**MCP server priority — Vaadin first.** Vaadin is a full-stack
application framework, not just a UI library. It owns or wraps
routing, security, server push, session handling, theming, and
many other application concerns with SPA-aware integration glue
around Spring. Defaulting to Spring directly for those concerns
defeats Vaadin's plumbing and produces broken or insecure
behavior. The canonical example is security: a Vaadin app is an
SPA, and Spring Security has no native concept of SPA route
transitions — Vaadin provides a helper that routes view/route
transitions through Spring Security so authorization works as
expected. Bypassing it (e.g., wiring servlet filters or
`@PreAuthorize` directly against Vaadin views) silently disables
view-level access control. Server push, theming, and session
handling have similar story arcs.

Therefore: every teammate consults the `vaadin` MCP server first
for any concern Vaadin owns or wraps. Use `spring-docs` or `java`
only when Vaadin doesn't cover the topic or its answer needs
augmenting. Don't rely on training data — framework APIs evolve
between versions and training data may describe deprecated
patterns. See "Documentation Sources (MCP Servers)" above for the
server list.

**Note for Claude Code:** Customize this section for the project's
framework. The principle — use the framework's idioms, not generic
web patterns — applies to any framework. Replace the Vaadin-specific
content above with the relevant paradigm, anti-patterns, and
documentation sources for the project's actual framework.

## CRITICAL: Requirements Are Not Negotiable
**Teammates must NEVER change project requirements to match their implementation.**

If a requirement specifies a version, library, framework, or approach:
- Use that exact version/library/framework, even with limited
  training data on it. Search documentation, read source, experiment.
- If you cannot make it work after a thorough attempt, MESSAGE THE
  LEAD with what you tried and what failed. Do NOT silently
  downgrade, substitute, or rewrite the requirement.
- When training data ("conventional wisdom") conflicts with project
  docs or code comments, THE PROJECT'S DOCS WIN. Always. Training
  data may be outdated, inapplicable, or wrong for this context.
- Before applying patterns from general knowledge, grep the
  project's docs, README, and code comments for "do not", "don't",
  "avoid", "WARNING", "NOTE" — patterns may be explicitly forbidden.

Silently changing requirements or ignoring in-project docs in favor
of general training is the highest-severity violation and must be
escalated to the Lead immediately by any teammate that notices.

## Requirements Ambiguity — Do Not Guess
Requirements are sometimes unclear, ambiguous, conflicting, or
insufficiently specified. Teammates must escalate, not guess.

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
- Fill gaps from training data or "common sense" — your
  assumptions may contradict the human's intent
- Pick the simplest interpretation because it's easier to implement
- Treat documentation silence as permission OR prohibition —
  silence on a requirement-level question (WHAT the system must
  do, WHAT constraints it must satisfy) is ambiguity, neither "go
  ahead" nor "don't". This covers only WHAT; silence at the HOW
  level (implementation within an existing requirement's bounds)
  is professional judgment, not ambiguity — see the Requirement
  Gate Workflow's refinement rule in `.claude/agents/lead.md`.
- Implement both interpretations and "let the human choose later"
  — creates dead code and doubles the test surface

**What you MUST do:**
- STOP implementation of the ambiguous part (continue with
  unambiguous parts of the same task)
- Escalate to the Architect with the specific ambiguity (see
  Requirements Clarification Escalation below)
- The Architect tries to resolve from existing docs; if not, the
  Lead escalates to the human
- Don't proceed with the ambiguous part until a resolution is in
  the task file

## Conventions
- Commit messages: conventional commits
  (https://www.conventionalcommits.org/) —
  `<type>(<scope>): <description>`. Types: `feat`, `fix`, `docs`,
  `test`, `refactor`, `chore`, `perf`. Scope is the feature or
  component affected (e.g., `feat(auth): add SSO integration`).
  The description should explain *what changed and why*, not
  itemize files or lines touched.
- All PRs require passing tests before merge.
- Do NOT commit directly to `{{DEV_BRANCH_NAME}}`.

## Status Tracking

### Requirement Status

Two-column `[D][C]` markers on every requirement and AC. See
"Requirement status convention" for the full format, key, and valid
combinations.

**Ownership:**
- All `[D][C]` marks — **Analyst** is the sole writer for all
  requirement docs. Other roles notify the Analyst when their portion
  of the lifecycle advances; the Analyst updates the docs:
  - **Coder** notifies Analyst: implementation begins → requirement C `[-]`;
    implementation committed → requirement C `[x]`.
  - **Tester** notifies Analyst: test being written → AC C `[-]`;
    test written → AC C `[x]`; test passing → AC C `[*]`.
  - **Analyst** marks requirement C `[*]` when all of that
    requirement's AC C statuses are `[*]`.

**Transitions:**
- `[ ]` → `[-]` (D): Analyst begins drafting.
- `[-]` → `[x]` (D): Analyst marks when draft is submitted for review.
- `[x]` → `[*]` (D): Human approves (relayed through Lead). Precondition:
  all AC D must already be `[*]`.
- `[*]` → `[x]` (D): Analyst resets when requirement text needs revision.
  Notifies Lead; Lead assesses impact on active and completed tasks.
- `[*]` → `[!]` (D): Analyst marks when approved text is known stale —
  the requirement no longer accurately describes the intended behavior.
- `[*]` → `[!]` (C): Analyst marks when a requirement change invalidates
  existing implementation.
- Adding or substantively changing any AC resets the parent requirement's
  D from `[*]` to `[x]` (Analyst). Analyst notifies Lead.
- Renaming or moving a requirement or AC does not reset status. Analyst
  must update all cross-references: `INDEX.md` and active task files in
  `.claude/.tasks/`.

### Task Plan Status
Each task file in `.claude/.tasks/<task-id>.md` tracks progress at the
plan-step level. Steps are role-assigned and use the same checkbox
notation. Each teammate marks their own steps as `[-]` when starting
and `[x]` when done.

### Project Status
`.claude/.progress.md` is a minimal dispatcher so the Lead can
recover current state after context compaction. It answers: "which
task am I working on, and what else is parked?" and "which
requirement branches are in flight?"

`progress.md` is gitignored local metadata. Branch operations
don't affect it; it persists across branch switches and task
suspension/resumption. Only IDs and one-line labels — all detail
lives in task files and requirement docs.

**Single writer:** Only the Integrator writes `.claude/.progress.md`.
On state changes (new active task, task suspended/resumed,
requirement branch created/merged), the Lead directs the
Integrator to update; the Integrator also updates proactively as
part of its workflows (e.g., Integration Merge Workflow).
Single-writer prevents concurrent writes.

Structure:
```markdown
# Progress

## Active Task
- <task-id>: <one-line description>
  <!-- Optional indented annotation if the task is resumed but held: -->
  - blocked on `{{DEV_BRANCH_NAME}}` health since <ISO 8601 UTC>

## Suspended Tasks
- <task-id>: Blocked by <prerequisite task-id or description>

## Requirement Branches
- requirement/<slug>: <status> — <one-line description>

## Pattern Branches
- pattern/<slug>: <status> — <one-line description>
```

Requirement and pattern branch statuses:
- `drafting` — Analyst/Architect is actively working on this branch
- `awaiting-approval` — draft submitted to the Lead for human review
- `approved` — human approved; ready to merge to `{{DEV_BRANCH_NAME}}`
- `merged` — merged to `{{DEV_BRANCH_NAME}}`; branch can be deleted

## Branching
- Development branch: `{{DEV_BRANCH_NAME}}` (e.g., `develop`)
- Requirement branches: `requirement/<slug>` — branched off
  `{{DEV_BRANCH_NAME}}` by the Integrator for the Analyst to draft
  requirement docs. One branch per topic or related group (e.g.,
  `requirement/authentication`, `requirement/dashboard-v2`), not
  per individual requirement; the Analyst freely splits, merges,
  and cross-references requirements within a group. Multiple
  requirement branches can be in flight simultaneously.
  Squash-merged back to `{{DEV_BRANCH_NAME}}` after human approval.
  Tracked in `.claude/.progress.md`.
- Pattern branches: `pattern/<slug>` — branched off
  `{{DEV_BRANCH_NAME}}` by the Integrator for the Architect to
  author `docs/patterns/` entries and `docs/glossary/technical.md`
  additions. One branch per topic or related group. Multiple
  pattern branches can be in flight simultaneously.
  Squash-merged back to `{{DEV_BRANCH_NAME}}` after human approval.
  Tracked in `.claude/.progress.md`.
- Task branches: `task/<task-id>` — branched off
  `{{DEV_BRANCH_NAME}}` by the Integrator for each implementation
  task.
- Teammate sub-branches: `task/<task-id>/<role>`:
  - `task/<task-id>/coder` (or `coder-a` … `coder-<n>` when the
    Lead splits a task across parallel Coders, up to
    `{{MAX_PARALLEL_CODERS}}` — see Parallel Subtask Coders in
    `.claude/agents/lead.md`)
  - `task/<task-id>/unit-tester`
  - `task/<task-id>/e2e-tester`
  - Analyst: no sub-branch — works on `requirement/<slug>` and
    commits status marks directly on the task branch.
  - Architect: no task sub-branch — reads code on others' branches
    but doesn't commit to task branches. Works on `pattern/<slug>`
    branches for `docs/patterns/` and `docs/glossary/technical.md`.
  - Tech Writer: no task sub-branch — works on `guide/<slug>` on
    the release cadence, not task cadence.
- Sub-branch operations: each teammate creates their sub-branch
  once and reuses it for all commit cycles within the task. Merge
  (not rebase) in both directions — FROM the task branch to stay
  current, INTO via the Task Branch Merge Protocol (below). No
  teammate commits to another's branch. Sub-branches are local
  only; only `{{DEV_BRANCH_NAME}}` interacts with the remote (via
  the Integration Merge Workflow).
- Merge strategy: squash for all branch-to-`{{DEV_BRANCH_NAME}}`
  merges — keeps history clean but loses per-commit granularity,
  so the squash commit message must capture key decisions and
  affected components (see Integration Merge Workflow T.5 in
  `.claude/agents/lead.md`).
- Merge method: `{{MERGE_METHOD}}`

### Cost report destinations

The Integrator computes a per-model token/cost delta report at
task conclusion (via `ccusage`, see T.6 in
`.claude/agents/lead.md`). It is **always reported verbally to
the human** at task wrap-up regardless of the settings below;
settings control durable recording.

- **Include cost report in commit message:** `{{COST_IN_COMMIT}}`
  (`yes` or `no`). When `yes`, the Integrator appends the report
  to the final squash-merge commit message so it persists in git.
- **Append cost report to project log:** `{{COST_IN_LOG}}` (`yes`
  or `no`). When `yes`, the Integrator appends the report (task
  ID, date, per-model breakdown) to `.claude/.cost-log.md` — a
  gitignored project-local log for cumulative developer-local
  cost tracking. Useful when the team wants cost visibility but
  not in git history.

Both `yes` (record everywhere), mixed (single destination), or
both `no` (verbal only) are valid. Other destinations may be
added later.

To change a setting later, ask the Lead ("change `Include cost
report in commit message` to `yes`/`no`"). The Lead delegates to
the Integrator, which creates a working branch off
`{{DEV_BRANCH_NAME}}`, updates `CLAUDE.md`, commits, and finalizes
per the project's merge method. **Caveat:** `CLAUDE.md` is read by
every role at task start, so a mid-engagement edit can collide
with any task/requirement branch in flight. Prefer a quiet period
(no active tasks, no open requirement branches) so the new value
applies uniformly. If a change must happen while work is in
flight, the Lead pauses new task creation until the edit is
merged and all teammates have pulled the latest
`{{DEV_BRANCH_NAME}}`.

### Workflow settings

These settings control what the kit *enforces*, not what the human
*can* do. Any setting can be changed by asking the Lead — the Lead
confirms the new value, explains the implication, and instructs
the Integrator to update `CLAUDE.md` via a working branch.

#### Existing code requirements: `{{EXISTING_CODE_REQS}}`

Governs how agents treat requirements relative to the codebase
when requirements are silent or incomplete. Requirements are
always consulted first regardless of this value — the difference
is what wins when requirements have gaps.

- **`explicit`** — Requirements are the authoritative
  specification. Any gap needed for the current task must be
  drafted and approved before coding proceeds; inferring intent
  from the code is not a substitute. Choose this when requirements
  are the primary source of intent for all team members.
- **`implicit`** — Requirements are consulted but gaps are
  acceptable. The code serves as supplementary guidance; agents
  infer intent from the existing implementation when requirements
  are silent. Choose this when the team wants flexibility to code
  without filling every requirement gap first.

Neither value implies requirements are complete. Changing this
setting adds or removes strictness going forward with no
implication that the codebase or requirements have changed:
- `implicit` → `explicit`: Code can no longer fill requirement
  gaps; any missing requirement must be drafted first.
- `explicit` → `implicit`: Code may fill requirement gaps without
  drafting new requirements first.

**Tension note:** `explicit` + `FEATURE_WORKFLOW=req-first` is the
strictest combination — any gap in requirements needed for the
current task must be explicitly drafted before coding proceeds.
`implicit` + `FEATURE_WORKFLOW=req-first` also creates friction
(changes hit the requirement gate for code with no formal
requirement), but agents may use the code to inform a draft before
proceeding. The Lead surfaces these tensions at kickoff.

#### Feature workflow: `{{FEATURE_WORKFLOW}}`

Controls whether new features and behavior changes require an
approved requirement before coding starts.

- **`req-first`** — Requirement Gate Workflow runs before the
  Coder starts: Analyst drafts, Architect pre-reviews, human
  approves, then implementation begins. Default for doc-centric
  projects.
- **`code-first`** — Coder implements from intent; Analyst
  backfills requirement docs after the fact.

#### Bug workflow: `{{BUG_WORKFLOW}}`

Controls whether documentation gaps are resolved before the code
fix. Code-level root-cause analysis is always expected regardless
of this setting; this setting controls only whether the *doc gap*
is diagnosed and fixed first.

- **`doc-first`** — Lead diagnoses what doc gap (missing
  requirement / AC / pattern / architecture entry) the bug exposes
  and fixes it before routing the code fix to the Coder. Slower
  per bug, but each bug strengthens the durable artifacts.
- **`fix-first`** — Lead routes straight to the Coder for the code
  fix. Faster per bug; doc gap harvesting happens separately or not
  at all.

## Team Coordination Procedures

These are the cross-cutting procedures every teammate may need to
execute or participate in. The Lead's standing instructions
(`.claude/agents/lead.md`) reference these from team-side workflows; teammates
load them via this file (which their CLAUDE.md auto-loads at
session start).

### Mid-Task Architect Escalation

When the Coder hits a problem that needs Architect involvement
before committing (see the Coder's DIAGNOSIS-FIRST FIX PROTOCOL
for triggers), use this procedure.

**Triggers** (Coder MUST escalate, not MAY):

- Failure classified as Structural
- Fix-attempt limit reached for the same root cause
  (`{{FIX_ATTEMPT_LIMIT}}` consecutive attempts; see FIX ATTEMPT
  LIMIT in the Coder's DIAGNOSIS-FIRST FIX PROTOCOL)
- Task requires modifying files or interfaces outside the task
  file's scope or the Architect's kickoff guidance
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

**Architect response.** Default to evaluating the Coder's diagnosis
before evaluating the approach. An exhausted Coder arriving with
`{{FIX_ATTEMPT_LIMIT}}` failed attempts creates pressure to revise
the architecture — that's the wrong instinct. Most fix-limit
escalations are diagnosis errors, not architectural ones. The
sound approach should stand unless the evidence shows it genuinely
cannot be implemented correctly. Drifting from a good architecture
poorly implemented to a worse architecture more easily implemented
is the failure mode this section guards against.

One of three outcomes:

1. **TARGETED GUIDANCE** *(default)* — the approach is sound; the
   Coder's diagnosis or fix is wrong; here is the correct fix with
   rationale. Coder proceeds. Most escalations land here.
2. **APPROACH REVISION** *(reserved)* — the approach itself cannot
   be made to work as designed. The Coder's struggle is not from
   oversight but from genuine architectural mismatch with the
   problem. Architect provides a revised implementation plan for
   the remaining work. Coder reverts uncommitted changes that
   conflict with the revised plan (see REVERT-BEFORE-REWORK in
   Coder rules), then proceeds with the new plan.
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
The Architect searches project documentation (`docs/`,
`CLAUDE.md`, code comments, commit messages) for evidence. If
the docs collectively make the answer clear — even if no single
doc states it explicitly — the Architect records the resolution
and rationale in the task file. Work proceeds.

**Step 3 — Lead escalates to human (if Architect can't resolve).**
The Architect routes the question to the Lead, who presents it
to the human using the teammate's original format plus the
Architect's research summary. The Lead records the human's
answer in the task file. If the answer reveals a docs gap, the
Lead assigns the Analyst to draft an update; the Analyst submits
the draft to the Architect for pre-review, incorporates feedback,
then submits to the Lead. Requirement docs are human-owned, so
the Architect-reviewed draft must be approved by the human
before commit (see Analyst rules in `.claude/agents/analyst.md`).

**While waiting:**

- Continue working on unambiguous parts of the task.
- Do NOT implement the ambiguous part with a guess, placeholder,
  or TODO comment — incomplete implementations create false
  progress and mislead others.
- If the ambiguous part blocks ALL remaining work, signal this
  in the escalation so the Lead can prioritize.

### Subtask Discovery

During implementation, satisfying in-scope requirements may need
work not in the original plan steps — but not enough to warrant
a separate task lifecycle. Examples: an additional validation
rule, a missing migration step, an unanticipated UI state.

**Procedure:**

1. Teammate reports the discovery to the Lead via `SendMessage`.
2. If the work maps to an existing documented requirement: the
   Lead asks the Integrator to add the requirement cross-reference
   to the task file's "Requirements in Scope" and adds plan steps.
3. If the work needs a new requirement: follow the ad-hoc
   discovery flow (Analyst drafts → Architect pre-reviews → human
   approves). Once approved, the Lead asks the Integrator to add
   the cross-reference and plan steps.
4. Analyst marks the newly in-scope ACs as `[-]` and rolls up
   the parent on the task branch.
5. Work continues on the same task branch — no suspension.

### Task Branch Merge Protocol

When any teammate merges their sub-branch into the task branch,
follow this protocol to prevent concurrent-merge conflicts:

1. **Announce:** Message all teammates on the task: "I'm merging
   to the task branch."
2. **Hold:** All other teammates hold off on merges until the
   release in step 5.
3. **Sync:** Merge from the task branch into your sub-branch
   first to pick up recent changes. Resolve conflicts.
4. **Merge:** Merge your sub-branch into the task branch.
5. **Release:** Message all teammates: "I'm done merging to the
   task branch."

This applies to all teammates that merge into the task branch
(Coder, Unit Tester, E2E Tester), not just parallel Coder work.
Teammates waiting proceed in the order they announced.

**Crash recovery:** If a teammate doesn't post the release
message within 5 minutes of announcing, the Lead investigates:

1. Check `git log` on the task branch for the merge commit.
2. If the merge completed: Lead posts release on behalf of the
   unresponsive teammate, then runs Teammate Recovery (resume via
   `SendMessage`; if resume fails, spawn a replacement from the
   same agent definition).
3. If the merge didn't complete or is partial: Lead reverts
   partial merge state, then runs Teammate Recovery.
4. Lead notifies all holding teammates before they proceed.

## What NOT to do
- Do not add new dependencies without messaging the Lead.
- Do not modify CI/CD pipeline files without explicit approval.
- Do not store secrets in code. Use environment variables.

## Technical Debt
See `docs/solutions/technical-debt.md` for known structural debt and
recommended resolutions.

## Non-Functional Requirements
See `docs/reqs/non-functional/` for performance, security, reliability,
usability, and other quality attribute requirements.

## Context Compaction Warning
<!-- SYNC NOTE: The file list below is duplicated in the Pre-Task
     Context Check in `.claude/agents/lead.md`. If you update one, update both. -->
This file is read at session start but may be LOST when context
compaction occurs. You cannot reliably detect compaction. Before
starting ANY task, verify you still have the needed context by
re-reading these files in order:

1. `CLAUDE.md` (this file) — stack, ownership rules, critical constraints
2. `docs/README.md` — master list of all requirement, glossary, and
   architecture documents
3. `docs/glossary/business.md` — agnostic business vocabulary referenced
   inline by requirement docs (Markdown links). Read this before any
   requirement doc so the linked terms make sense.
   `docs/glossary/technical.md` — technical implementation vocabulary
   curated by the Architect; read before any solutions or pattern docs.
4. Every file tagged NON-FUNCTIONAL, FUNCTIONAL-CROSS-CUTTING, or
   ARCHITECTURAL in `docs/README.md`, plus any TECHNICAL,
   ENVIRONMENTAL, or EXTERNAL-INTERFACE docs relevant to your
   current task. Also any PATTERN entry your role's Primary
   references list points at (see your agent definition).
5. `docs/solutions/technical-debt.md` — known structural debt
6. The FEATURE doc in `docs/README.md` matching your current task, plus
   all FEATURE-SUPPLEMENTAL docs linked from it. Follow inline
   Markdown links from the requirements into `docs/glossary/business.md`
   and `docs/solutions/` as you encounter them — those are part of the
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
