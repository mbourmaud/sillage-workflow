---
type: decision
status: stable
owners: [maintainers]
---

# 0001 — Human-governed portable core

## Context

Coding-agent workflows often bind their process to one forge, language,
runtime, issue tracker, or autonomous orchestration model. That prevents teams
from adopting the useful discipline independently of the surrounding tooling.

## Decision

Sillage defines a human-governed protocol independent of language, forge, task
store, CI provider, and agent runtime. Skills own judgment. Deterministic tools
validate facts. Project profiles and adapters bind portable concepts to local
systems.

Human approval remains required for product decisions, scope changes,
destructive actions, external writes, evidence waivers, merges, and deployments.

## Consequences

- Local-file operation is a first-class task-store mode.
- Git and remote forges are optional capabilities.
- Skills cannot require multi-agent support.
- The CLI can validate approval but cannot manufacture it.
- Provider-specific vocabulary stays outside the portable domain model.
