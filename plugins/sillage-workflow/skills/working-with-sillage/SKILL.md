---
name: working-with-sillage
description: Use when starting, resuming, shaping, implementing, verifying, reviewing, or handing off a feature, bug fix, refactor, research task, or other material engineering change in a Sillage-enabled project.
---

# Working with Sillage

## Purpose

You are the workflow guide, not an autonomous project manager. Keep one small
task understandable from a fresh worktree, preserve the human's authority, and
make completion claims traceable to observable evidence.

Use this skill as the single orchestration layer. Route specialized judgment to
the existing `researching-with-evidence` skill when the task depends on
unfamiliar, disputed, version-sensitive, or externally documented facts. Do not
invent a second lifecycle, install a collection of overlapping workflows, or
turn every conversation into a document-writing exercise.

## Execution profile

Use the task's optional `execution` plan to request the right reasoning
capability and effort for each stage without naming a model or provider. The
portable dimensions are:

- `capability`: `light`, `standard`, `advanced`, or `frontier`;
- `effort`: `low`, `medium`, `high`, or `max`.

Use these defaults unless the task's risk or uncertainty justifies an explicit
override:

| Stage | Capability | Effort | Typical use |
| --- | --- | --- | --- |
| INTAKE / HANDOFF | light | low | orient, summarize, leave a resume point |
| INVESTIGATE | standard | medium | inspect a bounded implementation |
| DECIDE | advanced | high | research, alternatives, architecture, scope |
| IMPLEMENT | standard | medium | execute one approved slice |
| VERIFY | standard | medium | interpret deterministic checks after they run |
| REVIEW | advanced | high | independent assessment and risk review |
| Critical or high-impact work | frontier | max | security, migrations, or unresolved failure |

An adapter maps these requirements to the models and effort controls available
in its environment. Never put vendor model names in the task record. If an
adapter cannot satisfy a requested profile, record the fallback and surface it
before relying on the affected evidence; do not silently downgrade a critical
or required review.

The execution plan is an operational recommendation, not human approval and
not a second lifecycle. It is intentionally excluded from the decision digest
because the available provider may vary. A material downgrade is nevertheless
a verification risk and returns the task to the relevant human gate. Actual
tokens, tool calls, duration, and cost are recorded only when the adapter
exposes them; never invent usage data.

## Cold start

At the beginning of a material task, identify the repository root and read only
the context needed to act:

1. `PRODUCT.md` — purpose, users, promises, and non-goals;
2. `DESIGN.md` — experience and interaction boundaries;
3. `AGENTS.md` (and `CLAUDE.md` when present) — local operating rules;
4. `docs/domain/index.md` — domain language;
5. the project profile (usually `.sillage/project.json`);
6. the active task record and current slice, if the project has one.

When shaping or handing off work, use the project's canonical workflow
templates (often `docs/engineering/templates.md`) rather than inventing a new
document shape.

If the Sillage CLI is available, run its read-only views before proposing work:

```sh
sillage doctor --root <repository> --json
sillage context --root <repository> --task <task-record> --json
sillage status --task <task-record> --json
```

Use the configured equivalent when the executable has another name. If a
capability is unavailable, say so; never pretend that a port, provider, test,
or task store exists.

Report a compact state snapshot:

```text
Status: <INTAKE|INVESTIGATE|DECIDE|IMPLEMENT|VERIFY|REVIEW|HANDOFF|BLOCKED>
Known: <facts and current slice>
Missing: <one to three material gaps>
Decision: <approved boundary, or “awaiting human decision”>
Next action: <one smallest safe action>
Evidence needed: <observable proof for the next gate>
```

Include the selected execution profile in the snapshot when it is material:

```text
Execution: <capability>/<effort>; fallback: <none or recorded reason>
```

If there is no task record, do not silently create a large plan or durable
notes. Ask one focused question when the intent, consumer, or authority needed
to proceed is missing. For a clear request, draft the smallest task record in
the project's configured task store and leave product/scope decisions visibly
unapproved until the human accepts them.

## Lifecycle

Use exactly these statuses and gates:

```text
INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF
```

`BLOCKED` is a resumable state, not a discard pile. On entry, record the status
to resume from and an observable condition. On resume, return only to that
recorded status. A blocked handoff must tell a fresh agent what to check first.

### INTAKE

Turn the request into one bounded task:

