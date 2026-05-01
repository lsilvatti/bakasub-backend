package routes

import (
	"bakasub-backend/internal/services"
	"bakasub-backend/internal/utils"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func APIRoutes(database *sql.DB, secretKey string) chi.Router {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, "success", "Health check completed", struct {
			OK         bool                      `json:"ok"`
			VideoTools services.VideoToolsStatus `json:"videoTools"`
		}{
			OK:         true,
			VideoTools: services.CheckVideoTools(),
		})
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
