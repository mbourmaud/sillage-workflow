---
type: standard
status: stable
owners: [maintainers]
---

# Durable knowledge

Task cards and optional task stores retain operational history. Repositories
retain only the knowledge needed by future work.

Promote a finding when losing it would force future contributors to rediscover
an important product boundary, domain invariant, architectural constraint,
standard, runbook, or landmark decision. Update a canonical page before adding
a new one.

Do not promote raw research, experiments, discarded approaches, temporary
plans, verification logs, screenshots, or task histories. A task card is the
default owner for local context; an ADR is reserved for a durable,
cross-cutting decision. Git history is not a reason to keep obsolete pages in
the active navigation surface.

Durable pages declare a type, lifecycle status, owner, sources when applicable,
and freshness policy when the knowledge can drift.
