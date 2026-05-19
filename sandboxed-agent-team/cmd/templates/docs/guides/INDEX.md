# Guides — Index

End-user, administrator, and operator-facing guides for using the
running system. Distinct from `docs/reqs/`, `docs/solutions/`,
and `docs/patterns/` — those trees are about *building* the system;
this one is about *using* it.

Owned by the **Tech Writer**. Updated when a release ships, not on
every implementation task — guide content tracks releases, not
internal development cadence.

## Audience-based organization

Guides are typically organized by reader role. Replace the
sub-directories below with what fits the project's actual user
roles:

| Path | Audience |
|------|----------|
| `user-guide/` | End users — day-to-day usage of the application |
| `admin-guide/` | Administrators — configuration, account management, settings |
| `operator-guide/` | Operators — deployment, monitoring, incident response |
| `release-notes/` | All audiences — what changed in each release |

Each subdirectory carries its own `INDEX.md` listing the topics
covered.

## Guide-writing conventions

- **Audience-first.** Each guide names its audience in the first
  line. A reader should know within seconds whether they're in the
  right place.
- **Task-oriented headings.** "Reset a user's password" beats
  "Password Reset Functionality" — guides answer "how do I do X?",
  not "what is X?".
- **Stable across releases.** Avoid embedding screenshots tied to
  internal-version-specific UI; prefer descriptive callouts.
  Screenshots are useful but they age — expect to refresh them per
  release.
- **Reference, don't restate, requirements.** When a guide describes
  intended behavior that originates in `docs/reqs/`, link to the
  requirement rather than repeating it. The guide explains *how to
  use* the behavior; the requirement defines *what the behavior
  is*.
- **Loose lockstep with code.** Guides update at release boundaries,
  not on every commit. A small UI tweak that doesn't affect user
  workflows usually doesn't trigger a guide update.

## Cross-references

| To link a guide entry to | Use |
|--------------------------|-----|
| A requirement (the *what*) | Relative link to `docs/reqs/...` |
| An architecture entry (the *how*) | Relative link to `docs/solutions/...` (mostly for developer-facing operator guides) |
| A pattern (style or recipe) | Usually unnecessary in guides; patterns are internal-facing |
