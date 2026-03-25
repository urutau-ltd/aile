GUIX ?= guix
GUIX_MANIFEST ?= ./manifest.scm
GUIX_PACKAGE_FILE ?= ./guix.scm
GUIX_SHELL = $(GUIX) shell --network -m $(GUIX_MANIFEST) --
PODMAN ?= podman
PODMAN_COMPOSE ?= podman-compose

GO_ENV = CGO_ENABLED=0

.PHONY: all test test-local vet vet-local check check-local env guix-env emacs \
	guix-test guix-vet guix-check podman-build podman-check podman-shell pkg ci \
	example-htmx example-htmx-local example-rest example-rest-local

all: ci

test-local:
	$(GO_ENV) go test -v ./...

vet-local:
	$(GO_ENV) go vet ./...

check-local: test-local vet-local

guix-test:
	$(GUIX_SHELL) $(MAKE) test-local

guix-vet:
	$(GUIX_SHELL) $(MAKE) vet-local

guix-check:
	$(GUIX_SHELL) $(MAKE) check-local

test: guix-test

vet: guix-vet

check: guix-check

podman-build:
	$(PODMAN_COMPOSE) --podman-path $(PODMAN) build ci

podman-check:
	$(PODMAN_COMPOSE) --podman-path $(PODMAN) run --rm ci

podman-shell:
	$(PODMAN_COMPOSE) --podman-path $(PODMAN) run --rm shell

env: guix-env

guix-env:
	$(GUIX) shell --network -m $(GUIX_MANIFEST)

emacs:
	$(GUIX_SHELL) emacs

pkg:
	$(GUIX) build -f $(GUIX_PACKAGE_FILE)

ci:
	$(MAKE) check

example-htmx-local:
	$(GO_ENV) go run ./examples/htmx-counter/main.go

example-rest-local:
	$(GO_ENV) go run ./examples/rest-api/main.go

example-htmx:
	$(GUIX_SHELL) $(MAKE) example-htmx-local

example-rest:
	$(GUIX_SHELL) $(MAKE) example-rest-local
