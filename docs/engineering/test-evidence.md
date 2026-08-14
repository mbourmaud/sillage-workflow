---
type: standard
status: stable
owners: [maintainers]
---

# Test evidence

Tests are a maintained proof portfolio, not an accumulating line-count target.

For each acceptance risk, name one primary owning test layer. Other layers test
their own seams without replaying the complete scenario. A regression test
names the invariant or failure it protects. Review may delete or consolidate a
test when another proof makes it redundant.

Projects define their own layers. Common categories include pure policy, wire
contract, composition, persistence, orchestration, user interaction, critical
journey, performance, security, compatibility, and manual verification.

Coverage thresholds are project policy. Exhaustive coverage is valuable for
closed critical algebras; blanket targets must not incentivize tests coupled to
implementation details.
