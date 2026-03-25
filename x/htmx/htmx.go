package htmx

import (
	"net/http"
	"strings"
)

const (
	headerRequest            = "HX-Request"
	headerBoosted            = "HX-Boosted"
	headerTarget             = "HX-Target"
	headerTrigger            = "HX-Trigger"
	headerTriggerName        = "HX-Trigger-Name"
	headerRedirect           = "HX-Redirect"
	headerLocation           = "HX-Location"
	headerPushURL            = "HX-Push-Url"
	headerReplaceURL         = "HX-Replace-Url"
	headerRefresh            = "HX-Refresh"
	headerReswap             = "HX-Reswap"
	headerRetarget           = "HX-Retarget"
	headerTriggerResponse    = "HX-Trigger"
	headerTriggerAfterSwap   = "HX-Trigger-After-Swap"
	headerTriggerAfterSettle = "HX-Trigger-After-Settle"
)

// IsRequest reports whether the request was issued by HTMX.
func IsRequest(r *http.Request) bool {
	return headerTrue(r, headerRequest)
}

// IsBoosted reports whether the request was issued by an HTMX boosted element.
func IsBoosted(r *http.Request) bool {
	return headerTrue(r, headerBoosted)
}

// Target returns the HTMX target element id or selector.
func Target(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get(headerTarget))
}

// Trigger returns the HTMX trigger element id.
func Trigger(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get(headerTrigger))
}

// TriggerName returns the HTMX trigger name.
func TriggerName(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get(headerTriggerName))
}

// TargetIs reports whether the request targets one of the provided ids.
// Candidates may be passed with or without a leading "#".
func TargetIs(r *http.Request, ids ...string) bool {
	return matchToken(Target(r), true, ids...)
}

// TriggerIs reports whether the request was triggered by one of the provided ids.
func TriggerIs(r *http.Request, ids ...string) bool {
	return matchToken(Trigger(r), false, ids...)
}

// TriggerNameIs reports whether the request trigger name matches one of the
// provided names.
func TriggerNameIs(r *http.Request, names ...string) bool {
	return matchToken(TriggerName(r), false, names...)
}

// Redirect instructs HTMX to perform a client-side redirect.
func Redirect(w http.ResponseWriter, url string) {
	setHeader(w, headerRedirect, url)
}

// Location instructs HTMX to perform a client-side boosted request to url.
func Location(w http.ResponseWriter, url string) {
	setHeader(w, headerLocation, url)
}

// PushURL instructs HTMX to push the provided URL into browser history.
func PushURL(w http.ResponseWriter, url string) {
	setHeader(w, headerPushURL, url)
}

// ReplaceURL instructs HTMX to replace the current browser URL.
func ReplaceURL(w http.ResponseWriter, url string) {
	setHeader(w, headerReplaceURL, url)
}

// Refresh instructs HTMX to refresh the current page.
func Refresh(w http.ResponseWriter) {
	w.Header().Set(headerRefresh, "true")
}

// Reswap instructs HTMX to use the provided swap strategy.
func Reswap(w http.ResponseWriter, strategy string) {
	setHeader(w, headerReswap, strategy)
}

// Retarget instructs HTMX to retarget the response swap.
func Retarget(w http.ResponseWriter, target string) {
	setHeader(w, headerRetarget, target)
}

// SetTrigger emits one or more client-side events after the response is
// processed.
func SetTrigger(w http.ResponseWriter, events ...string) {
	setEventsHeader(w, headerTriggerResponse, events...)
}

// SetTriggerAfterSwap emits one or more client-side events after the swap step.
func SetTriggerAfterSwap(w http.ResponseWriter, events ...string) {
	setEventsHeader(w, headerTriggerAfterSwap, events...)
}

// SetTriggerAfterSettle emits one or more client-side events after settle.
func SetTriggerAfterSettle(w http.ResponseWriter, events ...string) {
	setEventsHeader(w, headerTriggerAfterSettle, events...)
}

func headerTrue(r *http.Request, name string) bool {
	if r == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(r.Header.Get(name)), "true")
}

func matchToken(got string, trimHash bool, want ...string) bool {
	got = normalizeToken(got, trimHash)
	if got == "" {
		return false
	}

	for _, candidate := range want {
		if got == normalizeToken(candidate, trimHash) {
			return true
		}
	}

	return false
}

func normalizeToken(v string, trimHash bool) string {
	v = strings.TrimSpace(v)
	if trimHash {
		v = strings.TrimPrefix(v, "#")
	}
	return v
}

func setHeader(w http.ResponseWriter, name, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	w.Header().Set(name, value)
}

func setEventsHeader(w http.ResponseWriter, name string, events ...string) {
	clean := make([]string, 0, len(events))
	for _, event := range events {
		event = strings.TrimSpace(event)
		if event == "" {
			continue
		}
		clean = append(clean, event)
	}
	if len(clean) == 0 {
		return
	}
	w.Header().Set(name, strings.Join(clean, ", "))
}
