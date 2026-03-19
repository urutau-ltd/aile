package aile

import "net/http"

// State is the built runtime state for a given [aile/App]
type State struct {
	Config  Config
	Mux     *http.ServeMux
	Handler http.Handler
	Server  *http.Server
	Values  map[string]any
}
