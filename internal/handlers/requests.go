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

type AddLanguageRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdateLanguageRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type DeleteLanguageRequest struct {
	Code string `json:"code" validate:"required"`
}

type DeletePresetRequest struct {
	ID int `json:"id" validate:"required"`
}

type AddPresetRequest struct {
	Alias        string  `json:"alias" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	SystemPrompt string  `json:"system_prompt" validate:"required"`
	BatchSize    int     `json:"batch_size" validate:"required"`
	Temperature  float64 `json:"temperature" validate:"required"`
}

type UpdatePresetRequest struct {
	ID           int     `json:"id"`
	Alias        string  `json:"alias"`
	Name         string  `json:"name"`
	SystemPrompt string  `json:"system_prompt"`
	BatchSize    int     `json:"batch_size"`
	Temperature  float64 `json:"temperature"`
}
