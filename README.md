# Aile - Sorry but I can't let you have the model W!

[![Go Reference](https://pkg.go.dev/badge/codeberg.org/urutau-ltd/aile.svg)](https://pkg.go.dev/codeberg.org/urutau-ltd/aile)

`aile` is a small stdlib-first HTTP runtime for Go. It's design is inspired in
`chi` and `hono`, and it doesn't strive to be a framework, it's more like a tiny
library that just adds a small amount of structure on top of the standard
library API:

- A light application container
- Middleware Wiring
- Graceful shutdown and signal handling

## Installation

Install this library with:

```bash
$ go get -u codeberg.org/urutau-ltd/aile
```

## Quick Start

Here's a quick usage example:

```go
package main

import (
    "context"
    "log"
    "net/http"
    
    "codeberg.org/urutau-ltd/aile"
)

func main() {
    app, err := aile.New()
    if err != nil {
        log.Fatal(err)
    }
    
    app.Use(aile.Recovery())
    
    app.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
        aile.Text(w, http.StatusOk, "ok")
    })
    
    // Or use the convenience helpers
    app.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
        aile.Text(w, http.StatusOK, "hello")
    })
    
    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Routes

`aile` does **NOT** parse or re-implement routing. It passes the literal patters
you should use with `http.ServeMux`:

```go
app.HandleFunc("GET /api/v1/something", handler)
app.HandleFunc("POST /users/{id}", handler)
app.HandleFunc("/plain/pattern", handler)
```

There are also convenience helpers:

```go
app.GET("/users/{id}", handler)
app.POST("/users", handler)
app.DELETE("/users/{id}", handler)
```

## Runtime

`Run` does the usual work for you:

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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/combine"
	"codeberg.org/urutau-ltd/aile/x/requestid"
	"codeberg.org/urutau-ltd/aile/x/secureheaders"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/cors"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/requestid"
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

### Secure Headers

Various security headers:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/secureheaders"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/trailingslash"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/bearerauth"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/basicauth"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/compress"
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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/iprestriction"
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

## Examples

You can see more examples at the [`examples/`](./examples/) directory.

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

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/combine"
	"codeberg.org/urutau-ltd/aile/x/cors"
	"codeberg.org/urutau-ltd/aile/x/requestid"
	"codeberg.org/urutau-ltd/aile/x/secureheaders"
	xlogger "codeberg.org/urutau-ltd/aile/x/logger"
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

## GNU Guix compatiblity

This project is being developed in GNU Guix. And thus the supported Go version
is the latest available on GNU Guix. You may use some utilities present in this
repository if you wish to contribute using guix:

- `make env` should create a new guix shell environment with the necessary
  development dependencies.

- `make pkg` should test if the `guix.scm` package definition of this library
  builds and installs correctly under GNU Guix.

## COPYING

Where applicable this project's source code is under the terms of the GNU Affero
General Public License version 3 or at your option, any later version.
