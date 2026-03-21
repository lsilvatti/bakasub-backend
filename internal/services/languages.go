package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/utils"
	"database/sql"
)

type LanguageService struct {
	DB *sql.DB
}

func NewLanguageService(db *sql.DB) *LanguageService {
	return &LanguageService{DB: db}
}

func (s *LanguageService) GetLanguages() ([]models.Language, error) {
	rows, err := s.DB.Query("SELECT id, code, name FROM languages")
	if err != nil {
		utils.LogError("languages", "Failed to query languages", map[string]any{"error": err.Error()})
		return nil, err
	}
	defer rows.Close()

	var languages []models.Language
	for rows.Next() {
		var lang models.Language
		if err := rows.Scan(&lang.ID, &lang.Code, &lang.Name); err != nil {
			utils.LogError("languages", "Failed to scan language row", map[string]any{"error": err.Error()})
			return nil, err
		}
		languages = append(languages, lang)
	}
	return languages, nil
}

func (s *LanguageService) GetLanguageByCode(code string) (*models.Language, error) {
	var lang models.Language
	err := s.DB.QueryRow("SELECT id, code, name FROM languages WHERE code = ?", code).Scan(&lang.ID, &lang.Code, &lang.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.LogError("languages", "Failed to fetch language by code", map[string]any{
				"code":  code,
				"error": err.Error(),
			})
		}
		return nil, err
	}
	return &lang, nil
}

func (s *LanguageService) AddLanguage(lang models.Language) error {
	_, err := s.DB.Exec("INSERT INTO languages (code, name) VALUES (?, ?)", lang.Code, lang.Name)
	if err != nil {
		utils.LogError("languages", "Failed to add language", map[string]any{
			"code":  lang.Code,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("languages", "create", "Language added successfully", map[string]any{
		"code": lang.Code,
		"name": lang.Name,
	})

	return nil
}

func (s *LanguageService) UpdateLanguage(lang models.Language) error {
	_, err := s.DB.Exec("UPDATE languages SET name = ? WHERE code = ?", lang.Name, lang.Code)
	if err != nil {
		utils.LogError("languages", "Failed to update language", map[string]any{
			"code":  lang.Code,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("languages", "update", "Language updated successfully", map[string]any{
		"code": lang.Code,
		"name": lang.Name,
	})

	return nil
}

func (s *LanguageService) DeleteLanguage(code string) error {
	_, err := s.DB.Exec("DELETE FROM languages WHERE code = ?", code)
	if err != nil {
		utils.LogError("languages", "Failed to delete language", map[string]any{
			"code":  code,
			"error": err.Error(),
		})
		return err
	}

	utils.LogInfo("languages", "delete", "Language deleted successfully", map[string]any{
		"code": code,
	})

	return nil
}
