package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bakasub-backend/internal/utils"
)

type FolderHandler struct {
	DB *sql.DB
}

type FolderConfig struct {
	ID    int    `json:"id"`
	Alias string `json:"alias"`
	Path  string `json:"path"`
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

	query := "INSERT INTO favorite_folders (alias, path) VALUES (?, ?)"
	_, err := h.DB.Exec(query, req.Alias, req.Path)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to add favorite folder: "+err.Error())
		return
	}

	utils.Success(w, http.StatusOK, "Favorite folder added successfully")
}

func (h *FolderHandler) GetFavoriteFolders(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, alias, path FROM favorite_folders"
	rows, err := h.DB.Query(query)

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve favorite folders: "+err.Error())
		return
	}

	defer rows.Close()

	var folders []FolderConfig

	for rows.Next() {
		var folder FolderConfig

		if err := rows.Scan(&folder.ID, &folder.Alias, &folder.Path); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to parse favorite folders: "+err.Error())
			return
		}

		folders = append(folders, folder)
	}

	utils.JSON(w, http.StatusOK, "success", "", folders)
}
