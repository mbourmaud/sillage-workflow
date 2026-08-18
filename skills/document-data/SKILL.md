---
name: document-data
description: Use automatically for document databases, key-value stores, wide-column data, search indexes, event records, schema evolution, partition keys, or non-relational persistence reviews.
metadata:
  namespace: sillage
  qualified-name: "sillage:document-data"
---

# Design non-relational data deliberately

Start from read/write access patterns and consistency needs. A document or
key-value record is an owned aggregate-shaped representation, not a free-form
escape from domain rules.

## Review sequence

1. Name the owning context, record identity, partition/tenant key, and
   lifecycle.
2. List critical reads, writes, fan-out, size/cardinality, hot-key, and query
   patterns before choosing embedding, referencing, or multiple projections.
3. Define schema versioning, validation, defaults, compatibility, and how old
   records are read or migrated.
4. State consistency, idempotency, ordering, retry, and reconciliation rules;
   make eventual consistency visible to users and operators.
5. Choose indexes/secondary access paths with their write, cost, and rebuild
   limits. Protect sensitive fields and retention boundaries.
6. Verify a representative real read/write path, not only a serialization
   test. Record backfill, replay, or rollback evidence when data changes.

Return:

```text
Access pattern: <read/write and scale shape>
Record: <key, partition, embedded/referenced fields>
Schema: <version and compatibility rule>
Consistency: <guarantee, retry, reconciliation>
Operational risk: <hot key, size, index, retention, or migration limit>
Proof: <runtime, query, migration, or replay observation>
```

If the representation changes the business model, route to `sillage:ddd`; if
it changes product scope, route to `sillage:shape`. Keep this lens inside the
same Sillage lifecycle and prove the real persistence boundary.
