package models

import "time"

type TranslationJob struct {
	ID               string    `json:"id"`
	Status           string    `json:"status"`
	FilePath         string    `json:"file_path"`
	TargetLang       string    `json:"target_lang"`
	Preset           string    `json:"preset"`
	Model            string    `json:"model"`
	TotalLines       int       `json:"total_lines"`
	ProcessedLines   int       `json:"processed_lines"`
	CachedLines      int       `json:"cached_lines"`
	PromptTokens     int       `json:"prompt_tokens"`
	CompletionTokens int       `json:"completion_tokens"`
	CostUSD          float64   `json:"cost_usd"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
