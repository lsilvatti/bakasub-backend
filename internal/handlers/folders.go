package handlers

import (
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
)

type FolderProcessor interface {
	AddFolder(folder models.FolderConfig) error
	GetFolders() ([]models.FolderConfig, error)
}

type FolderHandler struct {
	Service FolderProcessor
}

func (h *FolderHandler) AddFavoriteFolder(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeAndValidate[AddFolderRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.AddFolder(req.ToModel()); err != nil {
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
