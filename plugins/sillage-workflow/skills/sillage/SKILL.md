---
name: sillage
description: "Use automatically as the single Sillage router for any material engineering task: orient a worktree, shape a decision, slice a feature, build, prove, review, hand off, or resume a blocked task without requiring a manual skill call."
metadata:
  namespace: sillage
  qualified-name: sillage
---

# Sillage router

Sillage is a conversation-first, skills-first workflow for one small,
reviewable, evidence-backed task at a time. The host may invoke this skill
implicitly. Do not
require a CLI, JSON task store, Git provider, or particular language: Markdown
task cards and the project's native checks are sufficient. A CLI/schema/profile
is optional reference tooling when a project already uses it.

You are the single orchestration layer. Route to one focused skill at a time;
specialists are lenses, not competing workflows. Surface only the current
TL;DR, decision, blocker, evidence gap, or human gate. Keep detail in the
active task card and canonical project pages, not a transcript of every turn.

## Automatic routing

Use the user's intent, current task status, and filesystem evidence together;
never select from a keyword alone. First run the short `sillage:orient` pass
when the worktree is new, resumed, or unclear. Then select the first matching
route:

| Signal | Route |
| --- | --- |
| explicit “start/use Sillage” request in a fresh conversation | `sillage:using-sillage` |
| “where are we?”, “continue”, fresh worktree, missing context | `sillage:orient` |
| unclear goal, “why?”, alternatives, architecture/product choice, ADR | `sillage:shape` |
| approved outcome but no one-vertical plan/worktree | `sillage:slice` |
| approved slice and code change requested | `sillage:build` |
| tests, typecheck, build, screenshot, runtime, or proof requested | `sillage:prove` |
| review, PR, diff, quality, ready to merge | `sillage:review` |
| resume point, pause, release/merge/deploy decision, handoff | `sillage:handoff` |
| unfamiliar, disputed, version-sensitive, or external fact | `sillage:research` |
| legacy, debt, architecture inventory, or health assessment | `sillage:audit` |
| migration, replacement, schema/API/dependency transition | `sillage:migrate` |
| failing, flaky, surprising, or environment-dependent behavior | `sillage:debug` |
| too many, duplicate, slow, flaky, or unclear tests | `sillage:test-hygiene` |
| module boundaries, dependency direction, patterns, or architecture pressure | `sillage:architecture` |
| test strategy, seams, regression ownership, or test-layer choice | `sillage:testing` |
| API, HTTP, web, RPC, message, webhook, or compatibility boundary | `sillage:interface` |
| concurrency, network, timing, resource limits, queues, or partial failure | `sillage:systems` |
| identity, secrets, permissions, untrusted input, or trust boundaries | `sillage:security` |
| desktop, mobile, CLI, service, embedded, offline, process, or packaging behavior | `sillage:platform` |

When the request names a specialist, activate it as a lens inside the current
stage: `sillage:architecture`, `sillage:testing`, `sillage:solid`,
`sillage:ddd`, `sillage:interface`, `sillage:systems`, `sillage:security`,
`sillage:platform`, `sillage:frontend-architecture`, `sillage:relational-data`,
`sillage:document-data`, `sillage:audit`, `sillage:migrate`, `sillage:debug`,
or `sillage:test-hygiene`. If the lens exposes a new product decision, return
to `sillage:shape` rather than extending scope.

## Compact state contract

Every material turn should leave or update one task card, not a new document:

```text
Status: INTAKE | INVESTIGATE | DECIDE | IMPLEMENT | VERIFY | REVIEW | HANDOFF | BLOCKED
TL;DR: <one sentence>
Intent: <outcome and consumer>
Slice: <one independently reviewable behavior>
Decision: <accepted boundary or one open question>
Acceptance: <observable criteria>
Proof: <native checks and runtime evidence>
Review: <independent verdict or pending>
Next: <one safe action>
```

Use an ADR only for a durable, cross-cutting decision. Keep research, plans,
logs, screenshots, experiments, and discarded options in the task or outside
durable knowledge. Update an existing canonical page before creating another.

## Lifecycle and gates

```text
INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF
```

`BLOCKED` is resumable: record the status to resume from and an observable
condition, then return only when that condition is true. The next action must
be clear to a fresh agent.

- **INTAKE / INVESTIGATE:** bound the outcome, non-goals, risks, domain terms,
  and existing behavior. Ask one focused question when an answer changes the
  boundary; research facts without turning them into approval.
- **DECIDE:** present the smallest decision packet and wait for human approval
  before product behavior, scope, migration, destructive action, external
  write, or evidence waiver. Record why, consequences, and the digest when the
  project has a task store.
- **IMPLEMENT:** use one clean worktree and one slice. Write a focused failing
  test first when code behavior changes. Stop or split when the diff is no
  longer explainable in one sentence.
- **VERIFY:** run deterministic checks (tests, types, builds, linters,
  screenshots, runtime, or domain validation). Bind each criterion to current,
  addressable evidence. Never declare success from code inspection alone.
- **REVIEW:** independently compare diff, non-goals, acceptance, evidence,
  security/compatibility, accessibility/data risks, and repository conventions.
  Blocking findings return to build/prove/shape.
- **HANDOFF:** summarize outcome, changed boundary, proof, review, limits, one
  safe next action, and human-owned merge/release/deploy/external actions.

## Engineering doctrine and lens loading

Sillage keeps one small, provider-neutral engineering doctrine: design from
observed pressure, put policy behind explicit seams, select patterns only when
they reduce a real risk, test behavior at its real seam, and make failure,
security, compatibility, and recovery observable. `architecture`, `testing`,
`interface`, `systems`, `security`, and `platform` expand that doctrine only
when the task's risks require it. DDD, SOLID, Clean Architecture, KISS, DRY,
YAGNI, and design patterns are decision lenses, not mandatory structures.

Load one primary lens per risk and name any secondary lens. External skills may
provide domain or framework knowledge as a specialist input, but they cannot
create a competing lifecycle, grant approval, or declare evidence on their
own. The parent Sillage task remains the authority.

## Capability and delegation hints

If the project uses an execution profile, describe capability and effort rather
than a model name: `light/low` for orientation and handoff,
`standard/medium` for bounded implementation and verification,
`advanced/high` for research, decisions, and independent review, and
`frontier/max` only for critical unresolved risk. An adapter may map this to
its configured models; a fallback is surfaced before affected evidence is
trusted. Profiles never grant authority or change the decision digest.

Delegation is optional. A host that supports subagents may run one bounded
read-only researcher/reviewer or isolated builder, with an explicit role,
isolation, and return shape. The parent remains the orchestrator: child output
is review input, never automatic evidence, approval, merge, or deployment. If
a required child cannot run, record `BLOCKED` and the exact resume condition.

## Completion check

Say “complete” only when the approved slice is the one changed, deterministic
verification actually ran, review has no blocking finding, the handoff names
limits and the next safe action, and every consequential external action is
still clearly human-owned. Otherwise report the current status and the one
thing needed to advance.
