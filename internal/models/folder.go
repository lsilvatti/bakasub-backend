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

type ExploreResponse struct {
	Items      []FileNode `json:"items"`
	ParentPath *string    `json:"parentPath"`
	FolderName string     `json:"folderName"`
}

type RootEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
