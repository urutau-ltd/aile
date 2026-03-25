# Aile - Sorry but I can't let you have the model W!

[![Go Reference](https://pkg.go.dev/badge/codeberg.org/urutau-ltd/aile/v2.svg)](https://pkg.go.dev/codeberg.org/urutau-ltd/aile/v2)

`aile` is a small HTTP runtime for Go built around the standard library. It
takes a few cues from `chi` and `hono`, but it is not trying to be a full
framework. It adds a thin layer on top of `net/http`:

- A light application container
- Middleware wiring
- Graceful shutdown and signal handling

## Installation

Install this library with:

```bash
$ go get -u codeberg.org/urutau-ltd/aile/v2
```

See the v2 migration guide in [`MIGRATING.md`](./MIGRATING.md) if you are
upgrading an existing app.

## Development Environment

`aile` is developed with GNU Guix first.

This repo uses Go 1.26.

The usual development flow is:

```bash
$ make guix-env
```

or directly:

```bash
$ guix shell --network -m ./manifest.scm
```

The default `make test`, `make vet`, `make check`, `make example-htmx`,
`make example-rest`, and `make example-html-admin` targets run through the Guix
manifest. `make guix-test`, `make guix-vet`, and `make guix-check` are there if
you want to call the Guix-backed variants explicitly.

## Onboarding

If this is your first time in the repo, use this flow:

1. Enter the Guix development environment with `make guix-env`.
2. Confirm the toolchain inside that shell with `go version` and
   `gopls version`.
3. Run `make check-local` inside that shell.
4. Use `make guix-check` from the host when you want the Makefile to open the
   Guix environment for you.
5. Start the editor with `make emacs` if you use Emacs and Eglot.
6. Run `make pkg` when you need to confirm the local Guix package definition.

Inside that shell you should see a Go 1.26.x toolchain.

## Quick Start

Here's a quick usage example:

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "codeberg.org/urutau-ltd/aile/v2"
)

func main() {
    app, err := aile.New()
    if err != nil {
        log.Fatal(err)
    }
    
    app.Use(aile.Recovery())

    app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
        aile.Text(w, http.StatusOK, "ok")
    })
    
    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Routes

`aile` does **NOT** parse or re-implement routing. It keeps `http.ServeMux`
semantics under the hood, but the public API registers routes through the HTTP
method helpers:

```go
app.GET("/api/v1/something", handler)
app.GET("/users/{id}", handler)
app.POST("/users", handler)
app.DELETE("/users/{id}", handler)
```

## Runtime

`Run` handles the common case for you:

- Builds an app
- Opens a listener
- Serves requests
- Listens for `SIGINT` and `SIGTERM`
- Performs graceful shutdown

For lower-level control:

- `Build`
- `Serve`
- `ListenAndServe`
- `Shutdown`
- `Addr`

## Hooks

You can register startup and shutdown hooks:

```go
app.OnStart(func(ctx context.Context, st *aile.State) error {
    return nil
})

app.OnShutdown(func(ctx context.Context, st *aile.State) error {
    return nil
})
```

## Values

You can store small shared values on the app:

```go
app.Set("name", "api")
v, ok := app.Value("name")
```

Built state receives a copy of those values.

## Version

The package exports `aile.Version` and `aile.ReleaseTag`. They are useful for
logs, diagnostics, and asset versioning.

## Static Files

`aile` keeps static serving very close to the standard library.

If you want app-level convenience:

```go
assets, err := fs.Sub(public, "public")
if err != nil {
    log.Fatal(err)
}

if err := app.Static("/assets", assets); err != nil {
    log.Fatal(err)
}
```

If you want the raw handler for plain `net/http` usage:

```go
h, err := aile.StaticHandler("/assets/", assets)
if err != nil {
    log.Fatal(err)
}
```

Under the hood this is still just `http.FileServerFS` plus prefix handling.

## JSON helpers

```go
payload, err := aile.DecodeJSON[MyRequest](r)

if err != nil {
    aile.Error(w, http.StatusBadRequest, "bad json")
    return
}

_ = aile.WriteJSON(w, http.StatusOK, payload)
```

