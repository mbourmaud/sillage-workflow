.PHONY: check test build

test:
	go test ./...

check:
	go test ./...
	go vet ./...
	go run ./cmd/sillage doctor --root .

build:
	go build -o dist/sillage ./cmd/sillage
