package handlers

import "bakasub-backend/internal/models"

type TranslateRequest struct {
	FilePath   string `json:"filePath" validate:"required"`
	TargetLang string `json:"targetLang" validate:"required"`
	Preset     string `json:"preset" validate:"required"`
	Model      string `json:"model" validate:"required"`
	RemoveSDH  bool   `json:"removeSDH"`
	Context    string `json:"context"`
}

type AddFolderRequest struct {
	Alias string `json:"alias" validate:"required"`
	Path  string `json:"path" validate:"required"`
}

type ListFilesRequest struct {
	Path string `json:"path" validate:"required"`
}

type RemoveFolderRequest struct {
	ID int `json:"id" validate:"required"`
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
	ID           int     `json:"id" validate:"required"`
	Alias        string  `json:"alias"`
	Name         string  `json:"name"`
	SystemPrompt string  `json:"system_prompt"`
	BatchSize    int     `json:"batch_size"`
	Temperature  float64 `json:"temperature"`
}

type UpdateConfigRequest struct {
	DefaultModel        string `json:"default_model" validate:"required"`
	DefaultPreset       string `json:"default_preset" validate:"required"`
	RemoveSdhDefault    bool   `json:"remove_sdh_default"`
	VideoTimeoutMinutes int    `json:"video_timeout_minutes" validate:"required"`
	LogRetentionDays    int    `json:"log_retention_days" validate:"required"`
}

func (r *AddPresetRequest) ToModel() models.TranslationPreset {
	return models.TranslationPreset{
		Alias:        r.Alias,
		Name:         r.Name,
		SystemPrompt: r.SystemPrompt,
		BatchSize:    r.BatchSize,
		Temperature:  r.Temperature,
	}
}

func (r *UpdatePresetRequest) ToModel() models.TranslationPreset {
	return models.TranslationPreset{
		ID:           r.ID,
		Alias:        r.Alias,
		Name:         r.Name,
		SystemPrompt: r.SystemPrompt,
		BatchSize:    r.BatchSize,
		Temperature:  r.Temperature,
	}
}

func (r *AddLanguageRequest) ToModel() models.Language {
	return models.Language{
		Code: r.Code,
		Name: r.Name,
	}
}

func (r *UpdateLanguageRequest) ToModel() models.Language {
	return models.Language{
		Code: r.Code,
		Name: r.Name,
	}
}

func (r *AddFolderRequest) ToModel() models.FolderConfig {
	return models.FolderConfig{
		Alias: r.Alias,
		Path:  r.Path,
	}
}

func (r *UpdateConfigRequest) ToModel() models.UserConfig {
	return models.UserConfig{
		DefaultModel:        r.DefaultModel,
		DefaultPreset:       r.DefaultPreset,
		RemoveSdhDefault:    r.RemoveSdhDefault,
		VideoTimeoutMinutes: r.VideoTimeoutMinutes,
		LogRetentionDays:    r.LogRetentionDays,
	}
}
