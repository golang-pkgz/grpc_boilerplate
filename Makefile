all: lint test docs

docs:
	gomarkdoc -o godoc.md ./...

lint:
	golangci-lint run

.PHONY: test
test: lint
	go test -coverprofile cover.out ./...
	go tool cover -func cover.out
	rm cover.out
