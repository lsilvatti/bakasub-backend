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
		utils.LogError("video_handler", "Missing 'path' query parameter for GetTrack", nil)
		utils.Error(w, http.StatusBadRequest, "Missing 'path' query parameter")
		return
	}

	tracks, err := h.Processor.ScanSubtitles(path)
	if err != nil {
		utils.LogError("video_handler", "Failed to scan subtitles", map[string]any{
			"path":  path,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error mapping file: "+err.Error())
		return
	}

	utils.LogInfo("video_handler", "success", "Successfully scanned subtitle tracks", map[string]any{
		"path":        path,
		"track_count": len(tracks),
	})

	utils.JSON(w, http.StatusOK, "success", "Tracks read successfully", map[string]interface{}{
		"tracks": tracks,
	})
}

func (h *VideoHandler) ExtractTrackHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[ExtractTrackRequest](r)
	if err != nil {
		utils.LogError("video_handler", "Invalid payload for ExtractTrack", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	srtPath, err := h.Processor.ExtractSubtitle(reqData.VideoPath, reqData.SubtitleId)
	if err != nil {
		utils.LogError("video_handler", "Failed to extract subtitle via processor", map[string]any{
			"videoPath":  reqData.VideoPath,
			"subtitleId": reqData.SubtitleId,
			"error":      err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error extracting subtitle: "+err.Error())
		return
	}

	utils.LogInfo("video_handler", "success", "Successfully extracted subtitle track", map[string]any{
		"videoPath": reqData.VideoPath,
		"srtPath":   srtPath,
	})

	utils.JSON(w, http.StatusOK, "success", "Subtitle extracted successfully", map[string]interface{}{
		"srtPath": srtPath,
	})
}

func (h *VideoHandler) MergeTrackHandler(w http.ResponseWriter, r *http.Request) {
	reqData, err := utils.DecodeAndValidate[MergeTrackRequest](r)
	if err != nil {
		utils.LogError("video_handler", "Invalid payload for MergeTrack", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	timeout := 20
	if cfg, err := h.Config.GetConfig(); err == nil && cfg != nil {
		timeout = cfg.VideoTimeoutMinutes
	} else {
		utils.LogInfo("video_handler", "warning", "Failed to fetch user config for timeout, using default 20m", map[string]any{
			"error": err.Error(),
		})
	}

	outVideoPath, err := h.Processor.MergeSubtitle(reqData.VideoPath, reqData.SrtPath, reqData.LangCode, timeout)
	if err != nil {
		utils.LogError("video_handler", "Failed to merge subtitle via processor", map[string]any{
			"videoPath": reqData.VideoPath,
			"srtPath":   reqData.SrtPath,
			"error":     err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Error merging subtitle: "+err.Error())
		return
	}

	utils.LogInfo("video_handler", "success", "Successfully merged subtitle track into video", map[string]any{
		"videoPath":    reqData.VideoPath,
		"outVideoPath": outVideoPath,
	})

	utils.JSON(w, http.StatusOK, "success", "Video generated successfully", map[string]interface{}{
		"outVideoPath": outVideoPath,
	})
}
