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

- `x/logger`
- `x/health`

These are convenience packages, not core concepts.

## COPYING

Where applicable this project's source code is under the terms of the GNU
Affero General Public License version 3 or at your option, any later version.
