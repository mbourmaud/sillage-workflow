# Sillage Full Workflow Implementation Plan

**Goal:** make the complete Sillage lifecycle usable from a fresh agent
conversation through a resumable handoff, with one canonical task record and
deterministic gates.

**Architecture:** The task JSON remains the machine contract and local task
store. The project entry files remain the human context contract. The CLI
provides read-only cold-start/status views plus explicit state mutation; skills
provide judgment and conversation routing, not hidden state or provider writes.

**Tech Stack:** Go standard library, JSON Schema Draft 2020-12, Markdown Agent
Skills, local files, GitHub Actions.

**Spec:** PRODUCT.md, DESIGN.md, docs/domain/index.md, and
docs/adr/0001-human-governed-portable-core.md.

## Global Constraints

- Humans control product decisions, scope changes, destructive actions,
  external writes, evidence waivers, merges, and deployments.
- One implementation worktree owns exactly one active slice.
- Commands validate by default; state mutation requires an explicit --write.
- Verification evidence must be addressable, observed after approval, and bound
  to the decision digest. A project-specific adapter may impose a stricter age
  window.
- Tasks may carry a provider-neutral execution profile: adapters choose the
  available capability and effort mapping, while fallbacks remain visible
  verification risks and never grant authority.
- Tasks may request bounded parent or subagent work by stage. The request names
  role, isolation, and return packet; the host owns spawning and model choice,
  while the parent owns synthesis and all human gates.
- Review is separate from verification; handoff is separate from review.
- The core remains independent of language, forge, tracker, package manager,
  and CI provider.
- Operational records stay in the configured task store; durable knowledge is
  promoted only to an existing canonical owner.

### Task 1: Extend the task contract with decisions, review, and handoff

**Files:**

- Modify: schemas/task.schema.json
- Modify: internal/workflow/state.go
- Test: internal/workflow/state_test.go
- Test: internal/contracts/contracts_test.go

**Interfaces:**

- Decision records an accepted or rejected product/technical question:
  id, question, answer, status, by, and at.
- Review records an independent assessment: status, by, at, summary, findings,
  and decision_digest.
- Handoff records outcome, next_action, at, and decision_digest.
- ValidateTransition gates VERIFY to REVIEW on evidence and REVIEW to HANDOFF
  on an accepted current review and a current handoff.

- [x] Write failing tests for accepted review requirements, blocking review
  findings, handoff digest binding, and malformed optional artifacts.
- [x] Run the focused tests and observe failures caused by missing types/gates.
- [x] Add the minimal schema definitions and Go models.
- [x] Include stable decisions in DecisionDigest; exclude operational status,
  evidence, review findings, and handoff text.
- [x] Run the focused tests, then go test ./....

### Task 2: Add deterministic cold-start and status views

**Files:**

- Modify: internal/project/doctor.go
- Create: internal/project/context.go
- Modify: cmd/sillage/main.go
- Test: internal/project/context_test.go
- Test: cmd/sillage/main_test.go

**Interfaces:**

- project.Context(root string, task workflow.Task) ContextReport returns
  canonical entry-point paths, project readiness, task identity/status,
  active slices, required gate, and the next safe action.
- sillage context --root path [--task task.json] [--json] is read-only.
- sillage status --task task.json [--json] is read-only.

- [x] Write failing tests for a ready project, a missing project contract, an
  active slice, a blocked task, and a task with no next legal transition.
- [x] Run focused tests and confirm the commands do not exist yet.
- [x] Implement reports using existing profile resolution and lifecycle policy.
- [x] Keep human output concise and JSON field names stable.
- [x] Run CLI and project tests.

### Task 3: Add explicit, atomic transition writes

**Files:**

- Modify: cmd/sillage/main.go
- Create: internal/taskstore/local.go
- Test: internal/taskstore/local_test.go
- Test: cmd/sillage/main_test.go

**Interfaces:**

- taskstore.WriteTransition(path string, task workflow.Task, target
  workflow.Status, expected []byte) error writes a validated task atomically
  with mode 0600 and checks optimistic concurrency against the bytes read by
  the caller.
