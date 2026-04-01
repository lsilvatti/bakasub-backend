package handlers

import (
	"net/http"

	"bakasub-backend/internal/utils"
)

type FavoritesService interface {
	GetFavorites() ([]string, error)
	UpdateFavorites(favorites []string) error
}

type FavoritesHandler struct {
	Service FavoritesService
}

func (h *FavoritesHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	favorites, err := h.Service.GetFavorites()

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to retrieve favorite models: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Favorite models retrieved successfully", favorites)
}

func (h *FavoritesHandler) UpdateFavorites(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[UpdateFavoritesRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	if err := h.Service.UpdateFavorites(reqData.FavoriteModels); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update favorite models: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Favorite models updated successfully", nil)
}