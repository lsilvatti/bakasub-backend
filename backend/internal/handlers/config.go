package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bakasub-backend/internal/utils"
)

type ConfigHandler struct {
	DB *sql.DB
}

type UserConfig struct {
	DefaultModel     string `json:"default_model"`
	DefaultPreset    string `json:"default_preset"`
	RemoveSdhDefault bool   `json:"remove_sdh_default"`
}

func (h *ConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	query := "SELECT default_model, default_preset, remove_sdh_default FROM user_configs WHERE id = 1"
	row := h.DB.QueryRow(query)

	var config UserConfig
	err := row.Scan(&config.DefaultModel, &config.DefaultPreset, &config.RemoveSdhDefault)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusOK, "success", "No user config found", UserConfig{})
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve user config: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", config)
}

func (h *ConfigHandler) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	var config UserConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate.Struct(config); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	query := `
	INSERT INTO user_configs (id, default_model, default_preset, remove_sdh_default)
	VALUES (1, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		default_model=excluded.default_model,
		default_preset=excluded.default_preset,
		remove_sdh_default=excluded.remove_sdh_default;
	`

	_, err := h.DB.Exec(query, config.DefaultModel, config.DefaultPreset, config.RemoveSdhDefault)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update user config: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "User config updated successfully", nil)
}
