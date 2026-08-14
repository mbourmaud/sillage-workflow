# Sillage

Sillage is a human-governed engineering workflow for coding agents. It keeps
work small, state resumable, decisions explicit, and completion claims backed
by evidence without requiring a particular language, forge, task tracker, or
CI provider.

Sillage combines:

- portable Agent Skills for judgment;
- JSON contracts for tasks, slices, projects, and worktrees;
- a standalone Go CLI for deterministic checks;
- project profiles that bind the protocol to local tools;
- an OKF-inspired durable knowledge model.

## Status

Sillage is pre-1.0. Version 0.1 establishes the public protocol, cold-start
contract, lifecycle validator, and first skill evaluation fixtures. Interfaces
may evolve as the workflow is exercised on real projects.

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
go run ./cmd/sillage transition --task task.json --to IMPLEMENT
```

Commands validate; they do not grant approval or mutate task records.
Approvals, evidence, and waivers carry a decision digest, so a later scope or
plan change requires fresh human authority and fresh verification.

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
go test ./...
go vet ./...
```

## License

[MIT](LICENSE)
