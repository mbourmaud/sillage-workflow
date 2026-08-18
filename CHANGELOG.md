# Changelog

All notable changes to Sillage are documented here. This project follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and uses
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.2.0] - 2026-08-18

### Added

- Full workflow kernel with decision-bound review and handoff artifacts.
- Read-only `context` and `status` views for fresh-agent cold starts.
- Explicit, atomic `transition --write` for local task records, with optimistic
  concurrency and symlink protections.
- Initial orchestration skill, compact workflow templates, and a complete
  executable pilot.
- A responsive public website with workflow, installation, and release pages.
- Full-workflow installation and deliberate-update guidance for Agent Skills
  clients.
- Provider-neutral execution profiles for matching reasoning capability and
  effort to each workflow stage without naming a model.
- Provider-neutral delegation requests for bounded parent or subagent work,
  with explicit isolation and return packets.
- Codex adapter profiles for decision research, isolated implementation, and
  read-only review.
- A canonical provenance policy documenting inspirations, naming boundaries,
  and MIT attribution rules.
- A skills-first V1 suite with explicit `using-sillage` entry and compact
  Markdown task-card guidance.
- Language-neutral architecture, testing, interface, systems, security, and
  platform lenses with a shared engineering doctrine for pressure, seams,
  patterns, failure, compatibility, and proof.

### Changed

- Release validation now requires a versioned changelog section before a tag
  can publish.
- Specialist lenses are progressive and non-prescriptive: DDD, SOLID, Clean
  Architecture, patterns, HTTP, networking, security, and platform practices
  are selected by observed risk rather than imposed on every project.

## [0.1.0] - 2026-08-14

### Added

- Human-governed, forge-neutral task lifecycle and JSON contracts.
- Canonical project entry points and deterministic `doctor`, `digest`, and
  read-only `transition` commands.
- Agent Plugins 1.0 manifests for Codex and Claude Code distribution.
- Initial evidence-led research Agent Skill with evaluation prompts and a real
  installation pilot.

[unreleased]: https://github.com/mbourmaud/sillage-workflow/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/mbourmaud/sillage-workflow/releases/tag/v0.2.0
[0.1.0]: https://github.com/mbourmaud/sillage-workflow/releases/tag/v0.1.0
