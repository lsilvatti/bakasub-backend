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
