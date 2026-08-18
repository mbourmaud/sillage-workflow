.PHONY: build check changelog-check contracts format lint pilot release-notes site-check test test-race

test:
	go test ./...

test-race:
	go test -race ./cmd/... ./internal/project ./internal/taskstore ./internal/workflow

format:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	@git diff --check 4b825dc642cb6eb9a060e54bf8d69288fbee4904 HEAD
	@git diff --cached --check
	@git diff --check

contracts:
	go test ./internal/contracts ./internal/release ./internal/skills

changelog-check:
	go run ./cmd/sillage changelog check --file CHANGELOG.md

site-check:
	go test ./internal/site

pilot:
	@set -eu; \
	 task_file=$$(mktemp "$${TMPDIR:-/tmp}/sillage-pilot.XXXXXX"); \
	 trap 'rm -f "$$task_file"' EXIT; \
	 cp examples/full-workflow/task.json "$$task_file"; \
	 go run ./cmd/sillage doctor --root . --json; \
	 go run ./cmd/sillage context --root . --task examples/full-workflow/task.json --json; \
	 go run ./cmd/sillage conformance --task examples/conformance/task.json --json; \
	 go run ./cmd/sillage status --task "$$task_file" --json; \
	 go run ./cmd/sillage transition --task "$$task_file" --to HANDOFF --write --json; \
	 go run ./cmd/sillage status --task "$$task_file" --json

release-notes:
	@test -n "$(VERSION)" || { echo "usage: make release-notes VERSION=vX.Y.Z" >&2; exit 2; }
	go run ./cmd/sillage changelog extract --file CHANGELOG.md --version "$(VERSION)"

lint:
	go vet ./...
	go tool actionlint

check: format lint test-race contracts changelog-check site-check
	go run ./cmd/sillage doctor --root .

build:
	go build -o dist/sillage ./cmd/sillage
