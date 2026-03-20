package services

import (
	"bakasub-backend/internal/models"
	"bakasub-backend/internal/parser"
	"database/sql"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type LLMProvider interface {
	TranslateText(text string, model string, apiKey string, targetLangName string, systemPrompt string) (string, error)
}

type TranslationFileSystemProvider interface {
	ReadFile(path string) (string, error)
	SaveFile(path string, content string) error
}

type TranslatorService struct {
	LLM LLMProvider
	FS  TranslationFileSystemProvider
	DB  *sql.DB
}

func NewTranslatorService(llm LLMProvider, fs TranslationFileSystemProvider, db *sql.DB) *TranslatorService {
	return &TranslatorService{
		LLM: llm,
		FS:  fs,
		DB:  db,
	}
}

var separatorRegex = regexp.MustCompile(`\s*---NEXT---\s*`)

func (s *TranslatorService) ProcessSubtitleFile(inputPath string, model string, outputPath string, apiKey string, targetLangCode string, presetAlias string, removeSDH bool) error {

	var systemPrompt string
	var maxChars int
	err := s.DB.QueryRow("SELECT system_prompt, batch_size FROM translation_presets WHERE alias = ?", presetAlias).Scan(&systemPrompt, &maxChars)
	if err != nil {
		return fmt.Errorf("error fetching preset '%s' from database: %w", presetAlias, err)
	}

	var targetLangName string
	err = s.DB.QueryRow("SELECT name FROM languages WHERE code = ?", targetLangCode).Scan(&targetLangName)
	if err != nil {
		return fmt.Errorf("error fetching language '%s' from database: %w", targetLangCode, err)
	}

	rawText, err := s.FS.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	var blocks []models.SubtitleBlock
	var assDoc *parser.ASSDocument
	var vttHeader string

	switch ext {
	case ".ass", ".ssa":
		assDoc, blocks = parser.ParseASS(rawText)
	case ".vtt":
		vttHeader, blocks = parser.ParseVTT(rawText)
	default:
		blocks = parser.ParseToBlocks(rawText) // Default: SRT
	}

	if len(blocks) == 0 {
		return fmt.Errorf("no subtitle found (or all lines were protected by heuristics)")
	}

	if removeSDH && ext != ".ass" {
		blocks = parser.RemoveSDH(blocks)
	}

	if len(blocks) == 0 {
		return fmt.Errorf("no valid subtitle found after filtering")
	}

	if apiKey == "" {
		return fmt.Errorf("apiKey is empty before starting translation")
	}

	const separator = "\n---NEXT---\n"
	separatorLen := len(separator)

	var batches [][]int
	var currentBatch []int
	currentChars := 0

	for i, block := range blocks {
		blockLen := len(block.Text)

		if currentChars+blockLen > maxChars && len(currentBatch) > 0 {
			batches = append(batches, currentBatch)
			currentBatch = nil
			currentChars = 0
		}

		currentBatch = append(currentBatch, i)
		currentChars += blockLen + separatorLen
	}

	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var translationErr error

	var limit = make(chan struct{}, 5)

	for _, batchIndices := range batches {
		wg.Add(1)
		limit <- struct{}{}

		go func(indices []int) {
			defer wg.Done()
			defer func() { <-limit }()

			mu.Lock()
			if translationErr != nil {
				mu.Unlock()
				return
			}
			mu.Unlock()

			var texts []string
			for _, idx := range indices {
				texts = append(texts, blocks[idx].Text)
			}

			joinedText := strings.Join(texts, separator)

			translatedText, err := s.LLM.TranslateText(joinedText, model, apiKey, targetLangName, systemPrompt)

			if err != nil {
				mu.Lock()
				if translationErr == nil {
					// Mostra o índice inicial em caso de erro para ajudar no debug
					translationErr = fmt.Errorf("error translating batch starting at line index %d: %w", indices[0], err)
				}
				mu.Unlock()
				return
			}

			translatedLines := separatorRegex.Split(translatedText, -1)

			if len(translatedLines) > 0 && strings.TrimSpace(translatedLines[len(translatedLines)-1]) == "" {
				translatedLines = translatedLines[:len(translatedLines)-1]
			}

			if len(translatedLines) != len(indices) {
				cleanText := separatorRegex.ReplaceAllString(translatedText, "\n\n")
				fallbackLines := strings.Split(strings.TrimSpace(cleanText), "\n\n")

				translatedLines = make([]string, len(indices))
				for i := range fallbackLines {
					if i < len(indices) {
						translatedLines[i] = fallbackLines[i]
					}
				}
			}

			mu.Lock()
			for i, idx := range indices {
				if i < len(translatedLines) {
					blocks[idx].Text = strings.TrimSpace(translatedLines[i])
				}
			}
			mu.Unlock()

		}(batchIndices)
	}

	wg.Wait()

	if translationErr != nil {
		return translationErr
	}

	var outputText string
	switch ext {
	case ".ass", ".ssa":
		outputText = parser.BuildASS(assDoc, blocks)
	case ".vtt":
		outputText = parser.BuildVTT(vttHeader, blocks)
	default:
		outputText = parser.BuildString(blocks)
	}

	if err := s.FS.SaveFile(outputPath, outputText); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
