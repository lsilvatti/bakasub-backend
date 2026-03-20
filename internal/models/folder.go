package models

type FolderConfig struct {
	ID    int    `json:"id"`
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
