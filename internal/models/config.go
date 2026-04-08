package models

type UserConfig struct {
	DefaultModel        string `json:"default_model"`
	DefaultPreset       string `json:"default_preset"`
	DefaultLanguage     string `json:"default_language"`
	RemoveSdhDefault    bool   `json:"remove_sdh_default"`
	VideoTimeoutMinutes int    `json:"video_timeout_minutes"`
	LogRetentionDays    int    `json:"log_retention_days"`
}
