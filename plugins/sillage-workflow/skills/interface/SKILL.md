---
name: interface
description: Use automatically when a change crosses an API, HTTP, web, RPC, message, webhook, or compatibility boundary; make contracts explicit and evolvable without privileging a framework.
metadata:
  namespace: sillage
  qualified-name: "sillage:interface"
---

# Design interfaces that can evolve

An interface is a promise between independently changing consumers. Treat
HTTP, web pages, RPC, messages, files, and command protocols as observable
contracts. Start from the consumer's behavior, not from a controller or schema
shape.

## Interface pass

1. Identify provider, consumer, trust boundary, ownership, and compatibility
   window.
2. Specify inputs, outputs, errors, ordering, idempotence, limits, and
   cancellation. State what is guaranteed and what is deliberately undefined.
3. For HTTP/web, consider method semantics, status codes, headers, cache and
   freshness, content negotiation, authentication, redirects, cookies, CORS,
   pagination, and browser-visible failure states when relevant.
4. For RPC or messages, consider deadlines, delivery semantics, correlation,
   schema evolution, replay, deduplication, ordering, and poison messages.
5. Preserve backward and forward compatibility through additive change,
   tolerant readers, versioning, or an explicitly staged migration.
6. Test the real boundary with contract or integration evidence; types or a
   document alone are not proof.

Return:

```text
Consumer: <who calls or observes it>
Contract: <input/output/error and semantic guarantees>
Boundary: <protocol, trust, ownership, and compatibility window>
Failure: <timeout, invalid input, auth, overload, or partial result>
Evolution: <compatible change, version, or migration>
Proof: <contract, integration, browser, or runtime observation>
```

If the interface carries sensitive data, use `sillage:security`; if latency,
retry, ordering, or partial failure is material, use `sillage:systems` too.
