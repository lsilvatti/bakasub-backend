package handlers

import (
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
)

type FolderProcessor interface {
	AddFolder(folder models.FolderConfig) error
	GetFolders() ([]models.FolderConfig, error)
	RemoveFolder(id int) bool
	IsFolder(path string) bool
	IsFile(path string) bool
	ListVideoFiles(path string) ([]string, error)
	ListSubtitleFiles(path string) ([]string, error)
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

	if !h.Service.IsFolder(req.Path) {
		utils.Error(w, http.StatusBadRequest, "Provided path is not a folder")
		return
	}

	if err := h.Service.AddFolder(req.ToModel()); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to add favorite folder: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, "Favorite folder added successfully")
}

func (h *FolderHandler) RemoveFavoriteFolder(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeAndValidate[RemoveFolderRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if !h.Service.RemoveFolder(req.ID) {
		utils.Error(w, http.StatusInternalServerError, "Failed to remove favorite folder")
		return
	}

	utils.Success(w, http.StatusOK, "Favorite folder removed successfully")
}

func (h *FolderHandler) GetFavoriteFolders(w http.ResponseWriter, r *http.Request) {
	folders, err := h.Service.GetFolders()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve favorite folders: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", folders)
}

func (h *FolderHandler) ListVideoFiles(w http.ResponseWriter, r *http.Request) {
	folderPath := r.URL.Query().Get("path")

	if folderPath == "" {
		utils.Error(w, http.StatusBadRequest, "Missing 'path' query parameter")
		return
	}

	files, err := h.Service.ListVideoFiles(folderPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to list files in folder: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", files)
}

func (h *FolderHandler) ListSubtitleFiles(w http.ResponseWriter, r *http.Request) {
	folderPath := r.URL.Query().Get("path")

	if folderPath == "" {
		utils.Error(w, http.StatusBadRequest, "Missing 'path' query parameter")
		return
	}

	files, err := h.Service.ListSubtitleFiles(folderPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to list files in folder: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "", files)
}
