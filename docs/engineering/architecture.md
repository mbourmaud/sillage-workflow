---
type: architecture
status: stable
owners: [maintainers]
sources: [docs/adr/0001-human-governed-portable-core.md]
---

# Architecture

Sillage has four layers:

1. **Protocol** — portable lifecycle, authority, task, slice, evidence,
   execution-profile, and knowledge contracts.
2. **Skills** — judgment about shaping, research, slicing, implementation,
   verification, review, and handoff.
3. **Deterministic core** — validation of facts and transitions.
4. **Adapters** — task stores, workspace isolation, version control,
   verification commands, agent capability/effort mappings, and publication
   mechanisms.

The protocol works without the CLI. The CLI works without a remote forge. An
adapter may add capability but cannot weaken the authority model silently.

The Go runtime core uses the standard library and produces standalone binaries.
Test-only libraries validate JSON Schema and development tooling. JSON is the
machine contract; Markdown remains the human-readable knowledge format.

## Execution profiles

Tasks may request a provider-neutral capability and effort profile for the
current lifecycle stage. The portable contract uses `light`, `standard`,
`advanced`, or `frontier` capability and `low`, `medium`, `high`, or `max`
effort. Adapters map those requirements to whatever models and controls are
available locally. The core never stores vendor model names, grants authority,
or claims usage data that an adapter did not expose. The requested profile is
operational and is therefore excluded from the product decision digest; a
material fallback must still be surfaced as a verification risk.

## Workflow authority

The Sillage protocol, its task/project schemas, repository skills, deterministic
CLI, and release checks are the only supported workflow surface. External
methods are neither runtime dependencies nor alternate lifecycle authorities.
Ideas may be adapted, but the resulting contract must be expressed in
Sillage-native statuses, artifacts, gates, and commands.
