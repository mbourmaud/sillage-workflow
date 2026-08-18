---
name: security
description: Use automatically when a change handles identity, secrets, permissions, untrusted input, sensitive data, dependencies, or a trust boundary; produce a threat-informed security packet.
metadata:
  namespace: sillage
  qualified-name: "sillage:security"
---

# Reduce security risk at the boundary

Security is a property of assets, actors, trust boundaries, and failure
behavior. A checklist cannot prove that a system is secure; it can expose the
risks that require design, testing, or specialist review.

## Security pass

1. Name assets, actors, trust boundaries, abuse cases, and impact (confidentiality,
   integrity, availability, privacy, or accountability).
2. Apply least privilege, deny-by-default authorization, explicit session and
   credential lifecycle, and safe secret handling.
3. Validate at the boundary and encode for the sink. Consider injection,
   confused deputy, path traversal, SSRF, deserialization, CSRF, XSS, replay,
   request smuggling, and resource exhaustion where applicable.
4. Protect data in transit and at rest; minimize collection, retention,
   exposure, logs, and error detail.
5. Review dependencies, build provenance, update paths, permissions, and
   operational access. Make security failures observable without leaking the
   protected data.
6. Assign each risk an owning test, review, scanner, or external assessment;
   record unknowns and residual risk honestly.

Return:

```text
Asset: <what needs protection>
Threat: <actor, path, and impact>
Boundary: <where trust changes>
Control: <smallest preventive/detective measure>
Residual risk: <known limit and owner>
Proof: <security test, review, scanner, or external assessment>
```

Never claim “secure” from inspection or a green linter. Route regulated,
high-impact, cryptographic, authentication, or supply-chain decisions to the
appropriate human or specialist authority.
