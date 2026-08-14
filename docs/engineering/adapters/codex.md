---
type: adapter-guide
status: stable
owners: [maintainers]
sources: [https://developers.openai.com/codex/subagents]
---

# Codex delegation adapter

Sillage keeps the task contract portable. Codex is the host that turns a
`delegation` request and an `execution` profile into a child thread. The
repository skill tells the parent when to delegate; Codex owns spawning,
permissions, model selection, and thread cleanup.

## Local configuration

Codex supports project-scoped custom agents in `.codex/agents/`. Keep these
files local to the project or user configuration when they contain model,
budget, connector, or permission choices. Do not put those choices in a
portable task record or in the Sillage plugin manifest.

The smallest useful setup is one narrow agent per delegated role:

```toml
# .codex/agents/sillage-reviewer.toml
name = "sillage_reviewer"
description = "Independent read-only review of one Sillage task slice."
sandbox_mode = "read-only"
model_reasoning_effort = "high"
developer_instructions = """
Act as a Sillage reviewer. Read PRODUCT.md, DESIGN.md, AGENTS.md,
docs/domain/index.md, the active task, and the current diff. Review only the
bounded slice supplied by the parent. Return review_findings with severity,
file or criterion references, concrete details, and disposition. Do not edit,
transition the task, approve a decision, merge, deploy, or delegate further.
"""
```

Ready-to-copy role files live under
[`adapters/codex/agents`](../../adapters/codex/agents). Copy only the roles a
project actually uses into its own `.codex/agents/` directory; the files omit a
model so the local conversation or `[agents]` defaults remain authoritative.

The same shape can define a read-only decision researcher and an
`isolated_worktree` builder. Omit `model` when the host should inherit the
conversation or configured default. Set a concrete model only in local adapter
configuration when a team has deliberately chosen that mapping.

## Mapping a Sillage request

1. Read the current task status and resolve its `execution` and `delegation`
   entries.
2. Select one matching custom agent or the host's built-in subagent role. The
   included examples map `decision_researcher` to
   `sillage_decision_researcher`, `builder` to `sillage_builder`, and
   `reviewer` to `sillage_reviewer`.
3. Pass only the project entry points, task intent, active slice, current
   decision digest, isolation, and expected return shape.
4. Wait for the bounded packet and inspect it in the parent thread.
5. Bind usable evidence or review findings to the current digest and continue
   through the normal Sillage gate. A child cannot supply human approval.

When agents are unavailable, an optional request is recorded as a surfaced
parent fallback. A required request moves the task to `BLOCKED` with a resume
condition such as “Codex subagent capability is enabled for this session”.

## Capability boundary

Agent Plugins packages skills and MCP servers; it does not define portable
custom-agent or subagent semantics. This guide is therefore an optional Codex
adapter, not part of Sillage's portable core. Other hosts may implement the
same request shape with their own delegation mechanism or continue in the
parent conversation.
