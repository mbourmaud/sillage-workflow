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

Do not batch untested skills. Mechanical constraints belong in the CLI or a
schema rather than increasingly forceful prose.

## Code changes

Write a focused failing test before changing observable CLI or policy behavior.

```sh
make check
```

This is the authoritative local and CI gate. Run the smallest affected Go test
first, then `make check` before requesting review.

Keep adapters optional and keep the core free of assumptions about language,
forge, task store, CI provider, or agent runtime.
