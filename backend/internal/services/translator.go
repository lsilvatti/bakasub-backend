package services

import (
	"bakasub-backend/internal/ai"
	"bakasub-backend/internal/fileio"
	"bakasub-backend/internal/parser"
	"fmt"
)

func ProcessSubtitleFile(inputPath string, outputPath string, apiKey string) error {
	rawText, err := fileio.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	blocks := parser.ParseToBlocks(rawText)
	if len(blocks) == 0 {
		return fmt.Errorf("nenhuma legenda encontrada")
	}

	for i := range blocks {
		translatedText, err := ai.TranslateText(blocks[i].Text, "google/gemini-2.5-flash", apiKey)
		if err != nil {
			return fmt.Errorf("erro ao traduzir bloco %s: %w", blocks[i].ID, err)
		}
		blocks[i].Text = translatedText
	}

	finalText := parser.BuildString(blocks)

	err = fileio.SaveFile(outputPath, finalText)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo final: %w", err)
	}

	return nil
}
