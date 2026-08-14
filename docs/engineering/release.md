---
type: release-process
status: stable
owners: [maintainers]
---

# Release process

Sillage uses one human-readable [`CHANGELOG.md`](../../CHANGELOG.md) as the
release contract. It follows Keep a Changelog categories and keeps unreleased
work separate from published versions.

## During development

Add user-visible changes under `## [Unreleased]` using one of:

- `Added` — new capability;
- `Changed` — behavior or contract change;
- `Fixed` — corrected behavior;
- `Removed` — removed capability;
- `Security` — security-relevant correction.

Do not create a version heading before the release is approved. Internal
refactors that have no user-visible effect do not need a changelog entry.

## Before a tag

1. Move the intended `[Unreleased]` entries into a dated `## [X.Y.Z]` section.
2. Leave a new, empty `## [Unreleased]` section at the top.
3. Update the versioned plugin/manifests when the release includes them.
4. Run `make check` and `make release-notes VERSION=vX.Y.Z`.
5. Review the extracted notes and authorize the tag as a human release
   decision.

The release workflow repeats the version-specific changelog check before
building. It extracts the exact section into the GitHub release notes, so a tag
cannot publish without a corresponding documented release.

## Local commands

```sh
make changelog-check
make release-notes VERSION=v0.2.0
go run ./cmd/sillage changelog check --version v0.2.0
go run ./cmd/sillage changelog extract --version v0.2.0
```

The changelog commands are read-only. They do not create tags, releases, or
external writes.
