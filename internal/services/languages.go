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

func (s *LanguageService) CreateLanguage(lang models.Language) error {
	_, err := s.DB.Exec("INSERT INTO languages (code, name) VALUES (?, ?)", lang.Code, lang.Name)
	return err
}

func (s *LanguageService) GetLanguages() ([]models.Language, error) {
	rows, err := s.DB.Query("SELECT code, name FROM languages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []models.Language
	for rows.Next() {
		var l models.Language
		if err := rows.Scan(&l.Code, &l.Name); err != nil {
			return nil, err
		}
		languages = append(languages, l)
	}

	if languages == nil {
		languages = []models.Language{}
	}

	return languages, nil
}

func (s *LanguageService) UpdateLanguage(lang models.Language) error {
	_, err := s.DB.Exec("UPDATE languages SET name = ? WHERE code = ?", lang.Name, lang.Code)
	return err
}

func (s *LanguageService) DeleteLanguage(code string) error {
	_, err := s.DB.Exec("DELETE FROM languages WHERE code = ?", code)
	return err
}
