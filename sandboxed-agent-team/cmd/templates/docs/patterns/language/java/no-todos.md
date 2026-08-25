# No TODOs, FIXMEs, or XXX in Committed Code

When committing code with unfinished work, do not mark it with `// TODO`, `// FIXME`,
or `// XXX` — those comments are rarely revisited and become permanent noise. Leave
the work broken and obvious (won't compile, throws, or fails a test) so it cannot be
missed, or record it in `docs/reqs/open-items.md` and reference the item in the commit
message.
