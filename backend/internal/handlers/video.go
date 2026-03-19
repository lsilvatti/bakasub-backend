package handlers

import (
	"encoding/json"
	"net/http"

	"bakasub-backend/internal/utils"
	"bakasub-backend/internal/video"
)

func GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var reqData GetTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := utils.Validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	tracks, err := video.ScanSubtitles(reqData.VideoPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao mapear arquivo: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Trilhas lidas com sucesso", map[string]interface{}{
		"tracks": tracks,
	})
}

func ExtractTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var reqData ExtractTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := utils.Validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	srtPath, err := video.ExtractSubtitle(reqData.VideoPath, reqData.SubtitleId)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao extrair legenda: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Legenda extraída com sucesso", map[string]interface{}{
		"srtPath": srtPath,
	})
}

func MergeTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var reqData MergeTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := utils.Validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	outVideoPath, err := video.MergeSubtitle(reqData.VideoPath, reqData.SrtPath, reqData.LangCode)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao mesclar legenda: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Vídeo gerado com sucesso", map[string]interface{}{
		"outVideoPath": outVideoPath,
	})
}
