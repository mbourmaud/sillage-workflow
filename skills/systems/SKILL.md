---
name: systems
description: Use automatically when behavior depends on networks, concurrency, processes, queues, retries, time, resource limits, or partial failure; design for failure as a normal state.
metadata:
  namespace: sillage
  qualified-name: "sillage:systems"
---

# Make failure and time explicit

Distributed and concurrent systems do not provide one shared state, one
clock, or one reliable operation. Model the failure modes before selecting a
library or topology.

## Systems pass

1. Name the resource, owner, deadline, consistency need, and failure domain.
2. Define timeout/deadline, cancellation, retry, backoff, and maximum work.
   Never retry a non-idempotent operation without an explicit deduplication or
   compensation strategy.
3. State delivery, ordering, duplication, replay, and poison-message behavior
   for queues and events.
4. Consider connection lifecycle, DNS, TLS, proxies, TCP/UDP behavior, flow
   control, backpressure, and resource exhaustion only where they affect the
   observed risk.
5. Design degraded, unavailable, overloaded, and recovered states. Preserve
   user intent where possible and make loss visible.
6. Add observability that distinguishes timeout, cancellation, overload,
   dependency failure, and application rejection.

Return:

```text
Failure domain: <resource and dependency>
Guarantee: <availability, ordering, consistency, durability, or none>
Timing: <deadline, timeout, cancellation, retry budget>
Recovery: <deduplication, replay, compensation, or explicit loss>
Limits: <queue, memory, connection, CPU, disk, or rate bound>
Proof: <fault injection, concurrency, integration, load, or runtime check>
```

Do not add a queue, circuit breaker, saga, cache, or retry loop merely because
it is a familiar pattern. If its invariant changes product behavior, return to
`sillage:shape` before implementation.
