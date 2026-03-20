package services

import (
	"bakasub-backend/internal/models"
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
		return nil, err
	}
	defer rows.Close()

	var languages []models.Language
	for rows.Next() {
		var lang models.Language
		if err := rows.Scan(&lang.ID, &lang.Code, &lang.Name); err != nil {
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
		return nil, err
	}
	return &lang, nil
}

func (s *LanguageService) AddLanguage(lang models.Language) error {
	_, err := s.DB.Exec("INSERT INTO languages (code, name) VALUES (?, ?)", lang.Code, lang.Name)
	return err
}

func (s *LanguageService) UpdateLanguage(lang models.Language) error {
	_, err := s.DB.Exec("UPDATE languages SET name = ? WHERE code = ?", lang.Name, lang.Code)
	return err
}

func (s *LanguageService) DeleteLanguage(code string) error {
	_, err := s.DB.Exec("DELETE FROM languages WHERE code = ?", code)
	return err
}
