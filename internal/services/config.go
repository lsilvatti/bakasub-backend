package services

import (
	"bakasub-backend/internal/models"
	"database/sql"
)

type ConfigService struct {
	DB *sql.DB
}

func NewConfigService(db *sql.DB) *ConfigService {
	return &ConfigService{DB: db}
}

func (s *ConfigService) UpdateConfig(cfg models.UserConfig) error {
	_, err := s.DB.Exec(`
		UPDATE user_config 
		SET default_model = $1, default_preset = $2, remove_sdh_default = $3, video_timeout_minutes = $4, log_retention_days = $5 
		WHERE id = 1`,
		cfg.DefaultModel, cfg.DefaultPreset, cfg.RemoveSdhDefault, cfg.VideoTimeoutMinutes, cfg.LogRetentionDays,
	)
	return err
}

func (s *ConfigService) GetConfig() (*models.UserConfig, error) {
	var cfg models.UserConfig
	err := s.DB.QueryRow(`
		SELECT default_model, default_preset, remove_sdh_default, video_timeout_minutes, log_retention_days 
		FROM user_config 
		LIMIT 1`).
		Scan(&cfg.DefaultModel, &cfg.DefaultPreset, &cfg.RemoveSdhDefault, &cfg.VideoTimeoutMinutes, &cfg.LogRetentionDays)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
