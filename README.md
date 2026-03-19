# Aile - Sorry but I can't let you have the model W!

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

- `x/logger`
- `x/health`

These are convenience packages, not core concepts.

## COPYING

You may choose to use this software under either the GNU Affero General Public
License v3.0 or later (AGPL-3.0+) or the GNU Lesser General Public License v3.0
or later (LGPL-3.0+).

### When choosing the GNU AGPL-3.0

Best for open-source applications and network services.

- _Requirement:_ Any modified version of this software, when used over a
  network, must have its source code made available to users.
- [Read the AGPL-3.0 License](LICENSE-AGPL3.0-or-later.txt)

### When choosing the GNU LGPL-3.0

Best for linking in proprietary or commercial applications.

- _Requirement:_ You may link to this library, but modifications to the library
  code itself must be shared.
- [Read the LGPL-3.0 License](LICENSE-LGPL3.0-or-later.txt)
