package health_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/health"
)

func ExampleMount() {
	app := aile.MustNew()
	health.Mount(app)

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 ok
}
