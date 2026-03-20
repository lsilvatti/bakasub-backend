package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
)

type ConfigProcessor interface {
	GetConfig() (models.UserConfig, error)
	UpdateConfig(config models.UserConfig) error
}

type ConfigHandler struct {
	Service ConfigProcessor
}

func (h *ConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	config, err := h.Service.GetConfig()

	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSON(w, http.StatusOK, "success", "No user config found", models.UserConfig{})
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve user config: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", config)
}

func (h *ConfigHandler) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	var config models.UserConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate.Struct(config); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.UpdateConfig(config); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update user config: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "User config updated successfully", nil)
}
