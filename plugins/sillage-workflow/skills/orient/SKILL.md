---
name: orient
description: Use automatically when entering or resuming a Sillage worktree, when the project or task context is unclear, or when a fresh agent needs a concise status snapshot before acting.
metadata:
  namespace: sillage
  qualified-name: "sillage:orient"
---

# Orient a fresh worktree

Orient before material work, after a context switch, and whenever the user
says “continue”, “where are we?”, or “what is left?”. This is a read-only
cartography pass, not a planning or implementation pass.

## Read in order

1. `PRODUCT.md` — purpose, users, promises, and non-goals.
2. `DESIGN.md` — interaction and architecture boundaries.
3. `AGENTS.md` and any `CLAUDE.md` — local operating rules.
4. `docs/domain/index.md` — canonical vocabulary and invariants.
5. The active task/slice record, if one exists.
6. Only the files and checks named by the active task.

Use the project's native commands and profile when available. A Sillage CLI,
JSON schema, or task store is optional reference tooling; never block a
Markdown-first workflow because one is absent. Do not scan the whole history
or create notes merely to prove that orientation happened.

## Return a state card

```text
Status: <INTAKE|INVESTIGATE|DECIDE|IMPLEMENT|VERIFY|REVIEW|HANDOFF|BLOCKED>
TL;DR: <one sentence describing the current slice>
Known: <two to five verified facts>
Missing: <at most three material gaps>
Decision: <accepted boundary, or awaiting human decision>
Next: <one smallest safe action>
Proof needed: <observable check for the next gate>
```

If no task exists, propose the smallest task card in the conversation and ask
one focused question only when the intent, consumer, or authority is missing.
Do not silently create a large plan, ADR, or task hierarchy.

## Route, do not duplicate

After the state card, hand off to exactly one next skill:

- unclear outcome, alternatives, or product boundary → `sillage:shape`;
- accepted task without an independent slice → `sillage:slice`;
- approved slice ready to change → `sillage:build`;
- implementation claiming proof → `sillage:prove`;
- proof ready for an independent assessment → `sillage:review`;
- accepted review ready to resume elsewhere → `sillage:handoff`.

Use `sillage:research` only for unfamiliar or version-sensitive facts. Use a
specialist (`sillage:solid`, `sillage:ddd`, `sillage:frontend-architecture`,
`sillage:relational-data`, `sillage:document-data`, `sillage:audit`,
`sillage:migrate`, `sillage:debug`, or `sillage:test-hygiene`) as a focused lens
inside the current stage, never as a competing lifecycle.
