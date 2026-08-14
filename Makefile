.PHONY: build check changelog-check contracts format lint release-notes site-check test test-race

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
