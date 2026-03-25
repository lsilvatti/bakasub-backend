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

func (s *PresetService) CreatePreset(preset models.TranslationPreset) error {
	_, err := s.DB.Exec(`
		INSERT INTO translation_presets (alias, name, system_prompt, batch_size, temperature) 
		VALUES ($1, $2, $3, $4, $5)`,
		preset.Alias, preset.Name, preset.SystemPrompt, preset.BatchSize, preset.Temperature,
	)
	return err
}

func (s *PresetService) UpdatePreset(preset models.TranslationPreset) error {
	_, err := s.DB.Exec(`
		UPDATE translation_presets 
		SET alias = COALESCE(NULLIF($1, ''), alias),
		    name = COALESCE(NULLIF($2, ''), name),
		    system_prompt = COALESCE(NULLIF($3, ''), system_prompt),
		    batch_size = COALESCE(NULLIF($4, 0), batch_size),
		    temperature = COALESCE(NULLIF($5, 0), temperature)
		WHERE id = $6`,
		preset.Alias, preset.Name, preset.SystemPrompt, preset.BatchSize, preset.Temperature, preset.ID,
	)
	return err
}

func (s *PresetService) GetPresets() ([]models.TranslationPreset, error) {
	rows, err := s.DB.Query("SELECT id, alias, name, system_prompt, batch_size, temperature FROM translation_presets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var presets []models.TranslationPreset
	for rows.Next() {
		var p models.TranslationPreset
		if err := rows.Scan(&p.ID, &p.Alias, &p.Name, &p.SystemPrompt, &p.BatchSize, &p.Temperature); err != nil {
			return nil, err
		}
		presets = append(presets, p)
	}

	if presets == nil {
		presets = []models.TranslationPreset{}
	}

	return presets, nil
}

func (s *PresetService) DeletePreset(id int) error {
	_, err := s.DB.Exec("DELETE FROM translation_presets WHERE id = $1", id)
	return err
}