## Extras under `x/`

Optional extras live under `x/`.

Current first-party extras:

### Combine Middleware

This helper allows you to combine several middlewares as if they were one:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/combine"
	requestid "codeberg.org/urutau-ltd/aile/v2/x/request_id"
	secureheaders "codeberg.org/urutau-ltd/aile/v2/x/secure_headers"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	stack := combine.Middleware(
		requestid.Middleware(requestid.Config{}),
		secureheaders.Middleware(secureheaders.Config{
			ContentTypeNosniff: true,
			FrameDeny:          true,
		}),
	)

	app.Use(stack)

	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	log.Fatal(app.Run(context.Background()))
}
```

### Cors

Easily setup your `CORS` headers:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/cors"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(cors.Middleware(cors.Config{
		AllowOrigins: []string{"https://app.example.com"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:       600,
	}))

	app.GET("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		_ = aile.WriteJSON(w, http.StatusOK, map[string]string{"message": "hello"})
	})

	log.Fatal(app.Run(context.Background()))
}
```

### Request ID

Generates or propagates request id's:

```go
 package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	requestid "codeberg.org/urutau-ltd/aile/v2/x/request_id"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(requestid.Middleware(requestid.Config{
		Header: "X-Request-ID",
	}))

	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		id, _ := requestid.FromContext(r.Context())
		_ = aile.WriteJSON(w, http.StatusOK, map[string]string{
			"request_id": id,
		})
	})

log.Fatal(app.Run(context.Background()))
}
```

### HTMX Helpers

If you are building an HTML app with HTMX and want to stop parsing `HX-*`
headers by hand:

```go
package main

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/htmx"
)

var pageTmpl = template.Must(template.New("editor").Parse(`<div id="os-editor">{{.}}</div>`))

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.GET("/os/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		if htmx.IsRequest(r) && htmx.TargetIs(r, "os-editor") {
			_ = pageTmpl.Execute(w, "Partial editor for "+r.PathValue("id"))
			return
		}

		_ = pageTmpl.Execute(w, "Full page editor for "+r.PathValue("id"))
	})

	app.POST("/os/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		htmx.SetTrigger(w, "os:changed")
		htmx.Redirect(w, "/os")
		aile.Status(w, http.StatusNoContent)
	})

	log.Fatal(app.Run(context.Background()))
}
```

### Secure Headers

Various security headers:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	secureheaders "codeberg.org/urutau-ltd/aile/v2/x/secure_headers"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(secureheaders.Middleware(secureheaders.Config{
		ContentTypeNosniff:    true,
		FrameDeny:             true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' 'unsafe-inline'",
		HSTSMaxAge:            31536000,
		HSTSIncludeSubdomains: true,
	}))

	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "secure")
	})

	log.Fatal(app.Run(context.Background()))
}
```

### Trailing Slash

Redirect adding or removing the trailing slash of a route:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	trailingslash "codeberg.org/urutau-ltd/aile/v2/x/trailing_slash"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(trailingslash.Middleware(trailingslash.RedirectTrim))

	app.GET("/users", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "users")
	})

	log.Fatal(app.Run(context.Background()))
}
```

### Bearer Token validator

Simple `Bearer` token authentication.

```go
package main

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	bearerauth "codeberg.org/urutau-ltd/aile/v2/x/bearer_auth"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(bearerauth.Middleware(func(token string) bool {
		return subtle.ConstantTimeCompare([]byte(token), []byte("super-token")) == 1
	}))

	app.GET("/private", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "private data")
	})

	log.Fatal(app.Run(context.Background()))}
```

### Basic Authentication

Classic basig auth everyone knows and loves:

```go
package main

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	basicauth "codeberg.org/urutau-ltd/aile/v2/x/basic_auth"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(basicauth.Middleware("admin", func(user, pass string) bool {
		userOK := subtle.ConstantTimeCompare([]byte(user), []byte("admin")) == 1
		passOK := subtle.ConstantTimeCompare([]byte(pass), []byte("secret")) == 1
		return userOK && passOK
	}))

	app.GET("/admin", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "welcome")
	})

	log.Fatal(app.Run(context.Background()))
}
```

