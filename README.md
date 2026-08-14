# Sillage

Sillage is a human-governed engineering workflow for coding agents. It keeps
work small, state resumable, decisions explicit, and completion claims backed
by evidence without requiring a particular language, forge, task tracker, or
CI provider.

Sillage combines:

- portable Agent Skills for judgment;
- an Agent Plugins 1.0 distribution manifest;
- JSON contracts for tasks, slices, projects, and worktrees;
- a standalone Go CLI for deterministic checks;
- project profiles that bind the protocol to local tools;
- an OKF-inspired durable knowledge model.

## Status

Sillage is pre-1.0. Version 0.1.0-rc.1 establishes the public protocol, cold-start
contract, lifecycle validator, and first skill evaluation fixtures. Interfaces
may evolve as the workflow is exercised on real projects.

## Install the release candidate in Codex

The first public test installs the `researching-with-evidence` skill through
the Codex plugin marketplace. It does not install the optional Go CLI.

```sh
codex plugin marketplace add mbourmaud/sillage-workflow --ref codex/sillage-v0.1
codex plugin add sillage-workflow@sillage
```

Start a new Codex task, then ask:

```text
Use the researching-with-evidence skill. Determine whether Agent Plugins 1.0
replaces Agent Skills or packages them. Use primary sources, distinguish facts
from inference, and return an evidence packet without creating repository docs.
```

The expected result names the skill, states the research question and consumer,
uses current primary sources with dates or versions, distinguishes inference
and unknowns, and proposes no durable document unless one has a clear owner.
Remove the test installation with:

```sh
codex plugin remove sillage-workflow@sillage
codex plugin marketplace remove sillage
```

## Core lifecycle

```text
INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF
```

`BLOCKED` is reachable from every state. Product decisions, scope changes,
destructive actions, external writes, merges, deployments, and evidence waivers
remain human decisions.

## Project entry contract

Sillage recommends four stable entry points:

```text
PRODUCT.md
DESIGN.md
AGENTS.md
CLAUDE.md -> AGENTS.md
docs/domain/index.md
```

Names are configurable. Their responsibilities are not: product, experience,
agent operating rules, and domain language must each have one canonical owner.

## CLI

The CLI has no runtime dependency beyond its compiled binary.

```sh
go run ./cmd/sillage doctor --root /path/to/project
go run ./cmd/sillage digest --task task.json
go run ./cmd/sillage transition --task task.json --to IMPLEMENT
```

Commands validate; they do not grant approval or mutate task records.
Approvals, evidence, and waivers carry a decision digest, so a later scope or
plan change requires fresh human authority and fresh verification.

[`examples/pilot/task.json`](examples/pilot/task.json) is an executable,
evidence-backed task record. Repository tests exercise it through the same CLI
boundary used by adopters.

## Agent Plugin

The root [`plugin.json`](plugin.json) packages Sillage skills according to
Agent Plugins 1.0. The Codex marketplace adapter lives in
[`plugins/sillage-workflow`](plugins/sillage-workflow); contract tests keep its
generated skill payload identical to the canonical source under `skills/`.
Agent Plugins is a distribution adapter, not the Sillage runtime contract: the
Go CLI, task schemas, and project documents remain independently usable. No MCP
server is bundled until a concrete tool-server need exists.

## Skills

Once released, install skills explicitly so projects do not inherit an
overlapping workflow:

```sh
npx skills add mbourmaud/sillage-workflow --skill researching-with-evidence
```

See [skills/README.md](skills/README.md) for the released capability set.

## Repository map

```text
cmd/sillage/          CLI
internal/             deterministic workflow policies
schemas/              portable JSON schemas
skills/               original Agent Skills
evals/                behavioral evaluation prompts
docs/domain/          Sillage's domain language
docs/engineering/     current system documentation
docs/adr/             landmark decisions
examples/             forge- and language-neutral examples
```

## Development

```sh
make check
```

The complete gate checks Go formatting, whitespace, `go vet`, GitHub Actions,
the race-enabled test suite, JSON Schema examples, Agent Skill structure, and
the Sillage project contract.

## License

[MIT](LICENSE)
