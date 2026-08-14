---
type: architecture
status: stable
owners: [maintainers]
sources: [docs/adr/0001-human-governed-portable-core.md]
---

# Architecture

Sillage has four layers:

1. **Protocol** — portable lifecycle, authority, task, slice, evidence, and
   knowledge contracts.
2. **Skills** — judgment about shaping, research, slicing, implementation,
   verification, review, and handoff.
3. **Deterministic core** — validation of facts and transitions.
4. **Adapters** — task stores, workspace isolation, version control,
   verification commands, and publication mechanisms.

The protocol works without the CLI. The CLI works without a remote forge. An
adapter may add capability but cannot weaken the authority model silently.

The Go runtime core uses the standard library and produces standalone binaries.
Test-only libraries validate JSON Schema and development tooling. JSON is the
machine contract; Markdown remains the human-readable knowledge format.

## Workflow authority

The Sillage protocol, its task/project schemas, repository skills, deterministic
CLI, and release checks are the only supported workflow surface. External
methods are neither runtime dependencies nor alternate lifecycle authorities.
Ideas may be adapted, but the resulting contract must be expressed in
Sillage-native statuses, artifacts, gates, and commands.
