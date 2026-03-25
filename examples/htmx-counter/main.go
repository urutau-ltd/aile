package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/health"
)

// This is a similar concept as React component props. But in this
// case this type acts as a page prop.
type counterPageData struct {
	// A simple integer used for counting.
	Count int
}

var pageTmpl *template.Template = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>aile + HTMX counter</title>
  <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js"></script>
  <style>
    body { font-family: system-ui, sans-serif; max-width: 48rem; margin: 3rem auto; padding: 0 1rem; }
    button { padding: 0.6rem 1rem; }
    .card { border: 1px solid #ddd; border-radius: 12px; padding: 1rem; }
  </style>
</head>
<body>
  <h1>aile + HTMX</h1>
  <div class="card">
    <p>This is a tiny HTMX counter served by aile.</p>
    <div id="counter">{{template "counter" .}}</div>
  </div>
</body>
</html>
{{define "counter"}}
  <div id="counter">
    <p><strong>Count:</strong> {{.Count}}</p>
    <button
      hx-post="/increment"
      hx-target="#counter"
      hx-swap="outerHTML">
      Increment
    </button>
  </div>
{{end}}`))

type counter struct {
	mu    sync.Mutex
	value int
}

func (c *counter) get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func (c *counter) inc() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

func main() {
	c := &counter{}

	app, err := aile.New(aile.WithAddr(":9094"), aile.WithMiddleware(aile.Recovery()))
	if err != nil {
		log.Fatal(err)
	}

	health.Mount(app)

	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		if err := pageTmpl.Execute(w, counterPageData{Count: c.get()}); err != nil {
			aile.Error(w, http.StatusInternalServerError, err.Error())
		}
	})

	app.POST("/increment", func(w http.ResponseWriter, r *http.Request) {
		if err := pageTmpl.ExecuteTemplate(w, "counter", counterPageData{Count: c.inc()}); err != nil {
			aile.Error(w, http.StatusInternalServerError, err.Error())
		}
	})

	log.Println("listening on http://localhost:9094")
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
