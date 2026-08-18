---
name: test-hygiene
description: Use automatically when tests are growing too quickly, duplicated, flaky, slow, unowned, or confusing; reduce the proof portfolio while preserving risk coverage and deterministic evidence.
metadata:
  namespace: sillage
  qualified-name: "sillage:test-hygiene"
---

# Keep tests useful

Test quality is the confidence-to-cost ratio, not line count or a coverage
badge. Every risk gets one primary owning layer; supporting tests must add a
distinct contract or failure path.

## Hygiene pass

1. Inventory tests by behavior, layer, runtime cost, owner, and flake history.
2. Map each acceptance/risk to its strongest test or runtime observation.
3. Mark duplicates, implementation-detail tests, dead fixtures, broad setup,
   and tests that pass while the real boundary is broken.
4. Keep one clear test at the narrowest stable layer; promote to integration,
   browser, or domain checks only when the risk crosses that boundary.
5. Quarantine or repair flakes with a reproducible cause; never hide failures
   by retrying indefinitely or weakening assertions.
6. Delete obsolete tests only with evidence that their behavior is covered or
   intentionally removed by the approved decision.

```text
TL;DR: <what confidence is preserved and what noise is removed>
Risk map: <risk → primary proof>
Remove/merge: <tests and reason>
Keep/add: <missing boundary or failure path>
Runtime cost: <before/after observation>
Proof: <test lane and real-boundary check>
```

If the cleanup changes product behavior or removes a contractual assertion,
route to `sillage:shape`. Implement and prove the hygiene slice like any other
vertical; never claim success from a smaller test count alone.
