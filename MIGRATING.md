# Migrating from 1.x to 2.x

`aile` v2 follows the same idea as v1: stay close to `net/http`, add a small
application container and runtime helpers, and stop there.

The breaking changes are few.

## 1. Import path

The module path is now:

```go
import "codeberg.org/urutau-ltd/aile/v2"
```

Update every first-party extra import the same way:

```go
import requestid "codeberg.org/urutau-ltd/aile/v2/x/request_id"
import secureheaders "codeberg.org/urutau-ltd/aile/v2/x/secure_headers"
```

## 2. Manual route patterns were removed from the public API

The public `App.Handle` and `App.HandleFunc` API is gone.

Use the method helpers instead:

```go
// v1.x
app.HandleFunc("GET /users/{id}", getUser)
app.HandleFunc("POST /users", createUser)

// v2.x
app.GET("/users/{id}", getUser)
app.POST("/users", createUser)
```

This keeps route registration explicit while still using `http.ServeMux`
semantics under the hood.

## 3. Static file convenience

v2 also adds a small helper for static files:

```go
if err := app.Static("/assets", assetsFS); err != nil {
	log.Fatal(err)
}
```

If you want the raw handler instead of app registration, use:

```go
h, err := aile.StaticHandler("/assets/", assetsFS)
```

Both utilities are small wrappers around `http.FileServerFS` and
`http.StripPrefix`.

## 4. Version constant

The package now exports:

```go
aile.Version
aile.ReleaseTag
```

These are handy for diagnostics, asset versioning, and templates.

## 5. Runtime behavior fixes

No migration is required here, but the following behaviors were corrected:

- `Serve` and `Run` now return shutdown hook errors correctly.
- `Serve` closes the listener if startup hooks fail.
- `x/combine` now actually composes middleware.
- `x/compress` now actually wraps the response writer it creates.

## 6. Upgrade checklist

1. Change imports to `/v2`.
2. Replace every `Handle` / `HandleFunc` call with the appropriate method
   helper.
3. Run `go test ./...`.
4. If you serve assets, consider moving those routes to `app.Static(...)`.
