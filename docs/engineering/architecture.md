---
type: architecture
status: stable
owners: [maintainers]
sources: [docs/adr/0001-human-governed-portable-core.md]
---

# Architecture

Sillage has four layers:

1. **Protocol** — portable lifecycle, authority, task, slice, evidence,
   execution-profile, delegation, and knowledge contracts.
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

## Delegation plans

Tasks may request a separate agent thread for a bounded lifecycle stage. A
delegation request names only a mode (`parent` or `subagent`), a portable role,
an isolation boundary, and an expected return packet. `read_only` and
`isolated_worktree` make the two safe child shapes explicit; a required request
blocks rather than silently falling back when the host cannot spawn a child.

The parent agent remains the orchestrator and human-facing authority. The
deterministic core validates the request and exposes it through status/context;
it does not spawn processes. A host adapter maps the request and the execution
profile to its own agent threads, models, permissions, and worktree mechanism.
The child receives a bounded context and returns a packet to the parent. That
packet becomes evidence, approval, or a lifecycle transition only after the
parent applies the normal Sillage gates.

## Workflow authority

The Sillage protocol, its task/project schemas, repository skills, deterministic
CLI, and release checks are the only supported workflow surface. External
methods are neither runtime dependencies nor alternate lifecycle authorities.
Ideas may be adapted, but the resulting contract must be expressed in
Sillage-native statuses, artifacts, gates, and commands.
