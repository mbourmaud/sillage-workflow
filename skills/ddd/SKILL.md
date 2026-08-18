---
name: ddd
description: Use automatically when a feature, refactor, or review involves domain language, bounded contexts, aggregates, entities, value objects, invariants, domain services, or business events; model only the domain complexity that is actually present.
metadata:
  namespace: sillage
  qualified-name: "sillage:ddd"
---

# Model the domain that exists

Domain-driven design is a language and boundary discipline, not a reason to
wrap every value in a class. Start from the business outcome and observed
rules; keep technical plumbing out of the domain core.

## Modeling pass

1. Write the ubiquitous language in the project's canonical domain page; flag
   synonyms and overloaded words.
2. Identify bounded contexts and their ownership. Do not share an entity just
   because two contexts use the same noun.
3. Choose the smallest model: value object for an immutable concept, entity
   for identity and lifecycle, aggregate for a consistency boundary, domain
   service for a rule that has no natural owner.
4. State invariants and valid transitions before choosing persistence or an
   event shape.
5. Keep application orchestration, ports, adapters, and UI outside the domain
   core. Translate at boundaries with explicit contracts.
6. Use events only when a real consumer needs a durable fact; do not add an
   outbox or ceremony speculatively.

Return:

```text
Language: <terms and definitions>
Boundary: <context and ownership>
Invariant: <rule and rejected state>
Model: <value/entity/aggregate/service, with reason>
Contract: <input/output or event if needed>
Unknown: <one question that changes the model>
```

If the model changes product behavior or scope, route to `sillage:shape` for
human approval. Verify invariants with domain tests and real boundary checks;
do not call a diagram or type definition proof by itself.
