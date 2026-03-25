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
	CreatePreset(preset models.TranslationPreset) error
	UpdatePreset(preset models.TranslationPreset) error
	DeletePreset(id int) error
}

func (h *PresetHandler) GetPresets(w http.ResponseWriter, r *http.Request) {
	presets, err := h.Service.GetPresets()

	if err != nil {
		utils.LogError("preset_handler", "Failed to retrieve presets list", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve presets: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Presets retrieved", map[string]interface{}{
		"presets": presets,
	})
}

func (h *PresetHandler) CreatePreset(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[AddPresetRequest](r)
	if err != nil {
		utils.LogError("preset_handler", "Invalid payload for CreatePreset", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.CreatePreset(reqData.ToModel()); err != nil {
		utils.LogError("preset_handler", "Failed to create preset via service", map[string]any{
			"alias": reqData.Alias,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to create preset: "+err.Error())
		return
	}

	utils.LogInfo("preset_handler", "success", "Successfully created preset", map[string]any{
		"alias": reqData.Alias,
	})

	utils.JSON(w, http.StatusOK, "success", "Preset created successfully", nil)
}

func (h *PresetHandler) UpdatePreset(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[UpdatePresetRequest](r)
	if err != nil {
		utils.LogError("preset_handler", "Invalid payload for UpdatePreset", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	presetModel := reqData.ToModel()
	if err := h.Service.UpdatePreset(presetModel); err != nil {
		utils.LogError("preset_handler", "Failed to update preset via service", map[string]any{
			"id":    reqData.ID,
			"alias": reqData.Alias,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to update preset: "+err.Error())
		return
	}

	utils.LogInfo("preset_handler", "success", "Successfully updated preset", map[string]any{
		"id":    reqData.ID,
		"alias": reqData.Alias,
	})

	utils.JSON(w, http.StatusOK, "success", "Preset updated successfully", nil)
}

func (h *PresetHandler) DeletePreset(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[DeletePresetRequest](r)
	if err != nil {
		utils.LogError("preset_handler", "Invalid payload for DeletePreset", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.DeletePreset(reqData.ID); err != nil {
		utils.LogError("preset_handler", "Failed to delete preset via service", map[string]any{
			"id":    reqData.ID,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to delete preset: "+err.Error())
		return
	}

	utils.LogInfo("preset_handler", "success", "Successfully deleted preset", map[string]any{
		"id": reqData.ID,
	})

	utils.JSON(w, http.StatusOK, "success", "Preset deleted successfully", nil)
}
