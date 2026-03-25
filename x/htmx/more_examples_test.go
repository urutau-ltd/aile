package htmx_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"codeberg.org/urutau-ltd/aile/v2/x/htmx"
)

func ExampleIsRequest() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Request", "true")

	fmt.Println(htmx.IsRequest(req))
	// Output: true
}

func ExampleIsBoosted() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Boosted", "true")

	fmt.Println(htmx.IsBoosted(req))
	// Output: true
}

func ExampleTarget() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Target", "#os-editor")

	fmt.Println(htmx.Target(req))
	// Output: #os-editor
}

func ExampleTrigger() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Trigger", "os-search")

	fmt.Println(htmx.Trigger(req))
	// Output: os-search
}

func ExampleTriggerName() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Trigger-Name", "q")

	fmt.Println(htmx.TriggerName(req))
	// Output: q
}

func ExampleTriggerIs() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Trigger", "os-search")

	fmt.Println(htmx.TriggerIs(req, "os-search"))
	// Output: true
}

func ExampleTriggerNameIs() {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Trigger-Name", "limit")

	fmt.Println(htmx.TriggerNameIs(req, "limit"))
	// Output: true
}

func ExampleRedirect() {
	rec := httptest.NewRecorder()
	htmx.Redirect(rec, "/os")

	fmt.Println(rec.Header().Get("HX-Redirect"))
	// Output: /os
}

func ExampleLocation() {
	rec := httptest.NewRecorder()
	htmx.Location(rec, "/os/42")

	fmt.Println(rec.Header().Get("HX-Location"))
	// Output: /os/42
}

func ExamplePushURL() {
	rec := httptest.NewRecorder()
	htmx.PushURL(rec, "/os")

	fmt.Println(rec.Header().Get("HX-Push-Url"))
	// Output: /os
}

func ExampleReplaceURL() {
	rec := httptest.NewRecorder()
	htmx.ReplaceURL(rec, "/os/42/edit")

	fmt.Println(rec.Header().Get("HX-Replace-Url"))
	// Output: /os/42/edit
}

func ExampleRefresh() {
	rec := httptest.NewRecorder()
	htmx.Refresh(rec)

	fmt.Println(rec.Header().Get("HX-Refresh"))
	// Output: true
}

func ExampleReswap() {
	rec := httptest.NewRecorder()
	htmx.Reswap(rec, "outerHTML")

	fmt.Println(rec.Header().Get("HX-Reswap"))
	// Output: outerHTML
}

func ExampleRetarget() {
	rec := httptest.NewRecorder()
	htmx.Retarget(rec, "#os-editor")

	fmt.Println(rec.Header().Get("HX-Retarget"))
	// Output: #os-editor
}

func ExampleSetTriggerAfterSwap() {
	rec := httptest.NewRecorder()
	htmx.SetTriggerAfterSwap(rec, "os:swap")

	fmt.Println(rec.Header().Get("HX-Trigger-After-Swap"))
	// Output: os:swap
}

func ExampleSetTriggerAfterSettle() {
	rec := httptest.NewRecorder()
	htmx.SetTriggerAfterSettle(rec, "os:settle")

	fmt.Println(rec.Header().Get("HX-Trigger-After-Settle"))
	// Output: os:settle
}
