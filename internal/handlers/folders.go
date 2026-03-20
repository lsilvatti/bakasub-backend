package handlers

import (
	"encoding/json"
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
)

type FolderProcessor interface {
	AddFolder(alias, path string) error
	GetFolders() ([]models.FolderConfig, error)
}

type FolderHandler struct {
	Service FolderProcessor
}

func (h *FolderHandler) AddFavoriteFolder(w http.ResponseWriter, r *http.Request) {
	var req AddFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.AddFolder(req.Alias, req.Path); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to add favorite folder: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, "Favorite folder added successfully")
}

func (h *FolderHandler) GetFavoriteFolders(w http.ResponseWriter, r *http.Request) {
	folders, err := h.Service.GetFolders()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve favorite folders: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", folders)
}
