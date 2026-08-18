# Sillage skills

Sillage is a skills-first protocol. Each folder contains one portable
`SKILL.md` and optional host metadata. Descriptions are deliberately specific
so a host can invoke the right skill from ordinary language; `sillage` remains
the only lifecycle router.

## Published suite

### Route and deliver

- `using-sillage` — explicit “start Sillage” handshake for a fresh
  conversation; returns a state card and delegates to the router.
- `sillage` — implicit router, lifecycle, authority, and completion contract.
- `orient` — fresh-worktree state card and next safe action.
- `shape` — goals, context, decisions, and durable ADR boundaries.
- `slice` — one vertical, worktree, acceptance set, and stop condition.
- `build` — one approved change with a focused failing test when relevant.
- `prove` — deterministic verification and current evidence.
- `review` — independent assessment sized for a human reviewer.
- `handoff` — compact resume point, limits, and human-owned next action.
- `research` — source authority, freshness, inference, and durable-knowledge discipline.

### Engineering lenses

- `architecture` — observed pressure, module seams, dependency direction, and
  justified design patterns.
- `testing` — behavioral seams, risk ownership, red-green loops, and proof
  portfolios.
- `solid` — pragmatic SOLID, Clean Architecture, KISS, DRY, YAGNI, GRASP, and
  MoSCoW questions about responsibility, contracts, and dependency direction.
- `ddd` — ubiquitous language, bounded contexts, aggregates, invariants, and
  domain/application boundaries.
- `interface` — HTTP, web, RPC, message, webhook, file, and command contracts.
- `systems` — concurrency, timing, network behavior, resource limits, queues,
  recovery, and partial failure.
- `security` — assets, actors, trust boundaries, abuse cases, controls, and
  residual risk.
- `platform` — desktop, mobile, CLI, service, embedded, process, filesystem,
  packaging, and operating-system behavior.
- `frontend-architecture` — user journeys, UI state completeness, route/data
  ownership, accessibility, responsive behavior, and browser proof.
- `relational-data` — relational invariants, constraints, transactions,
  indexes, query ownership, and reversible migrations.
- `document-data` — non-relational access patterns, record shape, schema
  versioning, consistency, partitioning, and reconciliation.
- `audit` — observed facts, ranked risk, and a bounded remediation path for
  legacy and technical debt.
- `migrate` — staged compatibility, rollback, and proof for legacy, API,
  dependency, architecture, or data changes.
- `debug` — reproducible, hypothesis-driven diagnosis before a fix.
- `test-hygiene` — a small proof portfolio with one primary test owner per
  risk, rather than a line-count target.

Install the whole suite from GitHub:

```sh
npx skills add mbourmaud/sillage-workflow \
  --skill using-sillage --skill sillage --skill orient --skill shape --skill slice \
  --skill build --skill prove --skill review --skill handoff \
  --skill research --skill architecture --skill testing --skill solid --skill ddd \
  --skill interface --skill systems --skill security --skill platform \
  --skill frontend-architecture --skill relational-data --skill document-data \
  --skill audit --skill migrate --skill debug --skill test-hygiene
```

Or install a single lens, for example:

```sh
npx skills add mbourmaud/sillage-workflow --skill ddd
```

The plugin bundle mirrors these files for Codex and Claude Code. Contract tests
fail if a mirror drifts. The Go CLI and JSON schemas are optional companions;
plain Markdown task cards and project-native checks are enough to use Sillage.

## Skill design contract

- One skill protects one habit or bounded workflow.
- The frontmatter description says what it does and when it triggers.
- `agents/openai.yaml` permits implicit invocation when the host supports it.
- Output is a compact packet that can feed the next skill.
- Specialists are lenses, never alternate lifecycle authorities.
- No skill declares success from inspection alone or invents human approval.
- Durable decisions go to the existing canonical owner; transient notes do not
  become repository pollution.
- Patterns, DDD, SOLID, Clean Architecture, KISS, DRY, HTTP, networking,
  security, and platform conventions are selected by observed pressure. They
  are never universal ceremony or proof by themselves.

Behavioral prompts and pilot results live under [`../evals`](../evals). Keep
claims about maturity honest until the same prompts have baseline, with-skill,
objective, and human review evidence.