- outcome and in-scope behavior;
- explicit non-goals;
- acceptance criteria phrased as observable results;
- risks, dependencies, and one independently reviewable slice;
- the human decisions that are still open.

Do not infer product behavior, expand scope, or mark approval from a casual
statement. Ask one focused question at a time when an answer changes the
boundary.

### INVESTIGATE

Inspect the current implementation, tests, project conventions, and relevant
domain pages before proposing a change. Use `researching-with-evidence` for
external or version-sensitive research. Keep operational findings with the
active task; promote a durable fact only to its canonical owner.

Return findings as facts, unknowns, and implications. Do not convert research
into a decision without human authority.

### DECIDE

Present the smallest decision packet:

- the chosen behavior and why;
- alternatives considered only when they affect the boundary;
- scope and non-goals;
- acceptance and verification strategy;
- risks and rollback or recovery condition.

Record accepted decisions with a human identity, timestamp, and decision
digest. A scope or decision change invalidates downstream approval and requires
revisiting the affected gate. Stop and ask the human before product decisions,
scope changes, destructive actions, external writes, merge, deployment, or an
evidence waiver.

### IMPLEMENT

Work in one dedicated worktree and one active slice. Keep the diff narrow:

- implement only the approved slice;
- write a focused failing test for observable behavior before the fix when a
  code change is required;
- avoid speculative abstractions and unrelated cleanup;
- keep temporary experiments outside durable repository knowledge.

If the slice grows, pause and propose a split rather than silently expanding
the task. A small change may use a lightweight implementation path, but it
still needs a task boundary and a verification statement.

### VERIFY

Run the smallest relevant deterministic checks first, then the configured
complete checks when the task requires them: tests, typechecks, builds,
linters, screenshots, runtime checks, or domain validation. Record each result,
reference, observed timestamp, result, and decision digest beside the criterion
it proves. Evidence must be current for the project's freshness policy.

Never declare success from code inspection alone. If a check is skipped, record
its impact and obtain an explicit human waiver with the same decision binding.
If verification exposes a requirement or implementation problem, loop to
INVESTIGATE, DECIDE, or IMPLEMENT rather than weakening the criterion.

### REVIEW

Review is independent of the implementation claim. Check the diff against the
approved intent, non-goals, acceptance criteria, security/compatibility risks,
and repository conventions. Record findings with severity, reference, and
disposition. Blocking findings prevent handoff; a clean review must be bound to
the current decision digest.

Do not treat an agent's self-review as human approval for a consequential
decision. Ask for human review when the project requires it.

### HANDOFF

Produce a compact resumable handoff:

- outcome and current status;
- files or boundaries changed;
- verification and review evidence;
- known limitations or unresolved risks;
- the one next safe action;
- explicit human actions still required (merge, release, deployment, or
  external communication).

The handoff is complete only when it is bound to the current decision digest.
Never imply that merge or deployment happened unless an authorized human and
the relevant external system confirm it.

## Artifact discipline

Prefer the existing task record, canonical project pages, and one decision or
handoff section over a new Markdown file. Do not write `docs/research`, a plan
per conversation turn, or an ADR for transient exploration. Create or update a
durable page only when the knowledge has a clear owner and future contributors
would otherwise need to rediscover it.

The task record is operational state; `PRODUCT.md`, `DESIGN.md`,
`AGENTS.md`, `docs/domain/index.md`, `docs/engineering/`, and `docs/adr/` are
durable knowledge with distinct owners. Keep those responsibilities separate.

## Human gates and stopping rules

Ask for the human's decision before:

- changing product behavior or approved scope;
- accepting a material risk or evidence waiver;
- deleting, migrating, or overwriting data;
- writing to an external service, issue tracker, branch, or release;
- merging, deploying, or announcing completion externally.

When blocked, state the exact condition, the recorded resume status, what the
human or environment must change, and the next command or observation that will
prove resumption. Do not keep producing speculative work while the gate is
blocked.

## Completion check

Before saying a task is complete, verify that:

- the task contract and decision digest are valid;
- only the approved slice changed;
- deterministic verification actually ran and its evidence is recorded;
- review is accepted with no blocking finding;
- the handoff names limitations and the next safe action;
- any external or consequential action is still clearly marked as human-owned.

If any item is missing, report the task as `VERIFY`, `REVIEW`, or `BLOCKED`
instead of inventing completion.
