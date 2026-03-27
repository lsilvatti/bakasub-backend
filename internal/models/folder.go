package models

type FolderConfig struct {
	ID    int    `json:"id"`
	Alias string `json:"alias"`
	Path  string `json:"path"`
}

type FileNode struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"isDirectory"`
}
