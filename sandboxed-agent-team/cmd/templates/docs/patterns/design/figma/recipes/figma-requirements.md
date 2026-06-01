---
name: figma-requirements
description: How to handle discrepancies between Figma frames and docs/reqs requirement documents before coding — surface, propose resolution, identify the authoritative source, and update the stale artifact.
---

# Figma ↔ Requirements: Resolving Discrepancies

When a Figma frame and a requirement doc under `docs/reqs/` disagree, surface the
discrepancy before coding and identify which artifact is authoritative so the stale
one can be updated.

Figma and the requirement docs under `docs/reqs/` are two views of the same design.
They will drift — a Figma frame may be a demo of one view (so its chrome includes
things that belong to that view, not the shell), or the requirement doc may describe
an older ordering that Figma has since revised. **When they disagree, stop and surface
the discrepancy. Do not silently pick one side.**

- Describe both sides concretely: what Figma shows (without citing node IDs) and what the requirement doc says.
- Propose a resolution and the reasoning.
- Ask which is authoritative, so the stale artifact can be updated.
- Apply the resolution to the implementation AND to whichever artifact was stale (usually the requirement doc, via the Analyst).

A Figma frame is often a **subset** of the full app's intended chrome — it shows what's
needed to mock a particular view. Elements missing from Figma are not necessarily removed
from the spec; elements present in Figma are not necessarily additions to the spec. Treat
Figma as authoritative for **styling** (gaps, fonts, colors, component choice) and the
requirement docs as authoritative for **composition** (what elements exist in a shell
region, their ordering, their conditions) — and when those roles conflict, surface and ask.
