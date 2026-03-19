package health

import (
	"net/http"

	"codeberg.org/urutau-ltd/aile"
)

func Mount(a *aile.App) {
	a.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
}
