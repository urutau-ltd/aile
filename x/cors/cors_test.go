package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewareSetsHeadersForAllowedOrigin(t *testing.T) {
	h := Middleware(Config{
		AllowOrigins:     []string{"https://app.example.com"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Request-ID"},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("unexpected allow origin: got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("unexpected allow credentials: got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Expose-Headers"); got != "X-Request-ID" {
		t.Fatalf("unexpected expose headers: got %q", got)
	}
}

func TestMiddlewareHandlesPreflight(t *testing.T) {
	h := Middleware(Config{
		AllowOrigins: []string{"https://app.example.com"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{"Authorization"},
		MaxAge:       600,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("preflight request should not reach next handler")
	}))

	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set("Origin", "https://app.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST" {
		t.Fatalf("unexpected allow methods: got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got != "Authorization" {
		t.Fatalf("unexpected allow headers: got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Max-Age"); got != "600" {
		t.Fatalf("unexpected max age: got %q", got)
	}
}
