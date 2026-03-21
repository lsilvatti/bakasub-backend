package handlers

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"net/http"
)

type LanguageService interface {
	GetLanguages() ([]models.Language, error)
	GetLanguageByCode(code string) (*models.Language, error)
	AddLanguage(lang models.Language) error
	UpdateLanguage(lang models.Language) error
	DeleteLanguage(code string) error
}

type LanguageHandler struct {
	Service LanguageService
}

func (h *LanguageHandler) GetLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	languages, err := h.Service.GetLanguages()
	if err != nil {
		utils.LogError("language_handler", "Failed to retrieve languages list", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error fetching languages: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Languages listed successfully", map[string]interface{}{
		"languages": languages,
	})
}

func (h *LanguageHandler) AddLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[AddLanguageRequest](r)
	if err != nil {
		utils.LogError("language_handler", "Invalid payload for AddLanguage", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid fields: "+err.Error())
		return
	}

	err = h.Service.AddLanguage(reqData.ToModel())
	if err != nil {
		utils.LogError("language_handler", "Failed to add language via service", map[string]any{
			"code":  reqData.Code,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error adding language: "+err.Error())
		return
	}

	utils.LogInfo("language_handler", "success", "Successfully added new language", map[string]any{
		"code": reqData.Code,
		"name": reqData.Name,
	})

	utils.JSON(w, http.StatusOK, "success", "Language added successfully", nil)
}

func (h *LanguageHandler) UpdateLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[UpdateLanguageRequest](r)
	if err != nil {
		utils.LogError("language_handler", "Invalid payload for UpdateLanguage", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid fields: "+err.Error())
		return
	}

	err = h.Service.UpdateLanguage(reqData.ToModel())
	if err != nil {
		utils.LogError("language_handler", "Failed to update language via service", map[string]any{
			"code":  reqData.Code,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error updating language: "+err.Error())
		return
	}

	utils.LogInfo("language_handler", "success", "Successfully updated language", map[string]any{
		"code": reqData.Code,
		"name": reqData.Name,
	})

	utils.JSON(w, http.StatusOK, "success", "Language updated successfully", nil)
}

func (h *LanguageHandler) DeleteLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[DeleteLanguageRequest](r)
	if err != nil {
		utils.LogError("language_handler", "Invalid payload for DeleteLanguage", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid fields: "+err.Error())
		return
	}

	err = h.Service.DeleteLanguage(reqData.Code)
	if err != nil {
		utils.LogError("language_handler", "Failed to delete language via service", map[string]any{
			"code":  reqData.Code,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error deleting language: "+err.Error())
		return
	}

	utils.LogInfo("language_handler", "success", "Successfully deleted language", map[string]any{
		"code": reqData.Code,
	})

	utils.JSON(w, http.StatusOK, "success", "Language deleted successfully", nil)
}
