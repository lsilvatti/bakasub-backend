package handlers

import (
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
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
		utils.Error(w, http.StatusBadRequest, "Invalid fields: "+err.Error())
		return
	}

	err = h.Service.AddLanguage(reqData.ToModel())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Error adding language: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Language added successfully", nil)
}

func (h *LanguageHandler) UpdateLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[UpdateLanguageRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid fields: "+err.Error())
		return
	}

	err = h.Service.UpdateLanguage(reqData.ToModel())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Error updating language: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Language updated successfully", nil)
}

func (h *LanguageHandler) DeleteLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[DeleteLanguageRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid fields: "+err.Error())
		return
	}

	err = h.Service.DeleteLanguage(reqData.Code)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Error deleting language: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Language deleted successfully", nil)
}