- Existing transition remains read-only by default.
- transition --write --task task.json --to STATUS performs the same validation,
  changes only status (blocker metadata must already be present for blocked
  transitions), and refuses to overwrite a task whose bytes changed since it
  was read.

- [x] Write failing tests for read-only behavior, successful atomic write,
  malformed task rejection, invalid transition rejection, and concurrent
  modification detection.
- [x] Run focused tests and observe failures.
- [x] Implement a same-directory temporary file plus rename, writing mode 0600
  and refusing symlink targets.
- [x] Run race-enabled task-store and CLI tests.

### Task 4: Add the single Sillage orchestration skill

**Files:**

- Create: skills/working-with-sillage/SKILL.md
- Create: skills/working-with-sillage/agents/openai.yaml
- Create: evals/working-with-sillage/evals.json
- Modify: skills/README.md
- Modify: plugins/sillage-workflow/skills/
- Test: internal/skills/contract_test.go

**Interfaces:**

- The skill triggers on feature, bug, refactor, architecture, research, and
  continue/resume/ship requests that can change a repository.
- It reads project context, finds the active task, names the current status,
  and chooses one bounded next action.
- It routes research to researching-with-evidence; it does not duplicate it.
- It outputs a compact state report: Status, Known, Missing, Decision, Next
  action, and Evidence needed.

- [x] Write baseline prompts for normal feature intake, blocked resume, and a
  tiny reversible change.
- [ ] Run the baseline without the skill and retain the result. (Requires a
  separate fresh-agent evaluation run; the deterministic pilot does not claim
  to replace this comparison.)
- [x] Write the skill with explicit human gates and no automatic repository
  documentation.
- [x] Copy the canonical skill into the distribution bundle and add a
  byte-for-byte contract test.
- [ ] Run the same prompts with the skill and record qualitative differences.
  (The evaluation fixture remains a release gate; the executable pilot proves
  the deterministic core and installation path.)

### Task 5: Add a complete executable pilot

**Files:**

- Create: examples/full-workflow/task.json
- Create: examples/full-workflow/README.md
- Create: docs/engineering/templates.md
- Modify: README.md
- Modify: docs/engineering/roadmap.md
- Test: cmd/sillage/main_test.go

- [x] Create a task that demonstrates DECIDE approval, IMPLEMENT, VERIFY
  evidence, REVIEW acceptance, and HANDOFF.
- [x] Exercise context, status, read-only transition checks, and explicit write
  transitions in a temporary copy.
- [x] Record commands and observed output in the example README.
- [x] Run make check. The optional `skills` CLI smoke test remains blocked by
  the local Homebrew Node/simdjson runtime and is not presented as a code pass.

### Task 6: Review and ship the full workflow slice

- [x] Review the diff for scope, portability, and durable-document pollution.
- [x] Run make check and the repository's plugin/skill contract checks. Claude
  strict validation and skills discovery from a clean checkout remain release
  follow-ups.
- [x] Run the full pilot in a fresh temporary task copy.
- [ ] Keep the release behind a human merge/tag decision.

### Task 7: Add bounded delegation requests

**Files:**

- Modify: schemas/task.schema.json
- Modify: internal/workflow/state.go
- Modify: internal/project/context.go
- Modify: cmd/sillage/main.go
- Modify: skills/working-with-sillage/SKILL.md
- Create: docs/engineering/adapters/codex.md
- Test: internal/workflow/state_test.go
- Test: internal/project/context_test.go

**Interfaces:**

- A task may carry a default delegation request and stage overrides.
- Requests distinguish parent work from a child, read-only review/investigation
  from an isolated implementation worktree, and the expected return packet.
- Required child work is resumably blocked when the host cannot provide it;
  optional child work surfaces a parent fallback.
- The deterministic core validates and reports requests but never launches a
  provider process.

- [x] Write failing tests for stage resolution, unsafe isolation, required
  parent requests, invalid roles, and decision-digest independence.
- [x] Add the schema, Go model, validation, and cold-start/status projection.
- [x] Document the Codex adapter without coupling the portable plugin to Codex.
- [ ] Run a real host subagent pilot and record model/runtime usage separately
  from the task contract.
