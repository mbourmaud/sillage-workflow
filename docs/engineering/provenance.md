---
type: provenance
status: stable
owners: [maintainers]
sources:
  - https://github.com/mattpocock/skills
  - https://skills.sh/
  - https://www.aihero.dev/skills
  - https://raw.githubusercontent.com/mattpocock/skills/main/LICENSE
  - https://agentskills.io/specification
  - https://agent-plugins.org/specification
---

# Sources and inspiration

Sillage is independently authored. It is informed by established engineering
practice and by other agent-workflow projects, but its names, contracts,
artifacts, authority model, and implementation are owned by Sillage.

This page is the single durable record of those influences. It is intentionally
short: operational research belongs in the active task record, while source
attribution belongs here.

## Provenance policy

- Reuse general ideas and engineering principles; do not copy another
  project's skill names, prompts, prose, or artifact contracts by default.
- Give every Sillage skill a portable name and a Sillage-owned description.
  Source-specific names remain references, not Sillage API names.
- Translate a useful pattern into Sillage's lifecycle, task schema, human
  gates, evidence rules, and handoff model. Do not introduce a competing
  lifecycle or status vocabulary.
- Cite the source that materially influenced a design. Keep citations in this
  page and in the public site rather than duplicating attribution in every
  skill.
- If Sillage later copies a substantial portion of MIT-licensed code or text,
  retain the original copyright and license notices in the copied material and
  add the required project notice. No third-party notice is needed for an
  independently written implementation.

## Influences and Sillage adaptations

| Source | Useful influence | Sillage-owned adaptation |
| --- | --- | --- |
| [Matt Pocock's skills](https://github.com/mattpocock/skills) and [AIHero overview](https://www.aihero.dev/skills) | Small, composable instructions; progressive disclosure; a chain in which one skill's output can feed the next; separate orchestration from reusable discipline; explicit design interrogation and test seams. | `sillage` is the sole orchestration skill, phase skills return compact packets, and architecture/testing/DDD/interface/systems/security/platform/data lenses stay focused. The task card, lifecycle gates, decision digest, evidence receipts, review brief, and resumable handoff are Sillage contracts. |
| [obra/superpowers](https://github.com/obra/superpowers) | Explicit approval before implementation, red/green discipline, fresh verification, worktree isolation, and independent review loops. | Sillage adopts the invariant and evidence intent, while keeping its own statuses, compact artifacts, human authority, optional delegation, and provider-neutral portability. |
| [skills.sh](https://skills.sh/) | A discoverable ecosystem where focused, technical leaf skills complement workflow skills; popularity favors concrete verbs and immediately useful specialist guidance. | Sillage keeps one router and exposes portable lenses with explicit ownership, progressive loading, and no automatic authority transfer from an installed third-party skill. |
| [Agent Skills specification](https://agentskills.io/specification) | A portable `SKILL.md` shape and client-discoverable naming constraints. | Sillage keeps the portable file format while keeping its own vocabulary and workflow semantics. |
| [Agent Plugins specification](https://agent-plugins.org/specification) | A thin package boundary for distributing skills to compatible hosts. | The plugin is an optional distribution adapter. It does not define Sillage's CLI, schemas, task store, or authority model. |

## What this means for contributors

When a new capability resembles an existing external skill, start by writing
the Sillage invariant it must protect and the observable packet it must return.
Choose a distinct name, implement only the smallest useful behavior, and link
the influence here if it materially shaped the design. A familiar idea is not
a reason to reproduce a familiar repository.

The current public skill names are `using-sillage`, `sillage`, `orient`, `shape`, `slice`,
`build`, `prove`, `review`, `handoff`, `research`, `architecture`, `testing`,
`solid`, `ddd`, `interface`, `systems`, `security`, `platform`,
`frontend-architecture`, `relational-data`, and `document-data`. They remain
distinct, portable, and meaningful without relying on a vendor's private
namespace. Maintenance lenses include `audit`, `migrate`, `debug`, and
`test-hygiene`.
