---
type: domain-index
status: stable
owners: [maintainers]
---

# Sillage domain language

- **Project** — the product and repository context in which work occurs.
- **Task** — one desired outcome with intent, scope, acceptance, and lifecycle.
- **Slice** — one cohesive, independently understandable part of a Task. An
  implementation Worktree has exactly one active Slice.
- **Worktree** — an isolated workspace bound to a Slice. A Git worktree is one
  adapter, not the definition.
- **Task Store** — the authoritative operational record. It may be local or
  provided by GitHub, GitLab, Jira, Linear, or another adapter.
- **Durable Knowledge** — current product, design, domain, architecture,
  standard, runbook, or landmark-decision knowledge needed by future work.
- **Evidence** — a fresh, addressable observation supporting an Acceptance
  Criterion. Test output, runtime state, screenshots, traces, and manual checks
  are possible evidence kinds. Transition evidence must be observed no earlier
  than the current human approval and carry the same decision digest.
- **Verification** — the argument that evidence proves the approved behavior.
- **Review** — an independent assessment of correctness, scope, safety, and
  maintainability.
- **Handoff** — a self-contained outcome or resume point.
- **Approval** — an explicit human authorization for a governed transition or
  consequential action. Decision approvals are bound to the deterministic
  digest of the intent, acceptance criteria, and slices they authorize.
- **Waiver** — explicit human acceptance of missing required evidence and its
  residual risk. Like evidence, it is invalidated when the decision digest
  changes.

## Lifecycle

`INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF`

`BLOCKED` records an external dependency, missing authority, or unavailable
evidence together with the exact resume condition.
