# Sillage

Sillage is a human-governed engineering workflow for coding agents. It keeps
work small, state resumable, decisions explicit, and completion claims backed
by evidence without requiring a particular language, forge, task tracker, or
CI provider.

**[Visit the Sillage website](https://mbourmaud.github.io/sillage-workflow/)**
for the visual protocol guide, installation paths, and release notes.

Sillage combines:

- portable Agent Skills for judgment;
- an Agent Plugins 1.0 distribution manifest;
- JSON contracts for tasks, slices, projects, and worktrees;
- a standalone Go CLI for deterministic checks;
- project profiles that bind the protocol to local tools;
- provider-neutral execution profiles that match capability and effort to each
  stage;
- an OKF-inspired durable knowledge model.

## Status

Sillage is pre-1.0. Version 0.1.0 establishes the public protocol, cold-start
contract, lifecycle validator, and first skill evaluation fixtures. Interfaces
may evolve as the workflow is exercised on real projects.

## Install Sillage 0.1.0

The release ships one canonical Agent Skill through several thin
distribution adapters. None of these commands installs the optional Go CLI.

### Codex

```sh
codex plugin marketplace add mbourmaud/sillage-workflow
codex plugin add sillage-workflow@sillage
```

### Claude Code

```sh
claude plugin marketplace add mbourmaud/sillage-workflow
claude plugin install sillage-workflow@sillage
```

### Any Agent Skills client

The `skills` CLI supports Claude Code, Codex, Cursor, OpenCode, GitHub Copilot,
and many other clients. Install from GitHub and name one or more agents explicitly:

```sh
npx skills add mbourmaud/sillage-workflow \
  --skill researching-with-evidence \
  --agent claude-code --agent codex
```

### Full Sillage workflow

Install the orchestration skill when you want the complete lifecycle rather
than the research capability alone:

```sh
npx skills add mbourmaud/sillage-workflow \
  --skill working-with-sillage \
  --global --agent codex --yes
```

The skill is usable for pilots and is still marked as a workflow candidate
until its first independent evaluation is completed.

### Updating

Skills do not update in the background. Update them deliberately at a task
boundary, then run the project's own checks:

```sh
npx skills update working-with-sillage --global --yes
make check
```

This updates agent guidance only. It does not update Sillage's repository
schemas, CLI, or project documents; those arrive through the normal repository
release process.

### First behavior test

Start a fresh agent task, then ask:

```text
Use the researching-with-evidence skill. Determine whether Agent Plugins 1.0
replaces Agent Skills or packages them. Use primary sources, distinguish facts
from inference, and return an evidence packet without creating repository docs.
```

The expected result names the skill, states the research question and consumer,
uses current primary sources with dates or versions, distinguishes inference
and unknowns, and proposes no durable document unless one has a clear owner.
Remove a marketplace installation with the matching client:

```sh
codex plugin remove sillage-workflow@sillage
codex plugin marketplace remove sillage

claude plugin uninstall sillage-workflow@sillage
claude plugin marketplace remove sillage
```

## Core lifecycle

```text
INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF
```

`BLOCKED` is reachable from every state. Product decisions, scope changes,
destructive actions, external writes, merges, deployments, and evidence waivers
remain human decisions.

Tasks may also carry an execution profile such as `standard/medium` for a
bounded implementation or `advanced/high` for a difficult decision or
independent review. These profiles do not name a model, grant permissions, or
change the decision digest; an adapter maps them to the models available in its
environment and must surface any fallback before affected evidence is trusted.

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
go run ./cmd/sillage context --root /path/to/project --task task.json --json
go run ./cmd/sillage status --task task.json --json
go run ./cmd/sillage changelog check --version v0.1.0
go run ./cmd/sillage changelog extract --version v0.1.0
go run ./cmd/sillage digest --task task.json
go run ./cmd/sillage transition --task task.json --to IMPLEMENT
go run ./cmd/sillage transition --task task.json --to IMPLEMENT --write
```

Commands validate and remain read-only by default. Only `transition --write`
mutates a local task record, after validation and an optimistic-concurrency
check; it does not grant approval.
Approvals, evidence, and waivers carry a decision digest, so a later scope or
plan change requires new human authority and new approval-bound verification.
Projects that need temporal freshness can add that policy in their verification
adapter; the portable core does not invent a universal age window.

[`examples/pilot/task.json`](examples/pilot/task.json) is an executable,
evidence-backed task record. Repository tests exercise it through the same CLI
boundary used by adopters.

## Agent Plugin

The root [`plugin.json`](plugin.json) packages Sillage skills according to
Agent Plugins 1.0. Codex and Claude Code marketplace manifests point to the same
[`plugins/sillage-workflow`](plugins/sillage-workflow) bundle; contract tests
keep its payload identical to the canonical source under `skills/` and its
version aligned across ecosystems. Agent Plugins is a distribution adapter, not
the Sillage runtime contract: the Go CLI, task schemas, and project documents
remain independently usable. No MCP server is bundled until a concrete
tool-server need exists.

## Skills

Install skills explicitly so projects do not inherit an
overlapping workflow:

```sh
npx skills add mbourmaud/sillage-workflow --skill researching-with-evidence
```

See [skills/README.md](skills/README.md) for the released capability set.

The `working-with-sillage` skill is the full workflow orchestrator: it guides
cold start, one bounded slice, human decision gates, deterministic verification,
independent review, blocked-task resumption, and handoff. It is accompanied by
the executable pilot in [examples/full-workflow](examples/full-workflow).

## Repository map

```text
cmd/sillage/          CLI
internal/             deterministic workflow policies
schemas/              portable JSON schemas
skills/               original Agent Skills
plugins/              generated multi-client distribution bundle
evals/                behavioral evaluation prompts
docs/domain/          Sillage's domain language
docs/engineering/     current system documentation
docs/adr/             landmark decisions
examples/             forge- and language-neutral examples
```

## Development

```sh
make check
make pilot
```

The complete gate checks Go formatting, whitespace, `go vet`, GitHub Actions,
the race-enabled test suite, JSON Schema examples, Agent Skill structure, and
the Sillage project contract.

`make pilot` runs the full local workflow example against a temporary task copy:
it reads the project contract, reports context and status, performs the explicit
`REVIEW → HANDOFF` write, and confirms the resulting handoff state. It never
changes the canonical example or performs an external write.

## Releases

Keep user-visible changes in [`CHANGELOG.md`](CHANGELOG.md) under
`[Unreleased]`. Before a human-approved tag, run:

```sh
make check
make release-notes VERSION=vX.Y.Z
```

CI requires the matching version section and publishes those exact notes with
the release.

## License

[MIT](LICENSE)
