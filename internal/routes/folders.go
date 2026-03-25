package routes

import (
	"bakasub-backend/internal/fileio"
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func FoldersRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	diskService := fileio.NewFileIOService()
	folderService := services.NewFolderService(database, diskService)

	foldersHandler := handlers.FolderHandler{Service: folderService}

	r.Get("/", foldersHandler.GetFolders)
	r.Post("/", foldersHandler.CreateFolder)
	r.Delete("/", foldersHandler.DeleteFolder)
	r.Get("/scan/videos", foldersHandler.ListVideoFiles)
	r.Get("/scan/subtitles", foldersHandler.ListSubtitleFiles)

	return r
}
