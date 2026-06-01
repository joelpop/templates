---
name: figma-code-comments
description: When writing comments in Figma-to-Vaadin generated code, keep behavior-explaining comments but omit Figma node IDs, frame names, and requirement-doc paths so comments remain valid as designs evolve.
---

# Comments in Figma-Generated Code

When writing comments in Figma-to-Vaadin generated code, keep comments that explain
non-obvious behavior or invariants, but omit Figma node IDs, frame names, and
requirement-doc paths so comments remain valid as designs evolve.

Design sources rot: Figma nodes get reorganized, requirement docs get renamed, and most
developers reading the code won't have Figma access. Describe the behavior directly instead.

- Bad: `// Spoofing rings per main-header.md — deferred`
- Good: `// Spoofing ring modes (green / amber / red) are absent until impersonation state exists`
- Bad: `// Matches Figma 38:5727 dashboard frame`
- Good: (no comment — the code speaks for itself)
