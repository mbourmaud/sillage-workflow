---
type: workflow-templates
status: stable
owners: [maintainers]
---

# Workflow templates

These are compact human-readable prompts for the operational task record. Keep
the authoritative machine state in the configured task store and promote only
durable knowledge to its canonical page.

## Task / spec

```markdown
## Task <id> — <title>

Status: INTAKE

### Intent
- Outcome: <observable result>
- Scope: <one bounded slice>
- Non-goals: <explicit exclusions>

### Acceptance
- [ ] <observable criterion> (risk: <what could be wrong>)

### Slice
- ID: <slice-id>
- Dependencies: <none or named dependency>

### Open decision
<one focused question, or “none”>
```

## Decision

```markdown
### Decision <id> — <question>

Status: proposed | accepted | rejected
Answer: <chosen behavior and boundary>
By: human:<identity>
At: <RFC3339 timestamp>

Why: <short reason tied to the task outcome>
Non-goals preserved: <what remains out of scope>
Risk / recovery: <material risk and observable recovery condition>
```

An accepted decision changes the task decision digest. Re-open the decision
gate when its answer, scope, or non-goals change.

## Verification

```markdown
### Verification — <criterion id>

Check: <test, typecheck, build, screenshot, runtime, or domain command>
Observed at: <RFC3339 timestamp>
Result: passed | observed
Reference: <addressable command result, artifact, or URL>
Decision digest: <sha256>
Notes: <limits, skipped checks, or none>
```

Evidence belongs beside the criterion it proves. A skipped check records its
impact and needs an explicit human waiver; inspection alone is not evidence of
completion.

## Review

```markdown
### Review

Status: accepted | changes_requested
Reviewer: reviewer:<identity>
At: <RFC3339 timestamp>
Decision digest: <sha256>
Summary: <independent assessment against intent and acceptance>

Findings:
- Severity: blocking | non_blocking
  Reference: <file, criterion, or check>
  Detail: <what was found and disposition>
```

An accepted review has no blocking finding and is bound to the current
decision digest. Verification and review are separate activities.

## Handoff

```markdown
### Handoff

Outcome: <what is true now>
Next action: <one safe action for a fresh agent>
At: <RFC3339 timestamp>
Decision digest: <sha256>
Known limits: <unresolved risk, external action, or none>
Human gates: <merge, deploy, external write, or none>
```

The handoff is the final local workflow artifact. It never implies that a
merge, deployment, release, or external communication happened without the
authorized system confirming it.
