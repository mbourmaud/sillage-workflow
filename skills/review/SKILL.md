---
name: review
description: Use automatically when a focused diff, pull request, change set, or implementation asks whether it is safe and complete; independently compare code and evidence with the approved Sillage slice.
metadata:
  namespace: sillage
  qualified-name: "sillage:review"
---

# Review the vertical

Review is an independent assessment, not a second implementation pass and not
human approval. Start from the task TL;DR, non-goals, acceptance, and current
decision digest; then inspect the diff, tests, runtime evidence, and nearby
contracts.

## Review order

1. Can a reviewer explain the vertical in one sentence?
2. Does every changed behavior belong to the approved slice?
3. Does each acceptance criterion have current, addressable evidence?
4. Are failure paths, security, compatibility, migration, accessibility,
   protocol, concurrency, platform, and operational limits covered for this
   risk?
5. Is there dead code, duplicated policy, speculative abstraction, or test
   volume that hides the behavior?
6. Does the handoff state what remains human-owned?

Return a compact packet:

```text
TL;DR: <one sentence verdict>
Status: accepted | changes_requested
Digest: <decision digest reviewed>
Blocking: <finding, reference, or none>
Non-blocking: <small findings or none>
Evidence gaps: <none or exact missing proof>
Next: <one safe action>
```

Use one primary specialist lens for architecture, testing, domain, interface,
systems, security, platform, frontend, or data risks when needed, but keep one
review owner and one lifecycle. A self-review may expose risk; it does not
waive the human gate for scope, migration, external writes, merge, deployment,
or release.
