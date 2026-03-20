package handlers

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"net/http"
)

type PresetHandler struct {
	Service PresetService
}

type PresetService interface {
	GetPresets() ([]models.TranslationPreset, error)
	GetPresetByAlias(alias string) (*models.TranslationPreset, error)
	CreatePreset(preset models.TranslationPreset) error
	UpdatePreset(preset models.TranslationPreset) error
	DeletePreset(id int) error
}

func (h *PresetHandler) GetPresetsHandler(w http.ResponseWriter, r *http.Request) {
	presets, err := h.Service.GetPresets()

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve presets: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", map[string]interface{}{
		"presets": presets,
	})
}

func (h *PresetHandler) CreatePresetHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[AddPresetRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.CreatePreset(reqData.ToModel()); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create preset: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Preset created successfully", nil)
}

func (h *PresetHandler) UpdatePresetHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[UpdatePresetRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	presetModel := reqData.ToModel()
	if err := h.Service.UpdatePreset(presetModel); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update preset: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Preset updated successfully", nil)
}

func (h *PresetHandler) DeletePresetHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[DeletePresetRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.DeletePreset(reqData.ID); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to delete preset: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Preset deleted successfully", nil)
}
