---
name: research
description: Use automatically when a task depends on unfamiliar, disputed, version-sensitive, or externally documented facts; read specifications, source code, papers, release notes, or tutorials and return a traceable evidence packet without polluting the repository.
metadata:
  namespace: sillage
  qualified-name: "sillage:research"
---

# Research with evidence

Research reduces decision uncertainty. It does not manufacture a decision,
approve scope, or create a document for every search.

## Establish the question

Before searching, state the question, consumer decision, relevant version or
environment, sufficient evidence, and freshness. If an essential input is
missing, ask one focused question and stop. If no task record exists, return a
self-contained packet in the conversation.

## Follow authority

Prefer primary sources, in this order: observed behavior in the installed/owned source, official
specification or standard, official documentation/release notes/examples,
primary research or institutional publication, expert tutorial, then
community discussion. Tutorials help discover terminology and edge cases; they
do not upgrade a weak claim into a primary fact. Record source, locator,
version, access date, and drift.

Use parallel investigation only when the runtime supports it and branches are
independent. Delegation is an optimization, never a requirement.

## Return one evidence packet

```markdown
### Research: <question>

**Consumer:** <decision or implementation>
**As of:** <date and relevant versions>

#### Conclusion
<short answer, confidence, and limits>

#### Claims
- **Claim:** <statement>
  - Authority: <observed | specification | official docs | tutorial | discussion>
  - Source: <addressable reference and locator>
  - Observed: <date/version/environment>
  - Implication: <why it matters>

#### Conflicts and inference
- <conflict, synthesis, or none>

#### Unknowns
- <remaining uncertainty or none>

#### Durable knowledge candidate
<existing canonical owner to update, or “none”>
```

Label fact, synthesis, inference, conflict, and unknown separately. Keep
quotations short. Put operational findings in the active task; promote durable
knowledge only to its canonical owner, with provenance and freshness.

## Stop at the authority boundary

Research may recommend a choice. It cannot approve product behavior, expand
scope, waive evidence, or authorize an external write. If findings change an
approved assumption, route to `sillage:shape`.

Before finishing, check that every material claim has suitable authority, the
original question is answered, versions/freshness are visible, and no raw
research, copied tutorial, temporary plan, log, or screenshot was promoted into
durable knowledge.
