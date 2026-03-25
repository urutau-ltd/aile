package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewareCompressesResponsesWhenAccepted(t *testing.T) {
	h := Middleware(Config{MinSize: 1})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "hello compressed world")
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("unexpected content encoding: got %q want %q", got, "gzip")
	}

	zr, err := gzip.NewReader(bytes.NewReader(rec.Body.Bytes()))
	if err != nil {
		t.Fatalf("gzip.NewReader returned error: %v", err)
	}
	defer zr.Close()

	body, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("io.ReadAll returned error: %v", err)
	}

	if string(body) != "hello compressed world" {
		t.Fatalf("unexpected body: got %q want %q", string(body), "hello compressed world")
	}
}

func TestMiddlewareSkipsCompressionWithoutAcceptEncoding(t *testing.T) {
	h := Middleware(Config{MinSize: 1})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "plain response")
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "" {
		t.Fatalf("unexpected content encoding: got %q want empty", got)
	}

	if rec.Body.String() != "plain response" {
		t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), "plain response")
	}
}

func TestMiddlewareLeavesSmallResponsesPlain(t *testing.T) {
	h := Middleware(Config{MinSize: 64})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "small body")
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Content-Encoding"); got != "" {
		t.Fatalf("unexpected content encoding: got %q want empty", got)
	}

	if rec.Body.String() != "small body" {
		t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), "small body")
	}
}
