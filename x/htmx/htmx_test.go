package htmx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHelpers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Boosted", "true")
	req.Header.Set("HX-Target", "#os-editor")
	req.Header.Set("HX-Trigger", "os-search")
	req.Header.Set("HX-Trigger-Name", "q")

	if !IsRequest(req) {
		t.Fatal("expected HTMX request")
	}
	if !IsBoosted(req) {
		t.Fatal("expected boosted request")
	}
	if got := Target(req); got != "#os-editor" {
		t.Fatalf("unexpected target: got %q want %q", got, "#os-editor")
	}
	if !TargetIs(req, "os-editor") {
		t.Fatal("expected target match without leading hash")
	}
	if !TargetIs(req, "#os-editor") {
		t.Fatal("expected target match with leading hash")
	}
	if !TriggerIs(req, "os-search") {
		t.Fatal("expected trigger match")
	}
	if !TriggerNameIs(req, "q") {
		t.Fatal("expected trigger name match")
	}
}

func TestRequestHelpersReturnFalseOnMissingValues(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)

	if IsRequest(req) {
		t.Fatal("expected non-HTMX request")
	}
	if IsBoosted(req) {
		t.Fatal("expected non-boosted request")
	}
	if TargetIs(req, "os-editor") {
		t.Fatal("expected missing target mismatch")
	}
	if TriggerIs(req, "os-search") {
		t.Fatal("expected missing trigger mismatch")
	}
	if TriggerNameIs(req, "q") {
		t.Fatal("expected missing trigger name mismatch")
	}
}

func TestResponseHeaderHelpers(t *testing.T) {
	rec := httptest.NewRecorder()

	Redirect(rec, "/os")
	Location(rec, "/os/42")
	PushURL(rec, "/os")
	ReplaceURL(rec, "/os/42/edit")
	Refresh(rec)
	Reswap(rec, "outerHTML")
	Retarget(rec, "#os-editor")
	SetTrigger(rec, "os:changed", "os:notify")
	SetTriggerAfterSwap(rec, "os:swap")
	SetTriggerAfterSettle(rec, "os:settle")

	tests := map[string]string{
		"HX-Redirect":             "/os",
		"HX-Location":             "/os/42",
		"HX-Push-Url":             "/os",
		"HX-Replace-Url":          "/os/42/edit",
		"HX-Refresh":              "true",
		"HX-Reswap":               "outerHTML",
		"HX-Retarget":             "#os-editor",
		"HX-Trigger":              "os:changed, os:notify",
		"HX-Trigger-After-Swap":   "os:swap",
		"HX-Trigger-After-Settle": "os:settle",
	}

	for header, want := range tests {
		if got := rec.Header().Get(header); got != want {
			t.Fatalf("unexpected %s: got %q want %q", header, got, want)
		}
	}
}

func TestResponseHeaderHelpersIgnoreEmptyValues(t *testing.T) {
	rec := httptest.NewRecorder()

	Redirect(rec, "")
	Location(rec, " ")
	PushURL(rec, "")
	ReplaceURL(rec, "")
	Reswap(rec, "")
	Retarget(rec, "")
	SetTrigger(rec, "", " ")

	for _, header := range []string{
		"HX-Redirect",
		"HX-Location",
		"HX-Push-Url",
		"HX-Replace-Url",
		"HX-Reswap",
		"HX-Retarget",
		"HX-Trigger",
	} {
		if got := rec.Header().Get(header); got != "" {
			t.Fatalf("expected empty %s header, got %q", header, got)
		}
	}
}
