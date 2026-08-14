# Sillage skills

Every skill in this directory is Sillage-owned and follows the Sillage
lifecycle and authority model. Skills are published independently so projects
install only the judgment they need; the deterministic CLI and schemas remain
optional companions.

## Released in 0.1.0

### `researching-with-evidence`

Researches unfamiliar or version-sensitive questions, preserves source
authority and freshness, stores operational findings with the task, and
promotes durable knowledge sparingly.

```sh
npx skills add mbourmaud/sillage-workflow --skill researching-with-evidence
```

### `working-with-sillage` (workflow candidate)

The orchestration skill for cold starts, bounded slices, human decision gates,
deterministic verification, independent review, resumable blocking, and
handoff. It also recommends provider-neutral capability and effort profiles per
stage, without naming models or changing human authority. It can also request
bounded parent or subagent work with explicit isolation and return packets;
host adapters own the actual thread and model selection. It routes external
research to `researching-with-evidence` instead of creating a competing
lifecycle. Its evaluation prompts are still draft until the full workflow has
been exercised on a real project.

```sh
npx skills add mbourmaud/sillage-workflow --skill working-with-sillage
```

Skills are not updated in the background. Refresh an installed skill at a
task boundary and run the project's deterministic checks afterwards:

```sh
npx skills update working-with-sillage --global --yes
```

Behavioral evaluation fixtures live under `evals/`. The first skill was
released after baseline comparison, a real Codex installation pilot, and human
review; the roadmap is not an inventory of future skills.
