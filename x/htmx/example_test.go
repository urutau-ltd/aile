package htmx_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"codeberg.org/urutau-ltd/aile/v2/x/htmx"
)

func ExampleTargetIs() {
	req := httptest.NewRequest(http.MethodGet, "/os/42/edit", nil)
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Target", "#os-editor")

	fmt.Println(htmx.IsRequest(req))
	fmt.Println(htmx.TargetIs(req, "os-editor"))
	// Output:
	// true
	// true
}

func ExampleSetTrigger() {
	rec := httptest.NewRecorder()
	htmx.SetTrigger(rec, "os:changed")
	htmx.Redirect(rec, "/os")

	fmt.Println(rec.Header().Get("HX-Trigger"))
	fmt.Println(rec.Header().Get("HX-Redirect"))
	// Output:
	// os:changed
	// /os
}
