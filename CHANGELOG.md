# Changelog

All notable changes to Sillage are documented here. This project follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and uses
[Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Full workflow kernel with decision-bound review and handoff artifacts.
- Read-only `context` and `status` views for fresh-agent cold starts.
- Explicit, atomic `transition --write` for local task records, with optimistic
  concurrency and symlink protections.
- `working-with-sillage` orchestration skill, compact workflow templates, and a
  complete executable pilot.
- A responsive public website with workflow, installation, and release pages.
- Full-workflow installation and deliberate-update guidance for Agent Skills clients.

### Changed

- Release validation now requires a versioned changelog section before a tag
  can publish.

## [0.1.0] - 2026-08-14

### Added

- Human-governed, forge-neutral task lifecycle and JSON contracts.
- Canonical project entry points and deterministic `doctor`, `digest`, and
  read-only `transition` commands.
- Agent Plugins 1.0 manifests for Codex and Claude Code distribution.
- `researching-with-evidence` Agent Skill with evaluation prompts and a real
  installation pilot.

[unreleased]: https://github.com/mbourmaud/sillage-workflow/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mbourmaud/sillage-workflow/releases/tag/v0.1.0
