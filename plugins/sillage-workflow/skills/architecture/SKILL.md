---
name: architecture
description: Use automatically when a change raises module boundaries, dependency direction, architectural patterns, or pressure to add abstractions; choose the smallest structure that protects an observed change.
metadata:
  namespace: sillage
  qualified-name: "sillage:architecture"
---

# Shape architecture around pressure

Architecture is the set of boundaries that makes important change safe. It is
not a layer count, a framework choice, or a pattern catalogue. Start from the
behavior, the owners of decisions, and the reasons the code may change.

## Architecture pass

1. Name the policy, data, delivery, and infrastructure concerns that are
   actually present.
2. Identify the smallest stable seams where consumers observe behavior.
3. Record dependency direction: policy must not accidentally depend on a
   delivery mechanism, persistence detail, vendor, or process boundary.
4. Select a pattern only when it removes a measured coupling or protects a
   likely change; state the complexity it adds.
5. Keep modules deep enough to hide implementation, but do not manufacture
   ports, repositories, events, or layers for imaginary variants.
6. Describe failure, ownership, compatibility, and migration consequences.

Useful patterns include composition, adapter, facade, strategy, command,
state, pipeline, ports-and-adapters, layered, plugin, event-driven, CQRS, and
saga. They are options, not requirements. A functional or data-oriented design
can satisfy the same boundary goals without artificial objects.

## SOLID and Clean Architecture relationship

Use SOLID, KISS, DRY, YAGNI, GRASP, and Clean Architecture as questions about
changeability, locality, and dependency direction. Do not rename code or add a
layer merely to display a principle. DRY removes duplicated knowledge, not all
repetition; Clean Architecture protects policy only where policy is real.

Return:

```text
Pressure: <observed change, coupling, or failure>
Boundary: <owner and stable seam>
Option: <smallest pattern or no new pattern>
Cost: <complexity, latency, migration, or operational cost>
Trade-off: <what is deliberately not protected>
Proof: <test, dependency check, runtime observation, or review>
```

If the pressure is a domain invariant, use `sillage:ddd`; if it is a protocol,
distributed system, security, or platform concern, use the matching lens.
Return to `sillage:shape` when the boundary changes product behavior or scope.
