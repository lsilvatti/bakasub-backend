package routes

import (
	"bakasub-backend/internal/handlers"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func FoldersRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()
	foldersHandler := handlers.FolderHandler{DB: database}

	r.Get("/", foldersHandler.GetFavoriteFolders)
	r.Post("/", foldersHandler.AddFavoriteFolder)

	return r
}
