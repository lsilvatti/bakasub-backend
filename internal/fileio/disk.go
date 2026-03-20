package fileio

import "os"

type FileIOService struct{}

func NewFileIOService() *FileIOService {
	return &FileIOService{}
}

func (s *FileIOService) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *FileIOService) SaveFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *FileIOService) DeleteFile(path string) error {
	return os.Remove(path)
}

func (s *FileIOService) ReadFolder(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (s *FileIOService) CreateFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

func (s *FileIOService) DeleteFolder(path string) error {
	return os.RemoveAll(path)
}
