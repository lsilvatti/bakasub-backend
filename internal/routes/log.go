package routes

import (
	"database/sql"

	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

func LogRoutes(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	logService := services.NewLogService(db)

	logHandler := &handlers.LogHandler{
		Service: logService,
	}

	r.Get("/", logHandler.GetLogsHandler)

	return r
}