### GZIP Compression

> [!IMPORTANT]
> For this to work, your clients should send the `Accept-Encoding: gzip` header
> along their requests to get compressed responses. Otherwise the compression
> will be skipped and the response will be served at normal size.

Compress your responses with `gzip`:

```go
package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/compress"
)

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(compress.Middleware(compress.Config{
		MinSize: 256,
	}))

	app.GET("/big", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, strings.Repeat("hello ", 200))
	})

	log.Fatal(app.Run(context.Background()))
}
```

### IP (v4) restriction

Allow or Deny by IP/Network:

```go
package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	iprestriction "codeberg.org/urutau-ltd/aile/v2/x/ip_restriction"
)

func mustCIDR(s string) *net.IPNet {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Use(iprestriction.Middleware(iprestriction.Config{
		Allow: []*net.IPNet{
			mustCIDR("127.0.0.0/8"),
			mustCIDR("10.0.0.0/8"),
		},
		TrustProxy: false,
	}))

	app.GET("/internal", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "internal")
	})

	log.Fatal(app.Run(context.Background()))
}
```

These are convenience packages, not core concepts.

### Resource Mounting

If your `main.go` is repeating the same CRUD-style route blocks over and over,
you can mount those conventions once:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/resource"
)

type providersHandler struct{}

func (providersHandler) Index(w http.ResponseWriter, r *http.Request)  { aile.Text(w, http.StatusOK, "providers index") }
func (providersHandler) New(w http.ResponseWriter, r *http.Request)    { aile.Text(w, http.StatusOK, "providers new") }
func (providersHandler) Create(w http.ResponseWriter, r *http.Request) { aile.Status(w, http.StatusCreated) }
func (providersHandler) Show(w http.ResponseWriter, r *http.Request)   { aile.Text(w, http.StatusOK, "provider "+r.PathValue("id")) }
func (providersHandler) Edit(w http.ResponseWriter, r *http.Request)   { aile.Text(w, http.StatusOK, "provider edit "+r.PathValue("id")) }
func (providersHandler) Update(w http.ResponseWriter, r *http.Request) { aile.Text(w, http.StatusOK, "provider update "+r.PathValue("id")) }
func (providersHandler) Delete(w http.ResponseWriter, r *http.Request) { aile.Status(w, http.StatusNoContent) }

type appSettingsHandler struct{}

func (appSettingsHandler) Show(w http.ResponseWriter, r *http.Request)   { aile.Text(w, http.StatusOK, "settings show") }
func (appSettingsHandler) Edit(w http.ResponseWriter, r *http.Request)   { aile.Text(w, http.StatusOK, "settings edit") }
func (appSettingsHandler) Update(w http.ResponseWriter, r *http.Request) { aile.Text(w, http.StatusOK, "settings update") }

func main() {
	app, err := aile.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := resource.MountCollection(app, "/providers", providersHandler{}); err != nil {
		log.Fatal(err)
	}

	if err := resource.MountSingleton(app, "/app-settings", appSettingsHandler{}); err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Run(context.Background()))
}
```

## Examples

You can see more examples at the [`examples/`](./examples/) directory.

### CRUD-style HTML app wiring

For an app like:

- `/providers`
- `/locations`
- `/os`
- `/app-settings`
- `/account-settings`

your `main.go` can stay flat:

```go
providersHandler := providers.NewHandler(logger, uiRoot, dbConn)
locationsHandler := locations.NewHandler(logger, uiRoot, dbConn)
osHandler := operatingsystems.NewHandler(logger, uiRoot, dbConn)
appSettingsHandler := appsettings.NewHandler(logger, uiRoot, dbConn)

if err := app.Static("/static", staticRoot); err != nil {
	log.Fatal(err)
}

