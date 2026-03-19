package routes

import (
	"bakasub-backend/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func APIRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/v1/translate", handlers.TranslateHandler)

	r.Mount("/v1/video", VideoRoutes())

	return r
}
