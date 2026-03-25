package resource_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/resource"
)

type providersHandler struct{}

func (providersHandler) Index(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "providers index")
}
func (providersHandler) New(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "providers new")
}
func (providersHandler) Create(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusCreated, "providers create")
}
func (providersHandler) Show(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "provider "+r.PathValue("id"))
}
func (providersHandler) Edit(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "provider edit "+r.PathValue("id"))
}
func (providersHandler) Update(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "provider update "+r.PathValue("id"))
}
func (providersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	aile.Status(w, http.StatusNoContent)
}

type appSettingsHandler struct{}

func (appSettingsHandler) Show(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "settings show")
}
func (appSettingsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "settings edit")
}
func (appSettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	aile.Text(w, http.StatusOK, "settings update")
}

func ExampleMountCollection() {
	app := aile.MustNew()
	if err := resource.MountCollection(app, "/providers", providersHandler{}); err != nil {
		panic(err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/providers/42/edit", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 provider edit 42
}

func ExampleMountSingleton() {
	app := aile.MustNew()
	if err := resource.MountSingleton(app, "/app-settings", appSettingsHandler{}); err != nil {
		panic(err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/app-settings/edit", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 settings update
}
