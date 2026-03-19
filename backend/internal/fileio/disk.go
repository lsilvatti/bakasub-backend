package fileio

import "os"

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SaveFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func ReadFolder(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func CreateFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

func DeleteFolder(path string) error {
	return os.RemoveAll(path)
}
