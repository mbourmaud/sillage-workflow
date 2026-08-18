---
name: testing
description: Use automatically when designing a test strategy, choosing a seam, writing a regression, or deciding which test layer owns a risk; optimize for behavioral confidence, not test count.
metadata:
  namespace: sillage
  qualified-name: "sillage:testing"
---

# Test the behavior at its real seam

Testing is a risk allocation decision. Begin with the observable behavior,
the consumer, and the failure that would matter. A test is useful when it can
disagree with the implementation and survive an internal refactor.

## Testing pass

1. State the acceptance behavior and the risk it protects.
2. Identify the public seam: domain operation, application use case, protocol,
   command, UI journey, or real persistence boundary.
3. Assign one primary owning layer to the risk. Add another layer only when it
   proves a different failure mode.
4. For code behavior, write one focused failing test, observe the failure, make
   the smallest change, observe green, then refactor only with proof intact.
5. Keep expected values independent from the implementation. Avoid tautology,
   private-method tests, internal mocks, snapshot noise, and horizontal test
   batches that encode imagined behavior.
6. Make time, randomness, network, processes, and concurrency deterministic at
   the seam; test cancellation, retry, timeout, and cleanup when relevant.

## Portfolio choices

- Domain tests protect invariants and transitions.
- Application tests protect orchestration and policy composition.
- Contract tests protect a provider or protocol boundary.
- Integration tests protect real storage, network, process, or browser risk.
- End-to-end tests protect a critical user journey; keep them few and real.
- Property, fuzz, load, and mutation checks are justified by a specific risk.

Return:

```text
Behavior: <observable outcome>
Risk: <failure and affected consumer>
Seam: <public boundary and why>
Owner: <primary test layer>
Cases: <normal, edge, failure, cancellation, or concurrency cases>
Proof: <command and expected observation>
```

`build` owns the implementation loop; `prove` owns fresh evidence; this lens
owns the test design. Never claim a regression is fixed from a passing unit
test when the risk lives at a real external boundary.
