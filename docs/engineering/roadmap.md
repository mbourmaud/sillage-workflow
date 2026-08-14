---
type: roadmap
status: draft
owners: [maintainers]
---

# Roadmap

## 0.1 — Protocol foundation

- Portable lifecycle and human-authority contract.
- Canonical cold-start project contract.
- Read-only project doctor and transition validator.
- Public task and project schemas.
- First research-skill candidate with evaluation fixtures.

## Full workflow kernel (pilot-ready)

- Deterministic cold-start and task-status views.
- Explicit, atomic local transition writes with optimistic concurrency checks.
- Decision-bound review and handoff artifacts.
- One orchestration skill that routes research and preserves human gates.
- Provider-neutral execution profiles for stage-specific capability and effort.
- An executable end-to-end pilot without merge, deployment, or external writes.

The kernel is exercised by `make pilot` before it is promoted to a release.
The orchestration skill remains a candidate until an independent baseline run
and human review are recorded; this is a release gate, not a missing runtime
capability.

## Next capabilities

Each capability is developed and evaluated independently:

1. Entering and resuming a worktree.
2. Shaping and investigating work.
3. Slicing approved work into independent review units.
4. Binding clean workspaces without assuming Git.
5. Planning risks, acceptance, tests, and runtime evidence.
6. Implementing with focused tests.
7. Designing and cleaning test evidence.
8. Separating verification from review.
9. Producing resumable handoffs.
10. Adapter mappings, usage receipts, and budget enforcement for execution
    profiles.
11. Local, GitHub, GitLab, and other task-store adapters.

The roadmap describes intended capability, not committed behavior. Released
commands and skills remain the authoritative product surface.
