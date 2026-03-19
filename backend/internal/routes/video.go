package routes

import (
	"bakasub-backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func VideoRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/get-tracks", handlers.GetTrackHandler)
	r.Post("/extract-track", handlers.ExtractTrackHandler)
	r.Post("/merge-subtitle", handlers.MergeTrackHandler)

	return r
}
