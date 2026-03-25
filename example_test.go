package aile_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing/fstest"

	aile "codeberg.org/urutau-ltd/aile/v2"
)

func ExampleApp_GET() {
	app := aile.MustNew()
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 ok
}

func ExampleApp_Static() {
	app := aile.MustNew()
	err := app.Static("/assets", fstest.MapFS{
		"app.css": &fstest.MapFile{Data: []byte("body{color:black}")},
	})
	if err != nil {
		panic(err)
	}

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/assets/app.css", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 body{color:black}
}

func ExampleWriteJSON() {
	type response struct {
		Status string `json:"status"`
	}

	rec := httptest.NewRecorder()
	if err := aile.WriteJSON(rec, http.StatusCreated, response{Status: "created"}); err != nil {
		panic(err)
	}

	fmt.Println(rec.Code)
	fmt.Print(rec.Body.String())
	// Output:
	// 201
	// {"status":"created"}
}

func ExampleDecodeJSON() {
	type request struct {
		Title string `json:"title"`
	}

	req := httptest.NewRequest(http.MethodPost, "/articles", strings.NewReader(`{"title":"Aile v2"}`))
	payload, err := aile.DecodeJSON[request](req)
	if err != nil {
		panic(err)
	}

	fmt.Println(payload.Title)
	// Output: Aile v2
}

func ExampleRecovery() {
	h := aile.Recovery()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 500
}

func ExampleVersion() {
	fmt.Println(aile.Version)
	fmt.Println(aile.ReleaseTag)
	// Output:
	// 2.1.0
	// v2.1.0
}
