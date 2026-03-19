package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"

	"bakasub-backend/internal/services"
	"bakasub-backend/internal/utils"
)

type TranslateRequest struct {
	FilePath   string `json:"filePath" validate:"required"`
	TargetLang string `json:"targetLang" validate:"required"`
	Preset     string `json:"preset" validate:"required"`
	Model      string `json:"model" validate:"required"`
}

var validate = validator.New()

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.Error(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var reqData TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := validate.Struct(reqData); err != nil {
		utils.Error(w, http.StatusBadRequest, "Campos obrigatórios ausentes: "+err.Error())
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		utils.Error(w, http.StatusInternalServerError, "Configuração de API ausente")
		return
	}

	inputPath := reqData.FilePath
	dir := filepath.Dir(inputPath)
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), ext)
	outputPath := filepath.Join(dir, fmt.Sprintf("%s_%s%s", base, reqData.TargetLang, ext))

	if err := services.ProcessSubtitleFile(inputPath, reqData.Model, outputPath, apiKey, reqData.TargetLang, reqData.Preset); err != nil {
		utils.Error(w, http.StatusInternalServerError, "Falha no processamento: "+err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "success", "Tradução concluída", map[string]string{
		"output_path": outputPath,
	})
}
