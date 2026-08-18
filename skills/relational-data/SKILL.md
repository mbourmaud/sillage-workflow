---
name: relational-data
description: Use automatically for SQL schema design, relational migrations, constraints, transactions, indexes, query performance, or persistence reviews; make invariants and rollback evidence explicit without over-modeling.
metadata:
  namespace: sillage
  qualified-name: "sillage:relational-data"
---

# Design relational data safely

Model business invariants first, then choose tables, keys, and queries. The
database is a consistency boundary, not merely a serialization detail.

## Review sequence

1. Name the owning bounded context and aggregate/row lifecycle.
2. State required, unique, foreign-key, check, and temporal invariants.
3. Choose normalization or denormalization per access pattern and update
   ownership; document the deliberate trade-off.
4. Design a forward-compatible migration: expand, backfill, switch, contract
   when needed. Define lock/volume/rollback limits before running it.
5. Add indexes from measured query shapes; check selectivity, ordering,
   pagination, and write cost rather than indexing every column.
6. Make transaction boundaries, isolation, retries, idempotency, and failure
   recovery observable in tests or a runtime check.

Return:

```text
Invariant: <rule the database must enforce>
Shape: <tables/keys/relations and ownership>
Access: <critical query and index reason>
Migration: <safe sequence, rollback, data-risk limit>
Proof: <constraint, query plan, integration, or backfill evidence>
```

Do not hide a domain decision in a migration. Route product or boundary
changes to `sillage:shape`, domain language to `sillage:ddd`, and implementation
to `sillage:build`; verify with real database behavior, not schema inspection
alone.
