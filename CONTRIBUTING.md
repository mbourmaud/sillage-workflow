# Contributing

Sillage develops workflow skills like production behavior: establish a failing
baseline, make the smallest guidance change, compare fresh runs, and preserve
human review of qualitative results.

## Before changing a skill

1. Add or update realistic prompts under `evals/<skill>/evals.json`.
2. Run the prompts without the proposed guidance and retain the failure.
3. Change one skill.
4. Run the same prompts with the skill.
5. Grade objective requirements and inspect qualitative output.
6. Record the skill's maturity honestly.

Do not batch untested skills. Keep the SKILL.md concise; mechanical constraints
belong in the project's deterministic checks rather than increasingly forceful
prose. A specialist must return an observable design/proof packet, not a
catalogue of principles.

## Code changes

Write a focused failing test before changing observable CLI or policy behavior;
for a Markdown-only skill, add an objective contract assertion or evaluation
prompt before changing its guidance.

```sh
make check
```

This is the authoritative local and CI gate. Run the smallest affected Go test
first, then `make check` before requesting review.

Keep adapters optional and keep the core free of assumptions about language,
forge, task store, CI provider, or agent runtime.

Sillage is the sole workflow authority for this repository. New skills and
automation must use its statuses, task contracts, gates, and evidence model;
do not add a competing orchestration layer.

## Releases

Record user-visible changes under the `[Unreleased]` section of
[`CHANGELOG.md`](CHANGELOG.md). Before tagging a release, move those entries
to a dated version section, leave a new `[Unreleased]` section at the top, and
run:

```sh
make check
make release-notes VERSION=vX.Y.Z
```

The release workflow rejects a tag without a matching non-empty changelog
section and uses that section as the published release notes. Tagging,
publishing, merging, and external announcements remain human decisions.
