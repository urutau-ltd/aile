# Contributing

This project is developed on Codeberg and mirrored to GitHub.

Day-to-day development assumes the Guix manifest, not whatever toolchain happens
to be installed on the host. This repo uses Go 1.26.

The project does not depend on hosted CI. Use one of these local workflows:

- Guix-based development shells and verification for most work
- Native local development when you need it
- Local containerized pipelines with Podman and Podman Compose
- Emacs as the primary IDE, with Eglot driving `gopls` for Go buffers

## Standards

- Keep the library close to `net/http`
- Prefer small APIs over framework-style abstractions
- Avoid new runtime dependencies
- Keep middleware and helpers composable and stdlib-shaped
- Add tests for behavior changes and bug fixes
- Update GoDoc and README examples when the public API changes

## Primary workflow: Guix

The usual flow is:

```bash
make guix-env
make check
```

The default `make` targets run through the Guix manifest:

- `make test`
- `make vet`
- `make check`
- `make guix-test`
- `make guix-vet`
- `make guix-check`
- `make example-htmx`
- `make example-rest`
- `make example-html-admin`

## Onboarding checklist

Use this checklist when preparing a machine or a fresh checkout:

1. Run `make guix-env`.
2. Inside that shell, confirm `go version` reports a Go 1.26 toolchain.
3. Confirm `gopls version` is available for Emacs and Eglot.
4. Run `make check-local` inside that shell.
5. Use `make guix-check` from the host when you want the Makefile to create the
   Guix environment for you.
6. Run `make pkg` if you are touching `guix.scm` or release packaging.
7. Run `make emacs` if Emacs is your active editor session for the project.

## Secondary workflow: native shell

If you want to use the current host environment directly, use the explicit
`-local` targets:

```bash
make check-local
```

That runs:

- `go test ./...`
- `go vet ./...`

## Guix workflow

The repository ships a development manifest and a package definition.

Start a development shell with:

```bash
make guix-env
```

That shell should provide:

- Go
- `gopls`
- Emacs
- Podman and `podman-compose`
- `make`
- `git`
- `ripgrep`

Build the local checkout through Guix with:

```bash
make pkg
```

The current `guix.scm` is set up for local development from the repository
checkout. If you later want to publish a Guix package from a release tag,
replace the local source with a tagged `git-fetch` origin and update the hash.

## Emacs and Eglot

Project-local Emacs settings live in [`.dir-locals.el`](./.dir-locals.el).

The editor setup is:

- Emacs as the main editor
- Eglot as the LSP client
- `gopls` as the Go language server for `go-mode`, `go-ts-mode`, `go-mod-mode`
  and `go-mod-ts-mode`
- project compile commands that go through the Guix-backed Makefile targets

The local configuration also sets:

- project compile commands
- formatting conventions
- buffer-local `gopls` settings for static analysis and completions
- Guix-friendly defaults for Scheme files

## Local pipelines with Podman

Codeberg does not give us hosted Actions we can rely on, so the repository
includes a local container pipeline.

Build and run the local CI job with:

```bash
make podman-check
```

Open an interactive development shell in the same container image with:

```bash
make podman-shell
```

Those commands use `podman-compose` and the repository's
[`compose.yml`](./compose.yml). If your `podman` binary is not on the default
`PATH`, pass `PODMAN=/absolute/path/to/podman`.

## Documentation expectations

When public behavior changes, update as needed:

- GoDoc comments
- GoDoc examples in `*_test.go`
- [`README.md`](./README.md)
- [`MIGRATING.md`](./MIGRATING.md) for breaking changes

## Release Checklist

Before freezing or tagging a release:

1. Run `make check`.
2. Run `make podman-check`.
3. Ensure examples compile and tests pass.
4. Ensure README, migration notes and version constants reflect the release.
