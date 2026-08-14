---
name: researching-with-evidence
description: Use when a task depends on unfamiliar, disputed, version-sensitive, or externally documented facts; when reading documentation, source code, specifications, papers, or tutorials; or when research findings may influence a product or engineering decision.
---

# Researching with evidence

## Core principle

Research reduces decision uncertainty. It does not manufacture a decision or a
new repository document.

Record operational findings in the active task record. Promote only knowledge
that future work would otherwise need to rediscover.

## Establish the question

Before searching, state:

- the question;
- the decision or implementation it informs;
- the project version or environment that matters;
- what would count as sufficient evidence;
- the required freshness.

If no active task record exists, return a self-contained evidence packet in the
conversation. Do not invent a provider or storage location.

If an essential input such as the subject, relevant version, or consuming
decision is missing, ask one focused question and stop. Do not render an empty
evidence packet or enumerate every possible unknown before research begins.

## Follow source authority

Prefer evidence in this order, adjusting when the subject has a more appropriate
authority:

1. behavior observed in the relevant installed version or owned source code;
2. official specification or standard;
3. official documentation, release notes, and first-party examples;
4. primary research or authoritative institutional publication;
5. expert tutorials and secondary explanations;
6. community discussion.

Prefer primary sources for material claims; use secondary sources to discover
and interpret them.

Tutorials are valuable for discovering terminology, examples, and edge cases.
Trace any material claim back to the source that owns it when possible. When
that is impossible, preserve the weaker authority instead of upgrading it by
paraphrase.

Use parallel investigation only when the runtime supports it and the question
has independent branches. Delegation is an optimization, never a requirement.

## Inspect the actual project

For version-sensitive technical research:

- identify the installed or deployed version;
- inspect local source and configuration before generic online guidance;
- compare documentation with runtime behavior when the claim matters;
- record applicable versions, access date, and environment;
- identify drift between the project and current upstream documentation.

## Build an evidence packet

Use this shape in the task record:

```markdown
### Research: <question>

**Consumer:** <decision or implementation>
**As of:** <date and relevant versions>

#### Conclusion
<short answer, including confidence and limits>

#### Claims
- **Claim:** <statement>
  - Authority: <installed behavior | specification | official docs | tutorial | discussion>
  - Source: <addressable reference and locator>
  - Observed: <date/version/environment>
  - Implication: <why this matters>

#### Conflicts and inference
- <conflicting evidence, explicit inference, or none>

#### Unknowns
- <remaining uncertainty or none>

#### Durable knowledge candidate
<canonical page to update and why, or "none">
```

Keep quotations short. Put evidence beside the claim it supports. Label
synthesis and inference explicitly; citations do not automatically support an
agent's conclusion.

## Promote durable knowledge sparingly

Update durable knowledge when losing the finding would force future
contributors to rediscover an important product boundary, domain invariant,
architecture constraint, standard, runbook, or landmark decision.

Prefer updating an existing canonical page. Create a new page only when the
knowledge has a distinct durable owner. Preserve provenance, lifecycle status,
and freshness where the project convention supports them.

Do not promote raw notes, copied tutorials, discarded approaches, experiments,
temporary plans, logs, screenshots, or task history.

## Stop at the authority boundary

Research may recommend a choice. It cannot approve product behavior, expand
scope, waive evidence, or authorize an external write. Return the task to its
decision gate when the findings change an approved assumption.

## Quality check

Before finishing, verify:

- every material claim has an appropriately authoritative source;
- installed versions and freshness are explicit where relevant;
- facts, synthesis, inference, conflict, and unknowns are distinguishable;
- the result answers the original question;
- operational material stayed in the task record;
- any durable promotion updates the correct canonical owner.
