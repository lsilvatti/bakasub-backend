package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bakasub-backend/internal/utils"
)

type SubtitleTranslator interface {
	ProcessSubtitleFile(inputPath string, model string, outputPath string, apiKey string, targetLang string, preset string, removeSDH bool) error
}

type TranslateHandler struct {
	Translator SubtitleTranslator
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
	outputPath := filepath.Join(dir, fmt.Sprintf("%s_%s%s", base, reqData.TargetLang, ext))

	if err := h.Translator.ProcessSubtitleFile(inputPath, reqData.Model, outputPath, apiKey, reqData.TargetLang, reqData.Preset, reqData.RemoveSDH); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Processing failed: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Translation completed", map[string]string{
		"output_path": outputPath,
	})
}
