package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"bakasub-backend/internal/utils"
)

type SubtitleTranslator interface {
	ProcessSubtitleFile(inputPath string, model string, outputPath string, apiKey string, targetLang string, preset string, removeSDH bool, context string) error
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
		utils.LogError("translate_handler", "Invalid translation request payload", map[string]any{
			"error": err.Error(),
		})
		utils.Error(w, http.StatusBadRequest, "Invalid data: "+err.Error())
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		utils.LogError("translate_handler", "OPENROUTER_API_KEY environment variable is missing", nil)
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

	utils.LogInfo("translate_handler", "info", "Translation request received", map[string]any{
		"input":  inputPath,
		"model":  reqData.Model,
		"target": reqData.TargetLang,
	})

	if err := h.Translator.ProcessSubtitleFile(inputPath, reqData.Model, outputPath, apiKey, reqData.TargetLang, reqData.Preset, reqData.RemoveSDH, reqData.Context); err != nil {
		utils.LogError("translate_handler", "Subtitle translation processing failed", map[string]any{
			"input": inputPath,
			"error": err.Error(),
		})
		utils.Error(w, http.StatusInternalServerError, "Processing failed: "+err.Error())
		return
	}

	utils.LogInfo("translate_handler", "success", "Translation request processed successfully", map[string]any{
		"output": outputPath,
	})

	utils.JSON(w, http.StatusOK, "success", "Translation completed", map[string]string{
		"output_path": outputPath,
	})
}
