---
name: debug
description: Use automatically when behavior is failing, flaky, surprising, or environment-dependent; perform a bounded, evidence-led diagnosis before proposing or applying a fix.
metadata:
  namespace: sillage
  qualified-name: "sillage:debug"
---

# Diagnose before fixing

Debugging narrows a symptom to a reproducible cause. Do not patch the first
plausible line, broaden the task, or hide an environmental failure as a
product fix.

## Diagnosis loop

1. State expected versus observed behavior, boundary, version, and impact.
2. Reproduce with the smallest command, fixture, request, or user path.
3. Separate facts, hypotheses, and unknowns; inspect logs/traces/configuration
   at the actual failing boundary.
4. Use a binary split or minimal instrumentation to eliminate hypotheses.
5. Add a regression test or deterministic reproduction before the fix when
   behavior is code-owned.
6. Propose the smallest fix, its failure modes, and the proof needed.

```text
TL;DR: <symptom and likely boundary>
Reproduction: <exact observation>
Facts: <verified evidence>
Hypotheses: <ranked, falsifiable>
Cause: <confirmed or still unknown>
Fix boundary: <one slice>
Proof: <regression and runtime checks>
```

If the cause is external or not reproducible, record a `BLOCKED` condition
instead of inventing certainty. Route a scope change to `sillage:shape`, the
fix to `sillage:build`, and proof to `sillage:prove`.
