package handlers

type TranslateRequest struct {
	FilePath   string `json:"filePath" validate:"required"`
	TargetLang string `json:"targetLang" validate:"required"`
	Preset     string `json:"preset" validate:"required"`
	Model      string `json:"model" validate:"required"`
	RemoveSDH  bool   `json:"removeSDH"`
}

type AddFolderRequest struct {
	Alias string `json:"alias" validate:"required"`
	Path  string `json:"path" validate:"required"`
}

type GetTrackRequest struct {
	VideoPath string `json:"videoPath" validate:"required"`
}

type ExtractTrackRequest struct {
	VideoPath  string `json:"videoPath" validate:"required"`
	SubtitleId int    `json:"subtitleId" validate:"required"`
}

type MergeTrackRequest struct {
	VideoPath string `json:"videoPath" validate:"required"`
	SrtPath   string `json:"srtPath" validate:"required"`
	LangCode  string `json:"langCode" validate:"required"`
}
