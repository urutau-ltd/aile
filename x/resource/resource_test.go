package resource

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	aile "codeberg.org/urutau-ltd/aile/v2"
)

type testCollectionHandler struct{}

func (testCollectionHandler) Index(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "index")
}
func (testCollectionHandler) New(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "new")
}
func (testCollectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusCreated, "create")
}
func (testCollectionHandler) Show(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "show:"+r.PathValue("id"))
}
func (testCollectionHandler) Edit(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "edit:"+r.PathValue("id"))
}
func (testCollectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "update:"+r.PathValue("id"))
}
func (testCollectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "delete:"+r.PathValue("id"))
}

type testSingletonHandler struct{}

func (testSingletonHandler) Show(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "show")
}
func (testSingletonHandler) Edit(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "edit")
}
func (testSingletonHandler) Update(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "update")
}

func TestMountCollectionRegistersRoutes(t *testing.T) {
	app := aile.MustNew()
	err := MountCollection(app, "/providers/", testCollectionHandler{})
	if err != nil {
		t.Fatalf("MountCollection returned error: %v", err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	tests := []struct {
		method string
		path   string
		status int
		body   string
	}{
		{method: http.MethodGet, path: "/providers", status: http.StatusOK, body: "index"},
		{method: http.MethodGet, path: "/providers/new", status: http.StatusOK, body: "new"},
		{method: http.MethodPost, path: "/providers", status: http.StatusCreated, body: "create"},
		{method: http.MethodGet, path: "/providers/7", status: http.StatusOK, body: "show:7"},
		{method: http.MethodGet, path: "/providers/7/edit", status: http.StatusOK, body: "edit:7"},
		{method: http.MethodPost, path: "/providers/7/edit", status: http.StatusOK, body: "update:7"},
		{method: http.MethodDelete, path: "/providers/7", status: http.StatusOK, body: "delete:7"},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			st.Handler.ServeHTTP(rec, req)

			if rec.Code != tc.status {
				t.Fatalf("unexpected status: got %d want %d", rec.Code, tc.status)
			}
			if rec.Body.String() != tc.body {
				t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), tc.body)
			}
		})
	}
}

func TestMountSingletonRegistersRoutes(t *testing.T) {
	app := aile.MustNew()
	err := MountSingleton(app, "/app-settings", testSingletonHandler{})
	if err != nil {
		t.Fatalf("MountSingleton returned error: %v", err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	tests := []struct {
		method string
		path   string
		body   string
	}{
		{method: http.MethodGet, path: "/app-settings", body: "show"},
		{method: http.MethodGet, path: "/app-settings/edit", body: "edit"},
		{method: http.MethodPost, path: "/app-settings/edit", body: "update"},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			st.Handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
			}
			if rec.Body.String() != tc.body {
				t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), tc.body)
			}
		})
	}
}

func TestMountCollectionRejectsInvalidInputs(t *testing.T) {
	app := aile.MustNew()
	var nilHandler *testCollectionHandler

	tests := []struct {
		name string
		app  *aile.App
		path string
		h    Collection
	}{
		{name: "nil app", app: nil, path: "/providers", h: testCollectionHandler{}},
		{name: "empty path", app: app, path: "", h: testCollectionHandler{}},
		{name: "missing leading slash", app: app, path: "providers", h: testCollectionHandler{}},
		{name: "root path", app: app, path: "/", h: testCollectionHandler{}},
		{name: "pattern path", app: app, path: "/providers/{id}", h: testCollectionHandler{}},
		{name: "nil handler", app: app, path: "/providers", h: nilHandler},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := MountCollection(tc.app, tc.path, tc.h); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestMountSingletonRejectsInvalidInputs(t *testing.T) {
	app := aile.MustNew()
	var nilHandler *testSingletonHandler

	tests := []struct {
		name string
		app  *aile.App
		path string
		h    Singleton
	}{
		{name: "nil app", app: nil, path: "/app-settings", h: testSingletonHandler{}},
		{name: "empty path", app: app, path: "", h: testSingletonHandler{}},
		{name: "missing leading slash", app: app, path: "app-settings", h: testSingletonHandler{}},
		{name: "root path", app: app, path: "/", h: testSingletonHandler{}},
		{name: "pattern path", app: app, path: "/app-settings/{id}", h: testSingletonHandler{}},
		{name: "nil handler", app: app, path: "/app-settings", h: nilHandler},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := MountSingleton(tc.app, tc.path, tc.h); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
