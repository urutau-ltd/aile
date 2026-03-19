
all: ci

test:
	CGO_ENABLED=0 go test -v ./...

env:
	guix shell --network -m ./manifest.scm

pkg:
	guix build -f ./guix.scm

ci:
	$(MAKE) test

example-htmx:
	CGO_ENABLED=0 go run ./examples/htmx-counter/main.go

example-rest:
	CGO_ENABLED=0 go run ./examples/rest-api/main.go
