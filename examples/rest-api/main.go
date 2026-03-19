package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"codeberg.org/urutau-ltd/aile"
	"codeberg.org/urutau-ltd/aile/x/health"
	xlogger "codeberg.org/urutau-ltd/aile/x/logger"
)

type article struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type createArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type store struct {
	mu     sync.Mutex
	nextID int
	items  map[int]article
}

func newStore() *store {
	return &store{nextID: 1, items: make(map[int]article)}
}

func (s *store) list() []article {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]article, 0, len(s.items))
	for i := 1; i < s.nextID; i++ {
		if item, ok := s.items[i]; ok {
			out = append(out, item)
		}
	}
	return out
}

func (s *store) create(req createArticleRequest) article {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := article{ID: s.nextID, Title: req.Title, Body: req.Body}
	s.items[item.ID] = item
	s.nextID++
	return item
}

func (s *store) get(id int) (article, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[id]
	return item, ok
}

func main() {
	st := newStore()

	app, err := aile.New(
		aile.WithAddr(":9094"),
		aile.WithMiddleware(
			aile.Recovery(),
			xlogger.Middleware(slog.Default()),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	health.Mount(app)

	app.HandleFunc("GET /api/articles", func(w http.ResponseWriter, r *http.Request) {
		_ = aile.WriteJSON(w, http.StatusOK, st.list())
	})

	app.HandleFunc("POST /api/articles", func(w http.ResponseWriter, r *http.Request) {
		var req createArticleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			aile.Error(w, http.StatusBadRequest, "invalid json")
			return
		}
		if req.Title == "" {
			aile.Error(w, http.StatusBadRequest, "title is required")
			return
		}
		item := st.create(req)
		_ = aile.WriteJSON(w, http.StatusCreated, item)
	})

	app.HandleFunc("GET /api/articles/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			aile.Error(w, http.StatusBadRequest, "invalid id")
			return
		}
		item, ok := st.get(id)
		if !ok {
			aile.Error(w, http.StatusNotFound, "not found")
			return
		}
		_ = aile.WriteJSON(w, http.StatusOK, item)
	})

	log.Println("listening on http://localhost:9094")
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
