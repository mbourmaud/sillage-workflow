.PHONY: build check contracts format lint test test-race

test:
	go test ./...

test-race:
	go test -race ./cmd/... ./internal/project ./internal/workflow

format:
	@test -z "$$(gofmt -l .)" || { gofmt -l .; exit 1; }
	@git diff --check 4b825dc642cb6eb9a060e54bf8d69288fbee4904 HEAD
	@git diff --cached --check
	@git diff --check

contracts:
	go test ./internal/contracts ./internal/skills

lint:
	go vet ./...
	go tool actionlint

check: format lint test-race contracts
	go run ./cmd/sillage doctor --root .

build:
	go build -o dist/sillage ./cmd/sillage
