---
name: solid
description: Use automatically as a focused architecture or review lens when code changes raise responsibility boundaries, extension pressure, substitutability, interface size, or dependency direction; apply SOLID pragmatically without forcing object-oriented ceremony.
metadata:
  namespace: sillage
  qualified-name: "sillage:solid"
---

# Apply SOLID pragmatically

This is the compatibility lens for responsibility and dependency questions.
Use `sillage:architecture` when the task needs a broader pattern, topology, or
Clean Architecture decision; keep one primary owner for the decision.

SOLID is a set of questions about changeability, not a class-count target.
Apply the smallest principle that removes a real source of coupling or
regression. A functional module, data-oriented package, or procedural program
can satisfy the intent without artificial objects.

## Five checks

- **Single Responsibility:** one reason for a module or function to change;
  split policy from I/O, mapping, rendering, and orchestration when they
  evolve independently.
- **Open/Closed:** extend through stable contracts or composition when new
  variants are expected; do not pre-build an abstraction for an imaginary
  variant.
- **Liskov Substitution:** an implementation must preserve the contract's
  preconditions, postconditions, errors, and observable semantics.
- **Interface Segregation:** consumers depend on the smallest capability they
  need; do not expose a “god” interface merely to share a type.
- **Dependency Inversion:** policy depends on stable domain/application ports;
  frameworks, databases, networks, and vendors stay at an adapter boundary.

## Adjacent heuristics

Use these as questions, never as slogans:

- **KISS:** is the simplest design that satisfies the current contract easier
  to understand and operate?
- **DRY:** is duplicated *knowledge* causing drift, or is repetition keeping
  two independent policies clear? Deduplicate the former only.
- **YAGNI:** is this abstraction, configuration, or extension point required by
  an accepted use case, or merely imaginable?
- **Law of Demeter:** can a policy talk to the object/capability it owns rather
  than traversing a train of unrelated internals?
- **GRASP:** does responsibility sit with the information expert, creator,
  controller, or polymorphic boundary that can change for a coherent reason?
- **MoSCoW:** when scope is under pressure, classify must/should/could/won't
  with the human before trading correctness for speed.

Test-driven development and short feedback loops are implementation tactics;
use them when they reduce risk, not as a ceremony gate for every text-only or
configuration-only change.

Clean Architecture is the boundary version of these questions: keep domain
policy independent from delivery, persistence, frameworks, and vendors; let
application orchestration depend on explicit ports; keep adapters at the edge.
Do not create layers with no policy or ownership to protect.

## Return a small design note

```text
Pressure: <observed change or coupling>
Principle: <one or more, only if evidenced>
Smallest move: <boundary/composition/test change>
Trade-off: <complexity introduced or deliberately avoided>
Proof: <test, typecheck, or review observation>
```

Do not rename code or add layers just to display SOLID. If the pressure is
domain language or an invariant, use `sillage:ddd`; if it is UI or persistence,
use the corresponding specialist and keep the Sillage lifecycle unchanged.
