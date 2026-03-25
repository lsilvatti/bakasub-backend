package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
)

type FolderFileSystemProvider interface {
	ReadFolder(path string) ([]os.DirEntry, error)
}

type FolderService struct {
	DB *sql.DB
	FS FolderFileSystemProvider
}

func NewFolderService(db *sql.DB, fs FolderFileSystemProvider) *FolderService {
	return &FolderService{
		DB: db,
		FS: fs,
	}
}

func (s *FolderService) CreateFolder(folder models.FolderConfig) error {
	query := "INSERT INTO folders (alias, path) VALUES ($1, $2)"
	_, err := s.DB.Exec(query, folder.Alias, folder.Path)

	if err != nil {
		utils.LogError("folder", "Failed to add favorite folder", map[string]any{
			"alias": folder.Alias,
			"path":  folder.Path,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("folder", "create", "Favorite folder added", map[string]any{
		"alias": folder.Alias,
		"path":  folder.Path,
	})

	return nil
}

func (s *FolderService) DeleteFolder(id int) bool {
	query := "DELETE FROM folders WHERE id = $1"
	result, err := s.DB.Exec(query, id)
	if err != nil {
		utils.LogError("folder", "Failed to execute delete folder query", map[string]any{
			"id":    id,
			"error": err.Error(),
		})
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError("folder", "Failed to get rows affected for folder deletion", map[string]any{
			"id":    id,
			"error": err.Error(),
		})
		return false
	}

	if rowsAffected == 0 {
		utils.LogError("folder", "Folder to remove not found in database", map[string]any{
			"id": id,
		})
		return false
	}

	utils.LogInfo("folder", "delete", "Favorite folder removed", map[string]any{
		"id": id,
	})

	return true
}

func (s *FolderService) GetFolders() ([]models.FolderConfig, error) {
	query := "SELECT id, alias, path FROM folders"
	rows, err := s.DB.Query(query)
	if err != nil {
		utils.LogError("folder", "Failed to get favorite folders from database", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}
	defer rows.Close()

	folders := make([]models.FolderConfig, 0)

	for rows.Next() {
		var folder models.FolderConfig
		if err := rows.Scan(&folder.ID, &folder.Alias, &folder.Path); err != nil {
			utils.LogError("folder", "Failed to scan folder row", map[string]any{
				"error": err.Error(),
			})
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (s *FolderService) ListFiles(folderPath string) ([]string, error) {
	entries, err := s.FS.ReadFolder(folderPath)
	if err != nil {
		utils.LogError("folder", "Failed to read directory for files", map[string]any{
			"path":  folderPath,
			"error": err.Error(),
		})
		return nil, err
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (s *FolderService) IsVideoFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	videoExtensions := []string{".mkv", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mpeg", ".mpg"}
	for _, vExt := range videoExtensions {
		if ext == vExt {
			return true
		}
	}
	return false
}

func (s *FolderService) IsSubtitleFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	subtitleExtensions := []string{".srt", ".ass", ".ssa", ".vtt", ".sub"}
	for _, sExt := range subtitleExtensions {
		if ext == sExt {
			return true
		}
	}
	return false
}

func (s *FolderService) IsFolder(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (s *FolderService) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (s *FolderService) ListVideoFiles(folderPath string) ([]string, error) {
	utils.SendSSE("info", "folder", "Scanning directory for video files...", map[string]any{"path": folderPath})

	entries, err := s.FS.ReadFolder(folderPath)
	if err != nil {
		utils.LogError("folder", "Failed to read directory for videos", map[string]any{
			"path":  folderPath,
			"error": err.Error(),
		})
		utils.SendSSE("error", "folder", "Failed to scan directory.", nil)
		return nil, err
	}

	videoFiles := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && s.IsVideoFile(entry.Name()) {
			videoFiles = append(videoFiles, entry.Name())
		}
	}

	utils.LogInfo("folder", "scan", "Video directory scan completed", map[string]any{
		"path":  folderPath,
		"count": len(videoFiles),
	})
	utils.SendSSE("success", "folder", "Directory scan complete.", map[string]any{"found": len(videoFiles)})

	return videoFiles, nil
}

func (s *FolderService) ListSubtitleFiles(folderPath string) ([]string, error) {
	utils.SendSSE("info", "folder", "Scanning directory for subtitle files...", map[string]any{"path": folderPath})

	entries, err := s.FS.ReadFolder(folderPath)
	if err != nil {
		utils.LogError("folder", "Failed to read directory for subtitles", map[string]any{
			"path":  folderPath,
			"error": err.Error(),
		})
		utils.SendSSE("error", "folder", "Failed to scan directory.", nil)
		return nil, err
	}

	subtitleFiles := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && s.IsSubtitleFile(entry.Name()) {
			subtitleFiles = append(subtitleFiles, entry.Name())
		}
	}

	utils.LogInfo("folder", "scan", "Subtitle directory scan completed", map[string]any{
		"path":  folderPath,
		"count": len(subtitleFiles),
	})
	utils.SendSSE("success", "folder", "Directory scan complete.", map[string]any{"found": len(subtitleFiles)})

	return subtitleFiles, nil
}
