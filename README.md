# Sillage

Sillage is a skills-first engineering workflow for coding agents. It keeps
work small, decisions human-owned, context resumable, and completion claims
backed by observable evidence — without turning a repository into a diary.

**[Read the Sillage guide](https://mbourmaud.github.io/sillage-workflow/)** for
the visual protocol, installation paths, and changelog.

## The idea

The agent should be able to enter a fresh worktree and answer: what is this
project, what language does its domain use, what task is active, what decision
was accepted, and what is the next safe action? The reviewer should be able to
summarize a vertical without reading hundreds of lines of plans or tests.

Sillage provides one implicit router and small composable skills:

```text
orient → shape → slice → build → prove → review → handoff
             ↘ research and specialist lenses ↗
```

The lifecycle remains:

```text
INTAKE → INVESTIGATE → DECIDE → IMPLEMENT → VERIFY → REVIEW → HANDOFF
```

`BLOCKED` is a resumable state with an explicit condition and resume status.
Product decisions, scope changes, destructive actions, evidence waivers,
external writes, merges, releases, and deployments remain human-owned.

## Install the skills

The portable product is `SKILL.md` plus compact Markdown task records. The
Sillage CLI, JSON schemas, task store, Git provider, and plugin host are
optional adapters; a project can use its own commands and files.

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

Install the complete suite, or select only the skills a project needs:

```sh
npx skills add mbourmaud/sillage-workflow \
  --skill using-sillage --skill sillage --skill orient --skill shape --skill slice \
  --skill build --skill prove --skill review --skill handoff \
  --skill research --skill architecture --skill testing --skill solid --skill ddd \
  --skill interface --skill systems --skill security --skill platform \
  --skill frontend-architecture --skill relational-data --skill document-data \
  --skill audit --skill migrate --skill debug --skill test-hygiene
```

The host may invoke skills implicitly from their descriptions. You can still
call one explicitly as `$sillage:review` or `$sillage:ddd` when you want a
focused lens. The host decides whether implicit invocation is supported; the
workflow remains useful when you describe the desired stage in ordinary words.

### Start a task

In a fresh conversation, use the simple entry command:

```text
$using-sillage
```

Or say “start Sillage” / “démarre Sillage”. The entry skill reads the minimal
project context, returns a short state card, and hands off to the right phase.
It does not implement or create a documentation pile during the handshake.

Update deliberately at a task boundary and run the project's own checks:

```sh
npx skills update sillage-workflow --global --yes
```

Skill updates change agent guidance only. They do not update project documents,
schemas, or the optional CLI.

## What the skills do

| Skill | Use it for |
| --- | --- |
| `sillage` | The single implicit router and lifecycle authority |
| `using-sillage` | Simple fresh-conversation entry point and state card |
| `orient` | Fresh-worktree context and a concise state card |
| `shape` | Goals, context, alternatives, decisions, and durable ADRs |
| `slice` | One vertical, worktree, acceptance set, and stop condition |
| `build` | One approved implementation with focused tests |
| `prove` | Deterministic checks and evidence packets |
| `review` | Independent, reviewer-sized quality assessment |
| `handoff` | A compact, resumable outcome and next safe action |
| `research` | Version-sensitive facts and source-traceable evidence |
| `architecture` | Patterns, boundaries, dependency direction, and Clean Architecture questions |
| `testing` | Behavioral seams, test ownership, red-green loops, and proof portfolios |
| `solid` | Pragmatic SOLID, Clean Architecture, KISS, DRY, YAGNI, and GRASP boundaries |
| `ddd` | Language, bounded contexts, invariants, and domain models |
| `interface` | HTTP, web, RPC, message contracts, failure semantics, and compatibility |
| `systems` | Networks, concurrency, timing, resources, queues, and partial failure |
| `security` | Assets, trust boundaries, abuse cases, controls, and residual risk |
| `platform` | Desktop, mobile, CLI, services, embedded, process, and packaging constraints |
| `frontend-architecture` | UI states, route/data boundaries, accessibility, browser proof |
| `relational-data` | Constraints, transactions, indexes, and safe migrations |
| `document-data` | Access patterns, schema evolution, consistency, and reconciliation |
| `audit` | Evidence-led codebase, architecture, legacy, and debt assessment |
| `migrate` | Staged, compatible, observable migration paths |
| `debug` | Reproducible diagnosis before a bounded fix |
| `test-hygiene` | Smaller, clearer proof portfolios without losing risk coverage |

Specialists are progressive lenses inside the same lifecycle. Sillage's
engineering doctrine is immutable, but no lens mandates a language, framework,
class model, pattern, test pyramid, or Clean Architecture. They do not create
parallel statuses, duplicate plans, or authorize changes.

## Markdown-first project setup

Keep the repository entry points small and stable:

```text
PRODUCT.md
DESIGN.md
AGENTS.md
CLAUDE.md → AGENTS.md
docs/domain/index.md
docs/engineering/templates.md
```

For each vertical, keep one task card with a TL;DR, intent, decision, slice,
acceptance/proof, review, and handoff. Store transient research, logs,
screenshots, experiments, and discarded plans outside durable knowledge. Add
an ADR only when a decision is durable and cross-cutting.

## Optional reference CLI

This repository includes a standalone Go CLI and JSON contracts for projects
that want deterministic validation, profiles, or local task transitions. They
are useful reference adapters, not a prerequisite for the skills:

```sh
go run ./cmd/sillage doctor --root .
go run ./cmd/sillage status --task task.json
go run ./cmd/sillage conformance --task task.json --json
go run ./cmd/sillage transition --task task.json --to VERIFY
```

`conformance` is the stricter opt-in check for a task that is ready to enter
engineering work: it requires a `probe`, `bounded`, or `cross-cutting`
classification, one primary lens, optional distinct secondary lenses, and a
named owning test or review layer for every acceptance criterion. Existing
0.2 task records remain valid without these optional fields.

Use the project's native test, typecheck, build, linter, screenshot, runtime,
and domain checks as the source of verification evidence.

## Development and release

```sh
make check
make pilot       # optional CLI reference pilot
```

`make check` validates Go code, contracts, all skill mirrors, changelog, site,
and the project contract. Keep user-visible changes under `[Unreleased]` in
[`CHANGELOG.md`](CHANGELOG.md), then run `make release-notes VERSION=vX.Y.Z`
before a human-approved tag.

## Sources and license

Sillage is independently authored and MIT licensed. Its progressive,
composable skill shape is informed by the [Agent Skills specification](https://agentskills.io/specification),
[Matt Pocock's skills](https://github.com/mattpocock/skills), and established
engineering practice. Names, prompts, lifecycle contracts, evidence rules,
and human authority remain Sillage-owned; the attribution policy is in
[`docs/engineering/provenance.md`](docs/engineering/provenance.md).
