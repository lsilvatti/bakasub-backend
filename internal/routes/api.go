package routes

import (
	"bakasub-backend/internal/utils"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func APIRoutes(database *sql.DB, secretKey string) chi.Router {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Mount("/openrouter", OpenRouterTranslateRoutes(database, secretKey))
	r.Mount("/video", VideoRoutes(database, secretKey))
	r.Mount("/config", ConfigRoutes(database, secretKey))
	r.Mount("/favorites", FavoritesRoutes(database))
	r.Mount("/folders", FoldersRoutes(database))
	r.Mount("/presets", PresetRoutes(database))
	r.Mount("/languages", LanguageRoutes(database))
	r.Mount("/logs", LogRoutes(database))
	r.Mount("/jobs", JobRoutes(database))

	r.Get("/events", utils.Broker.ServeHTTP)

	return r
}
