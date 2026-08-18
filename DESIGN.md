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

### Engineering doctrine before technology

Sillage keeps a small invariant doctrine: design from observed pressure, make
ownership and seams explicit, select patterns only when they reduce a real
risk, test behavior at its real boundary, and make failure, security,
compatibility, and recovery observable. Architecture, testing, DDD, interface,
systems, security, platform, frontend, and data skills are progressive lenses;
they load only when the task's risks require them. They do not mandate a
language, framework, deployment topology, class model, or Clean Architecture.

## Optional reference tooling

The portable product is the skills and Markdown protocol. A project may add a
CLI, JSON schema, task store, or host adapter for deterministic checks. Such
tooling must remain optional: missing capabilities degrade explicitly and never
silently simulate a check, approval, or external write. Human-readable output
is concise; machine-readable output is stable when a tool provides it.
