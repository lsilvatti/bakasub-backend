package services

import (
	"database/sql"
	"encoding/json"
)

type FavoritesService struct {
	DB *sql.DB
}

func NewFavoritesService(db *sql.DB) *FavoritesService {
	return &FavoritesService{DB: db}
}

func (s *FavoritesService) GetFavorites() ([]string, error) {
	var favModelsJSON []byte

	err := s.DB.QueryRow(`
		SELECT COALESCE(favorite_models, '[]')
		FROM user_config
		LIMIT 1`).Scan(&favModelsJSON)

	if err != nil {
		return []string{}, err
	}

	var favorites []string
	if len(favModelsJSON) > 0 {
		json.Unmarshal(favModelsJSON, &favorites)
	}

	if favorites == nil {
		favorites = []string{}
	}

	return favorites, nil
}

func (s *FavoritesService) UpdateFavorites(favorites []string) error {
	if favorites == nil {
		favorites = []string{}
	}

	favJSON, err := json.Marshal(favorites)
	if err != nil {
		favJSON = []byte("[]")
	}

	_, err = s.DB.Exec(`UPDATE user_config SET favorite_models = $1 WHERE id = 1`, favJSON)
	return err
}
