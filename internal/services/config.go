package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"database/sql"
)

type ConfigService struct {
	DB        *sql.DB
	SecretKey string
}

func NewConfigService(db *sql.DB, secretKey string) *ConfigService {
	return &ConfigService{DB: db, SecretKey: secretKey}
}

func (s *ConfigService) UpdateConfig(cfg models.UserConfig) error {
	encOpenRouter, err := utils.Encrypt(cfg.OpenRouterApiKey, s.SecretKey)
	if err != nil {
		return err
	}
	encTmdb, err := utils.Encrypt(cfg.TmdbAccessToken, s.SecretKey)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`
		UPDATE user_config
		SET default_model = ?, default_preset = ?, remove_sdh_default = ?,
		    video_timeout_minutes = ?, log_retention_days = ?, default_language = ?,
		    openrouter_api_key = ?, tmdb_access_token = ?,
		    concurrent_translations = ?, max_retries = ?, base_retry_delay = ?,
		    tmdb_metadata_enabled = ?
		WHERE id = 1`,
		cfg.DefaultModel, cfg.DefaultPreset, cfg.RemoveSdhDefault,
		cfg.VideoTimeoutMinutes, cfg.LogRetentionDays, cfg.DefaultLanguage,
		encOpenRouter, encTmdb,
		cfg.ConcurrentTranslations, cfg.MaxRetries, cfg.BaseRetryDelay,
		cfg.TmdbMetadataEnabled,
	)
	return err
}

func (s *ConfigService) GetConfig() (*models.UserConfig, error) {
	var cfg models.UserConfig
	var encOpenRouter, encTmdb string

	err := s.DB.QueryRow(`
		SELECT default_model, default_preset, remove_sdh_default, video_timeout_minutes,
		       log_retention_days, default_language, openrouter_api_key, tmdb_access_token,
		       concurrent_translations, max_retries, base_retry_delay, tmdb_metadata_enabled
		FROM user_config
		LIMIT 1`).
		Scan(
			&cfg.DefaultModel, &cfg.DefaultPreset, &cfg.RemoveSdhDefault,
			&cfg.VideoTimeoutMinutes, &cfg.LogRetentionDays, &cfg.DefaultLanguage,
			&encOpenRouter, &encTmdb,
			&cfg.ConcurrentTranslations, &cfg.MaxRetries, &cfg.BaseRetryDelay,
			&cfg.TmdbMetadataEnabled,
		)
	if err != nil {
		return nil, err
	}

	cfg.OpenRouterApiKey, err = utils.Decrypt(encOpenRouter, s.SecretKey)
	if err != nil {
		return nil, err
	}
	cfg.TmdbAccessToken, err = utils.Decrypt(encTmdb, s.SecretKey)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
