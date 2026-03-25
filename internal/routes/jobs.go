package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func JobRoutes(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	jobService := services.NewJobService(db)
	jobHandler := &handlers.JobHandler{Service: jobService}

	r.Get("/", jobHandler.ListJobs)
	r.Get("/{id}", jobHandler.GetJob)

	return r
}
