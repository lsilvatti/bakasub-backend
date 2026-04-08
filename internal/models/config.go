package models

type UserConfig struct {
	DefaultModel           string `json:"default_model"`
	DefaultPreset          string `json:"default_preset"`
	DefaultLanguage        string `json:"default_language"`
	RemoveSdhDefault       bool   `json:"remove_sdh_default"`
	VideoTimeoutMinutes    int    `json:"video_timeout_minutes"`
	LogRetentionDays       int    `json:"log_retention_days"`
	OpenRouterApiKey       string `json:"openrouter_api_key"`
	TmdbAccessToken        string `json:"tmdb_access_token"`
	TmdbMetadataEnabled    bool   `json:"tmdb_metadata_enabled"`
	ConcurrentTranslations int    `json:"concurrent_translations"`
	MaxRetries             int    `json:"max_retries"`
	BaseRetryDelay         int    `json:"base_retry_delay"`
}
