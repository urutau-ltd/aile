package aile

import (
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestStaticHandlerServesFiles(t *testing.T) {
	h, err := StaticHandler("/assets/", fstest.MapFS{
		"app.css": &fstest.MapFile{Data: []byte("body{color:black}")},
	})
	if err != nil {
		t.Fatalf("StaticHandler returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/assets/app.css", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}

	if rec.Body.String() != "body{color:black}" {
		t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), "body{color:black}")
	}
}

func TestStaticRegistersRedirectAndSubtreeRoute(t *testing.T) {
	app := MustNew()
	err := app.Static("/assets", fstest.MapFS{
		"app.css": &fstest.MapFile{Data: []byte("body{color:black}")},
	})
	if err != nil {
		t.Fatalf("Static returned error: %v", err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	redirectReq := httptest.NewRequest(http.MethodGet, "/assets", nil)
	redirectRec := httptest.NewRecorder()
	st.Handler.ServeHTTP(redirectRec, redirectReq)

	if redirectRec.Code != http.StatusPermanentRedirect {
		t.Fatalf("unexpected redirect status: got %d want %d", redirectRec.Code, http.StatusPermanentRedirect)
	}
	if got := redirectRec.Header().Get("Location"); got != "/assets/" {
		t.Fatalf("unexpected redirect location: got %q want %q", got, "/assets/")
	}

	fileReq := httptest.NewRequest(http.MethodGet, "/assets/app.css", nil)
	fileRec := httptest.NewRecorder()
	st.Handler.ServeHTTP(fileRec, fileReq)

	if fileRec.Code != http.StatusOK {
		t.Fatalf("unexpected file status: got %d want %d", fileRec.Code, http.StatusOK)
	}
	if fileRec.Body.String() != "body{color:black}" {
		t.Fatalf("unexpected file body: got %q want %q", fileRec.Body.String(), "body{color:black}")
	}
}

func TestStaticGetRouteAlsoServesHead(t *testing.T) {
	app := MustNew()
	err := app.Static("/assets/", fstest.MapFS{
		"app.css": &fstest.MapFile{Data: []byte("body{color:black}")},
	})
	if err != nil {
		t.Fatalf("Static returned error: %v", err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodHead, "/assets/app.css", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty HEAD body, got %q", rec.Body.String())
	}
}

func TestStaticRejectsInvalidInputs(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		fsys   fs.FS
	}{
		{name: "empty prefix", prefix: "", fsys: fstest.MapFS{}},
		{name: "missing leading slash", prefix: "assets", fsys: fstest.MapFS{}},
		{name: "nil fs", prefix: "/assets", fsys: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := StaticHandler(tc.prefix, tc.fsys); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestStaticRedirectPreservesQueryString(t *testing.T) {
	app := MustNew()
	err := app.Static("/assets", fstest.MapFS{
		"app.css": &fstest.MapFile{Data: []byte("body{color:black}")},
	})
	if err != nil {
		t.Fatalf("Static returned error: %v", err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/assets?v=42", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusPermanentRedirect {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusPermanentRedirect)
	}
	if got := rec.Header().Get("Location"); got != "/assets/?v=42" {
		t.Fatalf("unexpected redirect location: got %q want %q", got, "/assets/?v=42")
	}
}
