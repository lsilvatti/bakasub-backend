package handlers

import (
	"encoding/json"
	"net/http"

	"bakasub-backend/internal/services"
	"bakasub-backend/internal/utils"
)

type VideoProcessor interface {
	ScanSubtitles(videoPath string) ([]services.SubtitleTrack, error)
	ExtractSubtitle(videoPath string, subtitleId int) (string, error)
	MergeSubtitle(videoPath string, srtPath string, langCode string) (string, error)
}

type VideoHandler struct {
	Processor VideoProcessor
}

func (h *VideoHandler) GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	var reqData GetTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := utils.Validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	tracks, err := h.Processor.ScanSubtitles(reqData.VideoPath)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao mapear arquivo: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Trilhas lidas com sucesso", map[string]interface{}{
		"tracks": tracks,
	})
}

func (h *VideoHandler) ExtractTrackHandler(w http.ResponseWriter, r *http.Request) {
	var reqData ExtractTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := utils.Validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	srtPath, err := h.Processor.ExtractSubtitle(reqData.VideoPath, reqData.SubtitleId)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao extrair legenda: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Legenda extraída com sucesso", map[string]interface{}{
		"srtPath": srtPath,
	})
}

func (h *VideoHandler) MergeTrackHandler(w http.ResponseWriter, r *http.Request) {
	var reqData MergeTrackRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := utils.Validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	outVideoPath, err := h.Processor.MergeSubtitle(reqData.VideoPath, reqData.SrtPath, reqData.LangCode)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Erro ao mesclar legenda: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Vídeo gerado com sucesso", map[string]interface{}{
		"outVideoPath": outVideoPath,
	})
}
