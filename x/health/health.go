package health

import (
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Mount registers a basic GET /healthz endpoint on the app.
func Mount(a *aile.App) {
	a.GET("/healthz", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
}
