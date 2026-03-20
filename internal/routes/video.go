package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

func VideoRoutes() chi.Router {

	r := chi.NewRouter()

	videoService := services.NewVideoService()

	videoHandler := &handlers.VideoHandler{
		Processor: videoService,
	}

	r.Post("/get-tracks", videoHandler.GetTrackHandler)
	r.Post("/extract-track", videoHandler.ExtractTrackHandler)
	r.Post("/merge-subtitle", videoHandler.MergeTrackHandler)

	return r
}
