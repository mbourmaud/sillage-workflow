# Sillage agent instructions

## Start here

1. Read `PRODUCT.md` for the product contract.
2. Read `DESIGN.md` for the interaction contract.
3. Read `docs/domain/index.md` for canonical language.
4. Read the active task record before changing files.

## Work protocol

- Investigate before proposing durable changes.
- Do not implement before the relevant human approval exists.
- Keep one implementation worktree bound to one active slice.
- Write a focused failing test before observable behavior.
- Assign each risk one primary owning test layer.
- Never claim completion from code inspection alone.
- Keep operational notes, experiments, raw research, logs, and screenshots out
  of durable knowledge directories.
- Update an existing canonical page before creating another durable document.
- Ask before destructive actions, external writes, merges, deployments, or
  evidence waivers.

## Commands

```sh
go test ./...
go vet ./...
go run ./cmd/sillage doctor --root .
```

Run the smallest relevant lane first, then the complete repository checks.

## Code

- Keep the deterministic core dependent only on the Go standard library unless
  a dependency earns clear portability or correctness value.
- Public Go identifiers have concise documentation.
- Commands validate by default and mutate only through explicit verbs.
- Keep forge-, tracker-, language-, and CI-specific behavior behind profiles or
  adapters.
