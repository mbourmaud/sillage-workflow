---
type: workflow-templates
status: stable
owners: [maintainers]
---

# Workflow templates

These are compact Markdown prompts for the operational task card. A repository
may mirror them into JSON or a task store, but no CLI or machine format is
required. Keep one active card per vertical and promote only durable knowledge
to its canonical page.

## Task card / spec

```markdown
## Task <id> — <title>

Status: INTAKE

### TL;DR
<one sentence a reviewer can repeat>

### Intent
- Outcome: <observable result>
- Consumer: <who or what uses it>
- Class: probe | bounded | cross-cutting
- Scope: <one bounded slice>
- Non-goals: <explicit exclusions>

### Engineering lenses
- Primary risk lens: <architecture | testing | ddd | interface | systems |
  security | platform | frontend | relational-data | document-data | none>
- Secondary lenses: <only when a distinct risk needs them>

### Acceptance
- [ ] <observable criterion> (risk: <what could be wrong>)

### Decision
- Status: awaiting human | accepted | rejected
- Answer: <chosen behavior and why>
- By / at: <identity and timestamp>

### Slice / plan
- ID: <slice-id>
- Worktree: <path or “not applicable”>
- Dependencies: <none or named prerequisite>
- Stop condition: <when to split or return to shaping>

### Verification
- [ ] <criterion> — check: <native command or runtime observation>

### Review
Pending.

### Handoff
Pending.
```

The TL;DR, intent, decision, slice, proof, review, and handoff are the
reviewer's progressive-disclosure path. Keep detail in the relevant section;
do not create a separate file for every stage.

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

## Execution profile

The task may carry a provider-neutral reasoning recommendation. It describes
the capability and effort needed for the work, not a model name or permission
to bypass a human gate.

```json
{
  "execution": {
    "default": { "capability": "standard", "effort": "medium" },
    "overrides": {
      "DECIDE": { "capability": "advanced", "effort": "high" },
      "REVIEW": { "capability": "advanced", "effort": "high" }
    }
  }
}
```

Use `light/low` for orientation and handoff, `standard/medium` for bounded
implementation, `advanced/high` for research, decisions, and independent
review, and `frontier/max` only for critical or unresolved high-impact work.
Adapters choose the available provider mapping. A fallback is recorded and
surfaced before relying on affected evidence; actual usage is recorded only
when the adapter exposes it. Execution profiles are operational and do not
change the decision digest.

## Delegation plan

The task may also request a child agent for one stage. This is a host adapter
request, not a second lifecycle or a model selection:

```json
{
  "delegation": {
    "default": {
      "mode": "parent",
      "role": "orchestrator",
      "isolation": "same_context",
      "return": "summary"
    },
    "overrides": {
      "DECIDE": {
        "mode": "subagent",
        "role": "decision_researcher",
        "isolation": "read_only",
        "return": "decision_packet",
        "required": true
      },
      "IMPLEMENT": {
        "mode": "subagent",
        "role": "builder",
        "isolation": "isolated_worktree",
        "return": "implementation_patch"
      },
      "REVIEW": {
        "mode": "subagent",
        "role": "reviewer",
        "isolation": "read_only",
        "return": "review_findings",
        "required": true
      }
    }
  }
}
```

The parent supplies only the project context, task intent, active slice,
decision digest, requested isolation, and return shape. The host chooses the
actual model. Optional requests fall back to the parent with an explicit note;
required requests block until the host can provide a child. Child output is
review input, never automatic approval or evidence.

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
completion. Native project commands are the default proof surface; a Sillage
CLI is optional.

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
