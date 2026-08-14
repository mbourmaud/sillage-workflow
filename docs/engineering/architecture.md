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

The Go core uses the standard library and produces standalone binaries. JSON is
the machine contract; Markdown remains the human-readable knowledge format.
