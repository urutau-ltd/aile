package trailingslash

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareRedirectsTrim(t *testing.T) {
	h := Middleware(RedirectTrim)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("trim redirect should not reach next handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusPermanentRedirect {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusPermanentRedirect)
	}
	if got := rec.Header().Get("Location"); got != "/users" {
		t.Fatalf("unexpected location: got %q want %q", got, "/users")
	}
}

func TestMiddlewareRedirectsAppend(t *testing.T) {
	h := Middleware(RedirectAppend)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("append redirect should not reach next handler")
	}))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusPermanentRedirect {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusPermanentRedirect)
	}
	if got := rec.Header().Get("Location"); got != "/users/" {
		t.Fatalf("unexpected location: got %q want %q", got, "/users/")
	}
}
