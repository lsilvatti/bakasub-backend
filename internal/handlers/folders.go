package handlers

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"net/http"
)

type FolderProcessor interface {
	CreateFolder(folder models.FolderConfig) error
	GetFolders() ([]models.FolderConfig, error)
	DeleteFolder(id int) bool
	IsFolder(path string) bool
	IsFile(path string) bool
	ListVideoFiles(path string) ([]string, error)
	ListSubtitleFiles(path string) ([]string, error)
}

type FolderHandler struct {
	Service FolderProcessor
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeAndValidate[AddFolderRequest](r)
	if err != nil {
		utils.LogError("folder_handler", "Invalid payload for CreateFolder", map[string]any{"error": err.Error()})
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if !h.Service.IsFolder(req.Path) {
		utils.LogError("folder_handler", "Provided path is not a folder", map[string]any{"path": req.Path})
		utils.Error(w, http.StatusBadRequest, "Provided path is not a folder")
		return
	}

	if err := h.Service.CreateFolder(req.ToModel()); err != nil {
		utils.LogError("folder_handler", "Failed to create folder via service", map[string]any{
			"path":  req.Path,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to create folder: "+err.Error())
		return
	}

	utils.LogInfo("folder_handler", "success", "Successfully created folder", map[string]any{"path": req.Path})
	utils.Success(w, http.StatusOK, "Folder created successfully")
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeAndValidate[RemoveFolderRequest](r)
	if err != nil {
		utils.LogError("folder_handler", "Invalid payload for DeleteFolder", map[string]any{"error": err.Error()})
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if !h.Service.DeleteFolder(req.ID) {
		utils.LogError("folder_handler", "Failed to delete folder via service", map[string]any{"id": req.ID})
		utils.Error(w, http.StatusInternalServerError, "Failed to delete folder")
		return
	}

	utils.LogInfo("folder_handler", "success", "Successfully deleted folder", map[string]any{"id": req.ID})
	utils.Success(w, http.StatusOK, "Folder deleted successfully")
}

func (h *FolderHandler) GetFolders(w http.ResponseWriter, r *http.Request) {
	folders, err := h.Service.GetFolders()
	if err != nil {
		utils.LogError("folder_handler", "Failed to retrieve folders", map[string]any{"error": err.Error()})
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve folders: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Folders retrieved", folders)
}

func (h *FolderHandler) ListVideoFiles(w http.ResponseWriter, r *http.Request) {
	folderPath := r.URL.Query().Get("path")

	if folderPath == "" {
		utils.LogError("folder_handler", "Missing 'path' query parameter for ListVideoFiles", nil)
		utils.Error(w, http.StatusBadRequest, "Missing 'path' query parameter")
		return
	}

	files, err := h.Service.ListVideoFiles(folderPath)
	if err != nil {
		utils.LogError("folder_handler", "Failed to list video files via service", map[string]any{
			"path":  folderPath,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to list files in folder: "+err.Error())
		return
	}

	utils.LogInfo("folder_handler", "success", "Successfully listed video files", map[string]any{
		"path":  folderPath,
		"found": len(files),
	})

	utils.JSON(w, http.StatusOK, "success", "Video files listed", files)
}

func (h *FolderHandler) ListSubtitleFiles(w http.ResponseWriter, r *http.Request) {
	folderPath := r.URL.Query().Get("path")

	if folderPath == "" {
		utils.LogError("folder_handler", "Missing 'path' query parameter for ListSubtitleFiles", nil)
		utils.Error(w, http.StatusBadRequest, "Missing 'path' query parameter")
		return
	}

	files, err := h.Service.ListSubtitleFiles(folderPath)
	if err != nil {
		utils.LogError("folder_handler", "Failed to list subtitle files via service", map[string]any{
			"path":  folderPath,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Failed to list files in folder: "+err.Error())
		return
	}

	utils.LogInfo("folder_handler", "success", "Successfully listed subtitle files", map[string]any{
		"path":  folderPath,
		"found": len(files),
	})

	utils.JSON(w, http.StatusOK, "success", "Subtitle files listed", files)
}
