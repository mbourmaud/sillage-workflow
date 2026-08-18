---
name: migrate
description: Use automatically when changing a legacy boundary, schema, API, dependency, architecture, or data representation; create a staged, observable migration path with compatibility and rollback limits.
metadata:
  namespace: sillage
  qualified-name: "sillage:migrate"
---

# Migrate in reversible slices

Migration work is a sequence of small contracts, not a permission slip for a
big-bang rewrite. Preserve working behavior until the replacement is proven.

## Migration packet

```text
TL;DR: <target state and first safe step>
Current: <observed contract and consumers>
Target: <new contract and explicit non-goals>
Stages: <expand → backfill/dual-read → switch → contract, as applicable>
Compatibility: <old/new readers, writers, versions>
Rollback: <trigger, safe reversal, and data limit>
Proof per stage: <deterministic and runtime observations>
Cleanup: <when old code/data can be removed>
```

Inventory real consumers and ownership before editing. For data, state
invariants, idempotency, ordering, retention, locks, volume, and reconciliation.
For APIs/dependencies, define version negotiation and failure behavior. Keep
one stage per worktree/vertical and stop when a stage needs a new decision.

Never delete or contract an old path without human approval and evidence that
the replacement is live and recoverable. Route a changed target or risk to
`sillage:shape`; route each implemented stage to `sillage:build` and
`sillage:prove`.
