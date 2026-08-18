---
name: audit
description: Use automatically when a project needs a codebase, architecture, dependency, security, legacy, or technical-debt audit; separate observed facts from risk hypotheses and leave a ranked, bounded remediation path.
metadata:
  namespace: sillage
  qualified-name: "sillage:audit"
---

# Audit without creating a second backlog

An audit is a decision input, not permission for a rewrite. Inspect the stated
boundary, current behavior, tests, ownership, and operational evidence before
calling something debt.

## Audit packet

```text
TL;DR: <one-sentence health/risk summary>
Scope: <paths, runtime, and date>
Observed: <facts with references>
Risk: <impact, likelihood, and affected consumer>
Unknown: <what was not measured>
Keep: <working behavior worth preserving>
Remediate: <one smallest path, ranked by risk>
Defer: <explicit non-goals and why>
Proof: <check that would show the risk is reduced>
```

Use repository-native tests, dependency reports, runtime traces, and user
flows. Do not infer quality from line count, folder aesthetics, coverage
percentage, or a linter alone. Consolidate repeated findings into one owner
and one remediation slice. Keep raw inventories and screenshots out of durable
knowledge; promote only a stable risk or boundary to its canonical page.

If the audit changes product behavior, migration scope, or a destructive plan,
route to `sillage:shape`. If it reveals a safe sequence of changes, route to
`sillage:migrate`; otherwise leave the project at its current gate.
