package handlers

import (
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
	reqData, err := utils.DecodeAndValidate[GetTrackRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Dados inválidos: "+err.Error())
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
	reqData, err := utils.DecodeAndValidate[ExtractTrackRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Dados inválidos: "+err.Error())
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
	reqData, err := utils.DecodeAndValidate[MergeTrackRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Dados inválidos: "+err.Error())
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
