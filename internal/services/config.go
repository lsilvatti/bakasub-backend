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

func (s *ConfigService) GetConfig() (models.UserConfig, error) {
	query := "SELECT default_model, default_preset, remove_sdh_default FROM user_configs WHERE id = 1"
	row := s.DB.QueryRow(query)

	var config models.UserConfig
	err := row.Scan(&config.DefaultModel, &config.DefaultPreset, &config.RemoveSdhDefault)

	if err != nil {
		return config, err
	}

	return config, nil
}

func (s *ConfigService) UpdateConfig(config models.UserConfig) error {
	query := `
	INSERT INTO user_configs (id, default_model, default_preset, remove_sdh_default)
	VALUES (1, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		default_model=excluded.default_model,
		default_preset=excluded.default_preset,
		remove_sdh_default=excluded.remove_sdh_default;
	`
	_, err := s.DB.Exec(query, config.DefaultModel, config.DefaultPreset, config.RemoveSdhDefault)
	return err
}
