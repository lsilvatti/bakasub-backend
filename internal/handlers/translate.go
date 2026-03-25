package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"bakasub-backend/internal/models"
	"bakasub-backend/internal/services"
	"bakasub-backend/internal/utils"

	"github.com/google/uuid"
)

type SubtitleTranslator interface {
	ProcessSubtitleFile(jobID string, inputPath string, model string, outputPath string, apiKey string, targetLang string, preset string, removeSDH bool, context string) error
	PreFlight(inputPath string, model string, targetLang string, preset string, removeSDH bool) (*models.JobEstimate, error)
}

type TranslateHandler struct {
	Translator SubtitleTranslator
	JobService *services.JobService
}

func (h *TranslateHandler) PreFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	reqData, err := utils.DecodeAndValidate[PreFlightRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	estimate, err := h.Translator.PreFlight(reqData.FilePath, reqData.Model, reqData.TargetLang, reqData.Preset, reqData.RemoveSDH)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Analysis failed: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Pre-flight completed", estimate)
}

func (h *TranslateHandler) Translate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	reqData, err := utils.DecodeAndValidate[TranslateRequest](r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		utils.Error(w, http.StatusInternalServerError, "Missing API configuration")
		return
	}

	inputPath := reqData.FilePath
	dir := filepath.Dir(inputPath)
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), ext)
	langSuffixRegex := regexp.MustCompile(`(?i)_([a-z]{2,3}(-[a-z]{2,3})?)$`)
	base = langSuffixRegex.ReplaceAllString(base, "")
	outputPath := filepath.Join(dir, fmt.Sprintf("%s_%s%s", base, reqData.TargetLang, ext))

	jobID := uuid.New().String()
	err = h.JobService.CreateJob(jobID, reqData.FilePath, reqData.TargetLang, reqData.Preset, reqData.Model)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to initialize translation job")
		return
	}

	go func() {
		err := h.Translator.ProcessSubtitleFile(jobID, inputPath, reqData.Model, outputPath, apiKey, reqData.TargetLang, reqData.Preset, reqData.RemoveSDH, reqData.Context)
		if err != nil {
			h.JobService.UpdateStatus(jobID, "failed", err.Error())
		} else {
			h.JobService.UpdateStatus(jobID, "completed", "")
		}
	}()

	utils.JSON(w, http.StatusOK, "success", "Translation job started", map[string]string{
		"job_id":      jobID,
		"output_path": outputPath,
	})
}
