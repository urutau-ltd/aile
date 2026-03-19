// Package aile provides a small stdlib-first HTTP runtime for Go.
//
// aile keeps net/http front and center like chi does. it uses [http.ServeMux]
// [http.Server], [http.Handler] and literal stlib route patterns such as:
//
// "GET /api/v1/example"
// "POST /users/{id}"
//
// Aile mainly adds:
//   - A named build plan
//   - Patchable phases
//   - A small application container
//   - Runtime helpers for signal handling and graceful shutdown
package aile
