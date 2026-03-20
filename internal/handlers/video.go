package handlers

import (
	"net/http"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/services"
	"bakasub-backend/internal/utils"
)

type VideoProcessor interface {
	ScanSubtitles(videoPath string) ([]services.SubtitleTrack, error)
	ExtractSubtitle(videoPath string, subtitleId int) (string, error)
	MergeSubtitle(videoPath string, srtPath string, langCode string, timeoutMinutes int) (string, error)
}

type ConfigProvider interface {
	GetConfig() (*models.UserConfig, error)
}

type VideoHandler struct {
	Processor VideoProcessor
	Config    ConfigProvider
}

func (h *VideoHandler) GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		utils.Error(w, http.StatusBadRequest, "Missing 'path' query parameter")
		return
	}

	tracks, err := h.Processor.ScanSubtitles(path)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Error mapping file: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Tracks read successfully", map[string]interface{}{
		"tracks": tracks,
	})
}

func (h *VideoHandler) ExtractTrackHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[ExtractTrackRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	srtPath, err := h.Processor.ExtractSubtitle(reqData.VideoPath, reqData.SubtitleId)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Error extracting subtitle: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Subtitle extracted successfully", map[string]interface{}{
		"srtPath": srtPath,
	})
}

func (h *VideoHandler) MergeTrackHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[MergeTrackRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	timeout := 20
	if cfg, err := h.Config.GetConfig(); err == nil && cfg != nil {
		timeout = cfg.VideoTimeoutMinutes
	}

	outVideoPath, err := h.Processor.MergeSubtitle(reqData.VideoPath, reqData.SrtPath, reqData.LangCode, timeout)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Error merging subtitle: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Video generated successfully", map[string]interface{}{
		"outVideoPath": outVideoPath,
	})
}
