package aile

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()

	err := WriteJSON(rec, http.StatusCreated, map[string]string{"ok": "true"})
	if err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusCreated)
	}

	gotCT := rec.Header().Get("Content-Type")
	if gotCT != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type: got %q", gotCT)
	}
}

func TestDecodeJSON(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(`{"name":"nyctibius"}`))

	got, err := DecodeJSON[payload](req)
	if err != nil {
		t.Fatalf("DecodeJSON returned error: %v", err)
	}

	if got.Name != "nyctibius" {
		t.Fatalf("unexpected payload: got %#v", got)
	}
}
