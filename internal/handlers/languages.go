package handlers

import (
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
)

type LanguageService interface {
	GetLanguages() ([]models.Language, error)
	GetLanguageByCode(code string) (*models.Language, error)
	AddLanguage(code string, name string) error
	UpdateLanguage(code string, name string) error
	DeleteLanguage(code string) error
}

type LanguageHandler struct {
	Service LanguageService
}

func (h *LanguageHandler) GetLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	languages, err := h.Service.GetLanguages()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao buscar idiomas: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Idiomas listados com sucesso", map[string]interface{}{
		"languages": languages,
	})
}

func (h *LanguageHandler) AddLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[AddLanguageRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos inválidos: "+err.Error())
		return
	}

	err = h.Service.AddLanguage(reqData.Code, reqData.Name)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao adicionar idioma: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Idioma adicionado com sucesso", nil)
}

func (h *LanguageHandler) UpdateLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[UpdateLanguageRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos inválidos: "+err.Error())
		return
	}

	err = h.Service.UpdateLanguage(reqData.Code, reqData.Name)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao atualizar idioma: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Idioma atualizado com sucesso", nil)
}

func (h *LanguageHandler) DeleteLanguageHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[DeleteLanguageRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos inválidos: "+err.Error())
		return
	}

	err = h.Service.DeleteLanguage(reqData.Code)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao deletar idioma: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Idioma deletado com sucesso", nil)
}
