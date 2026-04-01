package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func FavoritesRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	favoritesService := services.NewFavoritesService(database)

	favoritesHandler := &handlers.FavoritesHandler{
		Service: favoritesService,
	}

	r.Get("/", favoritesHandler.GetFavorites)
	r.Put("/", favoritesHandler.UpdateFavorites)

	return r
}