package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/parser"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type LLMProvider interface {
	TranslateText(text string, model string, apiKey string, targetLang string, preset string) (string, error)
}

type TranslationFileSystemProvider interface {
	ReadFile(path string) (string, error)
	SaveFile(path string, content string) error
}

type TranslatorService struct {
	LLM LLMProvider
	FS  TranslationFileSystemProvider
}

func NewTranslatorService(llm LLMProvider, fs TranslationFileSystemProvider) *TranslatorService {
	return &TranslatorService{
		LLM: llm,
		FS:  fs,
	}
}

var separatorRegex = regexp.MustCompile(`\s*---NEXT---\s*`)

func (s *TranslatorService) ProcessSubtitleFile(inputPath string, model string, outputPath string, apiKey string, targetLang string, preset string, removeSDH bool) error {

	rawText, err := s.FS.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	blocks := parser.ParseToBlocks(rawText)

	if len(blocks) == 0 {
		return fmt.Errorf("nenhuma legenda encontrada")
	}

	if removeSDH {
		blocks = parser.RemoveSDH(blocks)
	}

	if len(blocks) == 0 {
		return fmt.Errorf("nenhuma legenda válida encontrada após filtragem")
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

	const separator = "\n---NEXT---\n"

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

			translatedText, err := s.LLM.TranslateText(joinedText, model, apiKey, targetLang, preset)

			if err != nil {
				mu.Lock()
				if translationErr == nil {
					translationErr = fmt.Errorf("erro ao traduzir lote %s: %w", batch[0].ID, err)
				}
				mu.Unlock()
				return
			}

			translatedLines := separatorRegex.Split(translatedText, -1)

			if len(translatedLines) > 0 && strings.TrimSpace(translatedLines[len(translatedLines)-1]) == "" {
				translatedLines = translatedLines[:len(translatedLines)-1]
			}

			if len(translatedLines) != len(batch) {
				cleanText := separatorRegex.ReplaceAllString(translatedText, "\n\n")
				fallbackLines := strings.Split(strings.TrimSpace(cleanText), "\n\n")
				translatedLines = make([]string, len(batch))
				for idx := range fallbackLines {
					if idx < len(batch) {
						translatedLines[idx] = fallbackLines[idx]
					}
				}
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

	if err := s.FS.SaveFile(outputPath, outputText); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}