if err := resource.MountCollection(app, "/providers", providersHandler); err != nil {
	log.Fatal(err)
}
if err := resource.MountCollection(app, "/locations", locationsHandler); err != nil {
	log.Fatal(err)
}
if err := resource.MountCollection(app, "/os", osHandler); err != nil {
	log.Fatal(err)
}
if err := resource.MountSingleton(app, "/app-settings", appSettingsHandler); err != nil {
	log.Fatal(err)
}
```

### Small Blog API example

Here's how to make a blog API with `aile` quickly:

```go
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/combine"
	"codeberg.org/urutau-ltd/aile/v2/x/cors"
	requestid "codeberg.org/urutau-ltd/aile/v2/x/request_id"
	secureheaders "codeberg.org/urutau-ltd/aile/v2/x/secure_headers"
	xlogger "codeberg.org/urutau-ltd/aile/v2/x/logger"
)

type article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	app, err := aile.New(aile.WithAddr(":8080"))
	if err != nil {
		log.Fatal(err)
	}

	// Example in-memory data.
	articles := []article{
		{ID: "1", Title: "Aile v1"},
		{ID: "2", Title: "Small APIs with Go"},
	}

	// Compose a practical middleware stack.
	app.Use(combine.Middleware(
		aile.Recovery(),
		requestid.Middleware(requestid.Config{
			Header: "X-Request-ID",
		}),
		xlogger.Middleware(slog.Default()),
		secureheaders.Middleware(secureheaders.Config{
			ContentTypeNosniff: true,
			FrameDeny:          true,
			ReferrerPolicy:     "strict-origin-when-cross-origin",
		}),
		cors.Middleware(cors.Config{
			AllowOrigins: []string{"http://localhost:3000"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders: []string{"Content-Type", "Authorization", "X-Request-ID"},
			MaxAge:       600,
		}),
	))

	app.GET("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = aile.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	app.GET("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		_ = aile.WriteJSON(w, http.StatusOK, articles)
	})

	app.POST("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		var in article
		in, err = aile.DecodeJSON[article](r)
		if err != nil {
			aile.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		if in.ID == "" || in.Title == "" {
			aile.Error(w, http.StatusBadRequest, "id and title are required")
			return
		}

		articles = append(articles, in)

		_ = aile.WriteJSON(w, http.StatusCreated, in)
	})

	app.GET("/api/meta", func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := requestid.FromContext(r.Context())

		_ = aile.WriteJSON(w, http.StatusOK, map[string]string{
			"service":    "articles",
			"request_id": reqID,
		})
	})

	log.Fatal(app.Run(context.Background()))
}
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Development

The repository lives on Codeberg and is mirrored to GitHub.

Most development happens through:

- Guix development shells via `make guix-env`
- Guix-backed verification via `make check`
- explicit Guix aliases via `make guix-test`, `make guix-vet`, and
  `make guix-check`
- local containerized verification via Podman and `podman-compose` as a
  secondary path
- Emacs as the primary IDE, with Eglot using `gopls` for Go files

Project-local Emacs settings live in [`.dir-locals.el`](./.dir-locals.el).
`make emacs` opens Emacs inside the Guix development environment.

If you want to bypass Guix in the current shell, use the explicit local targets:

```bash
$ make check-local
$ make test-local
$ make vet-local
```

### Local Podman Pipeline

Run the local containerized pipeline with:

```bash
$ make podman-check
```

Open an interactive shell in the same development image with:

```bash
$ make podman-shell
```

If `podman` lives outside your default `PATH`, override it explicitly:

```bash
$ make podman-check PODMAN=/absolute/path/to/podman
```

## GNU Guix Compatibility

This project is developed with GNU Guix. The manifest tracks the Go 1.26 series.
These helpers cover the usual Guix tasks:

- `make guix-env` opens a shell with the development tools used by the repo,
  including Emacs, `gopls`, Podman, and `podman-compose`.

- `make check` and `make guix-check` run the default verification flow through
  the Guix manifest, using Go 1.26.

- `make pkg` checks that the `guix.scm` package definition for the current
  checkout resolves and builds under GNU Guix.

## COPYING

Where applicable this project's source code is under the terms of the GNU Affero
General Public License version 3 or at your option, any later version.
