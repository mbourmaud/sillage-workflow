---
name: handoff
description: Use automatically when an accepted review must be summarized for a fresh agent, a task is being paused or resumed, or a vertical is ready for a human-owned merge, release, or deployment decision.
metadata:
  namespace: sillage
  qualified-name: "sillage:handoff"
---

# Hand off without context loss

Handoff is the smallest durable resume point. It is not a release claim and it
does not imply that an external action happened.

## Write one handoff card

```text
Status: HANDOFF | BLOCKED
TL;DR: <what is true now>
Changed: <files or boundaries, summarized>
Decision: <decision and digest>
Proof: <checks and references>
Review: <accepted verdict and reviewer, or why blocked>
Known limits: <unresolved risk or none>
Next safe action: <one action a fresh agent can take>
Human gate: <merge, release, deployment, external write, or none>
Resume condition: <required only for BLOCKED>
```

Read the task record and current diff before writing. Keep detail in the
existing task/knowledge owner; do not create a new report when an existing
handoff section is available. If review is not accepted or evidence is stale,
route back to `sillage:prove`, `sillage:review`, or the recorded blocked
status. A clean handoff always names its limits and the next safe action.
