package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/htmx"
	"codeberg.org/urutau-ltd/aile/v2/x/resource"
)

type provider struct {
	ID   int
	Name string
}

type appSettings struct {
	AppName string
}

type providerStore struct {
	mu     sync.Mutex
	nextID int
	items  map[int]provider
}

func newProviderStore() *providerStore {
	return &providerStore{
		nextID: 3,
		items: map[int]provider{
			1: {ID: 1, Name: "Hetzner"},
			2: {ID: 2, Name: "OVH"},
		},
	}
}

func (s *providerStore) list() []provider {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]provider, 0, len(s.items))
	for i := 1; i < s.nextID; i++ {
		if item, ok := s.items[i]; ok {
			out = append(out, item)
		}
	}
	return out
}

func (s *providerStore) create(name string) provider {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := provider{ID: s.nextID, Name: name}
	s.items[item.ID] = item
	s.nextID++
	return item
}

func (s *providerStore) get(id int) (provider, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	return item, ok
}

func (s *providerStore) update(id int, name string) (provider, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	if !ok {
		return provider{}, false
	}
	item.Name = name
	s.items[id] = item
	return item, true
}

func (s *providerStore) delete(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, id)
}

type providersHandler struct {
	store *providerStore
}

type providersPageData struct {
	Providers []provider
	Editor    template.HTML
}

type providerEditorData struct {
	Title    string
	Action   string
	Provider provider
}

var providersPageTmpl = template.Must(template.New("providers-page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>aile html admin</title>
  <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js"></script>
</head>
<body>
  <h1>Providers</h1>
  <p><a href="/app-settings">Open singleton settings</a></p>
  <button hx-get="/providers/new" hx-target="#providers-editor" hx-swap="innerHTML">New provider</button>
  <ul id="providers-list" hx-get="/providers" hx-trigger="load, providers:changed from:body" hx-target="#providers-list" hx-swap="outerHTML">
    {{range .Providers}}
      <li>
        {{.Name}}
        <button hx-get="/providers/{{.ID}}/edit" hx-target="#providers-editor" hx-swap="innerHTML">Edit</button>
        <button hx-delete="/providers/{{.ID}}" hx-target="#providers-editor" hx-swap="innerHTML">Delete</button>
      </li>
    {{end}}
  </ul>
  <section id="providers-editor">{{.Editor}}</section>
</body>
</html>`))

var providersListTmpl = template.Must(template.New("providers-list").Parse(`<ul id="providers-list" hx-get="/providers" hx-trigger="load, providers:changed from:body" hx-target="#providers-list" hx-swap="outerHTML">
{{range .}}
  <li>
    {{.Name}}
    <button hx-get="/providers/{{.ID}}/edit" hx-target="#providers-editor" hx-swap="innerHTML">Edit</button>
    <button hx-delete="/providers/{{.ID}}" hx-target="#providers-editor" hx-swap="innerHTML">Delete</button>
  </li>
{{end}}
</ul>`))

var providersEditorTmpl = template.Must(template.New("providers-editor").Parse(`<div>
  <h2>{{.Title}}</h2>
  <form method="post" action="{{.Action}}" hx-post="{{.Action}}" hx-target="#providers-editor" hx-swap="innerHTML">
    <label>Name <input type="text" name="name" value="{{.Provider.Name}}"></label>
    <button type="submit">Save</button>
  </form>
</div>`))

func (h providersHandler) Index(w http.ResponseWriter, r *http.Request) {
	writeHTML(w)
	items := h.store.list()
	if htmx.IsRequest(r) && htmx.TargetIs(r, "providers-list") {
		_ = providersListTmpl.Execute(w, items)
		return
	}

	editor := template.HTML(`<p>Select a provider or create a new one.</p>`)
	_ = providersPageTmpl.Execute(w, providersPageData{
		Providers: items,
		Editor:    editor,
	})
}

func (h providersHandler) New(w http.ResponseWriter, r *http.Request) {
	writeHTML(w)
	_ = providersEditorTmpl.Execute(w, providerEditorData{
		Title:  "New provider",
		Action: "/providers",
	})
}

func (h providersHandler) Create(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		writeHTML(w)
		_ = providersEditorTmpl.Execute(w, providerEditorData{
			Title:  "New provider",
			Action: "/providers",
		})
		return
	}

	created := h.store.create(name)
	htmx.SetTrigger(w, "providers:changed")
	writeHTML(w)
	_ = providersEditorTmpl.Execute(w, providerEditorData{
		Title:    "Edit provider",
		Action:   "/providers/" + strconv.Itoa(created.ID) + "/edit",
		Provider: created,
	})
}

func (h providersHandler) Show(w http.ResponseWriter, r *http.Request) {
	h.Edit(w, r)
}

func (h providersHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		aile.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, ok := h.store.get(id)
	if !ok {
		aile.Error(w, http.StatusNotFound, "not found")
		return
	}

	writeHTML(w)
	_ = providersEditorTmpl.Execute(w, providerEditorData{
		Title:    "Edit provider",
		Action:   "/providers/" + strconv.Itoa(item.ID) + "/edit",
		Provider: item,
	})
}

func (h providersHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		aile.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, ok := h.store.update(id, strings.TrimSpace(r.FormValue("name")))
	if !ok {
		aile.Error(w, http.StatusNotFound, "not found")
		return
	}

	htmx.SetTrigger(w, "providers:changed")
	writeHTML(w)
	_ = providersEditorTmpl.Execute(w, providerEditorData{
		Title:    "Edit provider",
		Action:   "/providers/" + strconv.Itoa(item.ID) + "/edit",
		Provider: item,
	})
}

func (h providersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		aile.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	h.store.delete(id)
	htmx.SetTrigger(w, "providers:changed")
	writeHTML(w)
	_, _ = w.Write([]byte(`<p>Select a provider or create a new one.</p>`))
}

type appSettingsHandler struct {
	mu       sync.Mutex
	settings appSettings
}

var settingsTmpl = template.Must(template.New("settings").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>aile singleton settings</title>
</head>
<body>
  <h1>Application settings</h1>
  <form method="post" action="/app-settings/edit">
    <label>App name <input type="text" name="app_name" value="{{.AppName}}"></label>
    <button type="submit">Save</button>
  </form>
  <p><a href="/providers">Back to providers</a></p>
</body>
</html>`))

func (h *appSettingsHandler) Show(w http.ResponseWriter, r *http.Request) {
	h.Edit(w, r)
}

func (h *appSettingsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	writeHTML(w)
	_ = settingsTmpl.Execute(w, h.settings)
}

func (h *appSettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.settings.AppName = strings.TrimSpace(r.FormValue("app_name"))
	snapshot := h.settings
	h.mu.Unlock()

	writeHTML(w)
	_ = settingsTmpl.Execute(w, snapshot)
}

func writeHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func main() {
	app := aile.MustNew(aile.WithAddr(":9094"), aile.WithMiddleware(aile.Recovery()))

	providers := providersHandler{store: newProviderStore()}
	settings := &appSettingsHandler{settings: appSettings{AppName: "Gavia"}}

	if err := resource.MountCollection(app, "/providers", providers); err != nil {
		log.Fatal(err)
	}
	if err := resource.MountSingleton(app, "/app-settings", settings); err != nil {
		log.Fatal(err)
	}

	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/providers", http.StatusSeeOther)
	})

	log.Println("listening on http://localhost:9094")
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
