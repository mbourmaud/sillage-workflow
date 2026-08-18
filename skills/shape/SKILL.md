---
name: shape
description: Use automatically when a request needs clarification, goals, context, alternatives, an architecture choice, an ADR, or an approved decision before implementation; ask one focused question at a time.
metadata:
  namespace: sillage
  qualified-name: "sillage:shape"
---

# Shape one decision

Shape turns an ambiguous request into a small, human-owned decision. It does
not implement, expand scope, or create a document for every thought.

## Produce a decision packet

```text
TL;DR: <chosen outcome in one sentence>
Intent: <problem and consumer>
Goals: <observable outcomes>
Non-goals: <explicit exclusions>
Context: <facts, constraints, and domain terms>
Options: <only material alternatives>
Decision: <proposed boundary and why>
Risks: <failure mode and recovery condition>
Proof: <checks that will make the decision believable>
Approval: awaiting human | accepted by <identity> at <time>
```

Classify the amount of ceremony before presenting options:

- `probe` — answer or throwaway investigation; no durable code or document.
- `bounded` — one existing behavior or boundary; a short in-chat decision.
- `cross-cutting` — shared contract, product behavior, or architecture; a
  structured decision and ADR only when the choice is durable.

The classification changes the size of the artifact, never the human gate.
When a task grows beyond its classification, stop and reshape it.

Read the project entry points and active task first. Use `sillage:research` for
external or version-sensitive facts; label facts, inference, conflicts, and
unknowns separately. Use `sillage:ddd` or a data/frontend specialist only when
that lens changes the decision.

## Durable decisions only

Update an existing canonical page or ADR only when the choice is durable,
cross-cutting, and worth rediscovering. Keep the decision packet in the active
task when it is local or reversible. An ADR must state context, decision,
consequences, alternatives, and the owner; it is never a transcript.

Stop at the human gate when the answer changes product behavior, scope,
destructive action, migration, external write, or an evidence waiver. Ask one
focused question, then wait. Once accepted, hand the exact boundary to
`sillage:slice`.
