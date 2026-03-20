package models

type UserConfig struct {
	DefaultModel     string `json:"default_model"`
	DefaultPreset    string `json:"default_preset"`
	RemoveSdhDefault bool   `json:"remove_sdh_default"`
}
