package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"bakasub-backend/internal/services"
)

type TranslateRequest struct {
	FilePath   string `json:"filePath"`
	TargetLang string `json:"targetLang"`
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		http.Error(w, "Erro de configuração da API", http.StatusInternalServerError)
		return
	}

	var reqData TranslateRequest
	err := json.NewDecoder(r.Body).Decode(&reqData)

	if err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if reqData.FilePath == "" {
		http.Error(w, "Caminho do arquivo é obrigatório", http.StatusBadRequest)
		return
	}

	if reqData.TargetLang == "" {
		http.Error(w, "Idioma de destino é obrigatório", http.StatusBadRequest)
		return
	}

	inputPath := reqData.FilePath
	dir := filepath.Dir(inputPath)
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), ext)

	outputPath := filepath.Join(dir, base+"_"+reqData.TargetLang+ext)

	err = services.ProcessSubtitleFile(inputPath, outputPath, apiKey)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Tradução salva em: %s", outputPath),
	})
}
