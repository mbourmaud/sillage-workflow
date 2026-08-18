---
name: build
description: Use automatically when an approved Sillage slice is ready for implementation, coding, refactoring, migration, or a focused test-first change in a clean worktree.
metadata:
  namespace: sillage
  qualified-name: "sillage:build"
---

# Build one approved slice

Implement only the accepted boundary. If no accepted decision and slice exist,
route back to `sillage:shape` or `sillage:slice`; do not infer approval from a
casual “go ahead”.

## Build loop

1. Re-read the task TL;DR, acceptance, non-goals, and domain terms.
2. Inspect the current code and tests at the narrowest useful boundary.
3. Write a focused failing test for observable behavior before the fix when a
   code change is required.
4. Implement the smallest clear change. Prefer existing primitives and
   contracts; avoid speculative abstractions, broad cleanup, dead code, and
   compatibility layers without a decision.
5. Run the smallest relevant check while iterating.
6. Stop when the slice is true, or split/return to `shape` when it grows.

Use a specialist as a review lens when the change concerns SOLID, domain
invariants, architecture, test seams, interfaces, systems behavior, security,
platform lifecycle, UI architecture, or data modeling. Select the smallest
pattern that answers an observed pressure; do not add layers, retries,
repositories, events, or abstractions for imaginary variants. The lens informs
code; it does not replace the task decision or invent a new workflow.

Keep temporary logs, experiments, generated files, and copied tutorials out of
durable knowledge. Record only material deviations in the active task. When
the implementation is ready, route to `sillage:prove`; never call code
inspection a completion proof.
