package services

import (
	"bakasub-backend/internal/models"
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

func (s *FolderService) AddFolder(folder models.FolderConfig) error {
	query := "INSERT INTO favorite_folders (alias, path) VALUES (?, ?)"
	_, err := s.DB.Exec(query, folder.Alias, folder.Path)
	return err
}

func (s *FolderService) RemoveFolder(id int) bool {
	query := "DELETE FROM favorite_folders WHERE id = ?"
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false
	}

	return rowsAffected > 0
}

func (s *FolderService) GetFolders() ([]models.FolderConfig, error) {
	query := "SELECT id, alias, path FROM favorite_folders"
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := make([]models.FolderConfig, 0)

	for rows.Next() {
		var folder models.FolderConfig
		if err := rows.Scan(&folder.ID, &folder.Alias, &folder.Path); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

func (s *FolderService) ListFiles(folderPath string) ([]string, error) {
	entries, err := s.FS.ReadFolder(folderPath)
	if err != nil {
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
	entries, err := s.FS.ReadFolder(folderPath)
	if err != nil {
		return nil, err
	}

	videoFiles := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && s.IsVideoFile(entry.Name()) {
			videoFiles = append(videoFiles, entry.Name())
		}
	}
	return videoFiles, nil
}

func (s *FolderService) ListSubtitleFiles(folderPath string) ([]string, error) {
	entries, err := s.FS.ReadFolder(folderPath)
	if err != nil {
		return nil, err
	}

	subtitleFiles := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && s.IsSubtitleFile(entry.Name()) {
			subtitleFiles = append(subtitleFiles, entry.Name())
		}
	}
	return subtitleFiles, nil
}
