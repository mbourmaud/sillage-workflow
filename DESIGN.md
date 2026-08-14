---
type: design
status: stable
owners: [maintainers]
---

# Sillage design

## Experience principles

### Conversation first

Humans speak naturally. Agents maintain structured state beneath the
conversation and surface only decisions, blockers, evidence gaps, and material
scope changes.

### Calm authority

Sillage never presents an agent recommendation as human approval. Commands say
what is known, missing, or blocked without theatrical urgency.

### Progressive disclosure

Cold start loads only the project entry points, active task, current slice, and
relevant domain pages. Historical detail remains reachable without dominating
the working context.

### Evidence near claims

Acceptance evidence appears beside the criterion it supports. Skipped checks
state their impact and waiver explicitly.

### Portable language

User-facing concepts are task, slice, evidence, review, handoff, and knowledge.
Provider terms such as pull request, merge request, GitHub issue, or Jira ticket
belong to adapters.

## CLI behavior

- Human-readable output is concise and actionable.
- `--json` provides stable machine-readable output.
- Validation commands are read-only unless their name explicitly denotes a
  mutation.
- Missing capabilities degrade explicitly; they are never silently simulated.
- Destructive commands are absent from the core CLI.
