---
name: slice
description: Use automatically when an approved task must become one small vertical, worktree, branch, or independently reviewable implementation plan with explicit acceptance and non-goals.
metadata:
  namespace: sillage
  qualified-name: "sillage:slice"
---

# Slice the work

Slice protects the reviewer and the next fresh agent. One active worktree owns
one behavior, one decision boundary, and one independently reviewable diff.

## Write the smallest plan

```text
TL;DR: <what this vertical changes>
Task: <approved intent and decision reference>
Slice: <one behavior from input to observable result>
Files/boundary: <likely ownership, not a speculative file list>
Non-goals: <what this slice will not touch>
Acceptance: <two to five observable criteria>
Test owner: <one primary layer per risk>
Runtime proof: <check, screenshot, build, or domain observation>
Dependencies: <none or named prerequisite>
Stop condition: <when to split or return to shape>
```

Split when the slice needs unrelated product decisions, crosses multiple
owners, or creates a diff a reviewer cannot summarize. Prefer a clean worktree
and a short branch; do not require Git if the host uses another forge. Do not
create a plan per conversation turn: update the active task card.

Before implementation, verify that the decision is accepted and the human
knows the boundary. If a dependency is unavailable, enter `BLOCKED` with the
status to resume from and an observable condition. Then route to
`sillage:build`.
