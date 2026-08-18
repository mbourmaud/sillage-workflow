---
name: prove
description: Use automatically when implementation needs tests, typechecks, builds, linters, screenshots, runtime checks, or domain validation before review; turn observable results into concise evidence.
metadata:
  namespace: sillage
  qualified-name: "sillage:prove"
---

# Prove the slice

Verification asks “does the approved behavior work?” It is separate from
review, and it never trusts code inspection alone.

## Evidence loop

1. Map each acceptance criterion to one primary owning check.
2. Run the smallest deterministic lane first, then the configured complete
   lane when the task requires it.
3. Exercise the real boundary when risk is runtime, UI, integration, or data;
   a unit test alone is not a substitute.
4. Record command/observation, result, observed time, addressable reference,
   limits, and the current decision digest beside the criterion.
5. If a check is skipped, state its impact and obtain an explicit human waiver.

Prefer a small proof portfolio over a large test count. Remove duplicate,
unowned, flaky, or obsolete tests when the task permits; do not weaken a
criterion to make a check pass. Freshness is the project's policy, not an
assumption: evidence from before the approved decision or after a material
scope change is invalid.

```text
Criterion: <observable requirement>
Check: <native command or runtime observation>
Result: passed | observed | failed | waived
Observed at: <timestamp>
Reference: <file, URL, artifact, or log locator>
Digest: <current decision digest>
Limits: <none or explicit impact>
```

Failures route to `sillage:build`, `sillage:shape`, or `BLOCKED` with a
resumable condition. Complete evidence routes to `sillage:review`.

Never make a completion claim from a previous run, a partial command, a green
linter, a child-agent report, or code inspection. Identify the exact claim,
run the fresh check, read its complete result, and report the observed status.
