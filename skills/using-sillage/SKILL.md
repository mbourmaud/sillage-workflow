---
name: using-sillage
description: Use when a developer says start Sillage, use Sillage, begin a task, continue in a fresh worktree, or asks how to work with the Sillage workflow; provide the simple entry ritual and hand off to the sillage router.
metadata:
  namespace: sillage
  qualified-name: "sillage:using-sillage"
---

# Start using Sillage

This is the explicit, human-friendly entry point. Invoke it as `$using-sillage`
in a new conversation, or let the host select it when the user says “start
Sillage”, “use Sillage”, “démarre Sillage”, or “on commence cette tâche”. It is
an onboarding handshake, not a second lifecycle.

## The three-line start

1. Confirm the repository/worktree root and whether the user wants to inspect
   first or resume an existing task.
2. Read the smallest stable context: `PRODUCT.md`, `DESIGN.md`, `AGENTS.md`
   (and `CLAUDE.md` when present), `docs/domain/index.md`, then the active task
   card if one exists.
3. Return the state card and route to exactly one next skill: `sillage:orient`
   for a fresh/unclear context, `sillage:shape` for an unclear outcome,
   `sillage:slice` for an accepted task without a vertical, or the current
   phase skill when the task is already in progress.

Use native project commands when the user asks for checks. A Sillage CLI,
plugin, JSON profile, Git worktree, or task tracker is optional; never make
installation of one a prerequisite for starting.

## First response

Keep the first response short and useful:

```text
Sillage: ready
Project: <name or unknown>
Status: <lifecycle status or INTAKE>
TL;DR: <one sentence, or “not shaped yet”>
Known: <two or three facts>
Missing: <at most two material gaps>
Next: <one safe action>
Human decision: <none, or one focused question>
```

Do not implement, create a large plan, generate an ADR, install tools, or
modify application code during this handshake unless the user explicitly asks
for that next action. If no task card exists and the request is clear, propose
one compact Markdown card; do not invent product approval.

## Hand off cleanly

After the state card, say which qualified skill will continue and why. Keep the
router authoritative:

- context missing → `sillage:orient`;
- intent/goal/architecture unclear → `sillage:shape`;
- approved work needs one vertical → `sillage:slice`;
- approved slice is being coded → `sillage:build`;
- proof/review/handoff requested → the corresponding phase skill.

When a specialist is needed, mention it as a lens inside that phase. Never
create a second status vocabulary or ask the developer to memorize the whole
suite before their first task.
