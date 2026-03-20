package services

import (
	"bakasub-backend/internal/models"
	"database/sql"
)

type PresetService struct {
	DB *sql.DB
}

func NewPresetService(db *sql.DB) *PresetService {
	return &PresetService{DB: db}
}

func (s *PresetService) GetPresets() ([]models.TranslationPreset, error) {
	rows, err := s.DB.Query("SELECT id, alias, name, system_prompt, batch_size, temperature FROM translation_presets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presets []models.TranslationPreset
	for rows.Next() {
		var preset models.TranslationPreset
		if err := rows.Scan(&preset.ID, &preset.Alias, &preset.Name, &preset.SystemPrompt, &preset.BatchSize, &preset.Temperature); err != nil {
			return nil, err
		}
		presets = append(presets, preset)
	}

	return presets, nil
}

func (s *PresetService) GetPresetByAlias(alias string) (*models.TranslationPreset, error) {
	row := s.DB.QueryRow("SELECT id, alias, name, system_prompt, batch_size, temperature FROM translation_presets WHERE alias = ?", alias)

	var preset models.TranslationPreset
	err := row.Scan(&preset.ID, &preset.Alias, &preset.Name, &preset.SystemPrompt, &preset.BatchSize, &preset.Temperature)
	if err != nil {
		return nil, err
	}

	return &preset, nil
}

func (s *PresetService) CreatePreset(preset models.TranslationPreset) error {
	_, err := s.DB.Exec("INSERT INTO translation_presets (alias, name, system_prompt, batch_size, temperature) VALUES (?, ?, ?, ?, ?)",
		preset.Alias, preset.Name, preset.SystemPrompt, preset.BatchSize, preset.Temperature)
	return err
}

func (s *PresetService) UpdatePreset(preset models.TranslationPreset) error {
	_, err := s.DB.Exec("UPDATE translation_presets SET alias = ?, name = ?, system_prompt = ?, batch_size = ?, temperature = ? WHERE id = ?",
		preset.Alias, preset.Name, preset.SystemPrompt, preset.BatchSize, preset.Temperature, preset.ID)
	return err
}

func (s *PresetService) DeletePreset(id int) error {
	_, err := s.DB.Exec("DELETE FROM translation_presets WHERE id = ?", id)
	return err
}
