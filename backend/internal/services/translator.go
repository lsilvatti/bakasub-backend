package services

import (
	"bakasub-backend/internal/ai"
	"bakasub-backend/internal/fileio"
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/parser"
	"fmt"
	"strings"
	"sync"
)

func ProcessSubtitleFile(inputPath string, model string, outputPath string, apiKey string, targetLang string, preset string) error {

	rawText, err := fileio.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	blocks := parser.ParseToBlocks(rawText)
	if len(blocks) == 0 {
		return fmt.Errorf("nenhuma legenda encontrada")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var translationErr error

	batchSize := models.Presets[preset].BatchSize
	if batchSize < 1 {
		batchSize = 1
	}

	var limit = make(chan struct{}, 5)

	if apiKey == "" {
		return fmt.Errorf("apiKey is empty before starting translation")
	}

	const separator = "\n---SEP---\n"

	for i := 0; i < len(blocks); i += batchSize {
		end := i + batchSize

		if end > len(blocks) {
			end = len(blocks)
		}

		currentBatch := blocks[i:end]

		wg.Add(1)
		limit <- struct{}{}

		go func(batch []models.SubtitleBlock) {
			defer wg.Done()
			defer func() { <-limit }()

			mu.Lock()
			if translationErr != nil {
				mu.Unlock()
				return
			}
			mu.Unlock()

			var texts []string
			for _, block := range batch {
				texts = append(texts, block.Text)
			}

			joinedText := strings.Join(texts, separator)

			translatedText, err := ai.TranslateText(joinedText, model, apiKey, targetLang, preset)

			if err != nil {
				mu.Lock()
				if translationErr == nil {
					translationErr = fmt.Errorf("erro ao traduzir bloco %s: %w", batch[0].ID, err)
				}
				mu.Unlock()
				return
			}

			translatedLines := strings.Split(translatedText, separator)

			if len(translatedLines) != len(batch) {
				translatedLines = strings.Split(strings.TrimSpace(translatedText), "\n\n")
			}

			mu.Lock()
			for i := range batch {
				if i < len(translatedLines) {
					batch[i].Text = strings.TrimSpace(translatedLines[i])
				}
			}
			mu.Unlock()

		}(currentBatch)
	}

	wg.Wait()

	if translationErr != nil {
		return translationErr
	}

	outputText := parser.BuildString(blocks)

	if err := fileio.SaveFile(outputPath, outputText); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}
