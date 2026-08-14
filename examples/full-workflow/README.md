# Full workflow pilot

This example is a task record at `REVIEW`. Its fields preserve the path already
completed through `DECIDE`, `IMPLEMENT`, and `VERIFY`: an accepted decision,
human approval, criterion evidence, a completed slice, an accepted review, and
a digest-bound handoff. The final gate is deliberately still explicit.

Run it from the repository root:

```sh
go run ./cmd/sillage doctor --root . --json
go run ./cmd/sillage context --root . --task examples/full-workflow/task.json --json
go run ./cmd/sillage status --task examples/full-workflow/task.json --json
go run ./cmd/sillage transition \
  --task examples/full-workflow/task.json --to HANDOFF --json
```

The read-only transition returns:

```json
{"ok":true,"code":"accepted"}
```

The example remains unchanged. To exercise the only state mutation in the
local core, copy it to a temporary task record and opt into the write:

```sh
pilot_task="$(mktemp)"
cp examples/full-workflow/task.json "$pilot_task"
go run ./cmd/sillage transition \
  --task "$pilot_task" --to HANDOFF --write --json
go run ./cmd/sillage status --task "$pilot_task" --json
rm "$pilot_task"
```

The command validates the task and transition, compares the bytes read before
the decision, writes a same-directory temporary file, and renames it into
place with private permissions. A concurrent edit, symlink target, malformed
record, or invalid gate refuses the write. Without `--write`, all commands
remain read-only.

The pilot intentionally does not merge, deploy, publish, or create a durable
note. Those actions remain human-owned gates after `HANDOFF`.
