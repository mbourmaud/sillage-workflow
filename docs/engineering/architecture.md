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
2. **Skills** — the portable product: implicit routing, shaping, research,
   slicing, implementation, verification, review, handoff, and progressive
   architecture, testing, DDD, interface, systems, security, platform,
   frontend, and data lenses.
3. **Deterministic core** — validation of facts and transitions.
4. **Adapters** — task stores, workspace isolation, version control,
   verification commands, agent capability/effort mappings, and publication
   mechanisms.

The protocol works with Markdown skills and the project's native commands. The
CLI, schemas, and task store are optional deterministic adapters. An adapter may
add capability but cannot weaken the authority model silently.

## Language-neutral engineering doctrine

The doctrine is intentionally smaller than a software architecture textbook.
It asks what pressure exists, who owns the decision, where the stable seam is,
what pattern or no-pattern reduces the risk, how failure and compatibility are
handled, and what observation would prove the claim. DDD, SOLID, Clean
Architecture, KISS, DRY, YAGNI, design patterns, HTTP, networking, security,
and platform conventions are invoked as lenses, not universal structures.

External technical skills can contribute a source, rule, or specialist review
through an adapter. The Sillage task remains the lifecycle authority: external
skills cannot create statuses, grant approval, or self-declare evidence.

The Go runtime core uses the standard library and produces standalone binaries.
Test-only libraries validate JSON Schema and development tooling. Markdown is
the portable human contract; JSON is an optional machine contract when a
project wants deterministic adapter validation.

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

The Sillage protocol, repository skills, Markdown artifacts, and release checks
are the only supported workflow surface. The optional CLI and schemas validate
that surface but do not replace it. External methods are neither runtime
dependencies nor alternate lifecycle authorities. Ideas may be adapted, but
the resulting contract must be expressed in Sillage-native statuses, artifacts,
gates, and commands.
