package routes

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func APIRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Get("/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Mount("/openrouter", OpenRouterTranslateRoutes(database))
	r.Mount("/video", VideoRoutes())
	r.Mount("/config", ConfigRoutes(database))
	r.Mount("/folders", FoldersRoutes(database))
	r.Mount("/presets", PresetRoutes(database))
	r.Mount("/languages", LanguageRoutes(database))

	return r
}
