---
type: product
status: stable
owners: [maintainers]
---

# Sillage product

## Purpose

Sillage helps a human and coding agents deliver small, understandable changes
without losing decisions, evidence, or the ability to resume in a fresh
context.

## Users

- Engineers using coding agents on existing or new projects.
- Product and domain experts who retain decision authority.
- Reviewers who need a concise link from intent to implementation evidence.
- Teams using GitHub, GitLab, Jira, local files, or other work systems.

## Product promises

1. A fresh agent can understand a worktree without previous chat history.
2. One active implementation worktree represents one task slice.
3. Operational history does not pollute durable repository knowledge.
4. Verification and review remain separate, evidence-backed activities.
5. No language, forge, CI system, or package manager is privileged.
6. Humans control decisions and consequential external actions.

7. A project can use the workflow with portable Markdown skills and its own
   native checks; no Sillage executable is required.

8. A reviewer can understand a vertical from a short TL;DR and expand into
   decisions and evidence only when needed.

9. A project can apply durable, language-neutral engineering disciplines —
   architecture, testing, domain modeling, interfaces, systems, security, and
   platform design — without adopting a framework or prescribed topology.

## Non-goals

- Autonomous product management or deployment.
- A replacement for issue trackers, Git, CI, or code review platforms.
- Mandatory multi-agent orchestration.
- Generating documentation for every thought or experiment.
- Maximizing test count or Markdown volume.
- Requiring a Sillage CLI, JSON task store, Git provider, or plugin host.
- Mandating DDD, SOLID, Clean Architecture, a design pattern, or a specific
  test pyramid for every project.
- Replacing specialist technical standards, security assessments, or platform
  documentation with a generic checklist.
