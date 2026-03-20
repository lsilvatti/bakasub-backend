package models

type TranslationPreset struct {
	ID           int     `json:"id"`
	Alias        string  `json:"alias"`
	Name         string  `json:"name"`
	SystemPrompt string  `json:"system_prompt"`
	BatchSize    int     `json:"batch_size"`
	Temperature  float64 `json:"temperature"`
}
