// Package aile provides a small HTTP runtime for Go built around net/http.
//
// It uses [http.ServeMux], [http.Server], and [http.Handler]. Routes are
// registered through helpers such as [App.GET], [App.POST], [App.PUT], and
// [App.DELETE].
//
// Paths keep the same [http.ServeMux] syntax under the hood, such as:
//
// "/api/v1/example"
// "/users/{id}"
//
// Aile mainly adds:
//   - a small application container
//   - middleware wiring
//   - runtime helpers for signals and graceful shutdown
package aile
