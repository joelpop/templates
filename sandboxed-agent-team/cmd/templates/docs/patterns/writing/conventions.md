# Writing Convention and Pattern Documents

Every convention and pattern document in this kit should follow two structural
rules so agents and developers know when and how to apply the practice.

## Open with a scope statement

The first sentence (before any section heading) should answer: "what should this
be applied to?" State the obligation directly:

> Every Vaadin view and custom component class should have a layout diagram in
> its class Javadoc.

Without a scope statement the document describes mechanics with no stated
audience or trigger. A reader scanning the file cannot tell whether it applies
to their current task.

If the document covers multiple unrelated practices, it should be split into
separate single-practice files — possibly grouped into a subdirectory — so each
file can carry its own concise, unambiguous obligation.

## Write INDEX.md entries that lead with the obligation

The INDEX.md description is the only thing an agent sees when scanning the index.
It should answer "when do I apply this?" at a glance, not enumerate topics:

**Avoid** — topic list only:
> Text-based layout diagrams in Javadoc `<pre>` blocks: placement, box labeling…

**Prefer** — obligation first, then topics:
> Every Vaadin view and custom component class should have a layout diagram in
> its class Javadoc. Defines placement, box labeling…