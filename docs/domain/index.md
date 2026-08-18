---
type: domain-index
status: stable
owners: [maintainers]
---

# Sillage domain language

- **Project** — the product and repository context in which work occurs.
- **Task** — one desired outcome with intent, scope, acceptance, and lifecycle.
- **Task Card** — the compact Markdown resume point for a Task: TL;DR, intent,
  decision, Slice, proof, review, handoff, and the next safe action. It is the
  portable source of context; a task-store adapter may mirror it.
- **Slice** — one cohesive, independently understandable part of a Task. An
  implementation Worktree has exactly one active Slice.
- **Worktree** — an isolated workspace bound to a Slice. A Git worktree is one
  adapter, not the definition.
- **Task Store** — an optional operational adapter for Task Cards. It may be
  local or provided by GitHub, GitLab, Jira, Linear, or another system; its
  absence never blocks the Markdown workflow.
- **Durable Knowledge** — current product, design, domain, architecture,
  standard, runbook, or landmark-decision knowledge needed by future work.
- **Evidence** — an addressable observation supporting an Acceptance Criterion.
  Test output, runtime state, screenshots, traces, and manual checks are
  possible evidence kinds. Transition evidence must be observed no earlier
  than the current human approval and carry the same decision digest. A project
  may add a stricter freshness window through its verification policy.
- **Verification** — the argument that evidence proves the approved behavior.
- **Skill Router** — the implicit `sillage` skill that selects one next phase
  from the user's intent, current status, and worktree evidence.
- **Specialist Lens** — a focused skill such as SOLID, DDD, frontend, data,
  audit, migration, debugging, or test hygiene. It informs the current phase
  and never creates a second lifecycle.
- **Task Classification** — the bounded context of the work: `probe`,
  `bounded`, or `cross-cutting`.
- **Primary Lens** — the one engineering lens that owns the task's dominant
  risk. Secondary lenses are allowed only for distinct additional risks.
- **Risk Owner** — the test or review layer that owns proof for one Acceptance
  Criterion. Every criterion has at most one primary owner in a conforming
  task.
- **Engineering Doctrine** — portable questions that apply across technology:
  observed pressure, ownership, seams, dependency direction, behavior,
  failure, security, compatibility, recovery, and proof.
- **Architecture Lens** — architecture and pattern guidance selected from a
  real change pressure; it never requires a layer, object model, or pattern.
- **Testing Lens** — test-seam and risk-ownership guidance; it assigns each
  risk a primary proof layer without optimizing for test count.
- **Interface Lens** — contracts across HTTP, web, RPC, messages, files, or
  commands, including semantics, failure, and compatible evolution.
- **Systems Lens** — concurrency, time, networks, queues, resources, partial
  failure, recovery, and observability.
- **Security Lens** — assets, actors, trust boundaries, abuse cases, controls,
  residual risk, and evidence; it never certifies security by inspection.
- **Platform Lens** — runtime, process, OS, filesystem, packaging, lifecycle,
  and resource constraints for desktop, mobile, CLI, service, or embedded work.
- **Native Check** — the project's own test, typecheck, build, linter,
  screenshot, runtime, or domain command used as verification evidence.
- **Review** — an independent assessment of correctness, scope, safety, and
  maintainability.
- **Handoff** — a self-contained outcome or resume point.
- **Approval** — an explicit human authorization for a governed transition or
  consequential action. Decision approvals are bound to the deterministic
  digest of the intent, acceptance criteria, and slices they authorize.
- **Execution Profile** — a provider-neutral recommendation for the reasoning
  capability and effort needed at a task stage. It may be mapped differently by
  each adapter and never grants authority or names a model in the portable
  task record.
- **Delegation Request** — a provider-neutral request for parent or child-agent
  work at one task stage. It names a role, isolation boundary, and expected
  return packet; it never names a model and never grants human authority.
- **Delegation Packet** — the bounded result returned by a delegated child to
  its parent. It is an input to review and verification, not approval or
  evidence until the normal Sillage gates accept it.
- **Decision Digest** — the lowercase SHA-256 of canonical JSON containing the
  Task intent, engineering context, acceptance criterion identities/statements/
  risks, and Slice identities/titles/acceptance/dependencies. Operational
  status, evidence, execution profiles, waivers, approvals, and blockers do
  not change it.
- **Waiver** — explicit human acceptance of missing required evidence and its
  residual risk. Like evidence, it is invalidated when the decision digest
  changes.

## Lifecycle

`INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF`

`BLOCKED` records an external dependency, missing authority, or unavailable
evidence together with the exact resume condition.
