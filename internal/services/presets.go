package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
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
		utils.LogError("presets", "Failed to fetch presets from database", map[string]any{"error": err.Error()})
		return nil, err
	}
	defer rows.Close()

	var presets []models.TranslationPreset
	for rows.Next() {
		var preset models.TranslationPreset
		if err := rows.Scan(&preset.ID, &preset.Alias, &preset.Name, &preset.SystemPrompt, &preset.BatchSize, &preset.Temperature); err != nil {
			utils.LogError("presets", "Failed to scan preset row", map[string]any{"error": err.Error()})
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
		// Só logamos se o erro for diferente de "não encontrado", pois buscar algo que não existe é normal.
		if err != sql.ErrNoRows {
			utils.LogError("presets", "Failed to fetch preset by alias", map[string]any{
				"alias": alias,
				"error": err.Error(),
			})
		}
		return nil, err
	}

	return &preset, nil
}

func (s *PresetService) CreatePreset(preset models.TranslationPreset) error {
	_, err := s.DB.Exec("INSERT INTO translation_presets (alias, name, system_prompt, batch_size, temperature) VALUES (?, ?, ?, ?, ?)",
		preset.Alias, preset.Name, preset.SystemPrompt, preset.BatchSize, preset.Temperature)

	if err != nil {
		utils.LogError("presets", "Failed to create translation preset", map[string]any{
			"alias": preset.Alias,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("presets", "create", "Translation preset created", map[string]any{
		"alias": preset.Alias,
		"name":  preset.Name,
	})

	return nil
}

func (s *PresetService) UpdatePreset(preset models.TranslationPreset) error {
	_, err := s.DB.Exec("UPDATE translation_presets SET alias = ?, name = ?, system_prompt = ?, batch_size = ?, temperature = ? WHERE id = ?",
		preset.Alias, preset.Name, preset.SystemPrompt, preset.BatchSize, preset.Temperature, preset.ID)

	if err != nil {
		utils.LogError("presets", "Failed to update translation preset", map[string]any{
			"id":    preset.ID,
			"alias": preset.Alias,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("presets", "update", "Translation preset updated", map[string]any{
		"id":    preset.ID,
		"alias": preset.Alias,
	})

	return nil
}

func (s *PresetService) DeletePreset(id int) error {
	_, err := s.DB.Exec("DELETE FROM translation_presets WHERE id = ?", id)

	if err != nil {
		utils.LogError("presets", "Failed to delete translation preset", map[string]any{
			"id":    id,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("presets", "delete", "Translation preset deleted", map[string]any{
		"id": id,
	})

	return nil
}
